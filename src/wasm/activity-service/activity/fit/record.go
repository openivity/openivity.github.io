package fit

import (
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/kit/semicircles"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/kit"
)

func NewRecord(mesg proto.Message) *activity.Record {
	rec := &activity.Record{}

	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.RecordTimestamp:
			rec.Timestamp = datetime.ToTime(field.Value)
		case fieldnum.RecordPositionLat:
			lat, ok := field.Value.(int32)
			if !ok || lat == basetype.Sint32Invalid {
				continue
			}
			rec.PositionLat = kit.Ptr(semicircles.ToDegrees(lat))
		case fieldnum.RecordPositionLong:
			long, ok := field.Value.(int32)
			if !ok || long == basetype.Sint32Invalid {
				continue
			}
			rec.PositionLong = kit.Ptr(semicircles.ToDegrees(long))
		case fieldnum.RecordAltitude:
			altitude, ok := field.Value.(uint16)
			if !ok || altitude == basetype.Uint16Invalid {
				continue
			}
			faltitude := (float64(altitude) / 5) - 500
			if faltitude < 0 {
				continue
			}
			rec.Altitude = &faltitude
		case fieldnum.RecordHeartRate:
			heartRate, ok := field.Value.(uint8)
			if !ok || heartRate == basetype.Uint8Invalid {
				continue
			}
			rec.HeartRate = &heartRate
		case fieldnum.RecordCadence:
			cadence, ok := field.Value.(uint8)
			if !ok || cadence == basetype.Uint8Invalid {
				continue
			}
			rec.Cadence = &cadence
		case fieldnum.RecordDistance:
			distance, ok := field.Value.(uint32)
			if !ok || distance == basetype.Uint32Invalid {
				continue
			}

			rec.Distance = kit.Ptr(float64(distance) / 100)
		case fieldnum.RecordSpeed:
			speed, ok := field.Value.(uint16)
			if !ok || speed == basetype.Uint16Invalid {
				continue
			}
			rec.Speed = kit.Ptr(float64(speed) / 1000)
		case fieldnum.RecordPower:
			power, ok := field.Value.(uint16)
			if !ok || power == basetype.Uint16Invalid {
				continue
			}
			rec.Power = &power
		case fieldnum.RecordTemperature:
			temperature, ok := field.Value.(int8)
			if !ok || temperature == basetype.Sint8Invalid {
				continue
			}
			rec.Temperature = &temperature
		}
	}

	return rec
}
