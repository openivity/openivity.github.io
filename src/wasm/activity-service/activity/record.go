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
)

// Record is a workout record. It wraps FIT SDK's mesgdef.Record as its base.
type Record struct {
	*mesgdef.Record

	SmoothedAltitude float64 // Smoothed Altitude using our preprocessor algorithm.
	Pace             float64
	Grade            float64
}

// CreateRecord creates new record.
func CreateRecord(rec *mesgdef.Record) Record {
	if rec == nil {
		rec = mesgdef.NewRecord(nil)
	}
	return Record{
		Record:           rec,
		SmoothedAltitude: math.NaN(),
		Pace:             math.NaN(),
		Grade:            math.NaN(),
	}
}

// MarshalAppendJSON appends the JSON format encoding of Record to b, returning the result.
func (r *Record) MarshalAppendJSON(b []byte) []byte {
	b = append(b, '{')

	if !r.Timestamp.IsZero() {
		b = append(b, `"timestamp":`...)
		b = strconv.AppendQuote(b, r.Timestamp.Format(time.RFC3339))
		b = append(b, ',')
	}
	if r.PositionLat != basetype.Sint32Invalid {
		b = append(b, `"positionLat":`...)
		b = strconv.AppendFloat(b, r.PositionLatDegrees(), 'g', -1, 64)
		b = append(b, ',')
	}
	if r.PositionLong != basetype.Sint32Invalid {
		b = append(b, `"positionLong":`...)
		b = strconv.AppendFloat(b, r.PositionLongDegrees(), 'g', -1, 64)
		b = append(b, ',')
	}
	if r.Distance != basetype.Uint32Invalid {
		b = append(b, `"distance":`...)
		b = strconv.AppendFloat(b, r.DistanceScaled(), 'g', -1, 64)
		b = append(b, ',')
	}

	altitude := r.SmoothedAltitude
	if math.IsNaN(altitude) {
		altitude = r.AltitudeScaled()
	}
	if math.IsNaN(altitude) {
		altitude = r.EnhancedAltitudeScaled()
	}
	if !math.IsNaN(altitude) {
		b = append(b, `"altitude":`...)
		b = strconv.AppendFloat(b, altitude, 'g', -1, 64)
		b = append(b, ',')
	}

	if r.HeartRate != basetype.Uint8Invalid {
		b = append(b, `"heartRate":`...)
		b = strconv.AppendUint(b, uint64(r.HeartRate), 10)
		b = append(b, ',')
	}
	if r.Cadence != basetype.Uint8Invalid {
		b = append(b, `"cadence":`...)
		b = strconv.AppendUint(b, uint64(r.Cadence), 10)
		b = append(b, ',')
	}

	speed := r.SpeedScaled()
	if math.IsNaN(speed) {
		speed = r.EnhancedAltitudeScaled()
	}
	if !math.IsNaN(speed) {
		b = append(b, `"speed":`...)
		b = strconv.AppendFloat(b, speed, 'g', -1, 64)
		b = append(b, ',')
	}

	if r.Power != basetype.Uint16Invalid {
		b = append(b, `"power":`...)
		b = strconv.AppendUint(b, uint64(r.Power), 10)
		b = append(b, ',')
	}
	if r.Temperature != basetype.Sint8Invalid {
		b = append(b, `"temperature":`...)
		b = strconv.AppendInt(b, int64(r.Temperature), 10)
		b = append(b, ',')
	}
	if !math.IsNaN(r.Pace) {
		b = append(b, `"pace":`...)
		b = strconv.AppendFloat(b, r.Pace, 'g', -1, 64)
		b = append(b, ',')
	}
	if !math.IsNaN(r.Grade) {
		b = append(b, `"grade":`...)
		b = strconv.AppendFloat(b, r.Grade, 'g', -1, 64)
	}

	if b[len(b)-1] == '{' {
		return b[:len(b)-1]
	}

	if b[len(b)-1] == ',' {
		b = b[:len(b)-1]
	}

	return append(b, '}')
}
