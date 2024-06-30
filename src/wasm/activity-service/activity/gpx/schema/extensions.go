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

	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/xmlutils"
)

// TrackPointExtension is a GPX extension for health-related data.
//
// We accept unmarshaling values from these following schema:
//  1. Garmin’s Track Point Extension v1: cad, atemp, hr (ref: http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd)
//  2. Cluetrust GPX extension: cadence, distance, hr, temp (ref: http://www.cluetrust.com/Schemas/gpxdata10.xsd)
//  3. Generic: cadence, distance, heartrate, temperature, power.
//
// However, we will marshal into Garmin’s Track Point Extension v1 schema so some fields will be omitted.
type TrackPointExtension struct {
	Cadence     uint8
	Distance    float64
	HeartRate   uint8
	Temperature int8
	Power       uint16
}

func (t *TrackPointExtension) reset() {
	t.Cadence = basetype.Uint8Invalid
	t.Distance = math.NaN()
	t.HeartRate = basetype.Uint8Invalid
	t.Temperature = basetype.Sint8Invalid
	t.Power = basetype.Uint16Invalid
}

func (t *TrackPointExtension) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
			return nil
		}
		if token.IsEndElement {
			continue
		}

		switch string(token.Name.Local) {
		case "cad", "cadence":
			val, err := strconv.ParseUint(string(token.Data), 10, 8)
			if err != nil {
				return err
			}
			t.Cadence = uint8(val)
		case "distance":
			val, err := strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return err
			}
			t.Distance = val
		case "hr", "heartrate":
			val, err := strconv.ParseUint(string(token.Data), 10, 8)
			if err != nil {
				return err
			}
			t.HeartRate = uint8(val)
		case "atemp", "temp", "temperature":
			val, err := strconv.ParseInt(string(token.Data), 10, 8)
			if err != nil {
				return err
			}
			t.Temperature = int8(val)
		case "power":
			val, err := strconv.ParseUint(string(token.Data), 10, 16)
			if err != nil {
				return err
			}
			t.Power = uint16(val)
		}
	}

	return nil
}

type garminTrackpoinExtensionV1 struct {
	Atemp int8  `xml:"gpxtpx:atemp,omitempty"`
	HR    uint8 `xml:"gpxtpx:hr,omitempty"`
	Cad   uint8 `xml:"gpxtpx:cad,omitempty"`
}

var _ xml.Marshaler = (*garminTrackpoinExtensionV1)(nil)

func (g *garminTrackpoinExtensionV1) MarshalXML(enc *xml.Encoder, se xml.StartElement) (err error) {
	if err = enc.EncodeToken(se); err != nil {
		return err
	}

	if g.Atemp != basetype.Sint8Invalid {
		if err = xmlutils.EncodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "gpxtpx:atemp"}},
			xml.CharData(strconv.FormatInt(int64(g.Atemp), 10))); err != nil {
			return fmt.Errorf("atemp: %w", err)
		}
	}

	if g.HR != basetype.Uint8Invalid {
		if err = xmlutils.EncodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "gpxtpx:hr"}},
			xml.CharData(strconv.FormatUint(uint64(g.HR), 10))); err != nil {
			return fmt.Errorf("hr: %w", err)
		}
	}

	if g.Cad != basetype.Uint8Invalid {
		if err = xmlutils.EncodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "gpxtpx:cad"}},
			xml.CharData(strconv.FormatUint(uint64(g.Cad), 10))); err != nil {
			return fmt.Errorf("cad: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

var _ xml.Marshaler = (*TrackPointExtension)(nil)

func (t *TrackPointExtension) MarshalXML(enc *xml.Encoder, se xml.StartElement) (err error) {
	m := garminTrackpoinExtensionV1{
		Atemp: t.Temperature,
		HR:    t.HeartRate,
		Cad:   t.Cadence,
	}

	if m.Atemp == basetype.Sint8Invalid && m.HR == basetype.Uint8Invalid && m.Cad == basetype.Uint8Invalid { // omit
		return nil
	}

	if err = enc.EncodeToken(se); err != nil {
		return err
	}

	if err = m.MarshalXML(enc, xml.StartElement{Name: xml.Name{Local: "gpxtpx:TrackPointExtension"}}); err != nil {
		return fmt.Errorf("gpxtpx: %w", err)
	}

	return enc.EncodeToken(se.End())
}
