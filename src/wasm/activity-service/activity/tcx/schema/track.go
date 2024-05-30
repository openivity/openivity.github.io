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
	"math"
	"strconv"
	"time"

	"github.com/muktihari/fit/kit/scaleoffset"
	"github.com/muktihari/fit/kit/semicircles"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/xmlutils"
)

type Track struct {
	Trackpoints []Trackpoint `xml:"Trackpoint,omitempty"`
}

var _ xml.Unmarshaler = (*Track)(nil)

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
					return fmt.Errorf("unmarshal Trackpoint: %w", err)
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

var _ xml.Marshaler = (*Track)(nil)

func (t *Track) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	for i := range t.Trackpoints {
		if err := t.Trackpoints[i].MarshalXML(enc, xmlutils.StartElement("Trackpoint")); err != nil {
			return fmt.Errorf("trackpoint[%d]: %w", i, err)
		}
	}

	return enc.EncodeToken(se.End())
}

type Trackpoint struct {
	Time           time.Time           `xml:"Time"`
	Position       Position            `xml:"Position,omitempty"`
	AltitudeMeters float64             `xml:"AltitudeMeters,omitempty"`
	DistanceMeters float64             `xml:"DistanceMeters,omitempty"`
	HeartRateBpm   uint8               `xml:"HeartRateBpm,omitempty"`
	Cadence        uint8               `xml:"Cadence,omitempty"`
	SensorState    SensorState         `xml:"SensorState,omitempty"`
	Extensions     TrackpointExtension `xml:"Extensions>TPX,omitempty"`
}

func (tp *Trackpoint) reset() {
	tp.Position.reset()
	tp.AltitudeMeters = math.NaN()
	tp.DistanceMeters = math.NaN()
	tp.HeartRateBpm = basetype.Uint8Invalid
	tp.Cadence = basetype.Uint8Invalid
	tp.Extensions.reset()
}

func (t *Trackpoint) ToRecord() activity.Record {
	rec := activity.CreateRecord(nil)

	rec.Timestamp = t.Time
	if !math.IsNaN(t.DistanceMeters) {
		rec.Distance = uint32(scaleoffset.Discard(t.DistanceMeters, 100, 0))
	}
	if !math.IsNaN(t.AltitudeMeters) {
		rec.Altitude = uint16(scaleoffset.Discard(t.AltitudeMeters, 5, 500))
	}
	rec.Cadence = t.Cadence
	rec.HeartRate = t.HeartRateBpm

	if !math.IsNaN(t.Position.LatitudeDegrees) && !math.IsNaN(t.Position.LongitudeDegrees) {
		rec.PositionLat = semicircles.ToSemicircles(t.Position.LatitudeDegrees)
		rec.PositionLong = semicircles.ToSemicircles(t.Position.LongitudeDegrees)
	}

	if !math.IsNaN(t.Extensions.Speed) {
		rec.Speed = uint16(scaleoffset.Discard(t.Extensions.Speed, 1000, 0))
	}

	return rec
}

var _ xml.Unmarshaler = (*Trackpoint)(nil)

