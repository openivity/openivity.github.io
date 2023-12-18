package schema

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/kit"
	kxml "github.com/muktihari/openactivity-fit/kit/xml"
)

type ActivityLap struct {
	StartTime           time.Time     `xml:"StartTime,attr"`
	TotalTimeSeconds    float64       `xml:"TotalTimeSeconds"`
	DistanceMeters      float64       `xml:"DistanceMeters"`
	MaximumSpeed        *float64      `xml:"MaximumSpeed,omitempty"`
	Calories            uint16        `xml:"Calories"`
	AverageHeartRateBpm *uint8        `xml:"AverageHeartRateBpm"`
	MaximumHeartRateBpm *uint8        `xml:"MaximumHeartRateBpm"`
	Intensity           Intensity     `xml:"Intensity,omitempty"`
	Cadence             *uint8        `xml:"Cadence,omitempty"`
	TriggerMethod       TriggerMethod `xml:"TriggerMethod,omitempty"`
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
				return fmt.Errorf("parse StartTime: %w", err)
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
					return fmt.Errorf("unmarshal Track: %w", err)
				}
				a.Tracks = append(a.Tracks, track)
			case "Value":
				targetCharData = targetCharData + "Value"
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "TotalTimeSeconds":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse TotalTimeSeconds: %w", err)
				}
				a.TotalTimeSeconds = f
			case "DistanceMeters":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse DistanceMeters: %w", err)
				}
				a.DistanceMeters = f
			case "MaximumSpeed":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse MaximumSpeed: %w", err)
				}
				a.MaximumSpeed = &f
			case "Calories":
				u, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return fmt.Errorf("parse Calories: %w", err)
				}
				a.Calories = uint16(u)
			case "AverageHeartRateBpm":
				continue
			case "AverageHeartRateBpmValue":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return fmt.Errorf("parse AverageHeartRateBpm: %w", err)
				}
				a.AverageHeartRateBpm = kit.Ptr(uint8(u))
			case "MaximumHeartRateBpm":
				continue
			case "MaximumHeartRateBpmValue":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return fmt.Errorf("parse MaximumHeartRateBpm: %w", err)
				}
				a.MaximumHeartRateBpm = kit.Ptr(uint8(u))
			case "Intensity":
				a.Intensity = Intensity(elem)
			case "Cadence":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return fmt.Errorf("parse Cadence: %w", err)
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

var _ xml.Marshaler = &ActivityLap{}

func (a *ActivityLap) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	se.Attr = append(se.Attr, xml.Attr{
		Name:  xml.Name{Local: "StartTime"},
		Value: a.StartTime.Format(time.RFC3339),
	})

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc,
		kxml.StartElement("TotalTimeSeconds"),
		xml.CharData(strconv.FormatFloat(a.TotalTimeSeconds, 'g', -1, 64))); err != nil {
		return fmt.Errorf("totalTimeSeconds: %w", err)
	}
	if err := kxml.EncodeElement(enc,
		kxml.StartElement("DistanceMeters"),
		xml.CharData(strconv.FormatFloat(a.DistanceMeters, 'g', -1, 64))); err != nil {
		return fmt.Errorf("distanceMeters: %w", err)
	}

	if a.MaximumSpeed != nil {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("MaximumSpeed"),
			xml.CharData(strconv.FormatFloat(*a.MaximumSpeed, 'g', -1, 64))); err != nil {
			return fmt.Errorf("maximumSpeed: %w", err)
		}
	}

	if err := kxml.EncodeElement(enc,
		kxml.StartElement("Calories"),
		xml.CharData(strconv.FormatUint(uint64(a.Calories), 10))); err != nil {
		return fmt.Errorf("calories: %w", err)
	}

	if a.AverageHeartRateBpm != nil {
		avgHR := kxml.StartElement("AverageHeartRateBpm")
		if err := enc.EncodeToken(avgHR); err != nil {
			return fmt.Errorf("averageHeartRateBpm start: %w", err)
		}
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("Value"),
			xml.CharData(strconv.FormatUint(uint64(*a.AverageHeartRateBpm), 10))); err != nil {
			return fmt.Errorf("averageHeartRateBpmValue: %w", err)
		}
		if err := enc.EncodeToken(avgHR.End()); err != nil {
			return fmt.Errorf("averageHeartRateBpm end: %w", err)
		}
	}
	if a.MaximumHeartRateBpm != nil {
		maxHR := kxml.StartElement("MaximumHeartRateBpm")
		if err := enc.EncodeToken(maxHR); err != nil {
			return fmt.Errorf("maximumHeartRateBpm start: %w", err)
		}
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("Value"),
			xml.CharData(strconv.FormatUint(uint64(*a.MaximumHeartRateBpm), 10))); err != nil {
			return fmt.Errorf("maximumHeartRateBpmValue: %w", err)
		}
		if err := enc.EncodeToken(maxHR.End()); err != nil {
			return fmt.Errorf("maximumHeartRateBpm end: %w", err)
		}
	}

	if len(a.Intensity) != 0 {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("Intensity"),
			xml.CharData(a.Intensity)); err != nil {
			return fmt.Errorf("intensity: %w", err)
		}
	}

	if a.Cadence != nil {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("Cadence"),
			xml.CharData(strconv.FormatUint(uint64(*a.Cadence), 10))); err != nil {
			return fmt.Errorf("cadence: %w", err)
		}

	}

	if len(a.TriggerMethod) != 0 {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("TriggerMethod"),
			xml.CharData(a.TriggerMethod)); err != nil {
			return fmt.Errorf("triggerMethod: %w", err)
		}
	}

	for i := range a.Tracks {
		if err := a.Tracks[i].MarshalXML(enc, kxml.StartElement("Track")); err != nil {
			return fmt.Errorf("track[%d]: %w", i, err)
		}
	}

	if len(a.Notes) != 0 {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("Notes"),
			xml.CharData(a.Notes)); err != nil {
			return fmt.Errorf("notes: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
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
