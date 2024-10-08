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

	"github.com/muktihari/fit/kit/scaleoffset"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/openivity/activity-service/geomath"
)

// Preprocessor is a preprocessor for improving activity's data such as smoothing the elevation and calculating slope/gradient.
type Preprocessor struct {
	options *options
}

type options struct {
	smoothingElevDistance  float64 // in meters
	calculateGradeDistance float64 // in meters
}

func defaultOptions() *options {
	return &options{
		smoothingElevDistance:  30,
		calculateGradeDistance: 100,
	}
}

type Option interface{ apply(*options) }

type fnApply func(*options)

func (f fnApply) apply(o *options) { f(o) }

func WithSmoothingElevationDistance(d float64) Option {
	return fnApply(func(o *options) {
		if d > 0 {
			o.smoothingElevDistance = d
		}
	})
}

func WithCalculateDistance(d float64) Option {
	return fnApply(func(o *options) {
		if d > 0 {
			o.calculateGradeDistance = d
		}
	})
}

// NewPreprocessor creates new preprocessor.
func NewPreprocessor(opts ...Option) *Preprocessor {
	options := defaultOptions()
	for i := range opts {
		opts[i].apply(options)
	}

	return &Preprocessor{options: options}
}

// AggregateByTimestamp aggregates fields with the same timestamp if any and return new slice of record.
// The FIT files produced by Strava is splitting values into multiple records with the same timestamp, so
// we think it's possible for other platforms/devices to produce similiar files.
func (p *Preprocessor) AggregateByTimestamp(records []Record) []Record {
	for i := 0; i < len(records); i++ {
		rec := &records[i]

		// Find forward
		for j := i + 1; j < len(records); j++ {
			next := records[j]
			if !rec.Timestamp.Equal(next.Timestamp) {
				records = append(records[:i+1], records[j:]...) // Remove aggregated records
				break
			}

			if rec.PositionLat == basetype.Sint32Invalid {
				rec.PositionLat = next.PositionLat
			}
			if rec.PositionLong == basetype.Sint32Invalid {
				rec.PositionLong = next.PositionLong
			}

			if rec.Altitude != basetype.Uint16Invalid && next.Altitude != basetype.Uint16Invalid {
				rec.Altitude = uint16((uint32(rec.Altitude) + uint32(next.Altitude)) / 2)
			} else if next.Altitude != basetype.Uint16Invalid {
				rec.Altitude = next.Altitude
			}

			if rec.EnhancedAltitude != basetype.Uint32Invalid && next.EnhancedAltitude != basetype.Uint32Invalid {
				rec.EnhancedAltitude = uint32((uint64(rec.EnhancedAltitude) + uint64(next.EnhancedAltitude)) / 2)
			} else if next.EnhancedAltitude != basetype.Uint32Invalid {
				rec.EnhancedAltitude = next.EnhancedAltitude
			}

			if rec.Cadence != basetype.Uint8Invalid && next.Cadence != basetype.Uint8Invalid {
				rec.Cadence = uint8((uint16(rec.Cadence) + uint16(next.Cadence)) / 2)
			} else if next.Cadence != basetype.Uint8Invalid {
				rec.Cadence = next.Cadence
			}

			if rec.Speed != basetype.Uint16Invalid && next.Speed != basetype.Uint16Invalid {
				rec.Speed = uint16((uint32(rec.Speed) + uint32(next.Speed)) / 2)
			} else if next.Speed != basetype.Uint16Invalid {
				rec.Speed = next.Speed
			}

			if rec.EnhancedSpeed != basetype.Uint32Invalid && next.Speed != basetype.Uint16Invalid {
				rec.EnhancedSpeed = uint32((uint64(rec.EnhancedSpeed) + uint64(next.EnhancedSpeed)) / 2)
			} else if next.EnhancedSpeed != basetype.Uint32Invalid {
				rec.EnhancedSpeed = next.EnhancedSpeed
			}

			if rec.Distance != basetype.Uint32Invalid && next.Distance != basetype.Uint32Invalid {
				rec.Distance = uint32((uint64(rec.Distance) + uint64(next.Distance)) / 2)
			} else if next.Distance != basetype.Uint32Invalid {
				rec.Distance = next.Distance
			}

			if rec.HeartRate != basetype.Uint8Invalid && next.HeartRate != basetype.Uint8Invalid {
				rec.HeartRate = uint8((uint16(rec.HeartRate) + uint16(next.HeartRate)) / 2)
			} else if next.HeartRate != basetype.Uint8Invalid {
				rec.HeartRate = next.HeartRate
			}

			if rec.Power != basetype.Uint16Invalid && next.Power != basetype.Uint16Invalid {
				rec.Power = uint16((uint32(rec.Power) + uint32(next.Power)) / 2)
			} else if next.Power != basetype.Uint16Invalid {
				rec.Power = next.Power
			}

			if rec.Temperature != basetype.Sint8Invalid && next.Temperature != basetype.Sint8Invalid {
				rec.Temperature = int8((int16(rec.Temperature) + int16(next.Temperature)) / 2)
			} else if next.Temperature != basetype.Sint8Invalid {
				rec.Temperature = next.Temperature
			}
		}
	}

	return records
}

