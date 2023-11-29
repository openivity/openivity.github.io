package activity

import (
	"time"

	"github.com/muktihari/openactivity-fit/kit"
)

type Record struct {
	Timestamp        time.Time
	PositionLat      *float64
	PositionLong     *float64
	Distance         *float64
	Altitude         *float64 // Original Altitude from file.
	SmoothedAltitude *float64 // Smoothed Altitude using our preprocessor algorithm.
	HeartRate        *uint8
	Cadence          *uint8
	Speed            *float64
	Power            *uint16
	Temperature      *int8
	Pace             *float64
	Grade            *float64
}

func (r *Record) ToMap() map[string]any {
	m := map[string]any{}

	if !r.Timestamp.IsZero() {
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
	if r.SmoothedAltitude != nil {
		m["altitude"] = *r.SmoothedAltitude // for better data representation.
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

func (r *Record) Clone() *Record {
	rec := &Record{
		Timestamp: r.Timestamp,
	}

	if r.PositionLat != nil {
		rec.PositionLat = kit.Ptr(*r.PositionLat)
	}
	if r.PositionLong != nil {
		rec.PositionLong = kit.Ptr(*r.PositionLong)
	}
	if r.Distance != nil {
		rec.Distance = kit.Ptr(*r.Distance)
	}
	if r.Altitude != nil {
		rec.Altitude = kit.Ptr(*r.Altitude)
	}
	if r.SmoothedAltitude != nil {
		rec.SmoothedAltitude = kit.Ptr(*r.SmoothedAltitude)
	}
	if r.HeartRate != nil {
		rec.HeartRate = kit.Ptr(*r.HeartRate)
	}
	if r.Cadence != nil {
		rec.Cadence = kit.Ptr(*r.Cadence)
	}
	if r.Speed != nil {
		rec.Speed = kit.Ptr(*r.Speed)
	}
	if r.Power != nil {
		rec.Power = kit.Ptr(*r.Power)
	}
	if r.Temperature != nil {
		rec.Temperature = kit.Ptr(*r.Temperature)
	}
	if r.Pace != nil {
		rec.Pace = kit.Ptr(*r.Pace)
	}
	if r.Grade != nil {
		rec.Grade = kit.Ptr(*r.Grade)
	}

	return rec
}
