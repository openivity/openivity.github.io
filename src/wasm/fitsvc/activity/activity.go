package activity

import (
	"time"

	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/kit/semicircles"
	"github.com/muktihari/fit/kit/typeconv"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/proto"
)

type ActivityFile struct {
	FileId   map[string]any `json:"fileId"`
	Sessions []any          `json:"sessions,omitempty"`
	Laps     []any          `json:"laps,omitempty"`
	Records  []any          `json:"records,omitempty"`
}

func (m *ActivityFile) ToMap() map[string]any {
	return map[string]any{
		"fileId":   m.FileId,
		"sessions": m.Sessions,
		"laps":     m.Laps,
		"records":  m.Records,
	}
}

type FileId struct {
	Manufacturer string    `json:"manufacturer"`
	Product      uint16    `json:"product"`
	TimeCreated  time.Time `json:"timeCreated"`
}

func (m *FileId) ToMap() map[string]any {
	return map[string]any{
		"manufacturer": m.Manufacturer,
		"product":      m.Product,
		"timeCreated":  m.TimeCreated.Format(time.RFC3339),
	}
}

func NewFileId(mesg proto.Message) map[string]any {
	m := map[string]any{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.FileIdManufacturer:
			m["manufacturer"] = typeconv.ToUint16[typedef.Manufacturer](field.Value).String()
		case fieldnum.FileIdProduct:
			m["product"], _ = field.Value.(uint16)
		case fieldnum.FileIdTimeCreated:
			m["timeCreated"] = datetime.ToTime(field.Value).Format(time.RFC3339)
		}
	}

	return m
}

func NewSession(mesg proto.Message) map[string]any {
	m := map[string]any{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.SessionSport:
			sport, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			m["sport"] = typedef.Sport(sport).String()
		case fieldnum.SessionSubSport:
			subSport, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			m["subSport"] = typedef.SubSport(subSport).String()
		case fieldnum.SessionTotalMovingTime:
			totalMovingTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["totalMovingTime"] = float64(totalMovingTime) / 1000
		case fieldnum.SessionTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["totalElapsedTime"] = float64(totalElapsedTime) / 1000
		case fieldnum.SessionTotalTimerTime:
			totalTimerTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["totalTimerTime"] = float64(totalTimerTime) / 1000
		case fieldnum.SessionTotalDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["totalDistance"] = float64(totalDistance) / 100
		case fieldnum.SessionTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["totalAscent"] = totalAscent
		case fieldnum.SessionTotalDescent:
			totalDescent, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["totalDescent"] = totalDescent
		case fieldnum.SessionTotalCycles:
			totalCycles, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			if totalCycles == basetype.Uint32Invalid {
				continue
			}
			m["totalCycles"] = totalCycles
		case fieldnum.SessionTotalCalories:
			totalCalories, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["totalCalories"] = totalCalories
		case fieldnum.SessionAvgSpeed:
			avgSpeed, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["avgSpeed"] = float64(avgSpeed) / 1000
		case fieldnum.SessionMaxSpeed:
			maxSpeed, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["maxSpeed"] = float64(maxSpeed) / 1000
		case fieldnum.SessionAvgHeartRate:
			avgHeartRate, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			m["avgHeartRate"] = avgHeartRate
		case fieldnum.SessionMaxHeartRate:
			maxHeartRate, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			m["maxHeartRate"] = maxHeartRate
		case fieldnum.SessionAvgCadence:
			avgCadence, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			if avgCadence == basetype.Uint8Invalid {
				continue
			}
			m["avgCadence"] = avgCadence
		case fieldnum.SessionMaxCadence:
			maxCadence, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			if maxCadence == basetype.Uint8Invalid {
				continue
			}
			m["maxCadence"] = maxCadence
		case fieldnum.SessionAvgPower:
			avgPower, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			if avgPower == basetype.Uint16Invalid {
				continue
			}
			m["avgPower"] = avgPower
		case fieldnum.SessionMaxPower:
			maxPower, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			if maxPower == basetype.Uint16Invalid {
				continue
			}
			m["maxPower"] = maxPower
		case fieldnum.SessionAvgTemperature:
			avgTemperature, ok := field.Value.(int8)
			if !ok {
				continue
			}
			m["avgTemperature"] = avgTemperature
		case fieldnum.SessionMaxTemperature:
			maxTemperature, ok := field.Value.(int8)
			if !ok {
				continue
			}
			m["maxTemperature"] = maxTemperature
		case fieldnum.SessionAvgAltitude:
			avgAltitude, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			fAvgAltitude := (float64(avgAltitude) / 5) - 500
			if fAvgAltitude < 0 {
				continue
			}
			m["avgAltitude"] = fAvgAltitude
		case fieldnum.SessionMaxAltitude:
			maxAltitude, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			fMaxAltitude := (float64(maxAltitude) / 5) - 500
			if fMaxAltitude < 0 {
				continue
			}
			m["maxAltitude"] = fMaxAltitude
		}
	}

	return m
}

func NewLap(mesg proto.Message) map[string]any {
	m := map[string]any{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.LapTimestamp:
			m["timestamp"] = datetime.ToTime(field.Value).Format(time.RFC3339)
		case fieldnum.LapTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["totalElapsedTime"] = float64(totalElapsedTime) / 1000
		case fieldnum.LapTotalTimerTime:
			totalTimerTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["totalTimerTime"] = float64(totalTimerTime) / 1000
		case fieldnum.TotalsDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["totalDistance"] = float64(totalDistance) / 100
		case fieldnum.LapTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["totalAscent"] = totalAscent
		}
	}

	return m
}

func NewRecord(mesg proto.Message) map[string]any {
	m := map[string]any{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.RecordTimestamp:
			m["timestamp"] = datetime.ToTime(field.Value).Format(time.RFC3339)
		case fieldnum.RecordPositionLat:
			lat, ok := field.Value.(int32)
			if !ok {
				continue
			}
			m["positionLat"] = semicircles.ToDegrees(lat)
		case fieldnum.RecordPositionLong:
			long, ok := field.Value.(int32)
			if !ok {
				continue
			}
			m["positionLong"] = semicircles.ToDegrees(long)
		case fieldnum.RecordAltitude:
			altitude, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			faltitude := (float64(altitude) / 5) - 500
			if faltitude < 0 {
				continue
			}
			m["altitude"] = faltitude
		case fieldnum.RecordHeartRate:
			heartRate, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			m["heartRate"] = heartRate
		case fieldnum.RecordCadence:
			cadence, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			m["cadence"] = cadence
		case fieldnum.RecordDistance:
			distance, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			m["distance"] = float64(distance) / 100
		case fieldnum.RecordSpeed:
			speed, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["speed"] = float64(speed) / 1000
		case fieldnum.RecordPower:
			power, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			m["power"] = power
		}
	}

	return m
}
