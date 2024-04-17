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
	"bytes"
	"encoding/json"

	"github.com/muktihari/openactivity-fit/activity/fit"
)

type ManufacturerList struct {
	Manufacturers []fit.Manufacturer
}

var _ json.Marshaler = &ManufacturerList{}

func (m *ManufacturerList) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')

	buf.WriteString("\"manufacturers\":[")
	for i := range m.Manufacturers {
		b, _ := m.Manufacturers[i].MarshalJSON()
		buf.Write(b)
		if i != len(m.Manufacturers)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	buf.WriteByte('}')
	return buf.Bytes(), nil
}
