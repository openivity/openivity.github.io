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
	"github.com/openivity/activity-service/strutils"
)

// Session is a workout session, it wraps FIT SDK's mesgdef.Session as its base.
type Session struct {
	*mesgdef.Session

	Laps    []Lap
	Records []Record
}

// CreateSession creates new session.
func CreateSession(session *mesgdef.Session) Session {
	if session == nil {
		session = mesgdef.NewSession(nil)
	}
	return Session{Session: session}
}

// EndTime returns sessions's EndTime (StartTime + TotalElapsedTime).
func (s *Session) EndTime() time.Time {
	if s.StartTime.IsZero() {
		return time.Time{}
	}
	if s.TotalElapsedTime == basetype.Uint32Invalid {
		return time.Time{}
	}
	return s.StartTime.Add(
		time.Duration(float64(s.TotalElapsedTime)/1000) * time.Second,
	)
}

// NewSessionFromLaps creates new session from Laps.
func NewSessionFromLaps(laps []Lap) Session {
	ses := CreateSession(
		mesgdef.NewSession(nil).
			SetStartTime(laps[0].StartTime))

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

	for i := range laps {
		lap := laps[i]

		if ses.Sport == typedef.SportInvalid {
			ses.Sport = lap.Sport
		}

		if lap.TotalElapsedTime != basetype.Uint32Invalid {
			ses.TotalElapsedTime += lap.TotalElapsedTime
		}
		if lap.TotalMovingTime != basetype.Uint32Invalid {
			ses.TotalMovingTime += lap.TotalMovingTime
		}
		if lap.TotalTimerTime != basetype.Uint32Invalid {
			ses.TotalTimerTime += lap.TotalTimerTime
		}
		if lap.TotalDistance != basetype.Uint32Invalid {
			ses.TotalDistance += lap.TotalDistance
		}
		if lap.TotalAscent != basetype.Uint16Invalid {
			ses.TotalAscent += lap.TotalAscent
		}
		if lap.TotalDescent != basetype.Uint16Invalid {
			ses.TotalDescent += lap.TotalDescent
		}
		if lap.TotalCalories != basetype.Uint16Invalid {
			ses.TotalCalories += lap.TotalCalories
		}
		if lap.AvgSpeed != basetype.Uint16Invalid {
			avgSpeed += uint64(lap.AvgSpeed)
			avgSpeedCount++
		}
		if lap.MaxSpeed != basetype.Uint16Invalid {
			if ses.MaxSpeed == basetype.Uint16Invalid || lap.MaxSpeed > ses.MaxSpeed {
				ses.MaxSpeed = lap.MaxSpeed
			}
		}
		if lap.AvgAltitude != basetype.Uint16Invalid {
			avgAltitude += uint64(lap.AvgAltitude)
			avgAltitudeCount++
		}
		if lap.MaxAltitude != basetype.Uint16Invalid {
			if ses.MaxAltitude == basetype.Uint16Invalid || lap.MaxAltitude > ses.MaxAltitude {
				ses.MaxAltitude = lap.MaxAltitude
			}
		}
		if lap.AvgCadence != basetype.Uint8Invalid {
			avgCadence += uint64(lap.AvgCadence)
			avgCadenceCount++
		}
		if lap.MaxCadence != basetype.Uint8Invalid {
			if ses.MaxCadence == basetype.Uint8Invalid || lap.MaxCadence > ses.MaxCadence {
				ses.MaxCadence = lap.MaxCadence
			}
		}
		if lap.AvgHeartRate != basetype.Uint8Invalid {
			avgHeartRate += uint64(lap.AvgHeartRate)
			avgHeartRateCount++
		}
		if lap.MaxHeartRate != basetype.Uint8Invalid {
			if ses.MaxHeartRate == basetype.Uint8Invalid || lap.MaxHeartRate > ses.MaxHeartRate {
				ses.MaxHeartRate = lap.MaxHeartRate
			}
		}
		if lap.AvgPower != basetype.Uint16Invalid {
			avgPower += uint64(lap.AvgPower)
			avgPowerCount++
		}
		if lap.MaxPower != basetype.Uint16Invalid {
			if ses.MaxPower == basetype.Uint16Invalid || lap.MaxPower > ses.MaxPower {
				ses.MaxPower = lap.MaxPower
			}
		}
		if lap.AvgTemperature != basetype.Sint8Invalid {
			avgTemperature += int64(lap.AvgTemperature)
			avgTemperatureCount++
		}
		if lap.MaxTemperature != basetype.Sint8Invalid {
			if ses.MaxTemperature == basetype.Sint8Invalid || lap.MaxTemperature > ses.MaxTemperature {
				ses.MaxTemperature = lap.MaxTemperature
			}
		}
	}

	if avgSpeedCount != 0 {
		ses.AvgSpeed = uint16(avgSpeed / avgSpeedCount)
	}
	if avgAltitudeCount != 0 {
		ses.AvgAltitude = uint16(avgAltitude / avgAltitudeCount)
	}
	if avgCadenceCount != 0 {
		ses.AvgCadence = uint8(avgCadence / avgCadenceCount)
	}
	if avgHeartRateCount != 0 {
		ses.AvgHeartRate = uint8(avgHeartRate / avgHeartRateCount)
	}
	if avgPowerCount != 0 {
		ses.AvgPower = uint16(avgPower / avgPowerCount)
	}
	if avgTemperatureCount != 0 {
		ses.AvgTemperature = int8(avgTemperature / avgTemperatureCount)
	}

	return ses
}