// CalculateDistanceAndSpeed calculates distance from latitude and longitude and speed when those values are missing.
func (p *Preprocessor) CalculateDistanceAndSpeed(records []Record) {
	for i := 1; i < len(records); i++ {
		rec := &records[i]
		prev := records[i-1]

		// Calculate distance from two coordinates
		var pointDistance float64
		if rec.Distance == basetype.Uint32Invalid {
			if rec.PositionLat != basetype.Sint32Invalid && rec.PositionLong != basetype.Sint32Invalid &&
				prev.PositionLat != basetype.Sint32Invalid && prev.PositionLong != basetype.Sint32Invalid {

				var prevDist float64
				if prev.Distance != basetype.Uint32Invalid {
					prevDist = prev.DistanceScaled()
				}
				pointDistance = geomath.HaversineDistance(
					rec.PositionLatDegrees(),
					rec.PositionLongDegrees(),
					prev.PositionLatDegrees(),
					prev.PositionLongDegrees(),
				)
				rec.Distance = uint32(scaleoffset.Discard(prevDist+pointDistance, 100, 0))
			}
		} else if rec.Distance != basetype.Uint32Invalid && prev.Distance != basetype.Uint32Invalid {
			pointDistance = rec.DistanceScaled() - prev.DistanceScaled()
		}

		// Speed
		if rec.Speed == basetype.Uint16Invalid && pointDistance > 0 {
			elapsed := rec.Timestamp.Sub(prev.Timestamp).Seconds()
			if elapsed > 0 {
				speed := pointDistance / elapsed
				rec.Speed = uint16(scaleoffset.Discard(speed, 1000, 0))
			}
		}
	}
}

// SmoothingElevation smoothing elevation values using simple moving average.
func (p *Preprocessor) SmoothingElevation(records []Record) {
	// Copy altitude value
	for i := range records {
		rec := &records[i]
		switch {
		case rec.Altitude != basetype.Uint16Invalid:
			rec.SmoothedAltitude = rec.AltitudeScaled()
		case rec.EnhancedAltitude != basetype.Uint32Invalid:
			rec.SmoothedAltitude = rec.EnhancedAltitudeScaled()
		}
	}

	for i := range records {
		rec := &records[i]
		if rec.Distance == basetype.Uint32Invalid || math.IsNaN(rec.SmoothedAltitude) {
			continue
		}

		var sum, counter float64
		for j := i; j >= 0; j-- {
			prev := records[j]
			if prev.Distance == basetype.Uint32Invalid || math.IsNaN(prev.SmoothedAltitude) {
				continue
			}

			d := rec.DistanceScaled() - prev.DistanceScaled()
			if d > p.options.smoothingElevDistance {
				break
			}
			sum += prev.SmoothedAltitude
			counter++
		}

		smoothedAltitude := sum / counter
		rec.SmoothedAltitude = smoothedAltitude
	}
}

// CalculateGrade calculates grade percentage.
func (p *Preprocessor) CalculateGrade(records []Record) {
	for i := range records {
		rec := &records[i]

		altitude := rec.SmoothedAltitude
		if math.IsNaN(altitude) {
			altitude = rec.AltitudeScaled()
		}
		if math.IsNaN(altitude) {
			altitude = rec.EnhancedAltitudeScaled()
		}

		if rec.Distance == basetype.Uint32Invalid || math.IsNaN(altitude) {
			continue
		}

		var run, rise float64
		for j := i + 1; j < len(records); j++ {
			next := records[j]

			nextAltitude := next.SmoothedAltitude
			if math.IsNaN(nextAltitude) {
				nextAltitude = next.AltitudeScaled()
			}
			if math.IsNaN(nextAltitude) {
				nextAltitude = next.EnhancedAltitudeScaled()
			}

			if next.Distance == basetype.Uint32Invalid || math.IsNaN(nextAltitude) {
				continue
			}

			d := next.DistanceScaled() - rec.DistanceScaled()
			if d > p.options.calculateGradeDistance {
				break
			}
			rise = nextAltitude - altitude
			run = d
		}

		if rise == 0 || run == 0 {
			continue
		}

		grade := rise / run * 100

		rec.Grade = grade
	}
}

// CalculatePace calculates pace time/distance (distance in km)
func (p *Preprocessor) CalculatePace(sport typedef.Sport, records []Record) {
	for i := 1; i < len(records); i++ {
		rec := &records[i]
		prev := records[i-1]

		if rec.Distance == basetype.Uint32Invalid || rec.Timestamp.IsZero() ||
			prev.Distance == basetype.Uint32Invalid || prev.Timestamp.IsZero() {
			continue
		}

		speed := rec.SpeedScaled()
		if math.IsNaN(speed) {
			speed = rec.EnhancedSpeedScaled()
		}
		if !IsConsideredMoving(sport, speed) {
			continue
		}

		if math.IsNaN(speed) {
			pointDistance := rec.DistanceScaled() - prev.DistanceScaled()
			elapsed := rec.Timestamp.Sub(prev.Timestamp).Seconds()
			pointDistanceInKM := pointDistance / 1000
			if pointDistanceInKM == 0 {
				continue
			}
			rec.Pace = elapsed / pointDistanceInKM
		} else {
			speedkph := speed * 3.6
			rec.Pace = (1 / speedkph) * 60 * 60
		}
	}
}
