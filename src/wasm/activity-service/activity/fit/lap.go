package fit

import (
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
)

func NewLap(mesg proto.Message) *activity.Lap {
	lap := &activity.Lap{}

	for i := range mesg.Fields {
		field := &mesg.Fields[i]
		switch field.Num {
		case fieldnum.LapTimestamp:
			lap.Timestamp = datetime.ToTime(field.Value)
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

func NewLapFromSession(session *activity.Session) *activity.Lap {
	return &activity.Lap{
		Timestamp:        session.Timestamp,
		TotalMovingTime:  session.TotalMovingTime,
		TotalElapsedTime: session.TotalElapsedTime,
		TotalDistance:    session.TotalDistance,
		TotalAscent:      session.TotalAscent,
		TotalDescent:     session.TotalDescent,
		TotalCalories:    session.TotalCalories,
		AvgSpeed:         session.AvgSpeed,
		MaxSpeed:         session.MaxSpeed,
		AvgHeartRate:     session.AvgHeartRate,
		MaxHeartRate:     session.MaxHeartRate,
		AvgCadence:       session.AvgCadence,
		MaxCadence:       session.MaxCadence,
		AvgPower:         session.AvgPower,
		MaxPower:         session.MaxPower,
		AvgTemperature:   session.AvgTemperature,
		MaxTemperature:   session.MaxTemperature,
		AvgAltitude:      session.AvgAltitude,
		MaxAltitude:      session.MaxAltitude,
	}
}
