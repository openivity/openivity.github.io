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

	"github.com/muktihari/fit/kit/scaleoffset"
	"github.com/muktihari/fit/kit/semicircles"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/xmlutils"
)

type Track struct {
	Trackpoints []Trackpoint `xml:"Trackpoint,omitempty"`
}

func (t *Track) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
		if token.IsEndElement() {
			continue
		}

		switch string(token.Name.Local) {
		case "Trackpoint":
			var trackpoint Trackpoint
			se := xmltokenizer.GetToken().Copy(token)
			err = trackpoint.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Trackpoint: %w", err)
			}
			t.Trackpoints = append(t.Trackpoints, trackpoint)
		}
	}

	return nil
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
	HeartRateBpm   uint8               `xml:"HeartRateBpm>Value,omitempty"`
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

func (t *Trackpoint) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	t.reset()

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
		if token.IsEndElement() {
			continue
		}

		switch string(token.Name.Local) {
		case "Position":
			se := xmltokenizer.GetToken().Copy(token)
			err = t.Position.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Position: %w", err)
			}
		case "Extensions":
			se := xmltokenizer.GetToken().Copy(token)
			err = t.Extensions.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Extensions: %w", err)
			}
		case "Time":
			t.Time, err = time.Parse(time.RFC3339, string(token.Data))
			if err != nil {
				return fmt.Errorf("parse Time %q: %w", token.Data, err)
			}
		case "AltitudeMeters":
			t.AltitudeMeters, err = strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse AltitudeMeters %q: %w", token.Data, err)
			}
		case "DistanceMeters":
			t.DistanceMeters, err = strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse DistanceMeters %q: %w", token.Data, err)
			}
		case "HeartRateBpm":
			token, err = getValueToken(tok)
			if err != nil {
				return err
			}
			u, err := strconv.ParseUint(string(token.Data), 10, 8)
			if err != nil {
				return fmt.Errorf("parse HeartRateBpm %q: %w", token.Data, err)
			}
			t.HeartRateBpm = uint8(u)
		case "Cadence":
			u, err := strconv.ParseUint(string(token.Data), 10, 8)
			if err != nil {
				return fmt.Errorf("parse Cadence %q: %w", token.Data, err)
			}
			t.Cadence = uint8(u)
		case "SensorState":
			t.SensorState = SensorState(token.Data)
		}
	}

	return nil
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

func (p *Position) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	p.reset()

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
		if token.IsEndElement() {
			continue
		}

		switch string(token.Name.Local) {
		case "LatitudeDegrees":
			p.LatitudeDegrees, err = strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse LatitudeDegrees: %w", err)
			}
		case "LongitudeDegrees":
			p.LongitudeDegrees, err = strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse LongitudeDegrees: %w", err)
			}
		}
	}

	return nil
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

func (t *TrackpointExtension) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	t.reset()

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
		if token.IsEndElement() {
			continue
		}

		switch string(token.Name.Local) {
		case "Speed":
			t.Speed, err = strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("parse Speed: %w", err)
			}
		}
	}

	return nil
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
