package schema

import (
	"encoding/xml"
	"fmt"
	"time"

	kxml "github.com/muktihari/openactivity-fit/kit/xml"
)

type ActivityList struct {
	Activity *Activity `xml:"Activity,omitempty"`
}

var _ xml.Unmarshaler = &ActivityList{}

func (a *ActivityList) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "Activity":
				var activity Activity
				if err := activity.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Activity: %w", err)
				}
				a.Activity = &activity
			}
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &ActivityList{}

func (a *ActivityList) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if a.Activity != nil {
		if err := a.Activity.MarshalXML(enc, kxml.StartElement("Activity")); err != nil {
			return fmt.Errorf("activity: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

type Activity struct {
	Sport   string        `xml:"Sport,attr"`
	ID      time.Time     `xml:"Id"`
	Laps    []ActivityLap `xml:"Lap"`
	Notes   string        `xml:"Notes,omitempty"`
	Creator *Device       `xml:"Creator,omitempty"`
}

var _ xml.Unmarshaler = &Activity{}

func (a *Activity) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for i := range se.Attr {
		attr := &se.Attr[i]

		switch attr.Name.Local {
		case "Sport":
			a.Sport = attr.Value
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
			switch elem.Name.Local {
			case "Lap":
				var activityLap ActivityLap
				if err := activityLap.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Lap: %w", err)
				}
				a.Laps = append(a.Laps, activityLap)
			case "Creator":
				var device Device
				if err := device.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Creator: %w", err)
				}
				a.Creator = &device
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "Id":
				t, err := time.Parse(time.RFC3339, string(elem))
				if err != nil {
					return fmt.Errorf("parse Id: %w", err)
				}
				a.ID = t
			case "Notes":
				a.Notes = string(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &Activity{}

func (a *Activity) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	se.Attr = append(se.Attr, xml.Attr{
		Name:  xml.Name{Local: "Sport"},
		Value: a.Sport,
	})

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc, kxml.StartElement("Id"), xml.CharData(a.ID.Format(time.RFC3339))); err != nil {
		return fmt.Errorf("id: %w", err)
	}

	for i := range a.Laps {
		if err := a.Laps[i].MarshalXML(enc, kxml.StartElement("Lap")); err != nil {
			return fmt.Errorf("lap[%d]: %w", i, err)
		}
	}

	if err := kxml.EncodeElement(enc, kxml.StartElement("Notes"), xml.CharData(a.Notes)); err != nil {
		return fmt.Errorf("notes: %w", err)
	}

	if a.Creator != nil {
		if err := a.Creator.MarshalXML(enc, kxml.StartElement("Creator")); err != nil {
			return fmt.Errorf("creator: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}
