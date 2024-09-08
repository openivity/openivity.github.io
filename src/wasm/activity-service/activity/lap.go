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

package activity

import (
	"math"
	"strconv"
	"time"

	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/openivity/activity-service/aggregator"
	"github.com/openivity/activity-service/strutils"
)

// Lap is a workout lap. It wraps FIT SDK's mesgdef.Lap as its base.
type Lap struct {
	*mesgdef.Lap
}

// CreateLap creates new lap.
func CreateLap(lap *mesgdef.Lap) Lap {
	if lap == nil {
		lap = mesgdef.NewLap(nil)
	}
	return Lap{Lap: lap}
}

// EndTime returns lap's EndTime (StartTime + TotalElapsedTime).
func (l *Lap) EndTime() time.Time {
	if l.StartTime.IsZero() {
		return time.Time{}
	}
	if l.TotalElapsedTime == basetype.Uint32Invalid {
		return time.Time{}
	}
	return l.StartTime.Add(
		time.Duration(float64(l.TotalElapsedTime)/1000) * time.Second,
	)
}

// CreateLap creates new lap from records.
func NewLapFromRecords(records []Record, sport typedef.Sport) Lap {
	lap := CreateLap(
		mesgdef.NewLap(nil).
			SetSport(sport).
			SetTimestamp(records[0].Timestamp).
			SetStartTime(records[0].Timestamp))

	var (
		avgSpeed            uint64 = 0
		avgSpeedCount       uint64 = 0
		avgAltitude         uint64 = 0
		avgAltitudeCount    uint64 = 0
		avgCadence          uint64 = 0
		avgCadenceCount     uint64 = 0
		avgHeartRate        uint64 = 0
		avgHeartRateCount   uint64 = 0
		avgPower            uint64 = 0
		avgPowerCount       uint64 = 0
		avgTemperature      int64  = 0
		avgTemperatureCount int64  = 0
	)

	for i := 0; i < len(records); i++ {
		rec := &records[i]

		if rec.Speed != basetype.Uint16Invalid {
			avgSpeed += uint64(rec.Speed)
			avgSpeedCount++
			if lap.MaxSpeed == basetype.Uint16Invalid || rec.Speed > lap.MaxSpeed {
				lap.MaxSpeed = rec.Speed
			}
		}
		if rec.Altitude != basetype.Uint16Invalid {
			avgAltitude += uint64(rec.Altitude)
			avgAltitudeCount++
			if lap.MaxAltitude == basetype.Uint16Invalid || rec.Altitude > lap.MaxAltitude {
				lap.MaxAltitude = rec.Altitude
			}
		}
		if rec.Cadence != basetype.Uint8Invalid {
			avgCadence += uint64(rec.Cadence)
			avgCadenceCount++
			if lap.MaxCadence == basetype.Uint8Invalid || rec.Cadence > lap.MaxCadence {
				lap.MaxCadence = rec.Cadence
			}
		}
		if rec.HeartRate != basetype.Uint8Invalid {
			avgHeartRate += uint64(rec.HeartRate)
			avgHeartRateCount++
			if lap.MaxHeartRate == basetype.Uint8Invalid || rec.HeartRate > lap.MaxHeartRate {
				lap.MaxHeartRate = rec.HeartRate
			}
		}
		if rec.Power != basetype.Uint16Invalid {
			avgPower += uint64(rec.Power)
			avgPowerCount++
			if lap.MaxPower == basetype.Uint16Invalid || rec.Power > lap.MaxPower {
				lap.MaxPower = rec.Power
			}
		}
		if rec.Temperature != basetype.Sint8Invalid {
			avgTemperature += int64(rec.Temperature)
			avgTemperatureCount++
			if lap.MaxTemperature == basetype.Sint8Invalid || rec.Temperature > lap.MaxTemperature {
				lap.MaxTemperature = rec.Temperature
			}
		}
	}

	if avgSpeedCount != 0 {
		lap.AvgSpeed = uint16(avgSpeed / avgSpeedCount)
	}
	if avgAltitudeCount != 0 {
		lap.AvgAltitude = uint16(avgAltitude / avgAltitudeCount)
	}
	if avgCadenceCount != 0 {
		lap.AvgCadence = uint8(avgCadence / avgCadenceCount)
	}
	if avgHeartRateCount != 0 {
		lap.AvgHeartRate = uint8(avgHeartRate / avgHeartRateCount)
	}
	if avgPowerCount != 0 {
		lap.AvgPower = uint16(avgPower / avgPowerCount)
	}
	if avgTemperatureCount != 0 {
		lap.AvgTemperature = int8(avgTemperature / avgTemperatureCount)
	}

	var startDistance, endDistance uint32 = basetype.Uint32Invalid, basetype.Uint32Invalid
	for i := 0; i < len(records); i++ {
		if records[i].Distance != basetype.Uint32Invalid {
			startDistance = records[i].Distance
			break
		}
	}
	for i := len(records) - 1; i >= 0; i-- {
		if records[i].Distance != basetype.Uint32Invalid {
			endDistance = records[i].Distance
			break
		}
	}
	lap.TotalDistance = endDistance - startDistance

	var startTimestamp, endTimestamp time.Time
	for i := 0; i < len(records); i++ {
		if !records[i].Timestamp.IsZero() {
			startTimestamp = records[i].Timestamp
			break
		}
	}
	for i := len(records) - 1; i >= 0; i-- {
		if !records[i].Timestamp.IsZero() {
			endTimestamp = records[i].Timestamp
			break
		}
	}
	if !startTimestamp.IsZero() && !endTimestamp.IsZero() {
		lap.TotalElapsedTime = uint32((endTimestamp.Sub(startTimestamp).Seconds() + 1) * 1000)
	}
	lap.TotalTimerTime = lap.TotalElapsedTime

	lap.TotalMovingTime = TotalMovingTime(records, lap.Sport)
	lap.TotalAscent, lap.TotalDescent = TotalAscentAndDescent(records)

	return lap
}

