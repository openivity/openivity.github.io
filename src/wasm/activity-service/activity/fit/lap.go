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

func NewLap(mesg proto.Message) *activity.Lap {
	lap := &activity.Lap{}

	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.LapTimestamp:
			lap.Timestamp = datetime.ToTime(field.Value)
		case fieldnum.LapStartTime:
			lap.StartTime = datetime.ToTime(field.Value)
		case fieldnum.LapTotalMovingTime:
			totalMovingTime, ok := field.Value.(uint32)
			if !ok || totalMovingTime == basetype.Uint32Invalid {
				continue
			}
			lap.TotalMovingTime = float64(totalMovingTime) / 1000
		case fieldnum.LapTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok || totalElapsedTime == basetype.Uint32Invalid {
				continue
			}
			lap.TotalElapsedTime = float64(totalElapsedTime) / 1000
		case fieldnum.LapTotalDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok || totalDistance == basetype.Uint32Invalid {
				continue
			}
			lap.TotalDistance = float64(totalDistance) / 100
		case fieldnum.LapTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok || totalAscent == basetype.Uint16Invalid {
				continue
			}
			lap.TotalAscent = totalAscent
		case fieldnum.LapTotalDescent:
			totalDescent, ok := field.Value.(uint16)
			if !ok || totalDescent == basetype.Uint16Invalid {
				continue
			}
			lap.TotalDescent = totalDescent
		case fieldnum.LapTotalCalories:
			totalCalories, ok := field.Value.(uint16)
			if !ok || totalCalories == basetype.Uint16Invalid {
				continue
			}
			lap.TotalCalories = totalCalories
		case fieldnum.LapAvgSpeed:
			avgSpeed, ok := field.Value.(uint16)
			if !ok || avgSpeed == basetype.Uint16Invalid {
				continue
			}
			fAvgSpeed := float64(avgSpeed) / 1000
			lap.AvgSpeed = &fAvgSpeed
		case fieldnum.LapMaxSpeed:
			maxSpeed, ok := field.Value.(uint16)
			if !ok || maxSpeed == basetype.Uint16Invalid {
				continue
			}
			fMaxSpeed := float64(maxSpeed) / 1000
			lap.MaxSpeed = &fMaxSpeed
		case fieldnum.LapAvgCadence:
			avgCadence, ok := field.Value.(uint8)
			if !ok || avgCadence == basetype.Uint8Invalid {
				continue
			}
			lap.AvgCadence = &avgCadence
		case fieldnum.LapMaxCadence:
			maxCadence, ok := field.Value.(uint8)
			if !ok || maxCadence == basetype.Uint8Invalid {
				continue
			}
			lap.MaxCadence = &maxCadence
		case fieldnum.LapAvgHeartRate:
			avgHeartRate, ok := field.Value.(uint8)
			if !ok || avgHeartRate == basetype.Uint8Invalid {
				continue
			}
			lap.AvgHeartRate = &avgHeartRate
		case fieldnum.LapMaxHeartRate:
			maxHeartRate, ok := field.Value.(uint8)
			if !ok || maxHeartRate == basetype.Uint8Invalid {
				continue
			}
			lap.MaxHeartRate = &maxHeartRate
		case fieldnum.LapAvgPower:
			avgPower, ok := field.Value.(uint16)
			if !ok || avgPower == basetype.Uint16Invalid {
				continue
			}
			lap.AvgPower = &avgPower
		case fieldnum.LapMaxPower:
			maxPower, ok := field.Value.(uint16)
			if !ok || maxPower == basetype.Uint16Invalid {
				continue
			}
			lap.MaxPower = &maxPower
		}
	}

	return lap
}

func convertLapToMesg(lap *activity.Lap) proto.Message {
	mesg := factory.CreateMesgOnly(mesgnum.Lap)

	if !lap.Timestamp.IsZero() {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapTimestamp)
		field.Value = datetime.ToUint32(lap.Timestamp)
		mesg.Fields = append(mesg.Fields, field)
	}
	if !lap.StartTime.IsZero() {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapStartTime)
		field.Value = datetime.ToUint32(lap.StartTime)
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.Sport != "" {
		sport := typedef.SportFromString(kit.FormatToLowerSnakeCase(lap.Sport))
		if sport == typedef.SportInvalid {
			sport = typedef.SportGeneric
		}
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapSport)
		field.Value = uint8(sport)
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.TotalMovingTime != 0 {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapTotalMovingTime)
		field.Value = scaleoffset.DiscardAny(lap.TotalMovingTime, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.TotalElapsedTime != 0 {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapTotalElapsedTime)
		field.Value = scaleoffset.DiscardAny(lap.TotalElapsedTime, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)

		totalTimerTimeField := factory.CreateField(mesgnum.Lap, fieldnum.LapTotalTimerTime)
		totalTimerTimeField.Value = field.Value
		mesg.Fields = append(mesg.Fields, totalTimerTimeField)
	}
	if lap.TotalDistance != 0 {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapTotalDistance)
		field.Value = scaleoffset.DiscardAny(lap.TotalDistance, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.TotalAscent != 0 {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapTotalAscent)
		field.Value = lap.TotalAscent
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.TotalDescent != 0 {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapTotalDescent)
		field.Value = lap.TotalDescent
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.TotalCalories != 0 {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapTotalCalories)
		field.Value = lap.TotalCalories
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.AvgSpeed != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapAvgSpeed)
		field.Value = scaleoffset.DiscardAny(*lap.AvgSpeed, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.MaxSpeed != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapMaxSpeed)
		field.Value = scaleoffset.DiscardAny(*lap.MaxSpeed, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.AvgHeartRate != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapAvgHeartRate)
		field.Value = *lap.AvgHeartRate
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.MaxHeartRate != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapMaxHeartRate)
		field.Value = *lap.MaxHeartRate
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.AvgCadence != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapAvgCadence)
		field.Value = *lap.AvgCadence
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.MaxCadence != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapMaxCadence)
		field.Value = *lap.MaxCadence
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.AvgPower != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapAvgPower)
		field.Value = *lap.AvgPower
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.MaxPower != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapMaxPower)
		field.Value = *lap.MaxPower
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.AvgTemperature != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapAvgTemperature)
		field.Value = *lap.AvgTemperature
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.MaxTemperature != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapMaxTemperature)
		field.Value = *lap.MaxTemperature
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.AvgAltitude != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapAvgAltitude)
		field.Value = scaleoffset.DiscardAny(*lap.AvgAltitude, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if lap.MaxAltitude != nil {
		field := factory.CreateField(mesgnum.Lap, fieldnum.LapMaxAltitude)
		field.Value = scaleoffset.DiscardAny(*lap.MaxAltitude, field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}

	eventField := factory.CreateField(mesgnum.Lap, fieldnum.LapEvent)
	eventField.Value = uint8(typedef.EventLap)
	mesg.Fields = append(mesg.Fields, eventField)

	eventTypeField := factory.CreateField(mesgnum.Lap, fieldnum.LapEventType)
	eventTypeField.Value = uint8(typedef.EventTypeStop)
	mesg.Fields = append(mesg.Fields, eventTypeField)

	return mesg
}
