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
)

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

var _ json.Marshaler = &Encode{}

func (e *Encode) MarshalJSON() ([]byte, error) {
	if e.Err != nil {
		return []byte(fmt.Sprintf("{%q:%q}", "err", e.Err)), nil
	}

	begin := time.Now()
	buf := new(bytes.Buffer)

	buf.WriteByte('{')
	buf.WriteString("\"fileName\":\"")
	buf.WriteString(e.FileName)
	buf.WriteString("\",")

	buf.WriteString("\"fileType\":\"")
	buf.WriteString(e.FileType)
	buf.WriteString("\",")

	buf.WriteString("\"encodeTook\":")
	buf.WriteString(strconv.FormatInt(e.EncodeTook.Milliseconds(), 10))
	buf.WriteByte(',')

	buf.WriteString("\"deserializeInputTook\":")
	buf.WriteString(strconv.FormatInt(e.DeserializeInputTook.Milliseconds(), 10))
	buf.WriteByte(',')

	buf.WriteString("\"filesBytes\":[")
	for i := range e.FilesBytes {
		buf.WriteByte('[')
		for j := range e.FilesBytes[i] {
			buf.WriteString(strconv.Itoa(int(e.FilesBytes[i][j]))) // keep it as number
			if j != len(e.FilesBytes[i])-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte(']')
		if i != len(e.FilesBytes)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString("],")

	e.SerializationTook = time.Since(begin)
	e.TotalElapsed = e.DeserializeInputTook + e.EncodeTook + e.SerializationTook

	buf.WriteString("\"serializationTook\":")
	buf.WriteString(strconv.FormatInt(e.DeserializeInputTook.Milliseconds(), 10))
	buf.WriteByte(',')

	buf.WriteString("\"totalElapsed\":")
	buf.WriteString(strconv.FormatInt(e.TotalElapsed.Milliseconds(), 10))

	buf.WriteByte('}')
	return buf.Bytes(), nil
}
