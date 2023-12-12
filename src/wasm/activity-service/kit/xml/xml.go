package xml

import (
	"bytes"
	"encoding/xml"
	"strings"
)

func StartElement(name string) xml.StartElement {
	return xml.StartElement{Name: xml.Name{Local: name}}
}

func EncodeElement(enc *xml.Encoder, se xml.StartElement, charData xml.CharData) error {
	if err := enc.EncodeToken(se); err != nil {
		return err
	}
	if charData != nil {
		if err := enc.EncodeToken(charData); err != nil {
			return err
		}
	}
	return enc.EncodeToken(se.End())
}

func Marshal(marshaler xml.Marshaler) (b []byte, err error) {
	w := bytes.NewBuffer(nil)
	w.Write([]byte(xml.Header))

	enc := xml.NewEncoder(w)
	defer func() {
		err = enc.Close()
		b = w.Bytes()
	}()

	enc.Indent("", strings.Repeat(" ", 2))

	if err = marshaler.MarshalXML(enc, xml.StartElement{}); err != nil {
		return nil, err
	}

	return
}
