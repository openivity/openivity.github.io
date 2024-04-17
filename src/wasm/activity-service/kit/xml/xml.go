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

package xml

import (
	"bytes"
	"encoding/xml"
	"strings"
)

func StartElement(name string) xml.StartElement {
	return xml.StartElement{Name: xml.Name{Local: name}}
}

func EncodeElement(enc *xml.Encoder, se xml.StartElement, charData xml.CharData) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}
	if charData != nil {
		if err := enc.EncodeToken(charData); err != nil {
			return err
		}
	}
	return enc.EncodeToken(se.End())
}

func Marshal(marshaler xml.Marshaler) (b []byte, err error) {
	w := bytes.NewBuffer(nil)
	w.Write([]byte(xml.Header))

	enc := xml.NewEncoder(w)
	defer func() {
		err = enc.Close()
		b = w.Bytes()
	}()

	enc.Indent("", strings.Repeat(" ", 1))

	if err = marshaler.MarshalXML(enc, xml.StartElement{}); err != nil {
		return nil, err
	}

	return
}
