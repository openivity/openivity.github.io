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
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/xmlutils"
)

type Track struct {
	Name          string         `xml:"name,omitempty"`
	Type          string         `xml:"type,omitempty"`
	TrackSegments []TrackSegment `xml:"trkseg,omitempty"`
}

var _ xml.Unmarshaler = &Track{}

func (t *Track) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "trkseg":
				var trkseg TrackSegment
				if err := trkseg.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				t.TrackSegments = append(t.TrackSegments, trkseg)
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "name":
				t.Name = string(elem)
			case "type":
				t.Type = string(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
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

var _ xml.Marshaler = &Track{}

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

var _ xml.Unmarshaler = &TrackSegment{}

func (t *TrackSegment) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "trkpt":
				var trkpt Waypoint
				if err := trkpt.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				t.Trackpoints = append(t.Trackpoints, trkpt)
			}
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
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

var _ xml.Marshaler = &TrackSegment{}

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

var _ xml.Unmarshaler = &Waypoint{}

func (w *Waypoint) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	w.reset()

	for i := range se.Attr {
		attr := se.Attr[i]

		switch attr.Name.Local {
		case "lat":
			f, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			w.Lat = f
		case "lon":
			f, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			w.Lon = f
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
			case "extensions":
				if err := w.TrackPointExtension.UnmarshalXML(dec, elem); err != nil {
					return err
				}
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "ele":
				f, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return fmt.Errorf("parsing ele:  %w", err)
				}
				w.Ele = f
			case "time":
				t, err := time.Parse(time.RFC3339, string(elem))
				if err != nil {
					return fmt.Errorf("parsing time:  %w", err)
				}
				w.Time = t
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
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

var _ xml.Marshaler = &Waypoint{}

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
