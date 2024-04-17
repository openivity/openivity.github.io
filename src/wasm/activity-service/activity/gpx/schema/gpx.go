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
	"strings"
)

const (
	xmlns       = "http://www.topografix.com/GPX/1/1"
	xmlnsxsi    = "http://www.w3.org/2001/XMLSchema-instance"
	xmlnsgpxtpx = "http://www.garmin.com/xmlschemas/TrackPointExtension/v1"
	xmlnsgpxx   = "http://www.garmin.com/xmlschemas/GpxExtensions/v3"
	Version     = "1.1"
)

var schemaLocations = [...]string{
	"http://www.topografix.com/GPX/1/1",
	"http://www.topografix.com/GPX/1/1/gpx.xsd",
	"http://www.garmin.com/xmlschemas/GpxExtensions/v3",
	"http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd",
	"http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
	"http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd",
}

// GPX is GPX schema (simplified).
//
// Note: Please define xml.Unmarshaler for each struct involved to avoid reflection as much as we can.
type GPX struct {
	XMLName xml.Name `xml:"gpx"`
	Creator string   `xml:"creator,attr"`
	Version string   `xml:"version,attr"`

	Metadata Metadata `xml:"metadata,omitempty"`
	Tracks   []Track  `xml:"trk,omitempty"`
}

var _ xml.Unmarshaler = &GPX{}

func (g *GPX) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	for i := range se.Attr {
		attr := se.Attr[i]

		switch attr.Name.Local {
		case "creator":
			g.Creator = attr.Value
		case "version":
			g.Version = attr.Value
		}
	}

	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "metadata":
				var metadata Metadata
				if err := metadata.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				g.Metadata = metadata
			case "trk":
				var track Track
				if err := track.UnmarshalXML(dec, elem); err != nil {
					return err
				}
				g.Tracks = append(g.Tracks, track)
			}

		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

func (g *GPX) Validate() error {
	if g == nil {
		return fmt.Errorf("%T is nil", g)
	}

	if err := g.Metadata.Validate(); err != nil {
		return fmt.Errorf("validate metadata: %w", err)
	}

	for i, track := range g.Tracks {
		if err := track.Validate(); err != nil {
			return fmt.Errorf("tracks[%d]: %w", i, err)
		}
	}
	if g.Creator == "" {
		return fmt.Errorf("creator is required but provided empty")
	}
	return nil
}

var _ xml.Marshaler = &GPX{}

func (g *GPX) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	version := g.Version
	if version == "" {
		version = Version
	}

	se.Name = xml.Name{Local: "gpx"}
	se.Attr = []xml.Attr{
		{Name: xml.Name{Local: "creator"}, Value: g.Creator},
		{Name: xml.Name{Local: "version"}, Value: version},
		{Name: xml.Name{Local: "xmlns"}, Value: xmlns},
		{Name: xml.Name{Local: "xmlns:xsi"}, Value: xmlnsxsi},
		{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: strings.Join(schemaLocations[:], " ")},
		{Name: xml.Name{Local: "xmlns:gpxtpx"}, Value: xmlnsgpxtpx},
		{Name: xml.Name{Local: "xmlns:gpxx"}, Value: xmlnsgpxx},
	}

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := g.Metadata.MarshalXML(enc, xml.StartElement{Name: xml.Name{Local: "metadata"}}); err != nil {
		return fmt.Errorf("metadata: %w", err)
	}

	for i := range g.Tracks {
		if err := g.Tracks[i].MarshalXML(enc, xml.StartElement{Name: xml.Name{Local: "trk"}}); err != nil {
			return fmt.Errorf("trk[%d]: %w", i, err)
		}
	}

	return enc.EncodeToken(se.End())
}
