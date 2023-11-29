package fit

import (
	"github.com/muktihari/fit/factory"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/kit/scaleoffset"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/kit"
)

func NewSession(mesg proto.Message) *activity.Session {
	ses := &activity.Session{}

	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.SessionTimestamp:
			ses.Timestamp = datetime.ToTime(field.Value)
		case fieldnum.SessionStartTime:
			ses.StartTime = datetime.ToTime(field.Value)
		case fieldnum.SessionSport:
			sport, ok := field.Value.(uint8)
			if !ok || sport == basetype.EnumInvalid {
				continue
			}
			ses.Sport = kit.FormatTitle(typedef.Sport(sport).String())
		case fieldnum.SessionTotalMovingTime:
			totalMovingTime, ok := field.Value.(uint32)
			if !ok || totalMovingTime == basetype.Uint32Invalid {
				continue
			}
			ses.TotalMovingTime = float64(totalMovingTime) / 1000
		case fieldnum.SessionTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok || totalElapsedTime == basetype.Uint32Invalid {
				continue
			}
			ses.TotalElapsedTime = float64(totalElapsedTime) / 1000
		case fieldnum.SessionTotalDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok || totalDistance == basetype.Uint32Invalid {
				continue
			}
			ses.TotalDistance = float64(totalDistance) / 100
		case fieldnum.SessionTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok || totalAscent == basetype.Uint16Invalid {
				continue
			}
			ses.TotalAscent = totalAscent
		case fieldnum.SessionTotalDescent:
			totalDescent, ok := field.Value.(uint16)
			if !ok || totalDescent == basetype.Uint16Invalid {
				continue
			}
			ses.TotalDescent = totalDescent
		case fieldnum.SessionTotalCalories:
			totalCalories, ok := field.Value.(uint16)
			if !ok || totalCalories == basetype.Uint16Invalid {
				continue
			}
			ses.TotalCalories = totalCalories
		case fieldnum.SessionAvgSpeed:
			avgSpeed, ok := field.Value.(uint16)
			if !ok || avgSpeed == basetype.Uint16Invalid {
				continue
			}
			ses.AvgSpeed = kit.Ptr(float64(avgSpeed) / 1000)
		case fieldnum.SessionMaxSpeed:
			maxSpeed, ok := field.Value.(uint16)
			if !ok || maxSpeed == basetype.Uint16Invalid {
				continue
			}
			ses.MaxSpeed = kit.Ptr(float64(maxSpeed) / 1000)
		case fieldnum.SessionAvgHeartRate:
			avgHeartRate, ok := field.Value.(uint8)
			if !ok || avgHeartRate == basetype.Uint8Invalid {
				continue
			}
			ses.AvgHeartRate = &avgHeartRate
		case fieldnum.SessionMaxHeartRate:
			maxHeartRate, ok := field.Value.(uint8)
			if !ok || maxHeartRate == basetype.Uint8Invalid {
				continue
			}
			ses.MaxHeartRate = &maxHeartRate
		case fieldnum.SessionAvgCadence:
			avgCadence, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			if avgCadence == basetype.Uint8Invalid {
				continue
			}
			ses.AvgCadence = &avgCadence
		case fieldnum.SessionMaxCadence:
			maxCadence, ok := field.Value.(uint8)
			if !ok || maxCadence == basetype.Uint8Invalid {
				continue
			}
			ses.MaxCadence = &maxCadence
		case fieldnum.SessionAvgPower:
			avgPower, ok := field.Value.(uint16)
			if !ok || avgPower == basetype.Uint16Invalid {
				continue
			}
			ses.AvgPower = &avgPower
		case fieldnum.SessionMaxPower:
			maxPower, ok := field.Value.(uint16)
			if !ok || maxPower == basetype.Uint16Invalid {
				continue
			}
			ses.MaxPower = &maxPower
		case fieldnum.SessionAvgTemperature:
			avgTemperature, ok := field.Value.(int8)
			if !ok || avgTemperature == basetype.Sint8Invalid {
				continue
			}
			ses.AvgTemperature = &avgTemperature
		case fieldnum.SessionMaxTemperature:
			maxTemperature, ok := field.Value.(int8)
			if !ok || maxTemperature == basetype.Sint8Invalid {
				continue
			}
			ses.MaxTemperature = &maxTemperature
		case fieldnum.SessionAvgAltitude:
			avgAltitude, ok := field.Value.(uint16)
			if !ok || avgAltitude == basetype.Uint16Invalid {
				continue
			}
			fAvgAltitude := (float64(avgAltitude) / 5) - 500
			if fAvgAltitude < 0 {
				continue
			}
			ses.AvgAltitude = &fAvgAltitude
		case fieldnum.SessionMaxAltitude:
			maxAltitude, ok := field.Value.(uint16)
			if !ok || maxAltitude == basetype.Uint16Invalid {
				continue
			}
			fMaxAltitude := (float64(maxAltitude) / 5) - 500
			if fMaxAltitude < 0 {
				continue
			}
			ses.MaxAltitude = &fMaxAltitude
		}
	}

	if ses.Sport == kit.FormatTitle(typedef.SportAll.String()) {
		ses.Sport = activity.SportGeneric
	}

	return ses
}

