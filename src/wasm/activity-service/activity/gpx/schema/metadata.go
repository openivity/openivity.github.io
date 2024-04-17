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
	"time"

	kxml "github.com/muktihari/openactivity-fit/kit/xml"
)

// Metadata is GPX's Metadata schema (simplified).
type Metadata struct {
	Name   string    `xml:"name,omitempty"`
	Desc   string    `xml:"desc,omitempty"`
	Author *Author   `xml:"author,omitempty"`
	Link   *Link     `xml:"link,omitempty"`
	Time   time.Time `xml:"time,omitempty"`
}

var _ xml.Unmarshaler = &Metadata{}

func (m *Metadata) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "author":
				var author Author
				if err := author.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				m.Author = &author
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "name":
				m.Name = string(elem)
			case "desc":
				m.Desc = string(elem)
			case "time":
				t, err := time.Parse(time.RFC3339, string(elem))
				if err != nil {
					return err
				}
				m.Time = t
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

func (m *Metadata) Validate() error {

	if err := m.Link.Validate(); err != nil {
		return fmt.Errorf("link: %w", err)
	}
	if err := m.Author.Validate(); err != nil {
		return fmt.Errorf("author: %w", err)
	}
	return nil
}

var _ xml.Marshaler = &Metadata{}

func (m *Metadata) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if len(m.Name) != 0 {
		if err := kxml.EncodeElement(enc, kxml.StartElement("name"), xml.CharData(m.Name)); err != nil {
			return fmt.Errorf("name: %w", err)
		}
	}

	if len(m.Desc) != 0 {
		if err := kxml.EncodeElement(enc, kxml.StartElement("desc"), xml.CharData(m.Desc)); err != nil {
			return fmt.Errorf("desc: %w", err)
		}
	}

	if m.Author != nil {
		if err := m.Author.MarshalXML(enc, kxml.StartElement("author")); err != nil {
			return fmt.Errorf("author: %w", err)
		}
	}

	if m.Link != nil {
		if err := m.Link.MarshalXML(enc, kxml.StartElement("link")); err != nil {
			return fmt.Errorf("link: %w", err)
		}
	}

	if !m.Time.IsZero() {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("time"),
			xml.CharData(m.Time.Format(time.RFC3339))); err != nil {
			return fmt.Errorf("time: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

// Author is Author schema (simplified).
type Author struct {
	Name string `xml:"name"`
	Link *Link  `xml:"link"`
}

var _ xml.Unmarshaler = &Author{}

func (a *Author) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "link":
				var link Link
				if err := link.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				a.Link = &link
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "name":
				a.Name = string(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

func (p *Author) Validate() error {
	if p == nil {
		return nil
	}
	if p.Name == "" {
		return fmt.Errorf("name should not be empty")
	}
	if p.Link == nil {
		return fmt.Errorf("link should not be empty")
	}
	if err := p.Link.Validate(); err != nil {
		return fmt.Errorf("link: %w", err)
	}
	return nil
}

var _ xml.Marshaler = &Author{}

func (a *Author) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc, kxml.StartElement("name"), xml.CharData(a.Name)); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	if a.Link != nil {
		if err := a.Link.MarshalXML(enc, kxml.StartElement("link")); err != nil {
			return fmt.Errorf("marshal link: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

// Link is Link schema.
type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`

	Text string `xml:"text,omitempty"`
	Type string `xml:"type,omitempty"`
}

var _ xml.Unmarshaler = &Link{}

func (a *Link) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for i := range se.Attr {
		attr := se.Attr[i]

		switch attr.Name.Local {
		case "href":
			a.Href = attr.Value
		}
	}

	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			targetCharData = elem.Name.Local
		case xml.CharData:
			switch targetCharData {
			case "text":
				a.Text = string(elem)
			case "type":
				a.Type = string(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

func (l *Link) Validate() error {
	if l == nil {
		return nil
	}
	if l.Href == "" {
		return fmt.Errorf("href should not be empty")
	}
	return nil
}

var _ xml.Marshaler = &Link{}

func (l *Link) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	se.Attr = append(se.Attr, xml.Attr{
		Name:  xml.Name{Local: "href"},
		Value: l.Href,
	})

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if len(l.Text) != 0 {
		if err := kxml.EncodeElement(enc, kxml.StartElement("text"), xml.CharData(l.Text)); err != nil {
			return fmt.Errorf("text: %w", err)
		}
	}

	if len(l.Type) != 0 {
		if err := kxml.EncodeElement(enc, kxml.StartElement("type"), xml.CharData(l.Type)); err != nil {
			return fmt.Errorf("type: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}
