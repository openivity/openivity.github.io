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
	"time"

	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/xmlutils"
)

type ActivityList struct {
	Activity Activity `xml:"Activity"`
}

func (a *ActivityList) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
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
		case "Activity":
			var activity Activity
			se := xmltokenizer.GetToken().Copy(token)
			err = activity.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Activity: %w", err)
			}
			a.Activity = activity
		}
	}
	return nil
}

var _ xml.Marshaler = (*ActivityList)(nil)

func (a *ActivityList) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := a.Activity.MarshalXML(enc, xmlutils.StartElement("Activity")); err != nil {
		return fmt.Errorf("activity: %w", err)
	}

	return enc.EncodeToken(se.End())
}

type Activity struct {
	Sport   string        `xml:"Sport,attr"`
	ID      time.Time     `xml:"Id"`
	Laps    []ActivityLap `xml:"Lap"`
	Notes   string        `xml:"Notes,omitempty"`
	Creator *Device       `xml:"Creator,omitempty"`
}

func (a *Activity) UnmarshalToken(tok *xmltokenizer.Tokenizer, se *xmltokenizer.Token) error {
	for i := range se.Attrs {
		attr := &se.Attrs[i]
		switch string(attr.Name.Local) {
		case "Sport":
			a.Sport = string(attr.Value)
		}
	}

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
		case "Lap":
			var activityLap ActivityLap
			se := xmltokenizer.GetToken().Copy(token)
			err = activityLap.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Lap: %w", err)
			}
			a.Laps = append(a.Laps, activityLap)
		case "Creator":
			var device Device
			se := xmltokenizer.GetToken().Copy(token)
			err = device.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return fmt.Errorf("unmarshal Creator: %w", err)
			}
			a.Creator = &device
		case "Id":
			t, err := time.Parse(time.RFC3339, string(token.Data))
			if err != nil {
				return fmt.Errorf("parse Id: %w", err)
			}
			a.ID = t
		case "Notes":
			a.Notes = string(token.Data)
		}
	}

	return nil
}

var _ xml.Marshaler = (*Activity)(nil)

func (a *Activity) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	se.Attr = append(se.Attr, xml.Attr{
		Name:  xml.Name{Local: "Sport"},
		Value: a.Sport,
	})

	if err := enc.EncodeToken(se); err != nil {
		return err
	}

	if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("Id"), xml.CharData(a.ID.Format(time.RFC3339))); err != nil {
		return fmt.Errorf("id: %w", err)
	}

	for i := range a.Laps {
		if err := a.Laps[i].MarshalXML(enc, xmlutils.StartElement("Lap")); err != nil {
			return fmt.Errorf("lap[%d]: %w", i, err)
		}
	}

	if err := xmlutils.EncodeElement(enc, xmlutils.StartElement("Notes"), xml.CharData(a.Notes)); err != nil {
		return fmt.Errorf("notes: %w", err)
	}

	if a.Creator != nil {
		if err := a.Creator.MarshalXML(enc, xmlutils.StartElement("Creator")); err != nil {
			return fmt.Errorf("creator: %w", err)
		}
	}

	return enc.EncodeToken(se.End())
}
