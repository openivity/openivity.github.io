package gpx

import (
	"github.com/muktihari/openactivity-fit/accumulator"
	"github.com/muktihari/openactivity-fit/activity"
)

func NewSession(laps []*activity.Lap, sport string) *activity.Session {
	if len(laps) == 0 {
		return nil
	}

	ses := &activity.Session{
		Timestamp: laps[0].Timestamp,
		Sport:     sport,
	}

	var (
		totalElapsedTimeAccumu = new(accumulator.Accumulator[float64])
		totalMovingTimeAccumu  = new(accumulator.Accumulator[float64])
		distanceAccumu         = new(accumulator.Accumulator[float64])
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
		distanceAccumu.Collect(&lap.TotalDistance)
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
	if value := distanceAccumu.Max(); value != nil {
		ses.TotalDistance = *value
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

	if activity.HasPace(sport) {
		var (
			paceAvgAccumu        = new(accumulator.Accumulator[float64])
			paceAvgElapsedAccumu = new(accumulator.Accumulator[float64])
		)
		for i := range laps {
			lap := laps[i]

			paceAvgAccumu.Collect(lap.AvgPace)
			paceAvgElapsedAccumu.Collect(lap.AvgElapsedPace)
		}
		ses.AvgPace = paceAvgAccumu.Avg()
		ses.AvgElapsedPace = paceAvgElapsedAccumu.Avg()
	}

	return ses
}
