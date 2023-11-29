package gpx

import (
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/gpx/schema"
)

func NewRecord(trkpt *schema.Waypoint) *activity.Record {
	rec := new(activity.Record)

	rec.Timestamp = trkpt.Time
	rec.PositionLat = trkpt.Lat
	rec.PositionLong = trkpt.Lon
	rec.Altitude = trkpt.Ele

	ext := trkpt.TrackPointExtension
	if ext != nil {
		rec.Distance = ext.Distance
		rec.Cadence = ext.Cadence
		rec.HeartRate = ext.HeartRate
		rec.Power = ext.Power
		rec.Temperature = ext.Temperature
	}

	return rec
}
