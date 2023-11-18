package preprocessor

import (
	"time"

	"github.com/muktihari/openactivity-fit/activity"
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

// SmoothingElev smoothing elevation values using simple moving average.
func (p *Preprocessor) SmoothingElev(records []*activity.Record) {
	for i := range records {
		rec := records[i]
		if rec.Distance == nil || rec.Altitude == nil {
			continue
		}

		var sum, counter float64
		for j := i; j >= 0; j-- {
			prev := records[j]
			if prev.Distance == nil || prev.Altitude == nil {
				continue
			}

			d := *rec.Distance - *prev.Distance
			if d > p.options.smoothingElevDistance {
				break
			}
			sum += *prev.Altitude
			counter++
		}

		altitude := sum / counter
		rec.Altitude = &altitude
	}
}

// CalculateGrade calculates grade percentage.
func (p *Preprocessor) CalculateGrade(records []*activity.Record) {
	for i := range records {
		rec := records[i]
		if rec.Distance == nil || rec.Altitude == nil {
			continue
		}

		var run, rise float64
		for j := i + 1; j < len(records); j++ {
			next := records[j]
			if next.Distance == nil || next.Altitude == nil {
				continue
			}

			d := *next.Distance - *rec.Distance
			if d > p.options.calculateGradeDistance {
				break
			}
			rise = *next.Altitude - *rec.Altitude
			run = d
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

		if rec.Distance == nil || rec.Timestamp == (time.Time{}) ||
			prev.Distance == nil || prev.Timestamp == (time.Time{}) {
			continue
		}

		if !activity.IsConsideredMoving(sport, rec.Speed) {
			continue
		}

		pointDistance := ((*rec.Distance) - (*prev.Distance))
		elapsed := rec.Timestamp.Sub(prev.Timestamp).Seconds()
		pointDistanceInKM := pointDistance / 1000
		if pointDistanceInKM == 0 {
			continue
		}

		rec.Pace = kit.Ptr(elapsed / pointDistanceInKM)
	}
}
