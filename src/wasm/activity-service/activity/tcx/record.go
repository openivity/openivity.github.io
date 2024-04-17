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
