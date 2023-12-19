package result

import (
	"bytes"
	"encoding/json"

	"github.com/muktihari/openactivity-fit/activity/fit"
)

type ManufacturerList struct {
	Manufacturers []fit.Manufacturer
}

var _ json.Marshaler = &ManufacturerList{}

func (m *ManufacturerList) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')

	buf.WriteString("\"manufacturers\":[")
	for i := range m.Manufacturers {
		b, _ := m.Manufacturers[i].MarshalJSON()
		buf.Write(b)
		if i != len(m.Manufacturers)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	buf.WriteByte('}')
	return buf.Bytes(), nil
}
