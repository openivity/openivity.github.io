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
	"io"
	"time"

	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/xmlutils"
)

// Metadata is GPX's Metadata schema (simplified).
type Metadata struct {
	Name   string    `xml:"name,omitempty"`
	Desc   string    `xml:"desc,omitempty"`
	Author *Author   `xml:"author,omitempty"`
	Link   *Link     `xml:"link,omitempty"`
	Time   time.Time `xml:"time,omitempty"`
}

func (m *Metadata) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	for {
		token, err := tok.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if token.IsEndElementOf(se) {
			return nil
		}
		if token.IsEndElement {
			continue
		}

		switch string(token.Name.Local) {
		case "name":
			m.Name = string(token.Data)
		case "desc":
			m.Desc = string(token.Data)
		case "author":
			m.Author = new(Author)
			se := xmltokenizer.GetToken().Copy(token)
			err = m.Author.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("author: %w", err)
			}
		case "link":
			m.Link = new(Link)
			se := xmltokenizer.GetToken().Copy(token)
			err = m.Link.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("link: %w", err)
			}
		case "time":
			m.Time, err = time.Parse(time.RFC3339, string(token.Data))
			if err != nil {
				return fmt.Errorf("time: %w", err)
			}
		}
	}

	return nil
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

var _ xml.Marshaler = (*Metadata)(nil)

func (m *Metadata) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if len(m.Name) != 0 {
		if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("name"), xml.CharData(m.Name)); err != nil {
			return fmt.Errorf("name: %w", err)
		}
	}

	if len(m.Desc) != 0 {
		if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("desc"), xml.CharData(m.Desc)); err != nil {
			return fmt.Errorf("desc: %w", err)
		}
	}

	if m.Author != nil {
		if err := m.Author.MarshalXML(enc, xmlutils.StartElement("author")); err != nil {
			return fmt.Errorf("author: %w", err)
		}
	}

	if m.Link != nil {
		if err := m.Link.MarshalXML(enc, xmlutils.StartElement("link")); err != nil {
			return fmt.Errorf("link: %w", err)
		}
	}

	if !m.Time.IsZero() {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("time"),
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

func (a *Author) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	for {
		token, err := tok.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("author: %w", err)
		}

		if token.IsEndElementOf(se) {
			return nil
		}
		if token.IsEndElement {
			continue
		}

		switch string(token.Name.Local) {
		case "name":
			a.Name = string(token.Data)
		case "link":
			a.Link = new(Link)
			se := xmltokenizer.GetToken().Copy(token)
			err := a.Link.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("link: %w", err)
			}
		}
	}

	return nil
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

var _ xml.Marshaler = (*Author)(nil)

func (a *Author) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("name"), xml.CharData(a.Name)); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	if a.Link != nil {
		if err := a.Link.MarshalXML(enc, xmlutils.StartElement("link")); err != nil {
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

func (a *Link) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	for i := range se.Attrs {
		attr := &se.Attrs[i]
		switch string(attr.Name.Local) {
		case "href":
			a.Href = string(attr.Value)
		}
	}

	for {
		token, err := tok.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("link: %w", err)
		}

		if token.IsEndElementOf(se) {
			return nil
		}
		if token.IsEndElement {
			continue
		}

		switch string(token.Name.Local) {
		case "text":
			a.Text = string(token.Data)
		case "type":
			a.Type = string(token.Data)
		}
	}

	return nil
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

var _ xml.Marshaler = (*Link)(nil)

func (l *Link) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	se.Attr = append(se.Attr, xml.Attr{
		Name:  xml.Name{Local: "href"},
		Value: l.Href,
	})

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if len(l.Text) != 0 {
		if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("text"), xml.CharData(l.Text)); err != nil {
			return fmt.Errorf("text: %w", err)
		}
	}

	if len(l.Type) != 0 {
		if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("type"), xml.CharData(l.Type)); err != nil {
			return fmt.Errorf("type: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}
