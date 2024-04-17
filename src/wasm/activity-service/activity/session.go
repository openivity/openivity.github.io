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
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/accumulator"
	"github.com/muktihari/openactivity-fit/kit"
)

type Session struct {
	Timestamp        time.Time
	StartTime        time.Time
	EndTime          time.Time
	Sport            string
	TotalMovingTime  float64
	TotalElapsedTime float64
	TotalDistance    float64
	TotalAscent      uint16
	TotalDescent     uint16
	TotalCalories    uint16
	AvgSpeed         *float64
	MaxSpeed         *float64
	AvgHeartRate     *uint8
	MaxHeartRate     *uint8
	AvgCadence       *uint8
	MaxCadence       *uint8
	AvgPower         *uint16
	MaxPower         *uint16
	AvgTemperature   *int8
	MaxTemperature   *int8
	AvgAltitude      *float64
	MaxAltitude      *float64
	AvgPace          *float64
	AvgElapsedPace   *float64

	WorkoutType WorkoutType
	Laps        []*Lap
	Records     []*Record
}

type WorkoutType byte

const (
	WorkoutTypeMoving WorkoutType = iota
	WorkoutTypeStationary
)

func NewSessionFromLaps(laps []*Lap, sport string) *Session {
	if len(laps) == 0 {
		return nil
	}

	ses := &Session{
		Timestamp: laps[0].Timestamp,
		StartTime: laps[0].StartTime,
		EndTime:   laps[len(laps)-1].EndTime,
		Sport:     sport,
	}

	var (
		totalElapsedTimeAccumu = new(accumulator.Accumulator[float64])
		totalMovingTimeAccumu  = new(accumulator.Accumulator[float64])
		totalDistanceAccumu    = new(accumulator.Accumulator[float64])
		totalAscentAccumu      = new(accumulator.Accumulator[uint16])
		totalDescentAccumu     = new(accumulator.Accumulator[uint16])
		totalCaloriesAccumu    = new(accumulator.Accumulator[uint16])
		speedAvgAccumu         = new(accumulator.Accumulator[float64])
		speedMaxAccumu         = new(accumulator.Accumulator[float64])
		altitudeAvgAccumu      = new(accumulator.Accumulator[float64])
		altitudeMaxAccumu      = new(accumulator.Accumulator[float64])
		cadenceAvgAccumu       = new(accumulator.Accumulator[uint8])
		cadenceMaxAccumu       = new(accumulator.Accumulator[uint8])
		heartRateAvgAccumu     = new(accumulator.Accumulator[uint8])
		heartRateMaxAccumu     = new(accumulator.Accumulator[uint8])
		powerAvgAccumu         = new(accumulator.Accumulator[uint16])
		powerMaxAccumu         = new(accumulator.Accumulator[uint16])
		temperatureAvgAccumu   = new(accumulator.Accumulator[int8])
		temperatureMaxAccumu   = new(accumulator.Accumulator[int8])
	)

	for i := range laps {
		lap := laps[i]

		totalElapsedTimeAccumu.Collect(&lap.TotalElapsedTime)
		totalMovingTimeAccumu.Collect(&lap.TotalMovingTime)
		totalDistanceAccumu.Collect(&lap.TotalDistance)
		totalAscentAccumu.Collect(&lap.TotalAscent)
		totalDescentAccumu.Collect(&lap.TotalDescent)
		totalCaloriesAccumu.Collect(&lap.TotalCalories)
		speedAvgAccumu.Collect(lap.AvgSpeed)
		speedMaxAccumu.Collect(lap.MaxSpeed)
		altitudeAvgAccumu.Collect(lap.AvgAltitude)
		altitudeMaxAccumu.Collect(lap.MaxAltitude)
		cadenceAvgAccumu.Collect(lap.AvgCadence)
		cadenceMaxAccumu.Collect(lap.MaxCadence)
		heartRateAvgAccumu.Collect(lap.AvgHeartRate)
		heartRateMaxAccumu.Collect(lap.MaxHeartRate)
		powerAvgAccumu.Collect(lap.AvgPower)
		powerMaxAccumu.Collect(lap.MaxPower)
		temperatureAvgAccumu.Collect(lap.AvgTemperature)
		temperatureMaxAccumu.Collect(lap.MaxTemperature)
	}

	if value := totalElapsedTimeAccumu.Sum(); value != nil {
		ses.TotalElapsedTime = *value
	}
	if value := totalMovingTimeAccumu.Sum(); value != nil {
		ses.TotalMovingTime = *value
	}
	if value := totalDistanceAccumu.Sum(); value != nil {
		ses.TotalDistance = *value
	}
	if value := totalAscentAccumu.Sum(); value != nil {
		ses.TotalAscent = *value
	}
	if value := totalDescentAccumu.Sum(); value != nil {
		ses.TotalDescent = *value
	}
	if value := totalCaloriesAccumu.Sum(); value != nil {
		ses.TotalCalories = *value
	}

	ses.AvgSpeed = speedAvgAccumu.Avg()
	ses.MaxSpeed = speedMaxAccumu.Max()
	ses.AvgAltitude = altitudeAvgAccumu.Avg()
	ses.MaxAltitude = altitudeMaxAccumu.Max()
	ses.AvgCadence = cadenceAvgAccumu.Avg()
	ses.MaxCadence = cadenceMaxAccumu.Max()
	ses.AvgHeartRate = heartRateAvgAccumu.Avg()
	ses.MaxHeartRate = heartRateMaxAccumu.Max()
	ses.AvgPower = powerAvgAccumu.Avg()
	ses.MaxPower = powerMaxAccumu.Max()
	ses.AvgTemperature = temperatureAvgAccumu.Avg()
	ses.MaxTemperature = temperatureMaxAccumu.Max()

	if HasPace(sport) {
		ses.AvgPace = kit.Ptr(ses.TotalMovingTime / (ses.TotalDistance / 1000))
		ses.AvgElapsedPace = kit.Ptr(ses.TotalElapsedTime / (ses.TotalDistance / 1000))
	}

	return ses
}

