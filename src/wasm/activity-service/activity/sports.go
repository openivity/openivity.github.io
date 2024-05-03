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

	"github.com/muktihari/fit/profile/typedef"
)

// Sport is a sport structure.
type Sport struct {
	ID                   typedef.Sport
	Name                 string
	ToleranceMovingSpeed float64
}

// MarshalAppendJSON appends the JSON format encoding of Sport to b, returning the result.
func (s *Sport) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')
	b = append(b, `"id":`...)
	b = strconv.AppendUint(b, uint64(s.ID), 10)
	b = append(b, ',')

	b = append(b, `"name":`...)
	b = strconv.AppendQuote(b, s.Name)
	b = append(b, ',')

	b = append(b, `"toleranceMovingSpeed":`...)
	b = strconv.AppendFloat(b, s.ToleranceMovingSpeed, 'g', -1, 64)
	b = append(b, '}')

	return b
}
