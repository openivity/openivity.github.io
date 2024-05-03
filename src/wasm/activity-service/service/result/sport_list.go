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

package result

import (
	"github.com/openivity/activity-service/activity"
)

type SportList struct {
	Sports []activity.Sport
}

// MarshalAppendJSON appends the JSON format encoding of SportList to b, returning the result.
func (s *SportList) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')
	b = append(b, `"sports":[`...)
	for i := range s.Sports {
		b = s.Sports[i].MarshalAppendJSON(b)
		if i != len(s.Sports)-1 {
			b = append(b, ',')
		}
	}
	b = append(b, ']')
	b = append(b, '}')
	return b
}
