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
	"github.com/muktihari/fit/factory"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
)

func convertEventToMesg(event *activity.Event) proto.Message {
	return factory.CreateMesgOnly(mesgnum.Event).WithFields(
		factory.CreateField(mesgnum.Event, fieldnum.EventTimestamp).WithValue(datetime.ToUint32(event.Timestamp)),
		factory.CreateField(mesgnum.Event, fieldnum.EventEvent).WithValue(event.Event),
		factory.CreateField(mesgnum.Event, fieldnum.EventEventType).WithValue(event.EventType),
	)
}
