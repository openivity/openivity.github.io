package gpx

import (
	"math"

	"github.com/muktihari/openactivity-fit/accumulator"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/kit"
)

func NewLap(records []*activity.Record, sport string) *activity.Lap {
	if len(records) == 0 {
		return nil
	}

	lap := &activity.Lap{
		Timestamp: records[0].Timestamp,
	}

	var (
		distanceAccumu    = new(accumulator.Accumulator[float64])
		speedAccumu       = new(accumulator.Accumulator[float64])
		altitudeAccumu    = new(accumulator.Accumulator[float64])
		cadenceAccumu     = new(accumulator.Accumulator[uint8])
		heartRateAccumu   = new(accumulator.Accumulator[uint8])
		powerAccumu       = new(accumulator.Accumulator[uint16])
		temperatureAccumu = new(accumulator.Accumulator[int8])
	)

	for i := 0; i < len(records); i++ {
		rec := records[i]

		distanceAccumu.Collect(rec.Distance)
		speedAccumu.Collect(rec.Speed)
		altitudeAccumu.Collect(rec.Altitude)
		cadenceAccumu.Collect(rec.Cadence)
		heartRateAccumu.Collect(rec.HeartRate)
		powerAccumu.Collect(rec.Power)
		temperatureAccumu.Collect(rec.Temperature)

		if i == 0 {
			continue
		}

		prev := records[i-1]

		// Calculate Total Elapsed and Total Moving Time
		if rec.Distance != nil && prev.Distance != nil {
			timeDiff := rec.Timestamp.Sub(prev.Timestamp).Seconds()
			lap.TotalElapsedTime += timeDiff

			if activity.IsConsideredMoving(sport, rec.Speed) {
				lap.TotalMovingTime += timeDiff
			}
		}

		// Calculate Total Ascent and Total Descent
		if rec.Altitude != nil && prev.Altitude != nil {
			delta := *rec.Altitude - *prev.Altitude
			if delta > 0 {
				lap.TotalAscent += uint16(delta)
			} else {
				lap.TotalDescent += uint16(math.Abs(delta))
			}
		}

	}

	if value := distanceAccumu.Max(); value != nil {
		lap.TotalDistance = *value
	}
	lap.AvgSpeed = speedAccumu.Avg()
	lap.MaxSpeed = speedAccumu.Max()
	lap.AvgAltitude = altitudeAccumu.Avg()
	lap.MaxAltitude = altitudeAccumu.Max()
	lap.AvgCadence = cadenceAccumu.Avg()
	lap.MaxCadence = cadenceAccumu.Max()
	lap.AvgHeartRate = heartRateAccumu.Avg()
	lap.MaxHeartRate = heartRateAccumu.Max()
	lap.AvgPower = powerAccumu.Avg()
	lap.MaxPower = powerAccumu.Max()
	lap.AvgTemperature = temperatureAccumu.Avg()
	lap.MaxTemperature = temperatureAccumu.Max()

	if activity.HasPace(sport) {
		lap.AvgPace = kit.Ptr(lap.TotalMovingTime / (lap.TotalDistance / 1000))
		lap.AvgElapsedPace = kit.Ptr(lap.TotalElapsedTime / (lap.TotalDistance / 1000))
	}

	return lap
}
