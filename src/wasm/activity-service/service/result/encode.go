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
)

// Encode is encode result.
type Encode struct {
	Err                  error
	DeserializeInputTook time.Duration
	EncodeTook           time.Duration
	SerializationTook    time.Duration
	TotalElapsed         time.Duration
	FileName             string
	FileType             string
	FilesBytes           [][]byte
}

// MarshalAppendJSON appends the JSON format encoding of Encode to b, returning the result.
func (e *Encode) MarshalAppendJSON(b []byte) []byte {
	if e.Err != nil {
		return []byte(fmt.Sprintf("{%q:%q}", "err", e.Err))
	}

	begin := time.Now()

	b = append(b, '{')

	b = append(b, `"fileName":`...)
	b = strconv.AppendQuote(b, e.FileName)
	b = append(b, ',')

	b = append(b, `"fileType":`...)
	b = strconv.AppendQuote(b, e.FileType)
	b = append(b, ',')

	b = append(b, `"encodeTook":`...)
	b = strconv.AppendInt(b, e.EncodeTook.Milliseconds(), 10)
	b = append(b, ',')

	b = append(b, `"deserializeInputTook":`...)
	b = strconv.AppendInt(b, e.DeserializeInputTook.Milliseconds(), 10)
	b = append(b, ',')

	b = append(b, `"filesBytes":[`...)
	for i := range e.FilesBytes {
		b = append(b, '[')
		for j := range e.FilesBytes[i] {
			b = strconv.AppendInt(b, int64(e.FilesBytes[i][j]), 10)
			if j != len(e.FilesBytes[i])-1 {
				b = append(b, ',')
			}
		}
		b = append(b, ']')
		if i != len(e.FilesBytes)-1 {
			b = append(b, ',')
		}
	}
	b = append(b, `],`...)

	e.SerializationTook = time.Since(begin)
	e.TotalElapsed = e.DeserializeInputTook + e.EncodeTook + e.SerializationTook

	b = append(b, `"serializationTook":`...)
	b = strconv.AppendInt(b, e.SerializationTook.Milliseconds(), 10)
	b = append(b, ',')

	b = append(b, `"totalElapsed":`...)
	b = strconv.AppendInt(b, e.TotalElapsed.Milliseconds(), 10)

	b = append(b, '}')

	return b
}
