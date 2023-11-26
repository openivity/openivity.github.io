package schema

import (
	"encoding/xml"
	"fmt"
)

// TODO: for implementing xml.Marshaler
const (
	xmlns    = "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2"
	xmlnsxsi = "http://www.w3.org/2001/XMLSchema-instance"
	xmlnsns5 = "http://www.garmin.com/xmlschemas/ActivityGoals/v1"
	xmlnsns3 = "http://www.garmin.com/xmlschemas/ActivityExtension/v2"
	xmlnsns2 = "http://www.garmin.com/xmlschemas/UserProfile/v2"
	xmlnstpx = "http://www.garmin.com/xmlschemas/GpxExtensions/v3"
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
