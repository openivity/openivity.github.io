package tcx

import (
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/tcx/schema"
)

func NewRecord(trackpoint *schema.Trackpoint) *activity.Record {
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

	if trackpoint.Extensions != nil {
		if trackpoint.Extensions.Speed != nil {
			rec.Speed = trackpoint.Extensions.Speed
		}
	}

	return rec
}
