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

func convertRecordToMesg(rec *activity.Record) proto.Message {
	mesg := factory.CreateMesgOnly(mesgnum.Record)

	if !rec.Timestamp.IsZero() {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordTimestamp)
		field.Value = datetime.ToUint32(rec.Timestamp)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.PositionLat != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordPositionLat)
		field.Value = semicircles.ToSemicircles(*rec.PositionLat)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.PositionLong != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordPositionLong)
		field.Value = semicircles.ToSemicircles(*rec.PositionLong)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Distance != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordDistance)
		field.Value = scaleoffset.DiscardAny(*rec.Distance, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Altitude != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordAltitude)
		field.Value = scaleoffset.DiscardAny(*rec.Altitude, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.HeartRate != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordHeartRate)
		field.Value = *rec.HeartRate
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Cadence != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordCadence)
		field.Value = *rec.Cadence
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Speed != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordSpeed)
		field.Value = scaleoffset.DiscardAny(*rec.Speed, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Power != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordPower)
		field.Value = *rec.Power
		mesg.Fields = append(mesg.Fields, field)
	}
	if rec.Temperature != nil {
		field := factory.CreateField(mesgnum.Record, fieldnum.RecordTemperature)
		field.Value = *rec.Temperature
		mesg.Fields = append(mesg.Fields, field)
	}

	return mesg
}
