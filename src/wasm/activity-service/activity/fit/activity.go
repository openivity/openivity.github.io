package fit

import (
	"time"

	"github.com/muktihari/fit/factory"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
)

func createActivityMesg(timestamp time.Time, timezone int, sessionCount uint16) proto.Message {
	mesg := factory.CreateMesgOnly(mesgnum.Activity)

	activityTypeField := factory.CreateField(mesgnum.Activity, fieldnum.ActivityType)
	activityTypeField.Value = uint8(typedef.ActivityAutoMultiSport)
	mesg.Fields = append(mesg.Fields, activityTypeField)

	timestampField := factory.CreateField(mesgnum.Activity, fieldnum.ActivityTimestamp)
	timestampField.Value = datetime.ToUint32(timestamp)
	mesg.Fields = append(mesg.Fields, timestampField)

	localTimestamp := datetime.ToLocalTime(timestampField.Value, timezone)
	localTimestampField := factory.CreateField(mesgnum.Activity, fieldnum.ActivityLocalTimestamp)
	localTimestampField.Value = datetime.ToUint32(localTimestamp)
	mesg.Fields = append(mesg.Fields, localTimestampField)

	sessionCountField := factory.CreateField(mesgnum.Activity, fieldnum.ActivityNumSessions)
	sessionCountField.Value = sessionCount

	return mesg
}