// IsBelongToThisLap check whether given t is belong to this session's time window.
func (s *Session) IsBelongToThisSession(t time.Time) bool {
	return isBelong(t, s.StartTime, s.EndTime())
}

// PutLaps puts given laps into session and return any remaining laps that doesn't belong to this session.
func (s *Session) PutLaps(laps ...Lap) (remainings []Lap) {
	var pos int
	for i := range laps {
		if s.IsBelongToThisSession(laps[i].StartTime) {
			laps[i], laps[pos] = laps[pos], laps[i]
			pos++
		}
	}
	s.Laps = append(s.Laps, laps[:pos]...)
	return laps[pos:]
}

// PutRecords puts given records into session and return any remaining records that doesn't belong to this session.
func (s *Session) PutRecords(records ...Record) (remainings []Record) {
	var pos int
	for i := range records {
		if s.IsBelongToThisSession(records[i].Timestamp) {
			records[i], records[pos] = records[pos], records[i]
			pos++
		}
	}
	s.Records = append(s.Records, records[:pos]...)
	return records[pos:]
}

// ReplaceValues replaces values with the corresponding values in the given session.
func (s *Session) ReplaceValues(session *Session) {
	if !session.StartTime.IsZero() {
		s.StartTime = session.StartTime
	}
	if session.TotalElapsedTime != basetype.Uint32Invalid {
		s.TotalElapsedTime = session.TotalElapsedTime
	}
	if session.TotalMovingTime != basetype.Uint32Invalid {
		s.TotalMovingTime = session.TotalMovingTime
	}
	if session.TotalTimerTime != basetype.Uint32Invalid {
		s.TotalTimerTime = session.TotalTimerTime
	}
	if session.TotalDistance != basetype.Uint32Invalid {
		s.TotalDistance = session.TotalDistance
	}
	if session.TotalCalories != basetype.Uint16Invalid {
		s.TotalCalories = session.TotalCalories
	}
	if session.TotalAscent != basetype.Uint16Invalid {
		s.TotalAscent = session.TotalAscent
	}
	if session.TotalDescent != basetype.Uint16Invalid {
		s.TotalDescent = session.TotalDescent
	}
	if session.AvgSpeed != basetype.Uint16Invalid {
		s.AvgSpeed = session.AvgSpeed
	}
	if session.MaxSpeed != basetype.Uint16Invalid {
		s.MaxSpeed = session.MaxSpeed
	}
	if session.AvgHeartRate != basetype.Uint8Invalid {
		s.AvgHeartRate = session.AvgHeartRate
	}
	if session.MaxHeartRate != basetype.Uint8Invalid {
		s.MaxHeartRate = session.MaxHeartRate
	}
	if session.AvgCadence != basetype.Uint8Invalid {
		s.AvgCadence = session.AvgCadence
	}
	if session.MaxCadence != basetype.Uint8Invalid {
		s.MaxCadence = session.MaxCadence
	}
	if session.AvgPower != basetype.Uint16Invalid {
		s.AvgPower = session.AvgPower
	}
	if session.MaxPower != basetype.Uint16Invalid {
		s.MaxPower = session.MaxPower
	}
	if session.AvgTemperature != basetype.Sint8Invalid {
		s.AvgTemperature = session.AvgTemperature
	}
	if session.MaxTemperature != basetype.Sint8Invalid {
		s.MaxTemperature = session.MaxTemperature
	}
	if session.AvgAltitude != basetype.Uint16Invalid {
		s.AvgAltitude = session.AvgAltitude
	}
	if session.MaxAltitude != basetype.Uint16Invalid {
		s.MaxAltitude = session.MaxAltitude
	}
}

