package fit

import (
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
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
			t, ok := mesg.Fields[i].Value.(uint32)
			if !ok {
				continue
			}
			timestamp = t
		case fieldnum.ActivityLocalTimestamp:
			t, ok := mesg.Fields[i].Value.(uint32)
			if !ok {
				continue
			}
			localDateTime = t
		}
	}

	if timestamp == basetype.Uint32Invalid || localDateTime == basetype.Uint32Invalid {
		return 0 // Default UTC
	}

	return datetime.TzOffsetHours(
		typedef.LocalDateTime(localDateTime),
		typedef.DateTime(timestamp),
	)
}
