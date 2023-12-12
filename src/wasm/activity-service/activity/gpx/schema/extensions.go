package schema

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/muktihari/openactivity-fit/kit"
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
	Cadence     *uint8
	Distance    *float64
	HeartRate   *uint8
	Temperature *int8
	Power       *uint16
}

var _ xml.Unmarshaler = &TrackPointExtension{}

func (t *TrackPointExtension) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
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
			case "cad", "cadence":
				val, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return err
				}
				t.Cadence = kit.Ptr(uint8(val))
			case "distance":
				val, err := strconv.ParseFloat(string(elem), 64)
				if err != nil {
					return err
				}
				t.Distance = &val
			case "hr", "heartrate":
				val, err := strconv.ParseUint(string(elem), 10, 8)
				if err != nil {
					return err
				}
				t.HeartRate = kit.Ptr(uint8(val))
			case "atemp", "temp", "temperature":
				val, err := strconv.ParseInt(string(elem), 10, 8)
				if err != nil {
					return err
				}
				t.Temperature = kit.Ptr(int8(val))
			case "power":
				val, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return err
				}
				t.Power = kit.Ptr(uint16(val))
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

type garminTrackpoinExtensionV1 struct {
	Atemp *int8  `xml:"gpxtpx:atemp,omitempty"`
	HR    *uint8 `xml:"gpxtpx:hr,omitempty"`
	Cad   *uint8 `xml:"gpxtpx:cad,omitempty"`
}

var _ xml.Marshaler = &garminTrackpoinExtensionV1{}

func (g *garminTrackpoinExtensionV1) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if g.Atemp != nil {
		if err := encodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "gpxtpx:atemp"}},
			xml.CharData(strconv.FormatInt(int64(*g.Atemp), 10))); err != nil {
			return fmt.Errorf("atemp: %w", err)
		}
	}

	if g.HR != nil {
		if err := encodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "gpxtpx:hr"}},
			xml.CharData(strconv.FormatUint(uint64(*g.HR), 10))); err != nil {
			return fmt.Errorf("hr: %w", err)
		}
	}

	if g.Cad != nil {
		if err := encodeElement(enc,
			xml.StartElement{Name: xml.Name{Local: "gpxtpx:cad"}},
			xml.CharData(strconv.FormatUint(uint64(*g.Cad), 10))); err != nil {
			return fmt.Errorf("cad: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

var _ xml.Marshaler = &TrackPointExtension{}

func (t *TrackPointExtension) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	m := garminTrackpoinExtensionV1{
		Atemp: t.Temperature,
		HR:    t.HeartRate,
		Cad:   t.Cadence,
	}

	if m.Atemp == nil && m.HR == nil && m.Cad == nil { // omit
		return nil
	}

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := m.MarshalXML(enc, xml.StartElement{Name: xml.Name{Local: "gpxtpx:TrackPointExtension"}}); err != nil {
		return fmt.Errorf("gpxtpx: %w", err)
	}

	return enc.EncodeToken(se.End())
}