func convertSessionToMesg(ses *activity.Session) proto.Message {
	mesg := factory.CreateMesgOnly(mesgnum.Session)

	if !ses.Timestamp.IsZero() {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTimestamp)
		field.Value = datetime.ToUint32(ses.Timestamp)
		mesg.Fields = append(mesg.Fields, field)
	}
	if !ses.StartTime.IsZero() {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionStartTime)
		field.Value = datetime.ToUint32(ses.StartTime)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.Sport != "" {
		sport := typedef.SportFromString(kit.FormatToLowerSnakeCase(ses.Sport))
		if sport == typedef.SportInvalid {
			sport = typedef.SportGeneric
		}
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionSport)
		field.Value = uint8(sport)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalMovingTime != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalMovingTime)
		field.Value = scaleoffset.DiscardAny(ses.TotalMovingTime, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalElapsedTime != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalElapsedTime)
		field.Value = scaleoffset.DiscardAny(ses.TotalElapsedTime, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)

		totalTimerTimeField := factory.CreateField(mesgnum.Lap, fieldnum.SessionTotalElapsedTime)
		totalTimerTimeField.Value = field.Value
		mesg.Fields = append(mesg.Fields, totalTimerTimeField)
	}
	if ses.TotalDistance != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalDistance)
		field.Value = scaleoffset.DiscardAny(ses.AvgAltitude, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalAscent != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalAscent)
		field.Value = ses.TotalAscent
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalDescent != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalDescent)
		field.Value = ses.TotalDescent
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalCalories != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalCalories)
		field.Value = ses.TotalCalories
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgSpeed != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgSpeed)
		field.Value = scaleoffset.DiscardAny(*ses.AvgSpeed, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxSpeed != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxSpeed)
		field.Value = scaleoffset.DiscardAny(*ses.MaxSpeed, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgHeartRate != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgHeartRate)
		field.Value = *ses.AvgHeartRate
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxHeartRate != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxHeartRate)
		field.Value = *ses.MaxHeartRate
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgCadence != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgCadence)
		field.Value = *ses.AvgCadence
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxCadence != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxCadence)
		field.Value = *ses.MaxCadence
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgPower != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgPower)
		field.Value = *ses.AvgPower
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxPower != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxPower)
		field.Value = *ses.MaxPower
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgTemperature != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgTemperature)
		field.Value = *ses.AvgTemperature
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxTemperature != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxTemperature)
		field.Value = *ses.MaxTemperature
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgAltitude != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgAltitude)
		field.Value = scaleoffset.DiscardAny(ses.AvgAltitude, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxAltitude != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxAltitude)
		field.Value = scaleoffset.DiscardAny(ses.MaxAltitude, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}

	eventField := factory.CreateField(mesgnum.Session, fieldnum.SessionEvent)
	eventField.Value = uint8(typedef.EventSession)
	mesg.Fields = append(mesg.Fields, eventField)

	eventTypeField := factory.CreateField(mesgnum.Session, fieldnum.SessionEventType)
	eventTypeField.Value = uint8(typedef.EventTypeStop)
	mesg.Fields = append(mesg.Fields, eventTypeField)

	return mesg
}
