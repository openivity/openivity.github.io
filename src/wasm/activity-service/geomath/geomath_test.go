// Copyright (C) 2024 Openivity

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

package geomath_test

import (
	"math"
	"testing"

	"github.com/openivity/activity-service/geomath"
)

func TestHaversineDistance(t *testing.T) {
	tt := []struct {
		lat1, lon1 float64
		lat2, lon2 float64
		expected   float64
	}{
		{
			lat1: -7.202760921791196, lon1: 109.93464292958379,
			lat2: -7.202771985903382, lon2: 109.93463496677577,
			expected: 1.51,
		},
		{
			lat1: -7.202783972024918, lon1: 109.93463094346225,
			lat2: -7.202786989510059, lon2: 109.93462993763387,
			expected: 0.35,
		},
	}

	for _, tc := range tt {
		distance := geomath.HaversineDistance(tc.lat1, tc.lon1, tc.lat2, tc.lon2)
		distance = math.Round(distance*100) / 100 // let's two decimals precision
		if distance != tc.expected {
			t.Fatalf("expected: %g, got: %g", tc.expected, distance)
		}
	}
}
