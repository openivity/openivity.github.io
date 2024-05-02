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
	"fmt"
	"strconv"
	"time"

	"github.com/openivity/activity-service/activity"
)

// Decode is decode result.
type Decode struct {
	Err               error
	DecodeTook        time.Duration
	SerializationTook time.Duration
	TotalElapsed      time.Duration
	Activities        []activity.Activity
}

// MarshalAppendJSON appends the JSON format encoding of Decode to b, returning the result.
func (d *Decode) MarshalAppendJSON(b []byte) []byte {
	if d.Err != nil {
		return []byte(fmt.Sprintf("{%q:%q}", "err", d.Err))
	}

	begin := time.Now()
	b = append(b, '{')

	b = append(b, `"err":null,`...)

	if len(d.Activities) != 0 {
		b = append(b, `"activities":[`...)
		for i := range d.Activities {
			n := len(b)
			b = d.Activities[i].MarshalAppendJSON(b)
			if len(b) != n && i != len(d.Activities)-1 {
				b = append(b, ',')
			}
		}
		b = append(b, ']')
		b = append(b, ',')
	}

	d.SerializationTook = time.Since(begin)
	d.TotalElapsed = d.DecodeTook + d.SerializationTook

	b = append(b, `"decodeTook":`...)
	b = append(b, strconv.FormatInt(d.DecodeTook.Milliseconds(), 10)...)
	b = append(b, ',')

	b = append(b, `"serializationTook":`...)
	b = append(b, strconv.FormatInt(d.SerializationTook.Milliseconds(), 10)...)
	b = append(b, ',')

	b = append(b, `"totalElapsed":`...)
	b = append(b, strconv.FormatInt(d.TotalElapsed.Milliseconds(), 10)...)

	b = append(b, '}')

	return b
}

// DecodeWorker is a decode worker.
type DecodeWorker struct {
	Err      error
	Index    int
	Activity *activity.Activity
}
