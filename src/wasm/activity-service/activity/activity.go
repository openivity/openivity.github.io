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
	"bytes"
	"strconv"
)

type Activity struct {
	Creator  Creator
	Timezone int
	Sessions []*Session
}

func (a *Activity) MarshalJSON() ([]byte, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('{')

	buf.WriteString("\"creator\":")
	b, _ := a.Creator.MarshalJSON()
	buf.Write(b)
	buf.WriteByte(',')

	buf.WriteString("\"timezone\":")
	buf.WriteString(strconv.FormatInt(int64(a.Timezone), 10))
	buf.WriteByte(',')

	if len(a.Sessions) != 0 {
		buf.WriteString("\"sessions\":[")
		for i := range a.Sessions {
			b, _ = a.Sessions[i].MarshalJSON()
			buf.Write(b)
			if i != len(a.Sessions)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte(']')
	}

	b = buf.Bytes()
	if b[len(b)-1] == ',' {
		b[len(b)-1] = '}'
		return b, nil
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (a *Activity) Clone() *Activity {
	act := &Activity{
		Creator:  *a.Creator.Clone(),
		Timezone: a.Timezone,
	}

	act.Sessions = make([]*Session, len(a.Sessions))
	for i := range a.Sessions {
		act.Sessions[i] = a.Sessions[i].Clone()
	}

	return act
}
