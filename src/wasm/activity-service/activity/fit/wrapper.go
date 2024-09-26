package fit

import (
	"github.com/muktihari/fit/factory"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
)

// wrapActivity wraps filedef.Activity so we can manipulate the messages according to our need.
// To maintain the messages arrival order, we injects timestamp for messages that does not have
// timestamp definition or it does have timestamp definition but it simply lack the timestamp
// with the previous message's timestamp.
//
// This injected timestamp (field) will be marked as ExpandedField so we can remove it before converting into
// proto.FIT. This way, the original messages is preserve with its arrival order as close as it was.
// There maybe slightly different order but those messages have same timestamp so it does not matter.
type wrapActivity struct {
	activity           *filedef.Activity
	timestampReference uint32
}

func (w *wrapActivity) Add(mesg proto.Message) {
	// Add missing timestamp to any message so we can restructure them as arrival order.
	if mesg.Num == mesgnum.FileId {
		w.timestampReference = mesg.FieldValueByNum(fieldnum.FileIdTimeCreated).Uint32()
	}

	if field := mesg.FieldByNum(proto.FieldNumTimestamp); field != nil {
		if field.Value.Uint32() == basetype.Uint32Invalid {
			field.Value = proto.Uint32(w.timestampReference) // Replace with previous timestamp reference
		} else {
			w.timestampReference = field.Value.Uint32()
		}
	} else {
		field := factory.CreateField(mesgnum.Session, proto.FieldNumTimestamp).WithValue(w.timestampReference)
		field.IsExpandedField = true // Mark for removal on ToFIT()
		mesg.Fields = append(mesg.Fields, field)
	}

	switch mesg.Num {
	case mesgnum.FileId:
		w.activity.FileId = *mesgdef.NewFileId(&mesg)
	case mesgnum.Sport:
		return // Let's discard this message, we will recreate them from Sessions.
	case mesgnum.Record:
		w.activity.Records = append(w.activity.Records, mesgdef.NewRecord(&mesg))
	case mesgnum.Lap:
		w.activity.Laps = append(w.activity.Laps, mesgdef.NewLap(&mesg))
	case mesgnum.SplitSummary:
		w.activity.SplitSummaries = append(w.activity.SplitSummaries, mesgdef.NewSplitSummary(&mesg))
	case mesgnum.Session:
		w.activity.Sessions = append(w.activity.Sessions, mesgdef.NewSession(&mesg))
	case mesgnum.Activity:
		w.activity.Activity = mesgdef.NewActivity(&mesg)
	default:
		mesg.Fields = append(mesg.Fields[:0:0], mesg.Fields...) // Clone
		w.activity.UnrelatedMessages = append(w.activity.UnrelatedMessages, mesg)
	}
}

func (w *wrapActivity) ToFIT(o *mesgdef.Options) proto.FIT {
	size := 2 + len(w.activity.Records) + len(w.activity.Laps) + len(w.activity.Sessions) +
		len(w.activity.SplitSummaries) + len(w.activity.Sports)

	fit := proto.FIT{Messages: make([]proto.Message, 0, size)}

	fit.Messages = append(fit.Messages, w.activity.FileId.ToMesg(o))

	// Append unrelated messages first, since there are messages that have the same timestamp
	// with first record's timestamp, but it might arrive before records. We want records order
	// to be placed after those messages. We use stable sort so the order will be preserved.
	fit.Messages = append(fit.Messages, w.activity.UnrelatedMessages...)

	// Create Sports from Sessions, add additional timestamp of the first record timestamp
	// so Sports will be put before Records, this field will be removed after sorting.
	w.activity.Sports = make([]*mesgdef.Sport, 0, len(w.activity.Sessions))
	var firstRecordTimestampField proto.Field
	for i := range w.activity.Records {
		if timestamp := datetime.ToUint32(w.activity.Records[i].Timestamp); timestamp != basetype.Uint32Invalid {
			firstRecordTimestampField = factory.CreateField(mesgnum.Record, fieldnum.RecordTimestamp).WithValue(timestamp)
			firstRecordTimestampField.IsExpandedField = true // Mark for removal after sorting.
			break
		}
	}

	for _, ses := range w.activity.Sessions {
		var ok bool
		for _, v := range w.activity.Sports {
			if v.Sport == ses.Sport {
				ok = true
				break
			}
		}
		if !ok {
			fit.Messages = append(fit.Messages, proto.Message{Num: mesgnum.Sport, Fields: []proto.Field{
				firstRecordTimestampField,
				factory.CreateField(mesgnum.Sport, fieldnum.SportSport).WithValue(ses.Sport.Byte()),
				factory.CreateField(mesgnum.Sport, fieldnum.SportSubSport).WithValue(ses.SubSport.Byte()),
				factory.CreateField(mesgnum.Sport, fieldnum.SportName).WithValue(ses.SportProfileName),
			}})
		}
	}

	for i := range w.activity.Records {
		fit.Messages = append(fit.Messages, w.activity.Records[i].ToMesg(o))
	}

	for i := range w.activity.Laps {
		fit.Messages = append(fit.Messages, w.activity.Laps[i].ToMesg(o))
	}

	// We use "Last Summary Order" pattern, so activity will have the latest timestamp.
	latestTimestampField := factory.CreateField(mesgnum.Session, fieldnum.SessionTimestamp).
		WithValue(datetime.ToUint32(w.activity.Activity.Timestamp))

	for i := range w.activity.SplitSummaries {
		// Split Summary actually does not have defined timestamp in Profile.xlsx, but Garmin Connect
		// enforce us to create one, otherwise, our file will be rejected.
		// Discussion: https://forums.garmin.com/developer/fit-sdk/f/discussion/385625/timestamp-field-in-split_summary-messages
		mesg := w.activity.SplitSummaries[i].ToMesg(o)
		mesg.RemoveFieldByNum(proto.FieldNumTimestamp)
		mesg.Fields = append(mesg.Fields, proto.Field{})
		copy(mesg.Fields[1:], mesg.Fields)
		mesg.Fields[0] = latestTimestampField
		fit.Messages = append(fit.Messages, mesg)
	}

	for i := range w.activity.Sessions {
		fit.Messages = append(fit.Messages, w.activity.Sessions[i].ToMesg(o))
	}

	fit.Messages = append(fit.Messages, w.activity.Activity.ToMesg(o))

	filedef.SortMessagesByTimestamp(fit.Messages[1:])

	// Since messages has been sorted, let's remove the additional timestamp field
	// (the one that we have marked with ExpandedField == true) from all messages.
	for i := range fit.Messages {
		mesg := &fit.Messages[i]
		var valid int
		for j := range mesg.Fields {
			field := &mesg.Fields[j]
			if field.Num == proto.FieldNumTimestamp && field.IsExpandedField {
				continue
			}
			if j != valid { // memmove
				mesg.Fields[j], mesg.Fields[valid] = mesg.Fields[valid], mesg.Fields[j]
			}
			valid++
		}
		mesg.Fields = mesg.Fields[:valid]
	}

	return fit
}
