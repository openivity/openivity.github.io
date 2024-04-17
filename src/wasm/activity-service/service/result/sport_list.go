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
	"strconv"
)

type SportList struct {
	Sports []Sport
}

var _ json.Marshaler = &SportList{}

func (s *SportList) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')

	buf.WriteString("\"sports\":[")
	for i := range s.Sports {
		b, _ := s.Sports[i].MarshalJSON()
		buf.Write(b)
		if i != len(s.Sports)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

type Sport struct {
	ID                   uint8
	Name                 string
	ToleranceMovingSpeed float64
}

var _ json.Marshaler = &Sport{}

func (s *Sport) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')
	buf.WriteString("\"id\":" + strconv.FormatUint(uint64(s.ID), 10) + ",")
	buf.WriteString("\"name\":\"" + s.Name + "\",")
	buf.WriteString("\"toleranceMovingSpeed\":" + strconv.FormatFloat(s.ToleranceMovingSpeed, 'g', -1, 64))
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
