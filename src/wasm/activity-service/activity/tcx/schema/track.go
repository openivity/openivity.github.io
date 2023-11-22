package schema

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/kit"
)

type Track struct {
	Trackpoints []Trackpoint `xml:"Trackpoint,omitempty"`
}

var _ xml.Unmarshaler = &Track{}

func (t *Track) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "Trackpoint":
				var trackpoint Trackpoint
				if err := trackpoint.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				t.Trackpoints = append(t.Trackpoints, trackpoint)
			}
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

type Trackpoint struct {
	Time           time.Time   `xml:"Time"`
	Position       *Position   `xml:"Position,omitempty"`
	AltitudeMeters *float64    `xml:"AltitudeMeters,omitempty"`
	DistanceMeters *float64    `xml:"DistanceMeters,omitempty"`
	HeartRateBpm   *uint8      `xml:"HeartRateBpm,omitempty"`
	Cadence        *uint8      `xml:"Cadence,omitempty"`
	SensorState    SensorState `xml:"SensorState,omitempty"`
}

var _ xml.Unmarshaler = &Trackpoint{}

func (tp *Trackpoint) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "Position":
				var position Position
				if err := position.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				tp.Position = &position
			case "Value":
				if targetCharData == "HeartRateBpm" { // check prev value
					targetCharData = "HeartRateBpmValue"
				}
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "Time":
				t, err := time.Parse(time.RFC3339, string(elem))
				if err != nil {
					return fmt.Errorf("parse Time %q: %w", elem, err)
				}
				tp.Time = t
			case "AltitudeMeters":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse AltitudeMeters %q: %w", elem, err)
				}
				tp.AltitudeMeters = &f
			case "DistanceMeters":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse DistanceMeters %q: %w", elem, err)
				}
				tp.DistanceMeters = &f
			case "HeartRateBpmValue":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return fmt.Errorf("parse HeartRateBpm %q: %w", elem, err)
				}
				tp.HeartRateBpm = kit.Ptr(uint8(u))
			case "Cadence":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return fmt.Errorf("parse Cadence %q: %w", elem, err)
				}
				tp.Cadence = kit.Ptr(uint8(u))
			case "SensorState":
				tp.SensorState = SensorState(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

type Position struct {
	LatitudeDegrees  float64 // -90.0 to 90.0
	LongitudeDegrees float64 // -180.0 to 180.0
}

var _ xml.Unmarshaler = &Position{}

func (p *Position) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
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
			case "LatitudeDegrees":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return err
				}
				p.LatitudeDegrees = f
			case "LongitudeDegrees":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return err
				}
				p.LongitudeDegrees = f
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

type SensorState string

const (
	SensorStatePresent SensorState = "Present"
	SensorStateAbsent  SensorState = "Absent"
)
