package schema

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
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

func (ts *TrackSegment) Validate() error {
	if ts == nil {
		return nil
	}
	for i, trackpoint := range ts.Trackpoints {
		if err := trackpoint.Validate(); err != nil {
			return fmt.Errorf("trackpoints[%d]: %w", i, err)
		}
	}
	return nil
}

type Waypoint struct {
	Lat                 *float64             `xml:"lat,attr,omitempty"`
	Lon                 *float64             `xml:"lon,attr,omitempty"`
	Ele                 *float64             `xml:"ele,omitempty"`
	Time                time.Time            `xml:"time,omitempty"`
	TrackPointExtension *TrackPointExtension `xml:"extensions>TrackPointExtension,omitempty"`
}

var _ xml.Unmarshaler = &Waypoint{}

func (w *Waypoint) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for i := range se.Attr {
		attr := se.Attr[i]

		// All attrs are float64
		f, err := strconv.ParseFloat(attr.Value, 64)
		if err != nil {
			return err
		}

		switch attr.Name.Local {
		case "lat":
			w.Lat = &f
		case "lon":
			w.Lon = &f
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
				var ext TrackPointExtension
				if err := ext.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				w.TrackPointExtension = &ext
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
				w.Ele = &f
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
	if w.Lat != nil {
		if *w.Lat < -90 || *w.Lat > 90 {
			return fmt.Errorf("lat %f is outside range -90.0 to 90.0", *w.Lat)
		}
	}
	if w.Lon != nil {
		if *w.Lon < -180 || *w.Lon > 180 {
			return fmt.Errorf("lon %f is outside range -180.0 to 180.0", *w.Lon)
		}
	}
	return nil
}
