package schema

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Metadata is GPX's Metadata schema (simplified).
type Metadata struct {
	Name   string    `xml:"name,omitempty"`
	Desc   string    `xml:"desc,omitempty"`
	Author *Author   `xml:"author,omitempty"`
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
	if err := m.Author.Validate(); err != nil {
		return fmt.Errorf("author: %w", err)
	}
	return nil
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

// Link is Link schema.
type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Value   string   `xml:",chardata"`

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
			default:
				a.Value = string(elem)
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
