package schema

import (
	"encoding/xml"
	"fmt"
	"strconv"

	kxml "github.com/muktihari/openactivity-fit/kit/xml"
)

type Application struct {
	Name       string `xml:"Name"`
	Build      *Build `xml:"Build"`
	LangID     string `xml:"LangID"`     // Specifies the two character ISO 693-1 language id that identifies the installed language of this application. see http://www.loc.gov/standards/iso639-2/ for appropriate ISO identifiers
	PartNumber string `xml:"PartNumber"` // The formatted XXX-XXXXX-XX Garmin part number of a PC application.
}

var _ xml.Unmarshaler = &Application{}

func (a *Application) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "Build":
				var build Build
				if err := build.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Build: %w", err)
				}
				a.Build = &build

			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "Name":
				a.Name = string(elem)
			case "LangID":
				a.LangID = string(elem)
			case "PartNumber":
				a.PartNumber = string(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &Application{}

func (a *Application) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc, kxml.StartElement("Name"), xml.CharData(a.Name)); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	if a.Build != nil {
		if err := a.Build.MarshalXML(enc, kxml.StartElement("Build")); err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	if err := kxml.EncodeElement(enc, kxml.StartElement("LangID"), xml.CharData(a.LangID)); err != nil {
		return fmt.Errorf("langID: %w", err)
	}
	if err := kxml.EncodeElement(enc, kxml.StartElement("PartNumber"), xml.CharData(a.PartNumber)); err != nil {
		return fmt.Errorf("partNumber: %w", err)
	}

	return enc.EncodeToken(se.End())
}

type Build struct {
	Type    BuildType
	Version *Version
}

var _ xml.Unmarshaler = &Build{}

func (b *Build) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "Version":
				var version Version
				if err := version.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Version: %w", err)
				}
				b.Version = &version
			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "Type":
				b.Type = BuildType(elem)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &Build{}

func (b *Build) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc, kxml.StartElement("Type"), xml.CharData(b.Type)); err != nil {
		return fmt.Errorf("type: %w", err)
	}

	if b.Version != nil {
		if err := b.Version.MarshalXML(enc, kxml.StartElement("Version")); err != nil {
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

var _ xml.Unmarshaler = &Device{}

func (d *Device) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
	var targetCharData string
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "Version":
				var version Version
				if err := version.UnmarshalXML(dec, elem); err != nil {
					return fmt.Errorf("unmarshal Version: %w", err)
				}
				d.Version = &version

			default:
				targetCharData = elem.Name.Local
			}
		case xml.CharData:
			switch targetCharData {
			case "Name":
				d.Name = string(elem)
			case "UnitId":
				u, err := strconv.ParseUint(string(elem), 10, 32)
				if err != nil {
					return fmt.Errorf("parse UnitId: %w", err)
				}
				d.UnitId = uint32(u)
			case "ProductId":
				u, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return fmt.Errorf("parse ProductId: %w", err)
				}
				d.ProductID = uint16(u)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &Device{}

func (d *Device) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc,
		kxml.StartElement("Name"),
		xml.CharData(d.Name)); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := kxml.EncodeElement(enc,
		kxml.StartElement("UnitId"),
		xml.CharData(strconv.FormatUint(uint64(d.UnitId), 10))); err != nil {
		return fmt.Errorf("unitId: %w", err)
	}
	if err := kxml.EncodeElement(enc,
		kxml.StartElement("ProductID"),
		xml.CharData(strconv.FormatUint(uint64(d.ProductID), 10))); err != nil {
		return fmt.Errorf("productID: %w", err)
	}

	if d.Version != nil {
		if err := d.Version.MarshalXML(enc, kxml.StartElement("Version")); err != nil {
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

var _ xml.Unmarshaler = &Version{}

func (v *Version) UnmarshalXML(dec *xml.Decoder, se xml.StartElement) error {
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
			case "VersionMajor":
				u, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return fmt.Errorf("parse VersionMajor: %w", err)
				}
				v.VersionMajor = uint16(u)
			case "VersionMinor":
				u, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return fmt.Errorf("parse VersionMinor: %w", err)
				}
				v.VersionMinor = uint16(u)
			case "BuildMajor":
				u, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return fmt.Errorf("parse BuildMajor: %w", err)
				}
				v.BuildMajor = uint16(u)
			case "BuildMinor":
				u, err := strconv.ParseUint(string(elem), 10, 16)
				if err != nil {
					return fmt.Errorf("parse BuildMinor: %w", err)
				}
				v.BuildMinor = uint16(u)
			}
			targetCharData = ""
		case xml.EndElement:
			if elem == se.End() {
				return nil
			}
		}
	}
}

var _ xml.Marshaler = &Version{}

func (v *Version) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := kxml.EncodeElement(enc,
		kxml.StartElement("VersionMajor"),
		xml.CharData(strconv.FormatUint(uint64(v.VersionMajor), 10))); err != nil {
		return fmt.Errorf("versionMajor:%w", err)
	}
	if err := kxml.EncodeElement(enc,
		kxml.StartElement("VersionMinor"),
		xml.CharData(strconv.FormatUint(uint64(v.VersionMinor), 10))); err != nil {
		return fmt.Errorf("versionMinor:%w", err)
	}
	if err := kxml.EncodeElement(enc,
		kxml.StartElement("BuildMajor"),
		xml.CharData(strconv.FormatUint(uint64(v.BuildMajor), 10))); err != nil {
		return fmt.Errorf("buildMajor:%w", err)
	}
	if err := kxml.EncodeElement(enc,
		kxml.StartElement("BuildMinor"),
		xml.CharData(strconv.FormatUint(uint64(v.BuildMinor), 10))); err != nil {
		return fmt.Errorf("buildMinor:%w", err)
	}

	return enc.EncodeToken(se.End())
}
