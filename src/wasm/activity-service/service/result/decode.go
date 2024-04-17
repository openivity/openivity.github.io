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
	"fmt"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/activity"
)

type Decode struct {
	Err               error
	DecodeTook        time.Duration
	SerializationTook time.Duration
	TotalElapsed      time.Duration
	Activities        []activity.Activity
}

var _ json.Marshaler = &Decode{}

func (d *Decode) MarshalJSON() ([]byte, error) {
	if d.Err != nil {
		return []byte(fmt.Sprintf("{%q:%q}", "err", d.Err)), nil
	}

	begin := time.Now()
	buf := new(bytes.Buffer)
	buf.WriteByte('{')

	buf.WriteString("\"err\":null,")

	if len(d.Activities) != 0 {
		buf.WriteString("\"activities\":[")
		for i := range d.Activities {
			b, _ := d.Activities[i].MarshalJSON()
			buf.Write(b)
			if i != len(d.Activities)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteString("],")
	}

	d.SerializationTook = time.Since(begin)
	d.TotalElapsed = d.DecodeTook + d.SerializationTook

	buf.WriteString("\"decodeTook\":")
	buf.WriteString(strconv.FormatInt(d.DecodeTook.Milliseconds(), 10))
	buf.WriteByte(',')

	buf.WriteString("\"serializationTook\":")
	buf.WriteString(strconv.FormatInt(d.SerializationTook.Milliseconds(), 10))
	buf.WriteByte(',')

	buf.WriteString("\"totalElapsed\":")
	buf.WriteString(strconv.FormatInt(d.TotalElapsed.Milliseconds(), 10))

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

type DecodeWorker struct {
	Err      error
	Index    int
	Activity *activity.Activity
}
