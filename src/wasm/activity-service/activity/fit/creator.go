package fit

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/muktihari/fit/factory"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
)

func NewCreator(mesg proto.Message) activity.Creator {
	m := activity.Creator{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.FileIdManufacturer:
			manufacturer := field.Value.Uint16()
			m.Manufacturer = &manufacturer
		case fieldnum.FileIdProduct:
			product := field.Value.Uint16()
			m.Product = &product
		case fieldnum.FileIdTimeCreated:
			m.TimeCreated = datetime.ToTime(field.Value)
		}
	}

	return m
}

func convertCreatorToMesg(cre *activity.Creator) proto.Message {
	mesg := factory.CreateMesgOnly(mesgnum.FileId)

	fileTypeField := factory.CreateField(mesgnum.FileId, fieldnum.FileIdType)
	fileTypeField.Value = proto.Uint8(uint8(typedef.FileActivity))
	mesg.Fields = append(mesg.Fields, fileTypeField)

	manufacturerField := factory.CreateField(mesgnum.FileId, fieldnum.FileIdManufacturer)
	manufacturerField.Value = proto.Uint16(*cre.Manufacturer)
	mesg.Fields = append(mesg.Fields, manufacturerField)

	productField := factory.CreateField(mesgnum.FileId, fieldnum.FileIdProduct)
	productField.Value = proto.Uint16(*cre.Product)
	mesg.Fields = append(mesg.Fields, productField)

	timeCreatedField := factory.CreateField(mesgnum.FileId, fieldnum.FileIdTimeCreated)
	timeCreatedField.Value = proto.Uint32(datetime.ToUint32(cre.TimeCreated))
	mesg.Fields = append(mesg.Fields, timeCreatedField)

	return mesg
}

type Manufacturer struct {
	ID       uint16
	Name     string
	Products []ManufacturerProduct
}

var _ json.Marshaler = &Manufacturer{}

func (m *Manufacturer) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteByte('{')
	buf.WriteString("\"id\":" + strconv.FormatUint(uint64(m.ID), 10) + ",")
	buf.WriteString("\"name\":\"" + m.Name + "\",")

	buf.WriteString("\"products\":[")
	for i := range m.Products {
		b, _ := m.Products[i].MarshalJSON()
		buf.Write(b)
		if i != len(m.Products)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

type ManufacturerProduct struct {
	ID   uint16
	Name string
}

var _ json.Marshaler = &ManufacturerProduct{}

func (p *ManufacturerProduct) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')
	buf.WriteString("\"id\":" + strconv.FormatUint(uint64(p.ID), 10) + ",")
	buf.WriteString("\"name\":\"" + p.Name + "\"")
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
