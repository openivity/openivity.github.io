package schema

import (
	"encoding/xml"
	"time"
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
					return err
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

type Activity struct {
	Sport   string        `xml:"Sport,attr"`
	ID      time.Time     `xml:"Id"`
	Laps    []ActivityLap `xml:"Lap"`
	Notes   string        `xml:"Notes,omitempty"`
	Creator *Device       `xml:"Creator,omitempty"`
}

var _ xml.Unmarshaler = &ActivityList{}

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
					return err
				}
				a.Laps = append(a.Laps, activityLap)
			case "Creator":
				var device Device
				if err := device.UnmarshalXML(dec, elem); err != nil {
					return err
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
					return err
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
