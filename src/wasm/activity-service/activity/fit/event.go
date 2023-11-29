package fit

import (
	"github.com/muktihari/fit/factory"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
)

func convertEventToMesg(event *activity.Event) proto.Message {
	return factory.CreateMesgOnly(mesgnum.Event).WithFields(
		factory.CreateField(mesgnum.Event, fieldnum.EventTimestamp).WithValue(datetime.ToUint32(event.Timestamp)),
		factory.CreateField(mesgnum.Event, fieldnum.EventEvent).WithValue(event.Event),
		factory.CreateField(mesgnum.Event, fieldnum.EventEventType).WithValue(event.EventType),
	)
}
