package fit

import (
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
			manufacturer, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m.Manufacturer = &manufacturer
		case fieldnum.FileIdProduct:
			product, ok := field.Value.(uint16)
			if !ok {
				continue
			}
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
	fileTypeField.Value = uint8(typedef.FileActivity)
	mesg.Fields = append(mesg.Fields, fileTypeField)

	manufacturerField := factory.CreateField(mesgnum.FileId, fieldnum.FileIdManufacturer)
	manufacturerField.Value = uint16(*cre.Manufacturer)
	mesg.Fields = append(mesg.Fields, manufacturerField)

	productField := factory.CreateField(mesgnum.FileId, fieldnum.FileIdProduct)
	productField.Value = uint16(*cre.Product)
	mesg.Fields = append(mesg.Fields, productField)

	timeCreatedField := factory.CreateField(mesgnum.FileId, fieldnum.FileIdTimeCreated)
	timeCreatedField.Value = datetime.ToUint32(cre.TimeCreated)
	mesg.Fields = append(mesg.Fields, timeCreatedField)

	return mesg
}

type Manufacturer struct {
	ID       uint16
	Name     string
	Products []ManufacturerProduct
}

func (m *Manufacturer) ToMap() map[string]any {
	products := make([]any, len(m.Products))
	for i := range m.Products {
		products[i] = m.Products[i].ToMap()
	}
	return map[string]any{
		"id":       uint16(m.ID),
		"name":     m.Name,
		"products": products,
	}
}

type ManufacturerProduct struct {
	ID   uint16
	Name string
}

func (p *ManufacturerProduct) ToMap() map[string]any {
	return map[string]any{
		"id":   p.ID,
		"name": p.Name,
	}
}
