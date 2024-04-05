package fit

import (
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/proto"
)

func CreateTimezone(mesg proto.Message) int {
	var (
		timestamp     = basetype.Uint32Invalid
		localDateTime = basetype.Uint32Invalid
	)

	for i := range mesg.Fields {
		switch mesg.Fields[i].Num {
		case fieldnum.ActivityTimestamp:
			timestamp = mesg.Fields[i].Value.Uint32()
		case fieldnum.ActivityLocalTimestamp:
			localDateTime = mesg.Fields[i].Value.Uint32()
		}
	}

	if timestamp == basetype.Uint32Invalid || localDateTime == basetype.Uint32Invalid {
		return 0 // Default UTC
	}

	return datetime.TzOffsetHoursFromUint32(localDateTime, timestamp)
}
