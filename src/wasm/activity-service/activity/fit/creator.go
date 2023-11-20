package fit

import (
	"strconv"

	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/kit/typeconv"
	"github.com/muktihari/fit/profile/typedef"
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
			m.Name = activity.FormatTitle(typeconv.ToUint16[typedef.Manufacturer](field.Value).String())
		case fieldnum.FileIdProduct:
			product, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m.Product = &product
			m.Name += " (" + strconv.FormatUint(uint64(product), 10) + ")"
		case fieldnum.FileIdTimeCreated:
			m.TimeCreated = datetime.ToTime(field.Value)
		}
	}

	return m
}
