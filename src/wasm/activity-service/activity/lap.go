package activity

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/accumulator"
	"github.com/muktihari/openactivity-fit/kit"
)

type Lap struct {
	Sport            string
	Timestamp        time.Time
	StartTime        time.Time
	EndTime          time.Time
	TotalMovingTime  float64
	TotalElapsedTime float64
	TotalDistance    float64
	TotalAscent      uint16
	TotalDescent     uint16
	TotalCalories    uint16
	AvgSpeed         *float64
	MaxSpeed         *float64
	AvgHeartRate     *uint8
	MaxHeartRate     *uint8
	AvgCadence       *uint8
	MaxCadence       *uint8
	AvgPower         *uint16
	MaxPower         *uint16
	AvgTemperature   *int8
	MaxTemperature   *int8
	AvgAltitude      *float64
	MaxAltitude      *float64
	AvgPace          *float64
	AvgElapsedPace   *float64
}

func NewLapFromRecords(records []*Record, sport string) *Lap {
	if len(records) == 0 {
		return nil
	}

	lap := &Lap{
		Timestamp: records[0].Timestamp,
		StartTime: records[0].Timestamp,
		EndTime:   records[len(records)-1].Timestamp,
	}

	var (
		distanceAccumu    = new(accumulator.Accumulator[float64])
		speedAccumu       = new(accumulator.Accumulator[float64])
		altitudeAccumu    = new(accumulator.Accumulator[float64])
		cadenceAccumu     = new(accumulator.Accumulator[uint8])
		heartRateAccumu   = new(accumulator.Accumulator[uint8])
		powerAccumu       = new(accumulator.Accumulator[uint16])
		temperatureAccumu = new(accumulator.Accumulator[int8])
	)

	for i := 0; i < len(records); i++ {
		rec := records[i]

		distanceAccumu.Collect(rec.Distance)
		speedAccumu.Collect(rec.Speed)
		altitudeAccumu.Collect(rec.Altitude)
		cadenceAccumu.Collect(rec.Cadence)
		heartRateAccumu.Collect(rec.HeartRate)
		powerAccumu.Collect(rec.Power)
		temperatureAccumu.Collect(rec.Temperature)
	}

	// Calculate Total Elapsed and Total Moving Time
	for i := 0; i < len(records); i++ {
		rec := records[i]
		if rec.Timestamp.IsZero() {
			continue
		}

		// Find next non-zero timestamp
		for j := i + 1; j < len(records); j++ {
			next := records[j]
			if !next.Timestamp.IsZero() {
				delta := next.Timestamp.Sub(rec.Timestamp).Seconds()
				lap.TotalElapsedTime += delta

				if IsConsideredMoving(sport, rec.Speed) {
					lap.TotalMovingTime += delta
				}
				i = j - 1 // move cursor
				break
			}
		}
	}

	// Calculate Total Ascent and Total Descent
	var totalAscent, totalDescent float64
	for i := 0; i < len(records)-1; i++ {
		rec := records[i]
		if rec.SmoothedAltitude == nil {
			continue
		}

		// Find next non-nil altitude
		for j := i + 1; j < len(records); j++ {
			next := records[j]
			if next.SmoothedAltitude != nil {
				delta := *next.SmoothedAltitude - *rec.SmoothedAltitude
				if delta > 0 {
					totalAscent += delta
				} else {
					totalDescent += math.Abs(delta)
				}
				i = j - 1 // move cursor
				break
			}
		}
	}

	lap.TotalAscent = uint16(math.Round(totalAscent))
	lap.TotalDescent = uint16(math.Round(totalDescent))

	if distanceAccumu.Min() != nil && distanceAccumu.Max() != nil {
		lap.TotalDistance = *distanceAccumu.Max() - *distanceAccumu.Min()
	}
	lap.AvgSpeed = speedAccumu.Avg()
	lap.MaxSpeed = speedAccumu.Max()
	lap.AvgAltitude = altitudeAccumu.Avg()
	lap.MaxAltitude = altitudeAccumu.Max()
	lap.AvgCadence = cadenceAccumu.Avg()
	lap.MaxCadence = cadenceAccumu.Max()
	lap.AvgHeartRate = heartRateAccumu.Avg()
	lap.MaxHeartRate = heartRateAccumu.Max()
	lap.AvgPower = powerAccumu.Avg()
	lap.MaxPower = powerAccumu.Max()
	lap.AvgTemperature = temperatureAccumu.Avg()
	lap.MaxTemperature = temperatureAccumu.Max()

	if HasPace(sport) {
		lap.AvgPace = kit.Ptr(lap.TotalMovingTime / (lap.TotalDistance / 1000))
		lap.AvgElapsedPace = kit.Ptr(lap.TotalElapsedTime / (lap.TotalDistance / 1000))
	}

	return lap
}