// Accumulate accumulates values between two sesions.
func (s *Session) Accumulate(session *Session) {
	gap := (session.StartTime.Sub(s.EndTime()).Seconds() * 1000)
	s.TotalElapsedTime += uint32(gap)
	s.TotalTimerTime += uint32(gap)

	s.TotalMovingTime += session.TotalMovingTime
	s.TotalElapsedTime += session.TotalElapsedTime
	s.TotalTimerTime += session.TotalTimerTime

	s.TotalDistance += session.TotalDistance
	s.TotalAscent += session.TotalAscent
	s.TotalDescent += session.TotalDescent
	s.TotalCalories += session.TotalCalories

	if s.AvgSpeed != basetype.Uint16Invalid && session.AvgSpeed != basetype.Uint16Invalid {
		s.AvgSpeed = uint16((uint32(s.AvgSpeed) + uint32(session.AvgSpeed)) / 2)
	} else if session.AvgSpeed != basetype.Uint16Invalid {
		s.AvgSpeed = session.AvgSpeed
	}
	if s.MaxSpeed != basetype.Uint16Invalid && session.MaxSpeed != basetype.Uint16Invalid {
		s.MaxSpeed = uint16((uint32(s.MaxSpeed) + uint32(session.MaxSpeed)) / 2)
	} else if session.MaxSpeed != basetype.Uint16Invalid {
		s.MaxSpeed = session.MaxSpeed
	}
	if s.AvgHeartRate != basetype.Uint8Invalid && session.AvgHeartRate != basetype.Uint8Invalid {
		s.AvgHeartRate = uint8((uint16(s.AvgHeartRate) + uint16(session.AvgHeartRate)) / 2)
	} else if session.AvgHeartRate != basetype.Uint8Invalid {
		s.AvgHeartRate = session.AvgHeartRate
	}
	if s.MaxHeartRate != basetype.Uint8Invalid && session.MaxHeartRate != basetype.Uint8Invalid {
		s.MaxHeartRate = uint8((uint16(s.MaxHeartRate) + uint16(session.MaxHeartRate)) / 2)
	} else if session.MaxHeartRate != basetype.Uint8Invalid {
		s.MaxHeartRate = session.MaxHeartRate
	}
	if s.AvgCadence != basetype.Uint8Invalid && session.AvgCadence != basetype.Uint8Invalid {
		s.AvgCadence = uint8((uint16(s.AvgCadence) + uint16(session.AvgCadence)) / 2)
	} else if session.AvgCadence != basetype.Uint8Invalid {
		s.AvgCadence = session.AvgCadence
	}
	if s.MaxCadence != basetype.Uint8Invalid && session.MaxCadence != basetype.Uint8Invalid {
		s.MaxCadence = uint8((uint16(s.MaxCadence) + uint16(session.MaxCadence)) / 2)
	} else if session.MaxCadence != basetype.Uint8Invalid {
		s.MaxCadence = session.MaxCadence
	}
	if s.AvgPower != basetype.Uint16Invalid && session.AvgPower != basetype.Uint16Invalid {
		s.AvgPower = uint16((uint32(s.AvgPower) + uint32(session.AvgPower)) / 2)
	} else if session.AvgPower != basetype.Uint16Invalid {
		s.AvgPower = session.AvgPower
	}
	if s.MaxPower != basetype.Uint16Invalid && session.MaxPower != basetype.Uint16Invalid {
		s.MaxPower = uint16((uint32(s.MaxPower) + uint32(session.MaxPower)) / 2)
	} else if session.MaxPower != basetype.Uint16Invalid {
		s.MaxPower = session.MaxPower
	}
	if s.AvgTemperature != basetype.Sint8Invalid && session.AvgTemperature != basetype.Sint8Invalid {
		s.AvgTemperature = int8((int16(s.AvgTemperature) + int16(session.AvgTemperature)) / 2)
	} else if session.AvgTemperature != basetype.Sint8Invalid {
		s.AvgTemperature = session.AvgTemperature
	}
	if s.MaxTemperature != basetype.Sint8Invalid && session.MaxTemperature != basetype.Sint8Invalid {
		s.MaxTemperature = int8((int16(s.MaxTemperature) + int16(session.MaxTemperature)) / 2)
	} else if session.MaxTemperature != basetype.Sint8Invalid {
		s.MaxTemperature = session.MaxTemperature
	}
	if s.AvgAltitude != basetype.Uint16Invalid && session.AvgAltitude != basetype.Uint16Invalid {
		s.AvgAltitude = uint16((uint32(s.AvgAltitude) + uint32(session.AvgAltitude)) / 2)
	} else if session.AvgAltitude != basetype.Uint16Invalid {
		s.AvgAltitude = session.AvgAltitude
	}
	if s.MaxAltitude != basetype.Uint16Invalid && session.MaxAltitude != basetype.Uint16Invalid {
		s.MaxAltitude = uint16((uint32(s.MaxAltitude) + uint32(session.MaxAltitude)) / 2)
	} else if session.MaxAltitude != basetype.Uint16Invalid {
		s.MaxAltitude = session.MaxAltitude
	}
}

