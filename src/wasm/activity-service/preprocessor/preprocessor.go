package preprocessor

import (
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/geomath"
	"github.com/muktihari/openactivity-fit/kit"
)

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

func New(opts ...Option) *Preprocessor {
	options := defaultOptions()
	for i := range opts {
		opts[i].apply(options)
	}

	return &Preprocessor{options: options}
}

// AggregateByTimestamp aggregates fields with the same timestamp if any and return new slice of record.
// The Fit file produced by Strava is splitting values into multiple records with the same timestamp, so
// I think it's possible for other platforms/devices to produce similiar files.
func (p *Preprocessor) AggregateByTimestamp(records []*activity.Record) []*activity.Record {
	newRecords := make([]*activity.Record, 0, len(records))

	for i := 0; i < len(records); i++ {
		rec := records[i]

		candidates := make([]*activity.Record, 0)
		for j := i + 1; j < len(records); j++ {
			next := records[j]
			if !rec.Timestamp.Equal(next.Timestamp) {
				i = j - 1
				break
			}
			candidates = append(candidates, next)
		}

		for j := range candidates {
			can := candidates[j]

			if rec.PositionLat == nil {
				rec.PositionLat = can.PositionLat
			}
			if rec.PositionLong == nil {
				rec.PositionLong = can.PositionLong
			}

			rec.Altitude = kit.Avg(rec.Altitude, can.Altitude)
			rec.Cadence = kit.Avg(rec.Cadence, can.Cadence)
			rec.Speed = kit.Avg(rec.Speed, can.Speed)
			rec.Distance = kit.Avg(rec.Distance, can.Distance)
			rec.HeartRate = kit.Avg(rec.HeartRate, can.HeartRate)
			rec.Power = kit.Avg(rec.Power, can.Power)
			rec.Temperature = kit.Avg(rec.Temperature, can.Temperature)
		}

		newRecords = append(newRecords, rec)
	}

	return newRecords
}

// CalculateDistanceAndSpeed calculates distance from latitude and longitude and speed when those values are missing.
func (p *Preprocessor) CalculateDistanceAndSpeed(records []*activity.Record) {
	for i := 1; i < len(records); i++ {
		rec := records[i]
		prev := records[i-1]

		// Calculate distance from two coordinates
		var pointDistance float64
		if rec.Distance == nil {
			if rec.PositionLat != nil && rec.PositionLong != nil &&
				prev.PositionLat != nil && prev.PositionLong != nil {

				var prevDist float64
				if prev.Distance != nil {
					prevDist = *prev.Distance
				}

				pointDistance = geomath.VincentyDistance(
					*rec.PositionLat,
					*rec.PositionLong,
					*prev.PositionLat,
					*prev.PositionLong,
				)

				rec.Distance = kit.Ptr(prevDist + pointDistance)
			}
		} else if rec.Distance != nil && prev.Distance != nil {
			pointDistance = *rec.Distance - *prev.Distance
		}

		// Speed
		if rec.Speed == nil && pointDistance > 0 {
			elapsed := rec.Timestamp.Sub(prev.Timestamp).Seconds()
			if elapsed > 0 {
				speed := pointDistance / elapsed
				rec.Speed = &speed
			}
		}
	}
}

// SmoothingElev smoothing elevation values using simple moving average.
func (p *Preprocessor) SmoothingElev(records []*activity.Record) {
	// Copy altitude value
	for i := range records {
		rec := records[i]
		if rec.Altitude != nil {
			rec.SmoothedAltitude = kit.Ptr(*rec.Altitude)
		}
	}

	for i := range records {
		rec := records[i]
		if rec.Distance == nil || rec.SmoothedAltitude == nil {
			continue
		}

		var sum, counter float64
		for j := i; j >= 0; j-- {
			prev := records[j]
			if prev.Distance == nil || prev.SmoothedAltitude == nil {
				continue
			}

			d := *rec.Distance - *prev.Distance
			if d > p.options.smoothingElevDistance {
				break
			}
			sum += *prev.SmoothedAltitude
			counter++
		}

		smoothedAltitude := sum / counter
		rec.SmoothedAltitude = &smoothedAltitude
	}
}

// CalculateGrade calculates grade percentage.
func (p *Preprocessor) CalculateGrade(records []*activity.Record) {
	for i := range records {
		rec := records[i]

		altitude := rec.SmoothedAltitude
		if altitude == nil {
			altitude = rec.Altitude
		}

		if rec.Distance == nil || altitude == nil {
			continue
		}

		var run, rise float64
		for j := i + 1; j < len(records); j++ {
			next := records[j]

			nextAltitude := next.SmoothedAltitude
			if nextAltitude == nil {
				nextAltitude = next.Altitude
			}

			if next.Distance == nil || nextAltitude == nil {
				continue
			}

			d := *next.Distance - *rec.Distance
			if d > p.options.calculateGradeDistance {
				break
			}
			rise = *nextAltitude - *altitude
			run = d
		}

		if rise == 0 || run == 0 {
			continue
		}

		grade := rise / run * 100

		rec.Grade = &grade
	}
}

// CalculatePace calculates pace time/distance (distance in km)
func (p *Preprocessor) CalculatePace(sport string, records []*activity.Record) {
	for i := 1; i < len(records); i++ {
		rec := records[i]
		prev := records[i-1]

		if rec.Distance == nil || rec.Timestamp.IsZero() ||
			prev.Distance == nil || prev.Timestamp.IsZero() {
			continue
		}

		if !activity.IsConsideredMoving(sport, rec.Speed) {
			continue
		}

		if rec.Speed == nil {
			pointDistance := ((*rec.Distance) - (*prev.Distance))
			elapsed := rec.Timestamp.Sub(prev.Timestamp).Seconds()
			pointDistanceInKM := pointDistance / 1000
			if pointDistanceInKM == 0 {
				continue
			}
			rec.Pace = kit.Ptr(elapsed / pointDistanceInKM)
		} else {
			speedkph := *rec.Speed * 3.6
			rec.Pace = kit.Ptr((1 / speedkph) * 60 * 60)
		}
	}
}
