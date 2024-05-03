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

type ManufacturerList struct {
	Manufacturers []activity.Manufacturer
}

// MarshalAppendJSON appends the JSON format encoding of ManufacturerList to b, returning the result.
func (m *ManufacturerList) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')
	b = append(b, `"manufacturers":[`...)
	for i := range m.Manufacturers {
		b = m.Manufacturers[i].MarshalAppendJSON(b)
		if i != len(m.Manufacturers)-1 {
			b = append(b, ',')
		}
	}
	b = append(b, ']')
	b = append(b, '}')
	return b
}
