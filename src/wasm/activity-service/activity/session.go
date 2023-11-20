package activity

import "time"

type Session struct {
	Timestamp        time.Time
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
}

func (s *Session) ToMap() map[string]any {
	m := map[string]any{}

	if s.Timestamp != (time.Time{}) {
		m["timestamp"] = s.Timestamp.Format(time.RFC3339)
	}
	if s.Sport != "" {
		m["sport"] = s.Sport
	}
	m["totalMovingTime"] = s.TotalMovingTime
	m["totalElapsedTime"] = s.TotalElapsedTime
	m["totalDistance"] = s.TotalDistance
	m["totalAscent"] = s.TotalAscent
	m["totalDescent"] = s.TotalDescent
	m["totalCalories"] = s.TotalCalories

	if s.AvgSpeed != nil {
		m["avgSpeed"] = *s.AvgSpeed
	}
	if s.MaxSpeed != nil {
		m["maxSpeed"] = *s.MaxSpeed
	}
	if s.AvgHeartRate != nil {
		m["avgHeartRate"] = *s.AvgHeartRate
	}
	if s.MaxHeartRate != nil {
		m["maxHeartRate"] = *s.MaxHeartRate
	}
	if s.AvgCadence != nil {
		m["avgCadence"] = *s.AvgCadence
	}
	if s.MaxCadence != nil {
		m["maxCadence"] = *s.MaxCadence
	}
	if s.AvgPower != nil {
		m["avgPower"] = *s.AvgPower
	}
	if s.MaxPower != nil {
		m["maxPower"] = *s.MaxPower
	}
	if s.AvgTemperature != nil {
		m["avgTemperature"] = *s.AvgTemperature
	}
	if s.MaxTemperature != nil {
		m["maxTemperature"] = *s.MaxTemperature
	}
	if s.AvgAltitude != nil {
		m["avgAltitude"] = *s.AvgAltitude
	}
	if s.MaxAltitude != nil {
		m["maxAltitude"] = *s.MaxAltitude
	}
	if s.AvgPace != nil {
		m["avgPace"] = *s.AvgPace
	}
	if s.AvgElapsedPace != nil {
		m["avgElapsedPace"] = *s.AvgElapsedPace
	}

	return m
}
