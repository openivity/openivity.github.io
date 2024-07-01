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
	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/xmlutils"
)

type Track struct {
	Name          string         `xml:"name,omitempty"`
	Type          string         `xml:"type,omitempty"`
	TrackSegments []TrackSegment `xml:"trkseg,omitempty"`
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
			return nil
		}
		if token.IsEndElement {
			continue
		}

		switch string(token.Name.Local) {
		case "name":
			t.Name = string(token.Data)
		case "type":
			t.Type = string(token.Data)
		case "trkseg":
			var trkseg TrackSegment
			se := xmltokenizer.GetToken().Copy(token)
			err = trkseg.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return err
			}
			t.TrackSegments = append(t.TrackSegments, trkseg)
		}
	}

	return nil
}

func (t *Track) Validate() error {
	if t == nil {
		return nil
	}
	for i, trackSegment := range t.TrackSegments {
		if err := trackSegment.Validate(); err != nil {
			return fmt.Errorf("trackSegments[%d]: %w", i, err)
		}
	}
	return nil
}

var _ xml.Marshaler = (*Track)(nil)

func (t *Track) MarshalXML(enc *xml.Encoder, se xml.StartElement) (err error) {
	if err = enc.EncodeToken(se); err != nil {
		return err
	}

	if len(t.Name) != 0 {
		if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("name"), xml.CharData(t.Name)); err != nil {
			return fmt.Errorf("name: %w", err)
		}
	}

	if len(t.Type) != 0 {
		if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("type"), xml.CharData(t.Name)); err != nil {
			return fmt.Errorf("type: %w", err)
		}
	}

	for i := range t.TrackSegments {
		if err := t.TrackSegments[i].MarshalXML(enc, xmlutils.StartElement("trkseg")); err != nil {
			return fmt.Errorf("trkseg[%d]: %w", i, err)
		}
	}

	return enc.EncodeToken(se.End())
}

type TrackSegment struct {
	Trackpoints []Waypoint `xml:"trkpt,omitempty"`
}

func (t *TrackSegment) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
		case "trkpt":
			var trkpt Waypoint
			se := xmltokenizer.GetToken().Copy(token)
			err = trkpt.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("trkpt: %w", err)
			}
			t.Trackpoints = append(t.Trackpoints, trkpt)
		}
	}

	return nil
}

func (t *TrackSegment) Validate() error {
	if t == nil {
		return nil
	}
	for i := range t.Trackpoints {
		if err := t.Trackpoints[i].Validate(); err != nil {
			return fmt.Errorf("trackpoints[%d]: %w", i, err)
		}
	}
	return nil
}

var _ xml.Marshaler = (*TrackSegment)(nil)

func (t *TrackSegment) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	for i := range t.Trackpoints {
		if err := t.Trackpoints[i].MarshalXML(enc, xmlutils.StartElement("trkpt")); err != nil {
			return fmt.Errorf("trkpt[%d]: %w", i, err)
		}
	}

	return enc.EncodeToken(se.End())
}

type Waypoint struct {
	Lat                 float64             `xml:"lat,attr,omitempty"`
	Lon                 float64             `xml:"lon,attr,omitempty"`
	Ele                 float64             `xml:"ele,omitempty"`
	Time                time.Time           `xml:"time,omitempty"`
	TrackPointExtension TrackPointExtension `xml:"extensions>TrackPointExtension,omitempty"`
}

func (w *Waypoint) reset() {
	w.Lat = math.NaN()
	w.Lon = math.NaN()
	w.Ele = math.NaN()
	w.Time = time.Time{}
	w.TrackPointExtension.reset()
}

func (w *Waypoint) ToRecord() activity.Record {
	rec := activity.CreateRecord(nil)

	rec.Timestamp = w.Time
	if !math.IsNaN(w.Lat) {
		rec.PositionLat = semicircles.ToSemicircles(w.Lat)
	}
	if !math.IsNaN(w.Lon) {
		rec.PositionLong = semicircles.ToSemicircles(w.Lon)
	}
	if !math.IsNaN(w.Ele) {
		rec.Altitude = uint16(scaleoffset.Discard(w.Ele, 5, 500))
	}

	ext := w.TrackPointExtension
	if !math.IsNaN(ext.Distance) {
		rec.Distance = uint32(scaleoffset.Discard(ext.Distance, 100, 0))
	}
	rec.Cadence = ext.Cadence
	rec.HeartRate = ext.HeartRate
	rec.Power = ext.Power
	rec.Temperature = ext.Temperature

	return rec
}

func (w *Waypoint) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	w.reset()

	var err error
	for i := range se.Attrs {
		attr := &se.Attrs[i]
		switch string(attr.Name.Local) {
		case "lat":
			w.Lat, err = strconv.ParseFloat(string(attr.Value), 64)
			if err != nil {
				return fmt.Errorf("lat: %w", err)
			}
		case "lon":
			w.Lon, err = strconv.ParseFloat(string(attr.Value), 64)
			if err != nil {
				return fmt.Errorf("lon: %w", err)
			}
		}
	}

	for {
		token, err := tok.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("waypoint: %w", err)
		}

		if token.IsEndElementOf(se) {
			return nil
		}
		if token.IsEndElement {
			continue
		}

		switch string(token.Name.Local) {
		case "ele":
			w.Ele, err = strconv.ParseFloat(string(token.Data), 64)
			if err != nil {
				return fmt.Errorf("ele: %w", err)
			}
		case "time":
			w.Time, err = time.Parse(time.RFC3339, string(token.Data))
			if err != nil {
				return fmt.Errorf("time: %w", err)
			}
		case "extensions":
			se := xmltokenizer.GetToken().Copy(token)
			err = w.TrackPointExtension.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("extensions: %w", err)
			}
		}
	}

	return nil
}

func (w *Waypoint) Validate() error {
	if !math.IsNaN(w.Lat) {
		if w.Lat < -90 || w.Lat > 90 {
			return fmt.Errorf("lat %f is outside range -90.0 to 90.0", w.Lat)
		}
	}
	if !math.IsNaN(w.Lon) {
		if w.Lon < -180 || w.Lon > 180 {
			return fmt.Errorf("lon %f is outside range -180.0 to 180.0", w.Lon)
		}
	}
	return nil
}

var _ xml.Marshaler = (*Waypoint)(nil)

func (w *Waypoint) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if !math.IsNaN(w.Lat) {
		se.Attr = append(se.Attr, xml.Attr{
			Name:  xml.Name{Local: "lat"},
			Value: strconv.FormatFloat(w.Lat, 'g', -1, 64)})
	}
	if !math.IsNaN(w.Lon) {
		se.Attr = append(se.Attr, xml.Attr{
			Name:  xml.Name{Local: "lon"},
			Value: strconv.FormatFloat(w.Lon, 'g', -1, 64)})
	}

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if !math.IsNaN(w.Ele) {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("ele"),
			xml.CharData(strconv.FormatFloat(w.Ele, 'g', -1, 64))); err != nil {
			return fmt.Errorf("ele: %w", err)
		}
	}

	if !w.Time.IsZero() {
		if err := xmlutils.EncodeElement(enc,
			xmlutils.StartElement("time"),
			xml.CharData(w.Time.Format(time.RFC3339))); err != nil {
			return fmt.Errorf("time: %w", err)
		}
	}

	if err := w.TrackPointExtension.MarshalXML(enc, xmlutils.StartElement("extensions")); err != nil {
		return fmt.Errorf("extensions: %w", err)
	}

	return enc.EncodeToken(se.End())
}
