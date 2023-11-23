package tcx

import (
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/tcx/schema"
)

func NewRecord(trackpoint *schema.Trackpoint, prevRec *activity.Record) *activity.Record {
	if trackpoint == nil {
		return nil
	}

	rec := &activity.Record{
		Timestamp: trackpoint.Time,
		Distance:  trackpoint.DistanceMeters,
		Altitude:  trackpoint.AltitudeMeters,
		Cadence:   trackpoint.Cadence,
		HeartRate: trackpoint.HeartRateBpm,
	}

	if trackpoint.Position != nil {
		rec.PositionLat = &trackpoint.Position.LatitudeDegrees
		rec.PositionLong = &trackpoint.Position.LongitudeDegrees
	}

	var pointDistance float64
	if prevRec != nil && prevRec.Distance != nil && rec.Distance != nil {
		pointDistance = *rec.Distance - *prevRec.Distance
	}

	if pointDistance != 0 {
		elapsed := rec.Timestamp.Sub(prevRec.Timestamp).Seconds()
		if elapsed > 0 {
			speed := pointDistance / elapsed
			rec.Speed = &speed
		}
	}

	return rec
}
