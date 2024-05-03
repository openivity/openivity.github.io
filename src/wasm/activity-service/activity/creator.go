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
	"strconv"
	"time"

	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
)

const Unknown = "Unknown"

type Creator struct {
	*mesgdef.FileId

	Name string
}

// CreateCreator creates new creator.
func CreateCreator(fileId *mesgdef.FileId) Creator {
	if fileId == nil {
		fileId = mesgdef.NewFileId(nil).
			SetType(typedef.FileActivity)
	}
	return Creator{FileId: fileId}
}

// MarshalAppendJSON appends the JSON format encoding of Creator to b, returning the result.
func (c *Creator) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')

	if c.Name == "" {
		c.Name = "Unknown"
	}

	b = append(b, `"name":`...)
	b = strconv.AppendQuote(b, c.Name)
	b = append(b, ',')

	if c.Manufacturer != typedef.ManufacturerInvalid {
		b = append(b, `"manufacturer":`...)
		b = strconv.AppendUint(b, uint64(c.Manufacturer), 10)
		b = append(b, ',')
	}
	if c.Product != basetype.Uint16Invalid {
		b = append(b, `"product":`...)
		b = strconv.AppendUint(b, uint64(c.Product), 10)
		b = append(b, ',')
	}
	if !c.TimeCreated.IsZero() {
		b = append(b, `"timeCreated":`...)
		b = strconv.AppendQuote(b, c.TimeCreated.Format(time.RFC3339))
	}

	if b[len(b)-1] == '{' {
		return b[:len(b)-1]
	}
	if b[len(b)-1] == ',' {
		b = b[:len(b)-1]
	}

	return append(b, '}')
}