func NewLapFromSession(session *Session) *Lap {
	return &Lap{
		Timestamp:        session.Timestamp,
		StartTime:        session.StartTime,
		EndTime:          session.EndTime,
		TotalMovingTime:  session.TotalMovingTime,
		TotalElapsedTime: session.TotalElapsedTime,
		TotalDistance:    session.TotalDistance,
		TotalAscent:      session.TotalAscent,
		TotalDescent:     session.TotalDescent,
		TotalCalories:    session.TotalCalories,
		AvgSpeed:         session.AvgSpeed,
		MaxSpeed:         session.MaxSpeed,
		AvgHeartRate:     session.AvgHeartRate,
		MaxHeartRate:     session.MaxHeartRate,
		AvgCadence:       session.AvgCadence,
		MaxCadence:       session.MaxCadence,
		AvgPower:         session.AvgPower,
		MaxPower:         session.MaxPower,
		AvgTemperature:   session.AvgTemperature,
		MaxTemperature:   session.MaxTemperature,
		AvgAltitude:      session.AvgAltitude,
		MaxAltitude:      session.MaxAltitude,
	}
}

var _ json.Marshaler = &Lap{}

func (l *Lap) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteByte('{')

	buf.WriteString("\"sport\":\"" + l.Sport + "\",")
	if !l.Timestamp.IsZero() {
		buf.WriteString("\"timestamp\":\"" + l.Timestamp.Format(time.RFC3339) + "\",")
	}
	if !l.StartTime.IsZero() {
		buf.WriteString("\"startTime\":\"" + l.StartTime.Format(time.RFC3339) + "\",")
	}
	if !l.EndTime.IsZero() {
		buf.WriteString("\"endTime\":\"" + l.EndTime.Format(time.RFC3339) + "\",")
	}

	buf.WriteString("\"totalMovingTime\":" + strconv.FormatFloat(l.TotalMovingTime, 'g', -1, 64) + ",")
	buf.WriteString("\"totalElapsedTime\":" + strconv.FormatFloat(l.TotalElapsedTime, 'g', -1, 64) + ",")
	buf.WriteString("\"totalDistance\":" + strconv.FormatFloat(l.TotalDistance, 'g', -1, 64) + ",")
	buf.WriteString("\"totalAscent\":" + strconv.FormatUint(uint64(l.TotalAscent), 10) + ",")
	buf.WriteString("\"totalDescent\":" + strconv.FormatUint(uint64(l.TotalDescent), 10) + ",")
	buf.WriteString("\"totalCalories\":" + strconv.FormatUint(uint64(l.TotalCalories), 10) + ",")

	if l.AvgSpeed != nil {
		buf.WriteString("\"avgSpeed\":" + strconv.FormatFloat(*l.AvgSpeed, 'g', -1, 64) + ",")
	}
	if l.MaxSpeed != nil {
		buf.WriteString("\"maxSpeed\":" + strconv.FormatFloat(*l.MaxSpeed, 'g', -1, 64) + ",")
	}
	if l.AvgHeartRate != nil {
		buf.WriteString("\"avgHeartRate\":" + strconv.FormatUint(uint64(*l.AvgHeartRate), 10) + ",")
	}
	if l.MaxHeartRate != nil {
		buf.WriteString("\"maxHeartRate\":" + strconv.FormatUint(uint64(*l.MaxHeartRate), 10) + ",")
	}
	if l.AvgCadence != nil {
		buf.WriteString("\"avgCadence\":" + strconv.FormatUint(uint64(*l.AvgCadence), 10) + ",")
	}
	if l.MaxCadence != nil {
		buf.WriteString("\"maxCadence\":" + strconv.FormatUint(uint64(*l.MaxCadence), 10) + ",")
	}
	if l.AvgPower != nil {
		buf.WriteString("\"avgPower\":" + strconv.FormatUint(uint64(*l.AvgPower), 10) + ",")
	}
	if l.MaxPower != nil {
		buf.WriteString("\"maxPower\":" + strconv.FormatUint(uint64(*l.MaxPower), 10) + ",")
	}
	if l.AvgTemperature != nil {
		buf.WriteString("\"avgTemperature\":" + strconv.FormatInt(int64(*l.AvgTemperature), 10) + ",")
	}
	if l.MaxTemperature != nil {
		buf.WriteString("\"maxTemperature\":" + strconv.FormatInt(int64(*l.MaxTemperature), 10) + ",")
	}
	if l.AvgAltitude != nil {
		buf.WriteString("\"avgAltitude\":" + strconv.FormatFloat(*l.AvgAltitude, 'g', -1, 64) + ",")
	}
	if l.MaxAltitude != nil {
		buf.WriteString("\"maxAltitude\":" + strconv.FormatFloat(*l.MaxAltitude, 'g', -1, 64) + ",")
	}
	if l.AvgPace != nil {
		buf.WriteString("\"avgPace\":" + strconv.FormatFloat(*l.AvgPace, 'g', -1, 64) + ",")
	}
	if l.AvgElapsedPace != nil {
		buf.WriteString("\"avgElapsedPace\":" + strconv.FormatFloat(*l.AvgElapsedPace, 'g', -1, 64))
	}

	b := buf.Bytes()
	if b[len(b)-1] == ',' {
		b[len(b)-1] = '}'
		return b, nil
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (l *Lap) IsBelongToThisLap(t time.Time) bool {
	return isBelong(t, l.StartTime, l.EndTime)
}

func (l *Lap) Clone() *Lap {
	lap := &Lap{
		Sport:            l.Sport,
		Timestamp:        l.Timestamp,
		StartTime:        l.StartTime,
		EndTime:          l.EndTime,
		TotalMovingTime:  l.TotalMovingTime,
		TotalElapsedTime: l.TotalElapsedTime,
		TotalDistance:    l.TotalDistance,
		TotalAscent:      l.TotalAscent,
		TotalDescent:     l.TotalDescent,
		TotalCalories:    l.TotalCalories,
	}

	if l.AvgSpeed != nil {
		lap.AvgSpeed = kit.Ptr(*l.AvgSpeed)
	}
	if l.MaxSpeed != nil {
		lap.MaxSpeed = kit.Ptr(*l.MaxSpeed)
	}
	if l.AvgHeartRate != nil {
		lap.AvgHeartRate = kit.Ptr(*l.AvgHeartRate)
	}
	if l.MaxHeartRate != nil {
		lap.MaxHeartRate = kit.Ptr(*l.MaxHeartRate)
	}
	if l.AvgCadence != nil {
		lap.AvgCadence = kit.Ptr(*l.AvgCadence)
	}
	if l.MaxCadence != nil {
		lap.MaxCadence = kit.Ptr(*l.MaxCadence)
	}
	if l.AvgPower != nil {
		lap.AvgPower = kit.Ptr(*l.AvgPower)
	}
	if l.MaxPower != nil {
		lap.MaxPower = kit.Ptr(*l.MaxPower)
	}
	if l.AvgTemperature != nil {
		lap.AvgTemperature = kit.Ptr(*l.AvgTemperature)
	}
	if l.MaxTemperature != nil {
		lap.MaxTemperature = kit.Ptr(*l.MaxTemperature)
	}
	if l.AvgAltitude != nil {
		lap.AvgAltitude = kit.Ptr(*l.AvgAltitude)
	}
	if l.MaxAltitude != nil {
		lap.MaxAltitude = kit.Ptr(*l.MaxAltitude)
	}
	if l.AvgPace != nil {
		lap.AvgPace = kit.Ptr(*l.AvgPace)
	}
	if l.AvgElapsedPace != nil {
		lap.AvgElapsedPace = kit.Ptr(*l.AvgElapsedPace)
	}

	return lap
}

// CombineLap combines lap's values into targetLap.
// Every zero value in targetLap will be replaced with the corresponding value in lap.
func CombineLap(targetLap, lap *Lap) {
	if targetLap == nil || lap == nil {
		return
	}

	if targetLap.StartTime.IsZero() {
		targetLap.StartTime = lap.StartTime
	}

	if targetLap.EndTime.IsZero() {
		targetLap.EndTime = lap.EndTime
	}

	targetLap.TotalElapsedTime = kit.PickNonZeroValue(targetLap.TotalElapsedTime, lap.TotalElapsedTime)
	targetLap.TotalMovingTime = kit.PickNonZeroValue(targetLap.TotalMovingTime, lap.TotalMovingTime)
	targetLap.TotalDistance = kit.PickNonZeroValue(targetLap.TotalDistance, lap.TotalDistance)
	targetLap.TotalCalories = kit.PickNonZeroValue(targetLap.TotalCalories, lap.TotalCalories)
	targetLap.TotalAscent = kit.PickNonZeroValue(targetLap.TotalAscent, lap.TotalAscent)
	targetLap.TotalDescent = kit.PickNonZeroValue(targetLap.TotalDescent, lap.TotalDescent)
	targetLap.AvgSpeed = kit.PickNonZeroValuePtr(targetLap.AvgSpeed, lap.AvgSpeed)
	targetLap.MaxSpeed = kit.PickNonZeroValuePtr(targetLap.MaxSpeed, lap.MaxSpeed)
	targetLap.AvgHeartRate = kit.PickNonZeroValuePtr(targetLap.AvgHeartRate, lap.AvgHeartRate)
	targetLap.MaxHeartRate = kit.PickNonZeroValuePtr(targetLap.MaxHeartRate, lap.MaxHeartRate)
	targetLap.AvgCadence = kit.PickNonZeroValuePtr(targetLap.AvgCadence, lap.AvgCadence)
	targetLap.MaxCadence = kit.PickNonZeroValuePtr(targetLap.MaxCadence, lap.MaxCadence)
	targetLap.AvgPower = kit.PickNonZeroValuePtr(targetLap.AvgPower, lap.AvgPower)
	targetLap.MaxPower = kit.PickNonZeroValuePtr(targetLap.MaxPower, lap.MaxPower)
	targetLap.AvgTemperature = kit.PickNonZeroValuePtr(targetLap.AvgTemperature, lap.AvgTemperature)
	targetLap.MaxTemperature = kit.PickNonZeroValuePtr(targetLap.MaxTemperature, lap.MaxTemperature)
	targetLap.AvgAltitude = kit.PickNonZeroValuePtr(targetLap.AvgAltitude, lap.AvgAltitude)
	targetLap.MaxAltitude = kit.PickNonZeroValuePtr(targetLap.MaxAltitude, lap.MaxAltitude)
	targetLap.AvgPace = kit.PickNonZeroValuePtr(targetLap.AvgPace, lap.AvgPace)
	targetLap.AvgElapsedPace = kit.PickNonZeroValuePtr(targetLap.AvgElapsedPace, lap.AvgElapsedPace)
}
