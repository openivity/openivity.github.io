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
	"math"
	"strconv"
	"time"

	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/xmlutils"
)

type ActivityLap struct {
	StartTime           time.Time     `xml:"StartTime,attr"`
	TotalTimeSeconds    float64       `xml:"TotalTimeSeconds"`
	DistanceMeters      float64       `xml:"DistanceMeters"`
	MaximumSpeed        float64       `xml:"MaximumSpeed,omitempty"`
	Calories            uint16        `xml:"Calories"`
	AverageHeartRateBpm uint8         `xml:"AverageHeartRateBpm>Value"`
	MaximumHeartRateBpm uint8         `xml:"MaximumHeartRateBpm>Value"`
	Intensity           Intensity     `xml:"Intensity,omitempty"`
	Cadence             uint8         `xml:"Cadence,omitempty"`
	TriggerMethod       TriggerMethod `xml:"TriggerMethod,omitempty"`
	Tracks              []Track       `xml:"Track,omitempty"`
	Notes               string        `xml:"Notes,omitempty"`
}

func (a *ActivityLap) reset() {
	a.TotalTimeSeconds = math.NaN()
	a.DistanceMeters = math.NaN()
	a.MaximumSpeed = math.NaN()
	a.Calories = basetype.Uint16Invalid
	a.AverageHeartRateBpm = basetype.Uint8Invalid
	a.MaximumHeartRateBpm = basetype.Uint8Invalid
	a.Cadence = basetype.Uint8Invalid
}

func (a *ActivityLap) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	a.reset()

	for i := range se.Attrs {
		attr := &se.Attrs[i]
		switch string(attr.Name.Local) {
		case "StartTime":
			t, err := time.Parse(time.RFC3339, string(attr.Value))
			if err != nil {
				return fmt.Errorf("parse StartTime: %w", err)
			}
			a.StartTime = t
		}
	}

	for {
		token, err := tok.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if token.IsEndElementOf(se) {
			break
		}
		if token.IsEndElement {
			continue
		}

		switch string(token.Name.Local) {
		case "Track":
			var track Track
			se := xmltokenizer.GetToken().Copy(token)
			err = track.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Track: %w", err)
			}
			a.Tracks = append(a.Tracks, track)
		case "TotalTimeSeconds":
			f, err := strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse TotalTimeSeconds: %w", err)
			}
			a.TotalTimeSeconds = f
		case "DistanceMeters":
			f, err := strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse DistanceMeters: %w", err)
			}
			a.DistanceMeters = f
		case "MaximumSpeed":
			f, err := strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse MaximumSpeed: %w", err)
			}
			a.MaximumSpeed = f
		case "Calories":
			u, err := strconv.ParseUint(string(token.Data), 10, 16)
			if err != nil {
				return fmt.Errorf("parse Calories: %w", err)
			}
			a.Calories = uint16(u)
		case "AverageHeartRateBpm":
			token, err = getValueToken(tok)
			if err != nil {
				return err
			}
			u, err := strconv.ParseUint(string(token.Data), 10, 8)
			if err != nil {
				return fmt.Errorf("parse AverageHeartRateBpm: %w", err)
			}
			a.AverageHeartRateBpm = uint8(u)
		case "MaximumHeartRateBpm":
			token, err = getValueToken(tok)
			if err != nil {
				return err
			}
			u, err := strconv.ParseUint(string(token.Data), 10, 8)
			if err != nil {
				return fmt.Errorf("parse MaximumHeartRateBpm: %w", err)
			}
			a.MaximumHeartRateBpm = uint8(u)
		case "Intensity":
			a.Intensity = Intensity(token.Data)
		case "Cadence":
			u, err := strconv.ParseUint(string(token.Data), 10, 8)
			if err != nil {
				return fmt.Errorf("parse Cadence: %w", err)
			}
			a.Cadence = uint8(u)
		case "TriggerMethod":
			a.TriggerMethod = TriggerMethod(token.Data)
		case "Notes":
			a.Notes = string(token.Data)
		}
	}

	return nil
}

var _ xml.Marshaler = (*ActivityLap)(nil)

func (a *ActivityLap) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	se.Attr = append(se.Attr, xml.Attr{
		Name:  xml.Name{Local: "StartTime"},
		Value: a.StartTime.Format(time.RFC3339),
	})

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("TotalTimeSeconds"),
		xml.CharData(strconv.FormatFloat(a.TotalTimeSeconds, 'g', -1, 64))); err != nil {
		return fmt.Errorf("totalTimeSeconds: %w", err)
	}
	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("DistanceMeters"),
		xml.CharData(strconv.FormatFloat(a.DistanceMeters, 'g', -1, 64))); err != nil {
		return fmt.Errorf("distanceMeters: %w", err)
	}

	if !math.IsNaN(a.MaximumSpeed) {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("MaximumSpeed"),
			xml.CharData(strconv.FormatFloat(a.MaximumSpeed, 'g', -1, 64))); err != nil {
			return fmt.Errorf("maximumSpeed: %w", err)
		}
	}

	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("Calories"),
		xml.CharData(strconv.FormatUint(uint64(a.Calories), 10))); err != nil {
		return fmt.Errorf("calories: %w", err)
	}

	if a.AverageHeartRateBpm != basetype.Uint8Invalid {
		avgHR := xmlutils.StartElement("AverageHeartRateBpm")
		if err := enc.EncodeToken(avgHR); err != nil {
			return fmt.Errorf("averageHeartRateBpm start: %w", err)
		}
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("Value"),
			xml.CharData(strconv.FormatUint(uint64(a.AverageHeartRateBpm), 10))); err != nil {
			return fmt.Errorf("averageHeartRateBpmValue: %w", err)
		}
		if err := enc.EncodeToken(avgHR.End()); err != nil {
			return fmt.Errorf("averageHeartRateBpm end: %w", err)
		}
	}
	if a.MaximumHeartRateBpm != basetype.Uint8Invalid {
		maxHR := xmlutils.StartElement("MaximumHeartRateBpm")
		if err := enc.EncodeToken(maxHR); err != nil {
			return fmt.Errorf("maximumHeartRateBpm start: %w", err)
		}
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("Value"),
			xml.CharData(strconv.FormatUint(uint64(a.MaximumHeartRateBpm), 10))); err != nil {
			return fmt.Errorf("maximumHeartRateBpmValue: %w", err)
		}
		if err := enc.EncodeToken(maxHR.End()); err != nil {
			return fmt.Errorf("maximumHeartRateBpm end: %w", err)
		}
	}

	if len(a.Intensity) != 0 {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("Intensity"),
			xml.CharData(a.Intensity)); err != nil {
			return fmt.Errorf("intensity: %w", err)
		}
	}

	if a.Cadence != basetype.Uint8Invalid {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("Cadence"),
			xml.CharData(strconv.FormatUint(uint64(a.Cadence), 10))); err != nil {
			return fmt.Errorf("cadence: %w", err)
		}

	}

	if len(a.TriggerMethod) != 0 {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("TriggerMethod"),
			xml.CharData(a.TriggerMethod)); err != nil {
			return fmt.Errorf("triggerMethod: %w", err)
		}
	}

	for i := range a.Tracks {
		if err := a.Tracks[i].MarshalXML(enc, xmlutils.StartElement("Track")); err != nil {
			return fmt.Errorf("track[%d]: %w", i, err)
		}
	}

	if len(a.Notes) != 0 {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("Notes"),
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
