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

package geomath

import (
	"math"
)

// HaversineDistance returns distance meters between two coordinates calculated using Haversine formula.
// The Haversine formula can result in an error of up to 0.5% since the Earth is not even a sphere — it’s an oblate ellipsoid.
//
// ref: https://en.wikipedia.org/wiki/Haversine_formula
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const radius = 6371.0

	lat1, lon1 = degreesToRadians(lat1), degreesToRadians(lon1)
	lat2, lon2 = degreesToRadians(lat2), degreesToRadians(lon2)

	// Haversine formula
	a := math.Pow(math.Sin((lat2-lat1)/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin((lon2-lon1)/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate the distance (in km)
	distance := radius * c

	return distance * 1000 // in meters
}

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
