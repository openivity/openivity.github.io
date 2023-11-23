package schema

import (
	"encoding/xml"
)

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
					return err
				}
				t.Activities = append(t.Activities, al)
			case "Author":
				var application Application
				if err := application.UnmarshalXML(dec, elem); err != nil {
					return err
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
