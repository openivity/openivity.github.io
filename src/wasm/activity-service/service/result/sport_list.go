package result

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type SportList struct {
	Sports []Sport
}

var _ json.Marshaler = &SportList{}

func (s *SportList) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')

	buf.WriteString("\"sports\":[")
	for i := range s.Sports {
		b, _ := s.Sports[i].MarshalJSON()
		buf.Write(b)
		if i != len(s.Sports)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

type Sport struct {
	ID                   uint8
	Name                 string
	ToleranceMovingSpeed float64
}

var _ json.Marshaler = &Sport{}

func (s *Sport) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')
	buf.WriteString("\"id\":" + strconv.FormatUint(uint64(s.ID), 10) + ",")
	buf.WriteString("\"name\":\"" + s.Name + "\",")
	buf.WriteString("\"toleranceMovingSpeed\":" + strconv.FormatFloat(s.ToleranceMovingSpeed, 'g', -1, 64))
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
