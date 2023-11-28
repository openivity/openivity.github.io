package scale

import (
	"time"

	"github.com/muktihari/openactivity-fit/accumulator"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/kit"
)

type Result struct {
	CombinedRecords []*activity.Record
	SessionRecords  [][]*activity.Record
}

type Scale struct {
	percentage float64
}

func New(percentage float64) *Scale {
	return &Scale{percentage: percentage}
}

func (s *Scale) Scale(activities []*activity.Activity) Result {
	var combinedRecords []*activity.Record
	var sessionRecords [][]*activity.Record

	prevSesDistance := 0.0
	for i := range activities {
		act := activities[i]

		for j := range act.Sessions {
			ses := act.Sessions[j]

			sessionRecords = append(sessionRecords, s.scaleRecords(ses.Records))

			for k := range ses.Records {
				rec := cloneRecord(ses.Records[k])
				if rec.Distance != nil {
					*rec.Distance += prevSesDistance
				}
				combinedRecords = append(combinedRecords, rec)
			}

			prevSesDistance = lastDistance(ses.Records)
		}
	}

	combinedRecords = s.scaleRecords(combinedRecords)

	return Result{
		CombinedRecords: combinedRecords,
		SessionRecords:  sessionRecords,
	}
}

func (s *Scale) scaleRecords(records []*activity.Record) []*activity.Record {
	summarizeByDistance := hasDistance(records)
	var threshold float64 // the unit value is either distance meters or duration seconds depends on conditions below:
	if summarizeByDistance {
		ld := lastDistance(records)
		threshold = ld * s.percentage
	} else {
		st := startTime(records)
		et := endTime(records)
		dur := et.Sub(st).Seconds()
		threshold = dur * s.percentage
	}

	var timestamp time.Time
	var positionLat *float64
	var positionLong *float64
	var distance *float64

	var (
		altitudeAccumu  = new(accumulator.Accumulator[float64])
		cadenceAccumu   = new(accumulator.Accumulator[uint8])
		heartrateAccumu = new(accumulator.Accumulator[uint8])
		speedAccumu     = new(accumulator.Accumulator[float64])
		powerAccumu     = new(accumulator.Accumulator[uint16])
		tempAccumu      = new(accumulator.Accumulator[int8])
		paceAccumu      = new(accumulator.Accumulator[float64])
		gradeAccumu     = new(accumulator.Accumulator[float64])
	)

	summarizedRecords := make([]*activity.Record, 0)

	var curIndex int
	for i := range records {
		rec := records[i]
		cur := records[curIndex]

		var delta float64
		if summarizeByDistance {
			if rec.Distance != nil && cur.Distance != nil {
				delta = *rec.Distance - *cur.Distance
			}
		} else {
			if !rec.Timestamp.IsZero() && !cur.Timestamp.IsZero() {
				delta = rec.Timestamp.Sub(cur.Timestamp).Seconds()
			}
		}

		if delta > threshold {
			newRec := &activity.Record{
				Timestamp:    timestamp,
				PositionLat:  positionLat,
				PositionLong: positionLong,
				Distance:     distance,
				Altitude:     altitudeAccumu.Avg(),
				Cadence:      cadenceAccumu.Avg(),
				HeartRate:    heartrateAccumu.Avg(),
				Speed:        speedAccumu.Avg(),
				Power:        powerAccumu.Avg(),
				Temperature:  tempAccumu.Avg(),
				Grade:        gradeAccumu.Avg(),
				Pace:         paceAccumu.Avg(),
			}

			summarizedRecords = append(summarizedRecords, newRec)

			// Reset
			timestamp = time.Time{}
			positionLat = nil
			positionLong = nil
			distance = nil
			altitudeAccumu = new(accumulator.Accumulator[float64])
			cadenceAccumu = new(accumulator.Accumulator[uint8])
			heartrateAccumu = new(accumulator.Accumulator[uint8])
			speedAccumu = new(accumulator.Accumulator[float64])
			powerAccumu = new(accumulator.Accumulator[uint16])
			tempAccumu = new(accumulator.Accumulator[int8])
			paceAccumu = new(accumulator.Accumulator[float64])
			gradeAccumu = new(accumulator.Accumulator[float64])

			curIndex = i
		}

		if timestamp.IsZero() {
			timestamp = rec.Timestamp
		}
		if positionLat == nil {
			positionLat = rec.PositionLat
		}
		if positionLong == nil {
			positionLong = rec.PositionLong
		}
		if distance == nil {
			distance = rec.Distance
		}

		altitudeAccumu.Collect(rec.Altitude)
		cadenceAccumu.Collect(rec.Cadence)
		heartrateAccumu.Collect(rec.HeartRate)
		speedAccumu.Collect(rec.Speed)
		powerAccumu.Collect(rec.Power)
		tempAccumu.Collect(rec.Temperature)
		paceAccumu.Collect(rec.Pace)
		gradeAccumu.Collect(rec.Grade)
	}

	return summarizedRecords
}

func cloneRecord(rec *activity.Record) *activity.Record {
	clone := new(activity.Record)

	clone.Timestamp = rec.Timestamp
	if rec.Distance != nil {
		clone.Distance = kit.Ptr(*rec.Distance)
	}

	clone.PositionLat = rec.PositionLat
	clone.PositionLong = rec.PositionLong
	clone.Altitude = rec.Altitude
	clone.HeartRate = rec.HeartRate
	clone.Cadence = rec.Cadence
	clone.Speed = rec.Speed
	clone.Power = rec.Power
	clone.Temperature = rec.Temperature
	clone.Pace = rec.Pace
	clone.Grade = rec.Grade

	return clone
}

func startTime(records []*activity.Record) time.Time {
	var t time.Time
	for i := range records {
		rec := records[i]
		if !t.IsZero() {
			break
		}
		t = rec.Timestamp
	}
	return t
}

func endTime(records []*activity.Record) time.Time {
	var t time.Time
	for i := len(records) - 1; i >= 0; i-- {
		rec := records[i]
		if !t.IsZero() {
			break
		}
		t = rec.Timestamp
	}
	return t
}

func lastDistance(records []*activity.Record) float64 {
	var d float64
	for i := len(records) - 1; i >= 0; i-- {
		rec := records[i]
		if d > 0 {
			break
		}
		if rec.Distance != nil {
			d = *rec.Distance
		}
	}
	return d
}

func hasDistance(records []*activity.Record) bool {
	for i := range records {
		rec := records[i]
		if rec.Distance != nil && *rec.Distance > 0 {
			return true
		}
	}
	return false

}
