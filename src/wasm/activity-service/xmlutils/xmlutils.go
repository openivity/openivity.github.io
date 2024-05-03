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

package xmlutils

import (
	"encoding/xml"
	"io"
)

// StartElement is syntax sugar for xml.StartElement{Name: xml.Name{Local: name}}
func StartElement(name string) xml.StartElement {
	return xml.StartElement{Name: xml.Name{Local: name}}
}

// EncodeElement is syntax sugar for encoding a full element: <se>charData<se.End()>.
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

// MarshalWrite marshals and write to w.
func MarshalWrite(w io.Writer, marshaler xml.Marshaler) error {
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", " ")
	if err := marshaler.MarshalXML(enc, xml.StartElement{}); err != nil {
		_ = enc.Close()
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}
	return nil
}
