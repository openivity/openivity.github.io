package fit

import (
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
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
			ses.Sport = activity.FormatTitle(typedef.Sport(sport).String())
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

	return ses
}
