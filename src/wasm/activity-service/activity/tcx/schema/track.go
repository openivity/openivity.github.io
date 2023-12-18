package schema

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/kit"
	kxml "github.com/muktihari/openactivity-fit/kit/xml"
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

var _ xml.Marshaler = &Track{}

func (t *Track) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	for i := range t.Trackpoints {
		if err := t.Trackpoints[i].MarshalXML(enc, kxml.StartElement("Trackpoint")); err != nil {
			return fmt.Errorf("trackpoint[%d]: %w", i, err)
		}
	}

	return enc.EncodeToken(se.End())
}

type Trackpoint struct {
	Time           time.Time            `xml:"Time"`
	Position       *Position            `xml:"Position,omitempty"`
	AltitudeMeters *float64             `xml:"AltitudeMeters,omitempty"`
	DistanceMeters *float64             `xml:"DistanceMeters,omitempty"`
	HeartRateBpm   *uint8               `xml:"HeartRateBpm,omitempty"`
	Cadence        *uint8               `xml:"Cadence,omitempty"`
	SensorState    SensorState          `xml:"SensorState,omitempty"`
	Extensions     *TrackpointExtension `xml:"Extensions>TPX,omitempty"`
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
					return fmt.Errorf("unmarshal Position: %w", err)
				}
				tp.Position = &position
			case "Value":
				targetCharData = targetCharData + "Value"
			case "Extensions":
				var trackpointExtension TrackpointExtension
				if err := trackpointExtension.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Extensions: %w", err)
				}
				tp.Extensions = &trackpointExtension
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
			case "HeartRateBpm":
				continue
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

var _ xml.Marshaler = &Trackpoint{}

func (tp *Trackpoint) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc,
		kxml.StartElement("Time"),
		xml.CharData(tp.Time.Format(time.RFC3339))); err != nil {
		return fmt.Errorf("time: %w", err)
	}

	if tp.Position != nil {
		if err := tp.Position.MarshalXML(enc, kxml.StartElement("Position")); err != nil {
			return fmt.Errorf("position: %w", err)
		}
	}
	if tp.AltitudeMeters != nil {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("AltitudeMeters"),
			xml.CharData(strconv.FormatFloat(*tp.AltitudeMeters, 'g', -1, 64))); err != nil {
			return fmt.Errorf("altitudeMeters: %w", err)
		}
	}
	if tp.DistanceMeters != nil {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("DistanceMeters"),
			xml.CharData(strconv.FormatFloat(*tp.DistanceMeters, 'g', -1, 64))); err != nil {
			return fmt.Errorf("distanceMeters: %w", err)
		}
	}
	if tp.HeartRateBpm != nil {
		hr := kxml.StartElement("HeartRateBpm")
		if err := enc.EncodeToken(hr); err != nil {
			return fmt.Errorf("heartRateBpm start: %w", err)
		}
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("Value"),
			xml.CharData(strconv.FormatUint(uint64(*tp.HeartRateBpm), 10))); err != nil {
			return fmt.Errorf("heartRateBpmValue: %w", err)
		}
		if err := enc.EncodeToken(hr.End()); err != nil {
			return fmt.Errorf("heartRateBpm end: %w", err)
		}
	}
	if tp.Cadence != nil {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("Cadence"),
			xml.CharData(strconv.FormatUint(uint64(*tp.Cadence), 10))); err != nil {
			return fmt.Errorf("cadence: %w", err)
		}
	}

	if len(tp.SensorState) != 0 {
		if err := kxml.EncodeElement(enc,
			kxml.StartElement("SensorState"),
			xml.CharData(tp.SensorState)); err != nil {
			return fmt.Errorf("sensorState: %w", err)
		}
	}

	if tp.Extensions != nil {
		if err := tp.Extensions.MarshalXML(enc, kxml.StartElement("Extensions")); err != nil {
			return fmt.Errorf("extension: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
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

var _ xml.Marshaler = &Position{}

func (p *Position) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc,
		xml.StartElement{Name: xml.Name{Local: "LatitudeDegrees"}},
		xml.CharData(strconv.FormatFloat(p.LatitudeDegrees, 'g', -1, 64))); err != nil {
		return fmt.Errorf("latitudeDegrees: %w", err)
	}

	if err := kxml.EncodeElement(enc,
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
	Speed *float64 `xml:"Speed,omitempty"`
}

var _ xml.Unmarshaler = &TrackpointExtension{}

func (tpe *TrackpointExtension) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
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
				tpe.Speed = &f
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &TrackpointExtension{}

func (tpe *TrackpointExtension) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	tpx := kxml.StartElement("ns1:TPX")
	if err := enc.EncodeToken(tpx); err != nil {
		return fmt.Errorf("tpx: %w", err)
	}

	if tpe.Speed != nil {
		if err := kxml.EncodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "ns1:Speed"}},
			xml.CharData(strconv.FormatFloat(*tpe.Speed, 'g', -1, 64))); err != nil {
			return fmt.Errorf("speed: %w", err)
		}
	}

	if err := enc.EncodeToken(tpx.End()); err != nil {
		return fmt.Errorf("tpx: %w", err)
	}

	return enc.EncodeToken(se.End())
}
