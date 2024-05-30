package activity

import (
	"math"

	"github.com/muktihari/fit/kit/scaleoffset"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
)

// AvgMaxSpeed returns average and max speed of the given records.
func AvgMaxSpeed(records []Record) (avgSpeed, maxSpeed uint16) {
	avgSpeed, maxSpeed = basetype.Uint16Invalid, basetype.Uint16Invalid
	if len(records) == 0 {
		return
	}
	var avg, count uint64
	for i := range records {
		rec := &records[i]
		if rec.Speed != basetype.Uint16Invalid {
			avg += uint64(rec.Speed)
			count++
			if maxSpeed == basetype.Uint16Invalid || rec.Speed > maxSpeed {
				maxSpeed = rec.Speed
			}
		}
	}
	if count == 0 {
		return
	}
	return uint16(avg / count), maxSpeed
}

// TotalAscentAndDescent calculate TotalAscent and TotalDescent from records.
func TotalAscentAndDescent(records []Record) (totalAscent, totalDescent uint16) {
	var ascent, descent float64
	var hasAltitude bool
	for i := 0; i < len(records)-1; i++ {
		rec := &records[i]
		if math.IsNaN(rec.SmoothedAltitude) {
			continue
		}
		hasAltitude = true
		// Find next non-nil altitude
		for j := i + 1; j < len(records); j++ {
			next := &records[j]
			if !math.IsNaN(rec.SmoothedAltitude) {
				delta := next.SmoothedAltitude - rec.SmoothedAltitude
				if delta > 0 {
					ascent += delta
				} else {
					descent += math.Abs(delta)
				}
				i = j - 1 // move cursor
				break
			}
		}
	}
	if !hasAltitude {
		return basetype.Uint16Invalid, basetype.Uint16Invalid
	}
	return uint16(math.Round(ascent)), uint16(math.Round(descent))
}

// TotalMovingTime calculates TotalMovingTime from records.
func TotalMovingTime(records []Record, sport typedef.Sport) (totalMovingTime uint32) {
	totalMovingTime = basetype.Uint32Invalid
	for i := 0; i < len(records); i++ {
		rec := &records[i]
		if rec.Timestamp.IsZero() {
			continue
		}
		// Find next non-zero timestamp
		for j := i + 1; j < len(records); j++ {
			next := &records[j]
			if !next.Timestamp.IsZero() {
				delta := next.Timestamp.Sub(rec.Timestamp).Seconds()
				if IsConsideredMoving(sport, rec.SpeedScaled()) {
					totalMovingTime += uint32(scaleoffset.Discard(delta, 1000, 0))
				}
				i = j - 1 // move cursor
				break
			}
		}
	}
	return totalMovingTime
}
