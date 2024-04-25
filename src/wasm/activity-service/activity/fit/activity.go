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
	activityTypeField.Value = proto.Uint8(uint8(typedef.ActivityAutoMultiSport))
	mesg.Fields = append(mesg.Fields, activityTypeField)

	timestampField := factory.CreateField(mesgnum.Activity, fieldnum.ActivityTimestamp)
	timestampField.Value = proto.Uint32(datetime.ToUint32(timestamp))
	mesg.Fields = append(mesg.Fields, timestampField)

	localTimestamp := datetime.ToLocalTime(timestampField.Value, timezone)
	localTimestampField := factory.CreateField(mesgnum.Activity, fieldnum.ActivityLocalTimestamp)
	localTimestampField.Value = proto.Uint32(datetime.ToUint32(localTimestamp))
	mesg.Fields = append(mesg.Fields, localTimestampField)

	sessionCountField := factory.CreateField(mesgnum.Activity, fieldnum.ActivityNumSessions)
	sessionCountField.Value = proto.Uint16(sessionCount)

	return mesg
}
