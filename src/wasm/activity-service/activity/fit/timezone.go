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
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/proto"
)

func CreateTimezone(mesg proto.Message) int {
	var (
		timestamp     = basetype.Uint32Invalid
		localDateTime = basetype.Uint32Invalid
	)

	for i := range mesg.Fields {
		switch mesg.Fields[i].Num {
		case fieldnum.ActivityTimestamp:
			timestamp = mesg.Fields[i].Value.Uint32()
		case fieldnum.ActivityLocalTimestamp:
			localDateTime = mesg.Fields[i].Value.Uint32()
		}
	}

	if timestamp == basetype.Uint32Invalid || localDateTime == basetype.Uint32Invalid {
		return 0 // Default UTC
	}

	return datetime.TzOffsetHoursFromUint32(localDateTime, timestamp)
}
