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
	"errors"
	"strconv"

	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/proto"
)

var ErrNoActivity = errors.New("no activity")

// Activity is an activity. It use FIT SDK's structure as its base since FIT is currently the most advance format.
type Activity struct {
	Creator  Creator
	Timezone int8
	Sessions []Session

	Sports         []*mesgdef.Sport
	SplitSummaries []*mesgdef.SplitSummary // required for FIT file; entries must be unique within each split_type
	Activity       *mesgdef.Activity       // required for FIT file.

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