func (t *Trackpoint) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	t.reset()

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
				if err := t.Position.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Position: %w", err)
				}
			case "Value":
				targetCharData = targetCharData + "Value"
			case "Extensions":
				if err := t.Extensions.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Extensions: %w", err)
				}
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "Time":
				_time, err := time.Parse(time.RFC3339, string(elem))
				if err != nil {
					return fmt.Errorf("parse Time %q: %w", elem, err)
				}
				t.Time = _time
			case "AltitudeMeters":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse AltitudeMeters %q: %w", elem, err)
				}
				t.AltitudeMeters = f
			case "DistanceMeters":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse DistanceMeters %q: %w", elem, err)
				}
				t.DistanceMeters = f
			case "HeartRateBpm":
				continue
			case "HeartRateBpmValue":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return fmt.Errorf("parse HeartRateBpm %q: %w", elem, err)
				}
				t.HeartRateBpm = uint8(u)
			case "Cadence":
				u, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return fmt.Errorf("parse Cadence %q: %w", elem, err)
				}
				t.Cadence = uint8(u)
			case "SensorState":
				t.SensorState = SensorState(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = (*Trackpoint)(nil)

func (t *Trackpoint) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("Time"),
		xml.CharData(t.Time.Format(time.RFC3339))); err != nil {
		return fmt.Errorf("time: %w", err)
	}

	if !math.IsNaN(t.Position.LatitudeDegrees) && !math.IsNaN(t.Position.LongitudeDegrees) {
		if err := t.Position.MarshalXML(enc, xmlutils.StartElement("Position")); err != nil {
			return fmt.Errorf("position: %w", err)
		}
	}
	if !math.IsNaN(t.AltitudeMeters) {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("AltitudeMeters"),
			xml.CharData(strconv.FormatFloat(t.AltitudeMeters, 'g', -1, 64))); err != nil {
			return fmt.Errorf("altitudeMeters: %w", err)
		}
	}
	if !math.IsNaN(t.DistanceMeters) {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("DistanceMeters"),
			xml.CharData(strconv.FormatFloat(t.DistanceMeters, 'g', -1, 64))); err != nil {
			return fmt.Errorf("distanceMeters: %w", err)
		}
	}
	if t.HeartRateBpm != basetype.Uint8Invalid {
		hr := xmlutils.StartElement("HeartRateBpm")
		if err := enc.EncodeToken(hr); err != nil {
			return fmt.Errorf("heartRateBpm start: %w", err)
		}
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("Value"),
			xml.CharData(strconv.FormatUint(uint64(t.HeartRateBpm), 10))); err != nil {
			return fmt.Errorf("heartRateBpmValue: %w", err)
		}
		if err := enc.EncodeToken(hr.End()); err != nil {
			return fmt.Errorf("heartRateBpm end: %w", err)
		}
	}
	if t.Cadence != basetype.Uint8Invalid {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("Cadence"),
			xml.CharData(strconv.FormatUint(uint64(t.Cadence), 10))); err != nil {
			return fmt.Errorf("cadence: %w", err)
		}
	}

	if len(t.SensorState) != 0 {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("SensorState"),
			xml.CharData(t.SensorState)); err != nil {
			return fmt.Errorf("sensorState: %w", err)
		}
	}

	if !math.IsNaN(t.Extensions.Speed) {
		if err := t.Extensions.MarshalXML(enc, xmlutils.StartElement("Extensions")); err != nil {
			return fmt.Errorf("extension: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

type Position struct {
	LatitudeDegrees  float64 // -90.0 to 90.0
	LongitudeDegrees float64 // -180.0 to 180.0
}

func (p *Position) reset() {
	p.LatitudeDegrees = math.NaN()
	p.LongitudeDegrees = math.NaN()
}

var _ xml.Unmarshaler = (*Position)(nil)

func (p *Position) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	p.reset()

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
					return fmt.Errorf("parse LatitudeDegrees: %w", err)
				}
				p.LatitudeDegrees = f
			case "LongitudeDegrees":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse LongitudeDegrees: %w", err)
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

var _ xml.Marshaler = (*Position)(nil)

func (p *Position) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if math.IsNaN(p.LatitudeDegrees) || math.IsNaN(p.LongitudeDegrees) { // omit
		return nil
	}

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc,
		xml.StartElement{Name: xml.Name{Local: "LatitudeDegrees"}},
		xml.CharData(strconv.FormatFloat(p.LatitudeDegrees, 'g', -1, 64))); err != nil {
		return fmt.Errorf("latitudeDegrees: %w", err)
	}

	if err := xmlutils.EncodeElement(enc,
		xml.StartElement{Name: xml.Name{Local: "LongitudeDegrees"}},
		xml.CharData(strconv.FormatFloat(p.LongitudeDegrees, 'g', -1, 64))); err != nil {
		return fmt.Errorf("longitudeDegrees: %w", err)
	}

	return enc.EncodeToken(se.End())
}

type SensorState string

const (
	SensorStatePresent SensorState = "Present"
	SensorStateAbsent  SensorState = "Absent"
)

type TrackpointExtension struct {
	Speed float64 `xml:"Speed,omitempty"`
}

func (t *TrackpointExtension) reset() {
	t.Speed = math.NaN()
}

var _ xml.Unmarshaler = (*TrackpointExtension)(nil)

func (t *TrackpointExtension) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	t.reset()

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
			case "Speed":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parse Speed: %w", err)
				}
				t.Speed = f
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = (*TrackpointExtension)(nil)

func (tpe *TrackpointExtension) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	tpx := xmlutils.StartElement("ns1:TPX")
	if err := enc.EncodeToken(tpx); err != nil {
		return fmt.Errorf("tpx: %w", err)
	}

	if !math.IsNaN(tpe.Speed) {
		if err := xmlutils.EncodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "ns1:Speed"}},
			xml.CharData(strconv.FormatFloat(tpe.Speed, 'g', -1, 64))); err != nil {
			return fmt.Errorf("speed: %w", err)
		}
	}

	if err := enc.EncodeToken(tpx.End()); err != nil {
		return fmt.Errorf("tpx: %w", err)
	}

	return enc.EncodeToken(se.End())
}