// Summarize summarizes the session such as updating StartPosition and EndPosition based on records.
func (s *Session) Summarize() {
	// Update GPS Positions
	for i := range s.Records {
		rec := &s.Records[i]
		if rec.PositionLat != basetype.Sint32Invalid && rec.PositionLong != basetype.Sint32Invalid {
			s.StartPositionLat = rec.PositionLat
			s.StartPositionLong = rec.PositionLong
			break
		}
	}
	for i := len(s.Records) - 1; i >= 0; i-- {
		rec := &s.Records[i]
		if rec.PositionLat != basetype.Sint32Invalid && rec.PositionLong != basetype.Sint32Invalid {
			s.EndPositionLat = rec.PositionLat
			s.EndPositionLong = rec.PositionLong
			break
		}
	}

	for i := len(s.Records) - 1; i >= 0; i-- {
		if !s.Records[i].Timestamp.IsZero() {
			s.TotalElapsedTime = uint32(s.Records[i].Timestamp.Sub(s.StartTime).Seconds() * 1000)
			break
		}
	}

	if s.TotalMovingTime == basetype.Uint32Invalid {
		s.TotalMovingTime = TotalMovingTime(s.Records, s.Sport)
	}

	if s.TotalAscent == basetype.Uint16Invalid || s.TotalDescent == basetype.Uint16Invalid {
		s.TotalAscent, s.TotalDescent = TotalAscentAndDescent(s.Records)
	}

	if s.AvgSpeed == basetype.Uint16Invalid || s.MaxSpeed == basetype.Uint16Invalid {
		s.AvgSpeed, s.MaxSpeed = AvgMaxSpeed(s.Records)
	}
}

