package fit

import (
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
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
