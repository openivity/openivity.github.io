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

	for i := range laps {
		aggregator.Aggregate(ses, laps[i].Lap)
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

	avgSpeed := s.AvgSpeedScaled()
	if math.IsNaN(avgSpeed) {
		avgSpeed = s.EnhancedAvgSpeedScaled()
	}
	if !math.IsNaN(avgSpeed) {
		b = append(b, `"avgSpeed":`...)
		b = strconv.AppendFloat(b, avgSpeed, 'g', -1, 64)
		b = append(b, ',')
	}

	maxSpeed := s.MaxSpeedScaled()
	if math.IsNaN(maxSpeed) {
		maxSpeed = s.EnhancedMaxSpeedScaled()
	}
	if !math.IsNaN(maxSpeed) {
		b = append(b, `"maxSpeed":`...)
		b = strconv.AppendFloat(b, maxSpeed, 'g', -1, 64)
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

	avgAltitude := s.AvgAltitudeScaled()
	if math.IsNaN(avgAltitude) {
		avgAltitude = s.EnhancedAvgAltitudeScaled()
	}
	if !math.IsNaN(avgAltitude) {
		b = append(b, `"avgAltitude":`...)
		b = strconv.AppendFloat(b, avgAltitude, 'g', -1, 64)
		b = append(b, ',')
	}

	maxAltitude := s.MaxAltitudeScaled()
	if math.IsNaN(maxAltitude) {
		maxAltitude = s.EnhancedMaxAltitudeScaled()
	}
	if !math.IsNaN(maxAltitude) {
		b = append(b, `"maxAltitude":`...)
		b = strconv.AppendFloat(b, maxAltitude, 'g', -1, 64)
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
