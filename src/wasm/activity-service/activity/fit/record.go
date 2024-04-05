package fit

import (
	"github.com/muktihari/fit/factory"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/kit/scaleoffset"
	"github.com/muktihari/fit/kit/semicircles"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
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
			lat := field.Value.Int32()
			if lat == basetype.Sint32Invalid {
				continue
			}
			rec.PositionLat = kit.Ptr(semicircles.ToDegrees(lat))
		case fieldnum.RecordPositionLong:
			long := field.Value.Int32()
			if long == basetype.Sint32Invalid {
				continue
			}
			rec.PositionLong = kit.Ptr(semicircles.ToDegrees(long))
		case fieldnum.RecordAltitude:
			altitude := field.Value.Uint16()
			if altitude == basetype.Uint16Invalid {
				continue
			}
			faltitude := (float64(altitude) / 5) - 500
			if faltitude < 0 {
				continue
			}
			rec.Altitude = &faltitude
		case fieldnum.RecordHeartRate:
			heartRate := field.Value.Uint8()
			if heartRate == basetype.Uint8Invalid {
				continue
			}
			rec.HeartRate = &heartRate
		case fieldnum.RecordCadence:
			cadence := field.Value.Uint8()
			if cadence == basetype.Uint8Invalid {
				continue
			}
			rec.Cadence = &cadence
		case fieldnum.RecordDistance:
			distance := field.Value.Uint32()
			if distance == basetype.Uint32Invalid {
				continue
			}

			rec.Distance = kit.Ptr(float64(distance) / 100)
		case fieldnum.RecordSpeed:
			speed := field.Value.Uint16()
			if speed == basetype.Uint16Invalid {
				continue
			}
			rec.Speed = kit.Ptr(float64(speed) / 1000)
		case fieldnum.RecordPower:
			power := field.Value.Uint16()
			if power == basetype.Uint16Invalid {
				continue
			}
			rec.Power = &power
		case fieldnum.RecordTemperature:
			temperature := field.Value.Int8()
			if temperature == basetype.Sint8Invalid {
				continue
			}
			rec.Temperature = &temperature
		}
	}

	return rec
}

func convertRecordToMesg(rec *activity.Record) proto.Message {
	mesg := factory.CreateMesgOnly(mesgnum.Record)

	if !rec.Timestamp.IsZero() {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordTimestamp)
		field.Value = proto.Uint32(datetime.ToUint32(rec.Timestamp))
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.PositionLat != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordPositionLat)
		field.Value = proto.Int32(semicircles.ToSemicircles(*rec.PositionLat))
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.PositionLong != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordPositionLong)
		field.Value = proto.Int32(semicircles.ToSemicircles(*rec.PositionLong))
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Distance != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordDistance)
		field.Value = scaleoffset.DiscardValue(*rec.Distance, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Altitude != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordAltitude)
		field.Value = scaleoffset.DiscardValue(*rec.Altitude, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.HeartRate != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordHeartRate)
		field.Value = proto.Uint8(*rec.HeartRate)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Cadence != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordCadence)
		field.Value = proto.Uint8(*rec.Cadence)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Speed != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordSpeed)
		field.Value = scaleoffset.DiscardValue(*rec.Speed, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Power != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordPower)
		field.Value = proto.Uint16(*rec.Power)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Temperature != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordTemperature)
		field.Value = proto.Int8(*rec.Temperature)
		mesg.Fields = append(mesg.Fields, field)
	}

	return mesg
}
