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

package activity

import (
	"math"
	"time"

	"github.com/muktihari/fit/profile/typedef"
)

const (
	ToleranceMovingSpeedSlowMovingSport  = 0.1388 // = 0.5km/h
	ToleranceMovingSpeedRunningLikeSport = 0.7916 // = 2.85 km/h
	ToleranceMovingSpeedCyclingLikeSport = 1.41   // = 5.07 km/h
)

// IsConsideredMoving check whether given speed in a given sport is considered moving.
func IsConsideredMoving(sport typedef.Sport, speed float64) bool {
	if math.IsNaN(speed) {
		return false
	}
	return speed > ToleranceMovingSpeed(sport)
}

// ToleranceMovingSpeed returns the tolerance moving speed of given record.
func ToleranceMovingSpeed(sport typedef.Sport) float64 {
	switch sport {
	case typedef.SportRunning:
		return ToleranceMovingSpeedRunningLikeSport
	case typedef.SportCycling:
		return ToleranceMovingSpeedCyclingLikeSport
	default:
		// Generic: since we don't know the specific sport, let's assume the slowest possible speed as our tolerance.
		// Included in this category: SportHiking, SportWalking, SportSwimming, SportTransition
		return ToleranceMovingSpeedSlowMovingSport
	}
}

// HasPace check whether given sport has pace for analytic.
func HasPace(sport typedef.Sport) bool {
	switch sport {
	case typedef.SportGeneric:
		// Since we don't know the specific sport, let's calculate the pace for now and allow the user to decide for themselves.
		// It's better to provide additional information than to have none. (TBD)
		fallthrough
	case typedef.SportHiking, typedef.SportWalking, typedef.SportRunning, typedef.SportSwimming, typedef.SportTransition:
		return true
	default:
		return false
	}
}

func isBelong(timestamp, startTime, endTime time.Time) bool {
	if timestamp.Equal(startTime) {
		return true
	}
	if endTime.IsZero() { // Last Lap or Session has no EndTime
		return timestamp.After(startTime)
	}
	if timestamp.Equal(endTime) {
		return true
	}
	return timestamp.After(startTime) && timestamp.Before(endTime)
}
