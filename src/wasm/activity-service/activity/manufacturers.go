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

// Manufacturer is manufacturer with its product list.
type Manufacturer struct {
	ID       typedef.Manufacturer
	Name     string
	Products []Product
}

// MarshalAppendJSON appends the JSON format encoding of Manufacturer to b, returning the result.
func (m *Manufacturer) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')
	b = append(b, `"id":`...)
	b = strconv.AppendUint(b, uint64(m.ID.Uint16()), 10)
	b = append(b, ',')

	b = append(b, `"name":`...)
	b = strconv.AppendQuote(b, m.Name)
	b = append(b, ',')

	b = append(b, `"products":[`...)
	for i := range m.Products {
		b = m.Products[i].MarshalAppendJSON(b)
		if i != len(m.Products)-1 {
			b = append(b, ',')
		}
	}
	b = append(b, ']')
	b = append(b, '}')
	return b
}

// Product is a manufacturer's product.
type Product struct {
	ID   uint16
	Name string
}

// // MarshalAppendJSON appends the JSON format encoding of Product to b, returning the result.
func (p *Product) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')
	b = append(b, `"id":`...)
	b = strconv.AppendUint(b, uint64(p.ID), 10)
	b = append(b, ',')

	b = append(b, `"name":`...)
	b = strconv.AppendQuote(b, p.Name)
	b = append(b, '}')

	return b
}
