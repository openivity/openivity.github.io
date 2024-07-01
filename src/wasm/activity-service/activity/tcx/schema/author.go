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
	"strconv"

	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/xmlutils"
)

type Application struct {
	Name       string `xml:"Name"`
	Build      *Build `xml:"Build"`
	LangID     string `xml:"LangID"`     // Specifies the two character ISO 693-1 language id that identifies the installed language of this application. see http://www.loc.gov/standards/iso639-2/ for appropriate ISO identifiers
	PartNumber string `xml:"PartNumber"` // The formatted XXX-XXXXX-XX Garmin part number of a PC application.
}

func (a *Application) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
		case "Build":
			var build Build
			se := xmltokenizer.GetToken().Copy(token)
			err = build.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Build: %w", err)
			}
			a.Build = &build
		case "Name":
			a.Name = string(token.Data)
		case "LangID":
			a.LangID = string(token.Data)
		case "PartNumber":
			a.PartNumber = string(token.Data)
		}
	}

	return nil
}

func (a *Application) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("Name"), xml.CharData(a.Name)); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	if a.Build != nil {
		if err := a.Build.MarshalXML(enc, xmlutils.StartElement("Build")); err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("LangID"), xml.CharData(a.LangID)); err != nil {
		return fmt.Errorf("langID: %w", err)
	}
	if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("PartNumber"), xml.CharData(a.PartNumber)); err != nil {
		return fmt.Errorf("partNumber: %w", err)
	}

	return enc.EncodeToken(se.End())
}

type Build struct {
	Type    BuildType
	Version *Version
}

func (b *Build) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
		case "Version":
			var version Version
			se := xmltokenizer.GetToken().Copy(token)
			err = version.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Version: %w", err)
			}
			b.Version = &version
		case "Type":
			b.Type = BuildType(token.Data)
		}
	}

	return nil
}

var _ xml.Marshaler = (*Build)(nil)

func (b *Build) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("Type"), xml.CharData(b.Type)); err != nil {
		return fmt.Errorf("type: %w", err)
	}

	if b.Version != nil {
		if err := b.Version.MarshalXML(enc, xmlutils.StartElement("Version")); err != nil {
			return fmt.Errorf("version: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

type BuildType string

const (
	BuildTypeInternal BuildType = "Internal"
	BuildTypeAlpha    BuildType = "Alpha"
	BuildTypeBeta     BuildType = "Beta"
	BuildTypeRelease  BuildType = "Release"
)

type Device struct {
	Name      string   `xml:"Name"`
	UnitId    uint32   `xml:"UnitId"`
	ProductID uint16   `xml:"ProductId"`
	Version   *Version `xml:"Version"`
}

func (d *Device) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
		case "Version":
			var version Version
			se := xmltokenizer.GetToken().Copy(token)
			err = version.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Version: %w", err)
			}
			d.Version = &version
		case "Name":
			d.Name = string(token.Data)
		case "UnitId":
			u, err := strconv.ParseUint(string(token.Data), 10, 32)
			if err != nil {
				return fmt.Errorf("parse UnitId: %w", err)
			}
			d.UnitId = uint32(u)
		case "ProductId":
			u, err := strconv.ParseUint(string(token.Data), 10, 16)
			if err != nil {
				return fmt.Errorf("parse ProductId: %w", err)
			}
			d.ProductID = uint16(u)
		}
	}

	return nil
}

var _ xml.Marshaler = (*Device)(nil)

func (d *Device) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("Name"),
		xml.CharData(d.Name)); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("UnitId"),
		xml.CharData(strconv.FormatUint(uint64(d.UnitId), 10))); err != nil {
		return fmt.Errorf("unitId: %w", err)
	}
	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("ProductID"),
		xml.CharData(strconv.FormatUint(uint64(d.ProductID), 10))); err != nil {
		return fmt.Errorf("productID: %w", err)
	}

	if d.Version != nil {
		if err := d.Version.MarshalXML(enc, xmlutils.StartElement("Version")); err != nil {
			return fmt.Errorf("version: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}

type Version struct {
	VersionMajor uint16 `xml:"VersionMajor"`
	VersionMinor uint16 `xml:"VersionMinor"`
	BuildMajor   uint16 `xml:"BuildMajor"`
	BuildMinor   uint16 `xml:"BuildMinor"`
}

func (v *Version) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
		case "VersionMajor":
			u, err := strconv.ParseUint(string(token.Data), 10, 16)
			if err != nil {
				return fmt.Errorf("parse VersionMajor: %w", err)
			}
			v.VersionMajor = uint16(u)
		case "VersionMinor":
			u, err := strconv.ParseUint(string(token.Data), 10, 16)
			if err != nil {
				return fmt.Errorf("parse VersionMinor: %w", err)
			}
			v.VersionMinor = uint16(u)
		case "BuildMajor":
			u, err := strconv.ParseUint(string(token.Data), 10, 16)
			if err != nil {
				return fmt.Errorf("parse BuildMajor: %w", err)
			}
			v.BuildMajor = uint16(u)
		case "BuildMinor":
			u, err := strconv.ParseUint(string(token.Data), 10, 16)
			if err != nil {
				return fmt.Errorf("parse BuildMinor: %w", err)
			}
			v.BuildMinor = uint16(u)
		}
	}

	return nil
}

var _ xml.Marshaler = (*Version)(nil)

func (v *Version) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("VersionMajor"),
		xml.CharData(strconv.FormatUint(uint64(v.VersionMajor), 10))); err != nil {
		return fmt.Errorf("versionMajor:%w", err)
	}
	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("VersionMinor"),
		xml.CharData(strconv.FormatUint(uint64(v.VersionMinor), 10))); err != nil {
		return fmt.Errorf("versionMinor:%w", err)
	}
	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("BuildMajor"),
		xml.CharData(strconv.FormatUint(uint64(v.BuildMajor), 10))); err != nil {
		return fmt.Errorf("buildMajor:%w", err)
	}
	if err := xmlutils.EncodeElement(enc,
		xmlutils.StartElement("BuildMinor"),
		xml.CharData(strconv.FormatUint(uint64(v.BuildMinor), 10))); err != nil {
		return fmt.Errorf("buildMinor:%w", err)
	}

	return enc.EncodeToken(se.End())
}