// MarshalAppendJSON appends the JSON format encoding of Session to b, returning the result.
func (s *Session) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')
	if s.Sport != typedef.SportInvalid {
		b = append(b, `"sport":`...)
		b = strconv.AppendQuote(b, strutils.ToTitle(s.Sport.String()))
		b = append(b, ',')
	}
	if !s.Timestamp.IsZero() {
		b = append(b, `"timestamp":`...)
		b = strconv.AppendQuote(b, s.Timestamp.Format(time.RFC3339))
		b = append(b, ',')
	}
	if !s.StartTime.IsZero() {
		b = append(b, `"startTime":`...)
		b = strconv.AppendQuote(b, s.StartTime.Format(time.RFC3339))
		b = append(b, ',')
	}
	if !s.EndTime().IsZero() {
		b = append(b, `"endTime":`...)
		b = strconv.AppendQuote(b, s.EndTime().Format(time.RFC3339))
		b = append(b, ',')
	}
	if s.TotalElapsedTime != basetype.Uint32Invalid {
		b = append(b, `"totalElapsedTime":`...)
		b = strconv.AppendFloat(b, s.TotalElapsedTimeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if s.TotalMovingTime != basetype.Uint32Invalid {
		b = append(b, `"totalMovingTime":`...)
		b = strconv.AppendFloat(b, s.TotalMovingTimeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if s.TotalTimerTime != basetype.Uint32Invalid {
		b = append(b, `"totalTimerTime":`...)
		b = strconv.AppendFloat(b, s.TotalTimerTimeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if s.TotalDistance != basetype.Uint32Invalid {
		b = append(b, `"totalDistance":`...)
		b = strconv.AppendFloat(b, s.TotalDistanceScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if s.TotalAscent != basetype.Uint16Invalid {
		b = append(b, `"totalAscent":`...)
		b = strconv.AppendUint(b, uint64(s.TotalAscent), 10)
		b = append(b, ',')
	}
	if s.TotalDescent != basetype.Uint16Invalid {
		b = append(b, `"totalDescent":`...)
		b = strconv.AppendUint(b, uint64(s.TotalDescent), 10)
		b = append(b, ',')
	}
	if s.TotalCalories != basetype.Uint16Invalid {
		b = append(b, `"totalCalories":`...)
		b = strconv.AppendUint(b, uint64(s.TotalCalories), 10)
		b = append(b, ',')
	}
	if s.AvgSpeed != basetype.Uint16Invalid {
		b = append(b, `"avgSpeed":`...)
		b = strconv.AppendFloat(b, s.AvgSpeedScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if s.MaxSpeed != basetype.Uint16Invalid {
		b = append(b, `"maxSpeed":`...)
		b = strconv.AppendFloat(b, s.MaxSpeedScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if s.AvgHeartRate != basetype.Uint8Invalid {
		b = append(b, `"avgHeartRate":`...)
		b = strconv.AppendUint(b, uint64(s.AvgHeartRate), 10)
		b = append(b, ',')
	}
	if s.MaxHeartRate != basetype.Uint8Invalid {
		b = append(b, `"maxHeartRate":`...)
		b = strconv.AppendUint(b, uint64(s.MaxHeartRate), 10)
		b = append(b, ',')
	}
	if s.AvgCadence != basetype.Uint8Invalid {
		b = append(b, `"avgCadence":`...)
		b = strconv.AppendUint(b, uint64(s.AvgCadence), 10)
		b = append(b, ',')
	}
	if s.MaxCadence != basetype.Uint8Invalid {
		b = append(b, `"maxCadence":`...)
		b = strconv.AppendUint(b, uint64(s.MaxCadence), 10)
		b = append(b, ',')
	}
	if s.AvgPower != basetype.Uint16Invalid {
		b = append(b, `"avgPower":`...)
		b = strconv.AppendUint(b, uint64(s.AvgPower), 10)
		b = append(b, ',')
	}
	if s.MaxPower != basetype.Uint16Invalid {
		b = append(b, `"maxPower":`...)
		b = strconv.AppendUint(b, uint64(s.MaxPower), 10)
		b = append(b, ',')
	}
	if s.AvgTemperature != basetype.Sint8Invalid {
		b = append(b, `"avgTemperature":`...)
		b = strconv.AppendInt(b, int64(s.AvgTemperature), 10)
		b = append(b, ',')
	}
	if s.MaxTemperature != basetype.Sint8Invalid {
		b = append(b, `"maxTemperature":`...)
		b = strconv.AppendInt(b, int64(s.MaxTemperature), 10)
		b = append(b, ',')
	}
	if s.AvgAltitude != basetype.Uint16Invalid {
		b = append(b, `"avgAltitude":`...)
		b = strconv.AppendFloat(b, s.AvgAltitudeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}
	if s.MaxAltitude != basetype.Uint16Invalid {
		b = append(b, `"maxAltitude":`...)
		b = strconv.AppendFloat(b, s.MaxAltitudeScaled(), 'g', -1, 64)
		b = append(b, ',')
	}

	if HasPace(s.Sport) {
		avgPace := s.TotalMovingTimeScaled() / (s.TotalDistanceScaled() / 1000)
		if !math.IsNaN(avgPace) && !math.IsInf(avgPace, 0) {
			b = append(b, `"avgPace":`...)
			b = strconv.AppendFloat(b, avgPace, 'g', -1, 64)
			b = append(b, ',')
		}
		avgElapsedPace := s.TotalElapsedTimeScaled() / (s.TotalDistanceScaled() / 1000)
		if !math.IsNaN(avgElapsedPace) && !math.IsInf(avgElapsedPace, 0) {
			b = append(b, `"avgElapsedPace":`...)
			b = strconv.AppendFloat(b, avgElapsedPace, 'g', -1, 64)
			b = append(b, ',')
		}
	}

	b = append(b, `"laps":[`...)
	for i := range s.Laps {
		n := len(b)
		b = s.Laps[i].MarshalAppendJSON(b)
		if len(b) != n && i != len(s.Laps)-1 {
			b = append(b, ',')
		}
	}
	b = append(b, ']')
	b = append(b, ',')

	b = append(b, `"records":[`...)
	for i := range s.Records {
		n := len(b)
		b = s.Records[i].MarshalAppendJSON(b)
		if len(b) != n && i != len(s.Records)-1 {
			b = append(b, ',')
		}
	}
	b = append(b, ']')
	b = append(b, ',')

	if b[len(b)-1] == '{' {
		return b[:len(b)-1]
	}
	if b[len(b)-1] == ',' {
		b = b[:len(b)-1]
	}

	return append(b, '}')
}
