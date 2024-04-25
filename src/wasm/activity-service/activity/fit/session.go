// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package fit

import (
	"time"

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
			sport := field.Value.Uint8()
			if sport == basetype.EnumInvalid {
				continue
			}
			ses.Sport = kit.FormatTitle(typedef.Sport(sport).String())
		case fieldnum.SessionTotalMovingTime:
			totalMovingTime := field.Value.Uint32()
			if totalMovingTime == basetype.Uint32Invalid {
				continue
			}
			ses.TotalMovingTime = float64(totalMovingTime) / 1000
		case fieldnum.SessionTotalElapsedTime:
			totalElapsedTime := field.Value.Uint32()
			if totalElapsedTime == basetype.Uint32Invalid {
				continue
			}
			ses.TotalElapsedTime = float64(totalElapsedTime) / 1000
			ses.EndTime = ses.StartTime.Add(time.Duration(ses.TotalElapsedTime) * time.Second)
		case fieldnum.SessionTotalDistance:
			totalDistance := field.Value.Uint32()
			if totalDistance == basetype.Uint32Invalid {
				continue
			}
			ses.TotalDistance = float64(totalDistance) / 100
		case fieldnum.SessionTotalAscent:
			totalAscent := field.Value.Uint16()
			if totalAscent == basetype.Uint16Invalid {
				continue
			}
			ses.TotalAscent = totalAscent
		case fieldnum.SessionTotalDescent:
			totalDescent := field.Value.Uint16()
			if totalDescent == basetype.Uint16Invalid {
				continue
			}
			ses.TotalDescent = totalDescent
		case fieldnum.SessionTotalCalories:
			totalCalories := field.Value.Uint16()
			if totalCalories == basetype.Uint16Invalid {
				continue
			}
			ses.TotalCalories = totalCalories
		case fieldnum.SessionAvgSpeed:
			avgSpeed := field.Value.Uint16()
			if avgSpeed == basetype.Uint16Invalid {
				continue
			}
			ses.AvgSpeed = kit.Ptr(float64(avgSpeed) / 1000)
		case fieldnum.SessionMaxSpeed:
			maxSpeed := field.Value.Uint16()
			if maxSpeed == basetype.Uint16Invalid {
				continue
			}
			ses.MaxSpeed = kit.Ptr(float64(maxSpeed) / 1000)
		case fieldnum.SessionAvgHeartRate:
			avgHeartRate := field.Value.Uint8()
			if avgHeartRate == basetype.Uint8Invalid {
				continue
			}
			ses.AvgHeartRate = &avgHeartRate
		case fieldnum.SessionMaxHeartRate:
			maxHeartRate := field.Value.Uint8()
			if maxHeartRate == basetype.Uint8Invalid {
				continue
			}
			ses.MaxHeartRate = &maxHeartRate
		case fieldnum.SessionAvgCadence:
			avgCadence := field.Value.Uint8()
			if avgCadence != basetype.Uint8Invalid {
				continue
			}
			if avgCadence == basetype.Uint8Invalid {
				continue
			}
			ses.AvgCadence = &avgCadence
		case fieldnum.SessionMaxCadence:
			maxCadence := field.Value.Uint8()
			if maxCadence == basetype.Uint8Invalid {
				continue
			}
			ses.MaxCadence = &maxCadence
		case fieldnum.SessionAvgPower:
			avgPower := field.Value.Uint16()
			if avgPower == basetype.Uint16Invalid {
				continue
			}
			ses.AvgPower = &avgPower
		case fieldnum.SessionMaxPower:
			maxPower := field.Value.Uint16()
			if maxPower == basetype.Uint16Invalid {
				continue
			}
			ses.MaxPower = &maxPower
		case fieldnum.SessionAvgTemperature:
			avgTemperature := field.Value.Int8()
			if avgTemperature == basetype.Sint8Invalid {
				continue
			}
			ses.AvgTemperature = &avgTemperature
		case fieldnum.SessionMaxTemperature:
			maxTemperature := field.Value.Int8()
			if maxTemperature == basetype.Sint8Invalid {
				continue
			}
			ses.MaxTemperature = &maxTemperature
		case fieldnum.SessionAvgAltitude:
			avgAltitude := field.Value.Uint16()
			if avgAltitude == basetype.Uint16Invalid {
				continue
			}
			fAvgAltitude := (float64(avgAltitude) / 5) - 500
			if fAvgAltitude < 0 {
				continue
			}
			ses.AvgAltitude = &fAvgAltitude
		case fieldnum.SessionMaxAltitude:
			maxAltitude := field.Value.Uint16()
			if maxAltitude == basetype.Uint16Invalid {
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
		field.Value = proto.Uint32(datetime.ToUint32(ses.Timestamp))
		mesg.Fields = append(mesg.Fields, field)
	}
	if !ses.StartTime.IsZero() {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionStartTime)
		field.Value = proto.Uint32(datetime.ToUint32(ses.StartTime))
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.Sport != "" {
		sport := typedef.SportFromString(kit.FormatToLowerSnakeCase(ses.Sport))
		if sport == typedef.SportInvalid {
			sport = typedef.SportGeneric
		}
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionSport)
		field.Value = proto.Uint8(uint8(sport))
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalMovingTime != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalMovingTime)
		field.Value = scaleoffset.DiscardValue(proto.Float64(ses.TotalMovingTime), field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalElapsedTime != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalElapsedTime)
		field.Value = scaleoffset.DiscardValue(proto.Float64(ses.TotalElapsedTime), field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)

		totalTimerTimeField := factory.CreateField(mesgnum.Lap, fieldnum.SessionTotalElapsedTime)
		totalTimerTimeField.Value = field.Value
		mesg.Fields = append(mesg.Fields, totalTimerTimeField)
	}
	if ses.TotalDistance != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalDistance)
		field.Value = scaleoffset.DiscardValue(proto.Float64(ses.TotalDistance), field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalAscent != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalAscent)
		field.Value = proto.Uint16(ses.TotalAscent)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalDescent != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalDescent)
		field.Value = proto.Uint16(ses.TotalDescent)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.TotalCalories != 0 {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionTotalCalories)
		field.Value = proto.Uint16(ses.TotalCalories)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgSpeed != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgSpeed)
		field.Value = scaleoffset.DiscardValue(proto.Float64(*ses.AvgSpeed), field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxSpeed != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxSpeed)
		field.Value = scaleoffset.DiscardValue(proto.Float64(*ses.MaxSpeed), field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgHeartRate != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgHeartRate)
		field.Value = proto.Uint8(*ses.AvgHeartRate)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxHeartRate != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxHeartRate)
		field.Value = proto.Uint8(*ses.MaxHeartRate)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgCadence != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgCadence)
		field.Value = proto.Uint8(*ses.AvgCadence)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxCadence != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxCadence)
		field.Value = proto.Uint8(*ses.MaxCadence)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgPower != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgPower)
		field.Value = proto.Uint16(*ses.AvgPower)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxPower != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxPower)
		field.Value = proto.Uint16(*ses.MaxPower)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgTemperature != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgTemperature)
		field.Value = proto.Int8(*ses.AvgTemperature)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxTemperature != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxTemperature)
		field.Value = proto.Int8(*ses.MaxTemperature)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.AvgAltitude != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionAvgAltitude)
		field.Value = scaleoffset.DiscardValue(proto.Float64(*ses.AvgAltitude), field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}
	if ses.MaxAltitude != nil {
		field := factory.CreateField(mesgnum.Session, fieldnum.SessionMaxAltitude)
		field.Value = scaleoffset.DiscardValue(proto.Float64(*ses.MaxAltitude), field.Type.BaseType(), field.Scale, field.Offset)
		mesg.Fields = append(mesg.Fields, field)
	}

	eventField := factory.CreateField(mesgnum.Session, fieldnum.SessionEvent)
	eventField.Value = proto.Uint8(uint8(typedef.EventSession))
	mesg.Fields = append(mesg.Fields, eventField)

	eventTypeField := factory.CreateField(mesgnum.Session, fieldnum.SessionEventType)
	eventTypeField.Value = proto.Uint8(uint8(typedef.EventTypeStop))
	mesg.Fields = append(mesg.Fields, eventTypeField)

	return mesg
}
