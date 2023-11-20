package gpx

import (
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/gpx/schema"
	"github.com/muktihari/openactivity-fit/geomath"
	"github.com/muktihari/openactivity-fit/kit"
)

func NewRecord(trkpt *schema.Waypoint, prevRec *activity.Record) *activity.Record {
	rec := new(activity.Record)

	rec.Timestamp = trkpt.Time
	rec.PositionLat = trkpt.Lat
	rec.PositionLong = trkpt.Lon
	rec.Altitude = trkpt.Ele

	if prevRec == nil {
		rec.Distance = kit.Ptr(0.0)
	}

	var pointDistance float64 // distance between two coordinates
	ext := trkpt.TrackPointExtension
	if ext != nil {
		if prevRec != nil && prevRec.Distance != nil && ext.Distance != nil {
			pointDistance = *ext.Distance - *prevRec.Distance
		}
		if ext.Distance != nil {
			rec.Distance = ext.Distance
		}
		rec.Cadence = ext.Cadence
		rec.HeartRate = ext.HeartRate
		rec.Power = ext.Power
		rec.Temperature = ext.Temperature
	}

	if rec.Distance == nil && prevRec != nil &&
		prevRec.PositionLat != nil && prevRec.PositionLong != nil &&
		rec.PositionLat != nil && rec.PositionLong != nil {

		// We got no distance from extension, let's calculate using coordinates.
		prevDist := 0.0
		if prevRec.Distance != nil {
			prevDist = *prevRec.Distance
		}

		pointDistance = geomath.VincentyDistance(
			*prevRec.PositionLat,
			*prevRec.PositionLong,
			*rec.PositionLat,
			*rec.PositionLong,
		)

		rec.Distance = kit.Ptr(prevDist + pointDistance)
	}

	if prevRec != nil && pointDistance > 0 {
		elapsed := rec.Timestamp.Sub(prevRec.Timestamp).Seconds()
		if elapsed > 0 {
			speed := pointDistance / elapsed
			rec.Speed = &speed
		}
	}

	return rec
}
