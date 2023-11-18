package activity

import (
	"time"
)

type Lap struct {
	Timestamp        time.Time
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
}

func (l *Lap) ToMap() map[string]any {
	m := map[string]any{}

	if l.Timestamp != (time.Time{}) {
		m["timestamp"] = l.Timestamp.Format(time.RFC3339)
	}

	m["totalMovingTime"] = l.TotalMovingTime
	m["totalElapsedTime"] = l.TotalElapsedTime
	m["totalDistance"] = l.TotalDistance
	m["totalAscent"] = l.TotalAscent
	m["totalDescent"] = l.TotalDescent
	m["totalCalories"] = l.TotalCalories

	if l.AvgSpeed != nil {
		m["avgSpeed"] = *l.AvgSpeed
	}
	if l.MaxSpeed != nil {
		m["maxSpeed"] = *l.MaxSpeed
	}
	if l.AvgHeartRate != nil {
		m["avgHeartRate"] = *l.AvgHeartRate
	}
	if l.MaxHeartRate != nil {
		m["maxHeartRate"] = *l.MaxHeartRate
	}
	if l.AvgCadence != nil {
		m["avgCadence"] = *l.AvgCadence
	}
	if l.MaxCadence != nil {
		m["maxCadence"] = *l.MaxCadence
	}
	if l.AvgPower != nil {
		m["avgPower"] = *l.AvgPower
	}
	if l.MaxPower != nil {
		m["maxPower"] = *l.MaxPower
	}
	if l.AvgTemperature != nil {
		m["avgTemperature"] = *l.AvgTemperature
	}
	if l.MaxTemperature != nil {
		m["maxTemperature"] = *l.MaxTemperature
	}
	if l.AvgAltitude != nil {
		m["avgAltitude"] = *l.AvgAltitude
	}
	if l.MaxAltitude != nil {
		m["maxAltitude"] = *l.MaxAltitude
	}
	if l.AvgPace != nil {
		m["avgPace"] = *l.AvgPace
	}
	if l.AvgElapsedPace != nil {
		m["maxPace"] = *l.AvgElapsedPace
	}

	return m
}
