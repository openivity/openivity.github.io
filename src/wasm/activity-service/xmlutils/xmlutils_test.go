// Copyright (C) 2024 Openivity

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

package xmlutils_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openivity/activity-service/xmlutils"
)

func TestStartElement(t *testing.T) {
	name := "name"
	se := xmlutils.StartElement(name)
	if diff := cmp.Diff(se, xml.StartElement{Name: xml.Name{Local: name}}); diff != "" {
		t.Fatal(diff)
	}
}

type fnWriter func(b []byte) (n int, err error)

func (f fnWriter) Write(b []byte) (n int, err error) { return f(b) }

func TestEncodeElement(t *testing.T) {
	tt := []struct {
		name     string
		w        io.Writer
		se       xml.StartElement
		charData xml.CharData
		expected string
		err      error
	}{
		{
			name:     "happy flow",
			w:        new(bytes.Buffer),
			se:       xml.StartElement{Name: xml.Name{Local: "speed"}},
			charData: xml.CharData("100"),
			expected: "<speed>100</speed>",
			err:      nil,
		},
		{
			name:     "happy error encode start element",
			w:        new(bytes.Buffer),
			se:       xml.StartElement{Name: xml.Name{}},
			charData: xml.CharData("100"),
			err:      fmt.Errorf("xml: start tag with no name"),
		},
	}

	for _, tc := range tt {
		enc := xml.NewEncoder(tc.w)
		err := xmlutils.EncodeElement(enc, tc.se, tc.charData)
		_ = enc.Flush()
		if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tc.err) {
			t.Fatalf("expected err: %v, got: %v", tc.err, err)
		}
		if err != nil {
			return
		}

		buf, ok := tc.w.(*bytes.Buffer)
		if !ok {
			return
		}
		if diff := cmp.Diff(buf.String(), tc.expected); diff != "" {
			t.Fatal(diff)
		}
	}
}

type fnMarshalXML func(enc *xml.Encoder, se xml.StartElement) error

func (m fnMarshalXML) MarshalXML(enc *xml.Encoder, se xml.StartElement) error {
	return m(enc, se)
}

func TestMarshalWrite(t *testing.T) {
	tt := []struct {
		name      string
		w         io.Writer
		marshaler xml.Marshaler
		expected  string
		err       error
	}{
		{
			name: "happy flow",
			w:    new(bytes.Buffer),
			marshaler: fnMarshalXML(func(enc *xml.Encoder, se xml.StartElement) error {
				return nil
			}),
		},
		{
			name: "marshalXML return error",
			w:    new(bytes.Buffer),
			marshaler: fnMarshalXML(func(enc *xml.Encoder, se xml.StartElement) error {
				return io.ErrClosedPipe
			}),
			err: io.ErrClosedPipe,
		},
		{
			name: "close return error",
			w: fnWriter(func(b []byte) (n int, err error) {
				return 0, io.ErrClosedPipe
			}),
			marshaler: fnMarshalXML(func(enc *xml.Encoder, se xml.StartElement) error {
				enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "speed"}}) // only to fill encoder so encoder.n != 0
				return nil
			}),
			err: io.ErrClosedPipe,
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			err := xmlutils.MarshalWrite(tc.w, tc.marshaler)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected err: %v, got: %v", tc.err, err)
			}
		})
	}
}
