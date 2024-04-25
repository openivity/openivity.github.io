// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
