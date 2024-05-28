// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package activity

import (
	"strconv"
	"time"

	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/proto"
)

// Activity is an activity. It use FIT SDK's structure as its base since FIT is currently the most advance format.
type Activity struct {
	Creator  Creator
	Timezone int8
	Sessions []Session

	// UnrelatedMessages contains all messages not used by our service
	// such as DeveloperDataIds, FieldDescriptions, Events, etc.
	// We will restore these messages as it is when we recreate the FIT files.
	UnrelatedMessages []proto.Message
}

// CreateActivity creates new activity.
func CreateActivity() Activity {
	return Activity{
		Creator: CreateCreator(nil),
	}
}

// ToFIT converts Activity into proto.FIT.
func (a *Activity) ToFIT(options *mesgdef.Options) proto.FIT {
	size := 1 + len(a.Sessions) + len(a.UnrelatedMessages)
	for i := range a.Sessions {
		ses := &a.Sessions[i]
		size += len(ses.Records) + len(ses.Laps)
	}

	fit := proto.FIT{Messages: make([]proto.Message, 0, size)}
	fit.Messages = append(fit.Messages, a.Creator.FileId.ToMesg(options))
	fit.Messages = append(fit.Messages, a.UnrelatedMessages...)

	var totalTimerTime uint32
	var lastTimestamp time.Time
	for i := range a.Sessions {
		ses := &a.Sessions[i]
		if ses.Timestamp.After(lastTimestamp) {
			lastTimestamp = ses.Timestamp
		}
		totalTimerTime += ses.TotalTimerTime

		for j := range ses.Records {
			rec := &ses.Records[j]
			fit.Messages = append(fit.Messages, rec.Record.ToMesg(options))
		}
		for j := range ses.Laps {
			lap := &ses.Laps[j]
			fit.Messages = append(fit.Messages, lap.Lap.ToMesg(options))
		}
		fit.Messages = append(fit.Messages, ses.Session.ToMesg(options))
	}

	activityMesg := mesgdef.NewActivity(nil).
		SetType(typedef.ActivityAutoMultiSport).
		SetTimestamp(lastTimestamp).
		SetLocalTimestamp(lastTimestamp.Add(time.Duration(a.Timezone) * time.Hour)).
		SetTotalTimerTime(totalTimerTime).
		SetNumSessions(uint16(len(a.Sessions)))

	fit.Messages = append(fit.Messages, activityMesg.ToMesg(options))

	filedef.SortMessagesByTimestamp(fit.Messages[1:]) // Exclude FileId

	return fit
}

// MarshalAppendJSON appends the JSON format encoding of Activity to b, returning the result.
func (a *Activity) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')

	b = append(b, `"creator":`...)
	b = a.Creator.MarshalAppendJSON(b)
	b = append(b, ',')

	b = append(b, `"timezone":`...)
	b = strconv.AppendInt(b, int64(a.Timezone), 10)
	b = append(b, ',')

	if len(a.Sessions) != 0 {
		b = append(b, `"sessions":[`...)
		for i := range a.Sessions {
			n := len(b)
			b = a.Sessions[i].MarshalAppendJSON(b)
			if len(b) != n && i != len(a.Sessions)-1 {
				b = append(b, ',')
			}
		}
		b = append(b, ']')
	}

	if b[len(b)-1] == '{' {
		return b[:len(b)-1]
	}
	if b[len(b)-1] == ',' {
		b = b[:len(b)-1]
	}

	return append(b, '}')
}
