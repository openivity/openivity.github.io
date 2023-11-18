package activity

import (
	"time"
)

type Record struct {
	Timestamp    time.Time
	PositionLat  *float64
	PositionLong *float64
	Distance     *float64
	Altitude     *float64
	HeartRate    *uint8
	Cadence      *uint8
	Speed        *float64
	Power        *uint16
	Temperature  *int8
	Pace         *float64
	Grade        *float64
}

func (r *Record) ToMap() map[string]any {
	m := map[string]any{}

	if r.Timestamp != (time.Time{}) {
		m["timestamp"] = r.Timestamp.Format(time.RFC3339)
	}
	if r.PositionLat != nil {
		m["positionLat"] = *r.PositionLat
	}
	if r.PositionLong != nil {
		m["positionLong"] = *r.PositionLong
	}
	if r.Distance != nil {
		m["distance"] = *r.Distance
	}
	if r.Altitude != nil {
		m["altitude"] = *r.Altitude
	}
	if r.HeartRate != nil {
		m["heartRate"] = *r.HeartRate
	}
	if r.Cadence != nil {
		m["cadence"] = *r.Cadence
	}
	if r.Speed != nil {
		m["speed"] = *r.Speed
	}
	if r.Power != nil {
		m["power"] = *r.Power
	}
	if r.Temperature != nil {
		m["temperature"] = *r.Temperature
	}
	if r.Pace != nil {
		m["pace"] = *r.Pace
	}
	if r.Grade != nil {
		m["grade"] = *r.Grade
	}

	return m
}