var _ json.Marshaler = &Session{}

func (s *Session) MarshalJSON() ([]byte, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('{')

	buf.WriteString("\"sport\":\"")
	buf.WriteString(s.Sport)
	buf.WriteString("\",")

	if !s.Timestamp.IsZero() {
		buf.WriteString("\"timestamp\":\"")
		buf.WriteString(s.Timestamp.Format(time.RFC3339))
		buf.WriteString("\",")
	}
	if !s.StartTime.IsZero() {
		buf.WriteString("\"startTime\":\"")
		buf.WriteString(s.StartTime.Format(time.RFC3339))
		buf.WriteString("\",")
	}
	if !s.EndTime.IsZero() {
		buf.WriteString("\"endTime\":\"")
		buf.WriteString(s.EndTime.Format(time.RFC3339))
		buf.WriteString("\",")
	}

	buf.WriteString("\"totalMovingTime\":")
	buf.WriteString(strconv.FormatFloat(s.TotalMovingTime, 'g', -1, 64))
	buf.WriteByte(',')

	buf.WriteString("\"totalElapsedTime\":")
	buf.WriteString(strconv.FormatFloat(s.TotalElapsedTime, 'g', -1, 64))
	buf.WriteByte(',')

	buf.WriteString("\"totalDistance\":")
	buf.WriteString(strconv.FormatFloat(s.TotalDistance, 'g', -1, 64))
	buf.WriteByte(',')

	buf.WriteString("\"totalAscent\":")
	buf.WriteString(strconv.FormatUint(uint64(s.TotalAscent), 10))
	buf.WriteByte(',')

	buf.WriteString("\"totalDescent\":")
	buf.WriteString(strconv.FormatUint(uint64(s.TotalDescent), 10))
	buf.WriteByte(',')

	buf.WriteString("\"totalCalories\":")
	buf.WriteString(strconv.FormatUint(uint64(s.TotalCalories), 10))
	buf.WriteByte(',')

	if s.AvgSpeed != nil && !math.IsInf(*s.AvgSpeed, 0) && !math.IsNaN(*s.AvgSpeed) {
		buf.WriteString("\"avgSpeed\":")
		buf.WriteString(strconv.FormatFloat(*s.AvgSpeed, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if s.MaxSpeed != nil && !math.IsInf(*s.MaxSpeed, 0) && !math.IsNaN(*s.MaxSpeed) {
		buf.WriteString("\"maxSpeed\":")
		buf.WriteString(strconv.FormatFloat(*s.MaxSpeed, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if s.AvgHeartRate != nil {
		buf.WriteString("\"avgHeartRate\":")
		buf.WriteString(strconv.FormatUint(uint64(*s.AvgHeartRate), 10))
		buf.WriteByte(',')
	}
	if s.MaxHeartRate != nil {
		buf.WriteString("\"maxHeartRate\":")
		buf.WriteString(strconv.FormatUint(uint64(*s.MaxHeartRate), 10))
		buf.WriteByte(',')
	}
	if s.AvgCadence != nil {
		buf.WriteString("\"avgCadence\":")
		buf.WriteString(strconv.FormatUint(uint64(*s.AvgCadence), 10))
		buf.WriteByte(',')
	}
	if s.MaxCadence != nil {
		buf.WriteString("\"maxCadence\":")
		buf.WriteString(strconv.FormatUint(uint64(*s.MaxCadence), 10))
		buf.WriteByte(',')
	}
	if s.AvgPower != nil {
		buf.WriteString("\"avgPower\":")
		buf.WriteString(strconv.FormatUint(uint64(*s.AvgPower), 10))
		buf.WriteByte(',')
	}
	if s.MaxPower != nil {
		buf.WriteString("\"maxPower\":")
		buf.WriteString(strconv.FormatUint(uint64(*s.MaxPower), 10))
		buf.WriteByte(',')
	}
	if s.AvgTemperature != nil {
		buf.WriteString("\"avgTemperature\":")
		buf.WriteString(strconv.FormatInt(int64(*s.AvgTemperature), 10))
		buf.WriteByte(',')
	}
	if s.MaxTemperature != nil {
		buf.WriteString("\"maxTemperature\":")
		buf.WriteString(strconv.FormatInt(int64(*s.MaxTemperature), 10))
		buf.WriteByte(',')
	}
	if s.AvgAltitude != nil && !math.IsInf(*s.AvgAltitude, 0) && !math.IsNaN(*s.AvgAltitude) {
		buf.WriteString("\"avgAltitude\":")
		buf.WriteString(strconv.FormatFloat(*s.AvgAltitude, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if s.MaxAltitude != nil && !math.IsInf(*s.MaxAltitude, 0) && !math.IsNaN(*s.MaxAltitude) {
		buf.WriteString("\"maxAltitude\":")
		buf.WriteString(strconv.FormatFloat(*s.MaxAltitude, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if s.AvgPace != nil && !math.IsInf(*s.AvgPace, 0) && !math.IsNaN(*s.AvgPace) {
		buf.WriteString("\"avgPace\":")
		buf.WriteString(strconv.FormatFloat(*s.AvgPace, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if s.AvgElapsedPace != nil && !math.IsInf(*s.AvgElapsedPace, 0) && !math.IsNaN(*s.AvgElapsedPace) {
		buf.WriteString("\"avgElapsedPace\":")
		buf.WriteString(strconv.FormatFloat(*s.AvgElapsedPace, 'g', -1, 64))
		buf.WriteByte(',')
	}

	if len(s.Laps) != 0 {
		buf.WriteString("\"laps\": [")
		for i := range s.Laps {
			b, _ := s.Laps[i].MarshalJSON()
			buf.Write(b)
			if i != len(s.Laps)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteString("],")
	}

	if len(s.Records) != 0 {
		buf.WriteString("\"records\": [")
		for i := range s.Records {
			b, _ := s.Records[i].MarshalJSON()
			buf.Write(b)
			if i != len(s.Records)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte(']')
	}

	b := buf.Bytes()
	if b[len(b)-1] == ',' {
		b[len(b)-1] = '}'
		return b, nil
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (s *Session) IsBelongToThisSession(t time.Time) bool {
	return isBelong(t, s.StartTime, s.EndTime)
}

// PutLaps puts given laps into session and return any remaining laps that doesn't belong to this session.
func (s *Session) PutLaps(laps ...*Lap) (remainings []*Lap) {
	remainings = make([]*Lap, 0, len(laps))
	for j := range laps {
		lap := laps[j]
		if s.IsBelongToThisSession(lap.StartTime) {
			s.Laps = append(s.Laps, lap)
		} else {
			remainings = append(remainings, lap)
		}
	}
	return remainings
}

// PutRecords puts given records into session and return any remaining records that doesn't belong to this session.
func (s *Session) PutRecords(records ...*Record) (remainings []*Record) {
	remainings = make([]*Record, 0, len(records))
	for j := range records {
		rec := records[j]
		if s.IsBelongToThisSession(rec.Timestamp) {
			s.Records = append(s.Records, rec)
		} else {
			remainings = append(remainings, rec)
		}
	}
	return remainings
}

func (s *Session) Clone() *Session {
	ses := &Session{
		Timestamp:        s.Timestamp,
		StartTime:        s.StartTime,
		EndTime:          s.EndTime,
		Sport:            s.Sport,
		TotalMovingTime:  s.TotalMovingTime,
		TotalElapsedTime: s.TotalElapsedTime,
		TotalDistance:    s.TotalDistance,
		TotalAscent:      s.TotalAscent,
		TotalDescent:     s.TotalDescent,
		TotalCalories:    s.TotalCalories,
	}

	if s.AvgSpeed != nil {
		ses.AvgSpeed = kit.Ptr(*s.AvgSpeed)
	}
	if s.MaxSpeed != nil {
		ses.MaxSpeed = kit.Ptr(*s.MaxSpeed)
	}
	if s.AvgHeartRate != nil {
		ses.AvgHeartRate = kit.Ptr(*s.AvgHeartRate)
	}
	if s.MaxHeartRate != nil {
		ses.MaxHeartRate = kit.Ptr(*s.MaxHeartRate)
	}
	if s.AvgCadence != nil {
		ses.AvgCadence = kit.Ptr(*s.AvgCadence)
	}
	if s.MaxCadence != nil {
		ses.MaxCadence = kit.Ptr(*s.MaxCadence)
	}
	if s.AvgPower != nil {
		ses.AvgPower = kit.Ptr(*s.AvgPower)
	}
	if s.MaxPower != nil {
		ses.MaxPower = kit.Ptr(*s.MaxPower)
	}
	if s.AvgTemperature != nil {
		ses.AvgTemperature = kit.Ptr(*s.AvgTemperature)
	}
	if s.MaxTemperature != nil {
		ses.MaxTemperature = kit.Ptr(*s.MaxTemperature)
	}
	if s.AvgAltitude != nil {
		ses.AvgAltitude = kit.Ptr(*s.AvgAltitude)
	}
	if s.MaxAltitude != nil {
		ses.MaxAltitude = kit.Ptr(*s.MaxAltitude)
	}
	if s.AvgPace != nil {
		ses.AvgPace = kit.Ptr(*s.AvgPace)
	}
	if s.AvgElapsedPace != nil {
		ses.AvgElapsedPace = kit.Ptr(*s.AvgElapsedPace)
	}

	ses.Records = make([]*Record, len(s.Records))
	for i := range s.Records {
		ses.Records[i] = s.Records[i].Clone()
	}

	ses.Laps = make([]*Lap, len(s.Laps))
	for i := range s.Laps {
		ses.Laps[i] = s.Laps[i].Clone()
	}

	return ses
}

// CombineSession combines sesssion's values into targetSession.
// Every zero value in targetSession will be replaced with the corresponding value in session.
func CombineSession(targetSession, session *Session) {
	if targetSession == nil || session == nil {
		return
	}

	if targetSession.EndTime.IsZero() {
		targetSession.EndTime = session.EndTime
	}

	targetSession.TotalElapsedTime = kit.PickNonZeroValue(targetSession.TotalElapsedTime, session.TotalElapsedTime)
	targetSession.TotalMovingTime = kit.PickNonZeroValue(targetSession.TotalMovingTime, session.TotalMovingTime)
	targetSession.TotalDistance = kit.PickNonZeroValue(targetSession.TotalDistance, session.TotalDistance)
	targetSession.TotalCalories = kit.PickNonZeroValue(targetSession.TotalCalories, session.TotalCalories)
	targetSession.TotalAscent = kit.PickNonZeroValue(targetSession.TotalAscent, session.TotalAscent)
	targetSession.TotalDescent = kit.PickNonZeroValue(targetSession.TotalDescent, session.TotalDescent)
	targetSession.AvgSpeed = kit.PickNonZeroValuePtr(targetSession.AvgSpeed, session.AvgSpeed)
	targetSession.MaxSpeed = kit.PickNonZeroValuePtr(targetSession.MaxSpeed, session.MaxSpeed)
	targetSession.AvgHeartRate = kit.PickNonZeroValuePtr(targetSession.AvgHeartRate, session.AvgHeartRate)
	targetSession.MaxHeartRate = kit.PickNonZeroValuePtr(targetSession.MaxHeartRate, session.MaxHeartRate)
	targetSession.AvgCadence = kit.PickNonZeroValuePtr(targetSession.AvgCadence, session.AvgCadence)
	targetSession.MaxCadence = kit.PickNonZeroValuePtr(targetSession.MaxCadence, session.MaxCadence)
	targetSession.AvgPower = kit.PickNonZeroValuePtr(targetSession.AvgPower, session.AvgPower)
	targetSession.MaxPower = kit.PickNonZeroValuePtr(targetSession.MaxPower, session.MaxPower)
	targetSession.AvgTemperature = kit.PickNonZeroValuePtr(targetSession.AvgTemperature, session.AvgTemperature)
	targetSession.MaxTemperature = kit.PickNonZeroValuePtr(targetSession.MaxTemperature, session.MaxTemperature)
	targetSession.AvgAltitude = kit.PickNonZeroValuePtr(targetSession.AvgAltitude, session.AvgAltitude)
	targetSession.MaxAltitude = kit.PickNonZeroValuePtr(targetSession.MaxAltitude, session.MaxAltitude)
	targetSession.AvgPace = kit.PickNonZeroValuePtr(targetSession.AvgPace, session.AvgPace)
	targetSession.AvgElapsedPace = kit.PickNonZeroValuePtr(targetSession.AvgElapsedPace, session.AvgElapsedPace)
}

// AccumulateSession combines sesssion's values into targetSession.
// While MergeSession is picking the first non-zero value, AccumulateSession accumulate values between two sesions.
func AccumulateSession(targetSession, session *Session) {
	targetSession.TotalMovingTime += session.TotalMovingTime
	targetSession.TotalElapsedTime += session.TotalElapsedTime
	gap := session.StartTime.Sub(targetSession.EndTime).Seconds()
	targetSession.TotalElapsedTime += gap
	targetSession.TotalDistance += session.TotalDistance
	targetSession.TotalAscent += session.TotalAscent
	targetSession.TotalDescent += session.TotalDescent
	targetSession.TotalCalories += session.TotalCalories
	targetSession.EndTime = session.EndTime

	targetSession.AvgSpeed = kit.Avg(targetSession.AvgSpeed, session.AvgSpeed)
	targetSession.MaxSpeed = kit.Max(targetSession.MaxSpeed, session.MaxSpeed)
	targetSession.AvgHeartRate = kit.Avg(targetSession.AvgHeartRate, session.AvgHeartRate)
	targetSession.MaxHeartRate = kit.Max(targetSession.MaxHeartRate, session.MaxHeartRate)
	targetSession.AvgCadence = kit.Avg(targetSession.AvgCadence, session.AvgCadence)
	targetSession.MaxCadence = kit.Max(targetSession.MaxCadence, session.MaxCadence)
	targetSession.AvgPower = kit.Avg(targetSession.AvgPower, session.AvgPower)
	targetSession.MaxPower = kit.Max(targetSession.MaxPower, session.MaxPower)
	targetSession.AvgTemperature = kit.Avg(targetSession.AvgTemperature, session.AvgTemperature)
	targetSession.MaxTemperature = kit.Max(targetSession.MaxTemperature, session.MaxTemperature)
	targetSession.AvgAltitude = kit.Avg(targetSession.AvgAltitude, session.AvgAltitude)
	targetSession.MaxAltitude = kit.Max(targetSession.MaxAltitude, session.MaxAltitude)
	targetSession.AvgPace = kit.Avg(targetSession.AvgPace, session.AvgPace)
	targetSession.AvgElapsedPace = kit.Avg(targetSession.AvgElapsedPace, session.AvgElapsedPace)
}
