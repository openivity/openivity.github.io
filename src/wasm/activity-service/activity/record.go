package activity

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/kit"
)

type Record struct {
	Timestamp        time.Time
	PositionLat      *float64
	PositionLong     *float64
	Distance         *float64
	Altitude         *float64 // Original Altitude from file.
	SmoothedAltitude *float64 // Smoothed Altitude using our preprocessor algorithm.
	HeartRate        *uint8
	Cadence          *uint8
	Speed            *float64
	Power            *uint16
	Temperature      *int8
	Pace             *float64
	Grade            *float64
}

var _ json.Marshaler = &Record{}

func (r *Record) MarshalJSON() ([]byte, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('{')

	if !r.Timestamp.IsZero() {
		buf.WriteString("\"timestamp\":\"")
		buf.WriteString(r.Timestamp.Format(time.RFC3339))
		buf.WriteString("\",")
	}
	if r.PositionLat != nil {
		buf.WriteString("\"positionLat\":")
		buf.WriteString(strconv.FormatFloat(*r.PositionLat, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if r.PositionLong != nil {
		buf.WriteString("\"positionLong\":")
		buf.WriteString(strconv.FormatFloat(*r.PositionLong, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if r.Distance != nil {
		buf.WriteString("\"distance\":")
		buf.WriteString(strconv.FormatFloat(*r.Distance, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if r.Altitude != nil {
		buf.WriteString("\"altitude\":")
		buf.WriteString(strconv.FormatFloat(*r.SmoothedAltitude, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if r.HeartRate != nil {
		buf.WriteString("\"heartRate\":")
		buf.WriteString(strconv.FormatUint(uint64(*r.HeartRate), 10))
		buf.WriteByte(',')
	}
	if r.Cadence != nil {
		buf.WriteString("\"cadence\":")
		buf.WriteString(strconv.FormatUint(uint64(*r.Cadence), 10))
		buf.WriteByte(',')
	}
	if r.Speed != nil {
		buf.WriteString("\"speed\":")
		buf.WriteString(strconv.FormatFloat(*r.Speed, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if r.Power != nil {
		buf.WriteString("\"power\":")
		buf.WriteString(strconv.FormatUint(uint64(*r.Power), 10))
		buf.WriteByte(',')
	}
	if r.Temperature != nil {
		buf.WriteString("\"temperature\":")
		buf.WriteString(strconv.FormatInt(int64(*r.Temperature), 10))
		buf.WriteByte(',')
	}
	if r.Pace != nil {
		buf.WriteString("\"pace\":")
		buf.WriteString(strconv.FormatFloat(*r.Pace, 'g', -1, 64))
		buf.WriteByte(',')
	}
	if r.Grade != nil {
		buf.WriteString("\"grade\":")
		buf.WriteString(strconv.FormatFloat(*r.Grade, 'g', -1, 64))
	}

	b := buf.Bytes()
	if len(b) == 1 { // only '{' means all fields is invalid
		return nil, nil
	}

	if b[len(b)-1] == ',' {
		b[len(b)-1] = '}'
		return b, nil
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (r *Record) Clone() *Record {
	rec := &Record{
		Timestamp: r.Timestamp,
	}

	if r.PositionLat != nil {
		rec.PositionLat = kit.Ptr(*r.PositionLat)
	}
	if r.PositionLong != nil {
		rec.PositionLong = kit.Ptr(*r.PositionLong)
	}
	if r.Distance != nil {
		rec.Distance = kit.Ptr(*r.Distance)
	}
	if r.Altitude != nil {
		rec.Altitude = kit.Ptr(*r.Altitude)
	}
	if r.SmoothedAltitude != nil {
		rec.SmoothedAltitude = kit.Ptr(*r.SmoothedAltitude)
	}
	if r.HeartRate != nil {
		rec.HeartRate = kit.Ptr(*r.HeartRate)
	}
	if r.Cadence != nil {
		rec.Cadence = kit.Ptr(*r.Cadence)
	}
	if r.Speed != nil {
		rec.Speed = kit.Ptr(*r.Speed)
	}
	if r.Power != nil {
		rec.Power = kit.Ptr(*r.Power)
	}
	if r.Temperature != nil {
		rec.Temperature = kit.Ptr(*r.Temperature)
	}
	if r.Pace != nil {
		rec.Pace = kit.Ptr(*r.Pace)
	}
	if r.Grade != nil {
		rec.Grade = kit.Ptr(*r.Grade)
	}

	return rec
}