// NewLapFromSession creates new lap from a session.
func NewLapFromSession(session *Session) Lap {
	lap := CreateLap(nil)
	aggregator.Replace(lap.Lap, session.Session)
	return lap
}

// IsBelongToThisLap check whether given t is belong to this lap's time window.
func (l *Lap) IsBelongToThisLap(t time.Time) bool {
	return isBelong(t, l.StartTime, l.EndTime())
}

// MarshalAppendJSON appends the JSON format encoding of Lap to b, returning the result.
func (l *Lap) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')

	b = append(b, `"sport":`...)
	b = strconv.AppendQuote(b, strutils.ToTitle(l.Sport.String()))
	b = append(b, ',')

	if !l.Timestamp.IsZero() {
		b = append(b, `"timestamp":`...)
		b = strconv.AppendQuote(b, l.Timestamp.Format(time.RFC3339))
		b = append(b, ',')
	}
	if !l.StartTime.IsZero() {
		b = append(b, `"startTime":`...)
		b = strconv.AppendQuote(b, l.StartTime.Format(time.RFC3339))
		b = append(b, ',')
	}
	if !l.EndTime().IsZero() {
		b = append(b, `"endTime":`...)
		b = strconv.AppendQuote(b, l.EndTime().Format(time.RFC3339))
		b = append(b, ',')
	}
	if l.TotalElapsedTime != basetype.Uint32Invalid {
		b = append(b, `"totalElapsedTime":`...)
		b = strconv.AppendFloat(b, l.TotalElapsedTimeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if l.TotalMovingTime != basetype.Uint32Invalid {
		b = append(b, `"totalMovingTime":`...)
		b = strconv.AppendFloat(b, l.TotalMovingTimeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if l.TotalTimerTime != basetype.Uint32Invalid {
		b = append(b, `"totalTimerTime":`...)
		b = strconv.AppendFloat(b, l.TotalTimerTimeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if l.TotalDistance != basetype.Uint32Invalid {
		b = append(b, `"totalDistance":`...)
		b = strconv.AppendFloat(b, l.TotalDistanceScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if l.TotalAscent != basetype.Uint16Invalid {
		b = append(b, `"totalAscent":`...)
		b = strconv.AppendUint(b, uint64(l.TotalAscent), 10)
		b = append(b, ',')
	}
	if l.TotalDescent != basetype.Uint16Invalid {
		b = append(b, `"totalDescent":`...)
		b = strconv.AppendUint(b, uint64(l.TotalDescent), 10)
		b = append(b, ',')
	}
	if l.TotalCalories != basetype.Uint16Invalid {
		b = append(b, `"totalCalories":`...)
		b = strconv.AppendUint(b, uint64(l.TotalCalories), 10)
		b = append(b, ',')
	}
	if l.AvgSpeed != basetype.Uint16Invalid {
		b = append(b, `"avgSpeed":`...)
		b = strconv.AppendFloat(b, l.AvgSpeedScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if l.MaxSpeed != basetype.Uint16Invalid {
		b = append(b, `"maxSpeed":`...)
		b = strconv.AppendFloat(b, l.MaxSpeedScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if l.AvgHeartRate != basetype.Uint8Invalid {
		b = append(b, `"avgHeartRate":`...)
		b = strconv.AppendUint(b, uint64(l.AvgHeartRate), 10)
		b = append(b, ',')
	}
	if l.MaxHeartRate != basetype.Uint8Invalid {
		b = append(b, `"maxHeartRate":`...)
		b = strconv.AppendUint(b, uint64(l.MaxHeartRate), 10)
		b = append(b, ',')
	}
	if l.AvgCadence != basetype.Uint8Invalid {
		b = append(b, `"avgCadence":`...)
		b = strconv.AppendUint(b, uint64(l.AvgCadence), 10)
		b = append(b, ',')
	}
	if l.MaxCadence != basetype.Uint8Invalid {
		b = append(b, `"maxCadence":`...)
		b = strconv.AppendUint(b, uint64(l.MaxCadence), 10)
		b = append(b, ',')
	}
	if l.AvgPower != basetype.Uint16Invalid {
		b = append(b, `"avgPower":`...)
		b = strconv.AppendUint(b, uint64(l.AvgPower), 10)
		b = append(b, ',')
	}
	if l.MaxPower != basetype.Uint16Invalid {
		b = append(b, `"maxPower":`...)
		b = strconv.AppendUint(b, uint64(l.MaxPower), 10)
		b = append(b, ',')
	}
	if l.AvgTemperature != basetype.Sint8Invalid {
		b = append(b, `"avgTemperature":`...)
		b = strconv.AppendInt(b, int64(l.AvgTemperature), 10)
		b = append(b, ',')
	}
	if l.MaxTemperature != basetype.Sint8Invalid {
		b = append(b, `"maxTemperature":`...)
		b = strconv.AppendInt(b, int64(l.MaxTemperature), 10)
		b = append(b, ',')
	}
	if l.AvgAltitude != basetype.Uint16Invalid {
		b = append(b, `"avgAltitude":`...)
		b = strconv.AppendFloat(b, l.AvgAltitudeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if l.MaxAltitude != basetype.Uint16Invalid {
		b = append(b, `"maxAltitude":`...)
		b = strconv.AppendFloat(b, l.MaxAltitudeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}

	if HasPace(l.Sport) {
		avgPace := l.TotalMovingTimeScaled() / (l.TotalDistanceScaled() / 1000)
		if !math.IsNaN(avgPace) && !math.IsInf(avgPace, 0) {
			b = append(b, `"avgPace":`...)
			b = strconv.AppendFloat(b, avgPace, 'g', -1, 64)
			b = append(b, ',')
		}
		avgElapsedPace := l.TotalElapsedTimeScaled() / (l.TotalDistanceScaled() / 1000)
		if !math.IsNaN(avgElapsedPace) && !math.IsInf(avgElapsedPace, 0) {
			b = append(b, `"avgElapsedPace":`...)
			b = strconv.AppendFloat(b, avgElapsedPace, 'g', -1, 64)
			b = append(b, ',')
		}
	}

	if b[len(b)-1] == '{' {
		return b[:len(b)-1]
	}

	if b[len(b)-1] == ',' {
		b = b[:len(b)-1]
	}

	return append(b, '}')
}
