package schema

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/kit"
)

type ActivityLap struct {
	StartTime           time.Time     `xml:"StartTime,attr"`
	TotalTimeSeconds    float64       `xml:"TotalTimeSeconds"`
	DistanceMeters      float64       `xml:"DistanceMeters"`
	MaximumSpeed        *float64      `xml:"MaximumSpeed,omitempty"`
	Calories            uint16        `xml:"Calories"`
	AverageHeartRateBpm *uint8        `xml:"AverageHeartRateBpm"`
	MaximumHeartRateBpm *uint8        `xml:"MaximumHeartRateBpm"`
	Intensity           Intensity     `xml:"Intensity"`
	Cadence             *uint8        `xml:"Cadence,omitempty"`
	TriggerMethod       TriggerMethod `xml:"TriggerMethod"`
	Tracks              []Track       `xml:"Track,omitempty"`
	Notes               string        `xml:"Notes,omitempty"`
}

var _ xml.Unmarshaler = &ActivityLap{}

func (a *ActivityLap) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for i := range se.Attr {
		attr := &se.Attr[i]

		switch attr.Name.Local {
		case "StartTime":
			t, err := time.Parse(time.RFC3339, attr.Value)
			if err != nil {
				return err
			}
			a.StartTime = t
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
			case "Track":
				var track Track
				if err := track.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				a.Tracks = append(a.Tracks, track)
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "TotalTimeSeconds":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return err
				}
				a.TotalTimeSeconds = f
			case "DistanceMeters":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return err
				}
				a.DistanceMeters = f
			case "MaximumSpeed":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return err
				}
				a.MaximumSpeed = &f
			case "Calories":
				u, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return err
				}
				a.Calories = uint16(u)
			case "AverageHeartRateBpm":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return err
				}
				a.AverageHeartRateBpm = kit.Ptr(uint8(u))
			case "MaximumHeartRateBpm":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return err
				}
				a.MaximumHeartRateBpm = kit.Ptr(uint8(u))
			case "Intensity":
				a.Intensity = Intensity(elem)
			case "Cadence":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return err
				}
				a.Cadence = kit.Ptr(uint8(u))
			case "TriggerMethod":
				a.TriggerMethod = TriggerMethod(elem)
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

type Intensity string

const (
	IntensityActive  Intensity = "Active"
	IntensityResting Intensity = "Resting"
)

type TriggerMethod string

const (
	TriggerMethodManual    TriggerMethod = "Manual"
	TriggerMethodDistance  TriggerMethod = "Distance"
	TriggerMethodLocation  TriggerMethod = "Location"
	TriggerMethodTime      TriggerMethod = "Time"
	TriggerMethodHeartRate TriggerMethod = "HeartRate"
)
