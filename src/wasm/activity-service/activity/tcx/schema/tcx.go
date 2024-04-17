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

package schema

import (
	"encoding/xml"
	"fmt"
	"strings"

	kxml "github.com/muktihari/openactivity-fit/kit/xml"
)

// TODO: for implementing xml.Marshaler
const (
	xmlns    = "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2"
	xmlnsxsi = "http://www.w3.org/2001/XMLSchema-instance"
	xmlnsns1 = "http://www.garmin.com/xmlschemas/ActivityExtension/v2"
)

var schemaLocations = [...]string{
	"http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2",
	"http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd",
}

// TCX simplified schema.
type TCX struct {
	Activities []ActivityList `xml:"Activities,omitempty"`
	Author     *Application   `xml:"Author,omitempty"`
}

var _ xml.Unmarshaler = &TCX{}

func (t *TCX) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "Activities":
				var al ActivityList
				if err := al.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Activities: %w", err)
				}
				t.Activities = append(t.Activities, al)
			case "Author":
				var application Application
				if err := application.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Author: %w", err)
				}
				t.Author = &application
			}
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &TCX{}

func (t *TCX) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	se.Name = xml.Name{Local: "TrainingCenterDatabase"}
	se.Attr = []xml.Attr{
		{Name: xml.Name{Local: "xmlns"}, Value: xmlns},
		{Name: xml.Name{Local: "xmlns:xsi"}, Value: xmlnsxsi},
		{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: strings.Join(schemaLocations[:], " ")},
		{Name: xml.Name{Local: "xmlns:ns1"}, Value: xmlnsns1},
	}

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	for i := range t.Activities {
		if err := t.Activities[i].MarshalXML(enc, kxml.StartElement("Activities")); err != nil {
			return fmt.Errorf("activities[%d]: %w", i, err)
		}
	}

	if t.Author != nil {
		if err := t.Author.MarshalXML(enc, kxml.StartElement("Author")); err != nil {
			return fmt.Errorf("author: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}
