package activity

import (
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/kit/semicircles"
	"github.com/muktihari/fit/kit/typeconv"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/proto"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ActivityFile struct {
	Creator  Creator
	Timezone int
	Sessions []any
	Laps     []any
	Records  []any
}

func (m *ActivityFile) ToMap() map[string]any {
	return map[string]any{
		"creator":  m.Creator.ToMap(),
		"timezone": m.Timezone,
		"sessions": m.Sessions,
		"laps":     m.Laps,
		"records":  m.Records,
	}
}

type Creator struct {
	Name         string
	Manufacturer uint16
	Product      uint16
	TimeCreated  time.Time
}

func (m *Creator) ToMap() map[string]any {
	if m.Name == "" {
		m.Name = "Unknown"
	}
	return map[string]any{
		"name":         m.Name,
		"manufacturer": m.Manufacturer,
		"product":      m.Product,
		"timeCreated":  m.TimeCreated.Format(time.RFC3339),
	}
}

func CreateTimezone(mesg proto.Message) int {
	var (
		timestamp     = basetype.Uint32Invalid
		localDateTime = basetype.Uint32Invalid
	)

	for i := range mesg.Fields {
		switch mesg.Fields[i].Num {
		case fieldnum.ActivityTimestamp:
			t, ok := mesg.Fields[i].Value.(uint32)
			if !ok {
				continue
			}
			timestamp = t
		case fieldnum.ActivityLocalTimestamp:
			t, ok := mesg.Fields[i].Value.(uint32)
			if !ok {
				continue
			}
			localDateTime = t
		}
	}

	if timestamp == basetype.Uint32Invalid || localDateTime == basetype.Uint32Invalid {
		return 0 // Default UTC
	}

	return datetime.TzOffsetHours(
		typedef.LocalDateTime(localDateTime),
		typedef.DateTime(timestamp),
	)
}

func NewCreator(mesg proto.Message) Creator {
	m := Creator{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.FileIdManufacturer:
			m.Manufacturer, _ = field.Value.(uint16)
			m.Name = title(typeconv.ToUint16[typedef.Manufacturer](field.Value).String())
		case fieldnum.FileIdProduct:
			m.Product, _ = field.Value.(uint16)
			m.Name += " (" + strconv.FormatUint(uint64(m.Product), 10) + ")"
		case fieldnum.FileIdTimeCreated:
			m.TimeCreated = datetime.ToTime(field.Value)
		}
	}

	return m
}

func NewSession(mesg proto.Message) map[string]any {
	m := map[string]any{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.SessionTimestamp:
			m["timestamp"] = datetime.ToTime(field.Value).Format(time.RFC3339)
		case fieldnum.SessionSport:
			sport, ok := field.Value.(uint8)
			if !ok || sport == basetype.EnumInvalid {
				continue
			}
			m["sport"] = title(typedef.Sport(sport).String())
		case fieldnum.SessionSubSport:
			subSport, ok := field.Value.(uint8)
			if !ok || subSport == basetype.EnumInvalid {
				continue
			}
			m["subSport"] = title(typedef.SubSport(subSport).String())
		case fieldnum.SessionTotalMovingTime:
			totalMovingTime, ok := field.Value.(uint32)
			if !ok || totalMovingTime == basetype.Uint32Invalid {
				continue
			}
			m["totalMovingTime"] = float64(totalMovingTime) / 1000
		case fieldnum.SessionTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok || totalElapsedTime == basetype.Uint32Invalid {
				continue
			}
			m["totalElapsedTime"] = float64(totalElapsedTime) / 1000
		case fieldnum.SessionTotalDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok || totalDistance == basetype.Uint32Invalid {
				continue
			}
			m["totalDistance"] = float64(totalDistance) / 100
		case fieldnum.SessionTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok || totalAscent == basetype.Uint16Invalid {
				continue
			}
			m["totalAscent"] = totalAscent
		case fieldnum.SessionTotalDescent:
			totalDescent, ok := field.Value.(uint16)
			if !ok || totalDescent == basetype.Uint16Invalid {
				continue
			}
			m["totalDescent"] = totalDescent
		case fieldnum.SessionTotalCycles:
			totalCycles, ok := field.Value.(uint32)
			if !ok || totalCycles == basetype.Uint32Invalid {
				continue
			}
			m["totalCycles"] = totalCycles
		case fieldnum.SessionTotalCalories:
			totalCalories, ok := field.Value.(uint16)
			if !ok || totalCalories == basetype.Uint16Invalid {
				continue
			}
			m["totalCalories"] = totalCalories
		case fieldnum.SessionAvgSpeed:
			avgSpeed, ok := field.Value.(uint16)
			if !ok || avgSpeed == basetype.Uint16Invalid {
				continue
			}
			m["avgSpeed"] = float64(avgSpeed) / 1000
		case fieldnum.SessionMaxSpeed:
			maxSpeed, ok := field.Value.(uint16)
			if !ok || maxSpeed == basetype.Uint16Invalid {
				continue
			}
			m["maxSpeed"] = float64(maxSpeed) / 1000
		case fieldnum.SessionAvgHeartRate:
			avgHeartRate, ok := field.Value.(uint8)
			if !ok || avgHeartRate == basetype.Uint8Invalid {
				continue
			}
			m["avgHeartRate"] = avgHeartRate
		case fieldnum.SessionMaxHeartRate:
			maxHeartRate, ok := field.Value.(uint8)
			if !ok || maxHeartRate == basetype.Uint8Invalid {
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
			if !ok || maxCadence == basetype.Uint8Invalid {
				continue
			}
			m["maxCadence"] = maxCadence
		case fieldnum.SessionAvgPower:
			avgPower, ok := field.Value.(uint16)
			if !ok || avgPower == basetype.Uint16Invalid {
				continue
			}
			m["avgPower"] = avgPower
		case fieldnum.SessionMaxPower:
			maxPower, ok := field.Value.(uint16)
			if !ok || maxPower == basetype.Uint16Invalid {
				continue
			}
			m["maxPower"] = maxPower
		case fieldnum.SessionAvgTemperature:
			avgTemperature, ok := field.Value.(int8)
			if !ok || avgTemperature == basetype.Sint8Invalid {
				continue
			}
			m["avgTemperature"] = avgTemperature
		case fieldnum.SessionMaxTemperature:
			maxTemperature, ok := field.Value.(int8)
			if !ok || maxTemperature == basetype.Sint8Invalid {
				continue
			}
			m["maxTemperature"] = maxTemperature
		case fieldnum.SessionAvgAltitude:
			avgAltitude, ok := field.Value.(uint16)
			if !ok || avgAltitude == basetype.Uint16Invalid {
				continue
			}
			fAvgAltitude := (float64(avgAltitude) / 5) - 500
			if fAvgAltitude < 0 {
				continue
			}
			m["avgAltitude"] = fAvgAltitude
		case fieldnum.SessionMaxAltitude:
			maxAltitude, ok := field.Value.(uint16)
			if !ok || maxAltitude == basetype.Uint16Invalid {
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
		case fieldnum.LapTotalMovingTime:
			totalMovingTime, ok := field.Value.(uint32)
			if !ok || totalMovingTime == basetype.Uint32Invalid {
				continue
			}
			m["totalMovingTime"] = float64(totalMovingTime) / 1000
		case fieldnum.LapTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok || totalElapsedTime == basetype.Uint32Invalid {
				continue
			}
			m["totalElapsedTime"] = float64(totalElapsedTime) / 1000
		case fieldnum.LapTotalDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok || totalDistance == basetype.Uint32Invalid {
				continue
			}
			m["totalDistance"] = float64(totalDistance) / 100
		case fieldnum.LapTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok || totalAscent == basetype.Uint16Invalid {
				continue
			}
			m["totalAscent"] = totalAscent
		case fieldnum.LapTotalDescent:
			totalDescent, ok := field.Value.(uint16)
			if !ok || totalDescent == basetype.Uint16Invalid {
				continue
			}
			m["totalDescent"] = totalDescent
		case fieldnum.LapTotalCalories:
			totalCalories, ok := field.Value.(uint16)
			if !ok || totalCalories == basetype.Uint16Invalid {
				continue
			}
			m["totalCalories"] = totalCalories
		case fieldnum.LapAvgSpeed:
			avgSpeed, ok := field.Value.(uint16)
			if !ok || avgSpeed == basetype.Uint16Invalid {
				continue
			}
			m["avgSpeed"] = float64(avgSpeed) / 1000
		case fieldnum.LapMaxSpeed:
			maxSpeed, ok := field.Value.(uint16)
			if !ok || maxSpeed == basetype.Uint16Invalid {
				continue
			}
			m["maxSpeed"] = float64(maxSpeed) / 1000
		case fieldnum.LapAvgCadence:
			avgCadence, ok := field.Value.(uint8)
			if !ok || avgCadence == basetype.Uint8Invalid {
				continue
			}
			m["avgCadence"] = avgCadence
		case fieldnum.LapMaxCadence:
			maxCadence, ok := field.Value.(uint8)
			if !ok || maxCadence == basetype.Uint8Invalid {
				continue
			}
			m["maxCadence"] = maxCadence
		case fieldnum.LapAvgHeartRate:
			avgHeartRate, ok := field.Value.(uint8)
			if !ok || avgHeartRate == basetype.Uint8Invalid {
				continue
			}
			m["avgHeartRate"] = avgHeartRate
		case fieldnum.LapMaxHeartRate:
			maxHeartRate, ok := field.Value.(uint8)
			if !ok || maxHeartRate == basetype.Uint8Invalid {
				continue
			}
			m["maxHeartRate"] = maxHeartRate
		case fieldnum.LapAvgPower:
			avgPower, ok := field.Value.(uint16)
			if !ok || avgPower == basetype.Uint16Invalid {
				continue
			}
			m["avgPower"] = avgPower
		case fieldnum.LapMaxPower:
			maxPower, ok := field.Value.(uint16)
			if !ok || maxPower == basetype.Uint16Invalid {
				continue
			}
			m["maxPower"] = maxPower
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
			if !ok || lat == basetype.Sint32Invalid {
				continue
			}
			m["positionLat"] = semicircles.ToDegrees(lat)
		case fieldnum.RecordPositionLong:
			long, ok := field.Value.(int32)
			if !ok || long == basetype.Sint32Invalid {
				continue
			}
			m["positionLong"] = semicircles.ToDegrees(long)
		case fieldnum.RecordAltitude:
			altitude, ok := field.Value.(uint16)
			if !ok || altitude == basetype.Uint16Invalid {
				continue
			}
			faltitude := (float64(altitude) / 5) - 500
			if faltitude < 0 {
				continue
			}
			m["altitude"] = faltitude
		case fieldnum.RecordHeartRate:
			heartRate, ok := field.Value.(uint8)
			if !ok || heartRate == basetype.Uint8Invalid {
				continue
			}
			m["heartRate"] = heartRate
		case fieldnum.RecordCadence:
			cadence, ok := field.Value.(uint8)
			if !ok || cadence == basetype.Uint8Invalid {
				continue
			}
			m["cadence"] = cadence
		case fieldnum.RecordDistance:
			distance, ok := field.Value.(uint32)
			if !ok || distance == basetype.Uint32Invalid {
				continue
			}
			m["distance"] = float64(distance) / 100
		case fieldnum.RecordSpeed:
			speed, ok := field.Value.(uint16)
			if !ok || speed == basetype.Uint16Invalid {
				continue
			}
			m["speed"] = float64(speed) / 1000
		case fieldnum.RecordPower:
			power, ok := field.Value.(uint16)
			if !ok || power == basetype.Uint16Invalid {
				continue
			}
			m["power"] = power
		case fieldnum.RecordTemperature:
			temperature, ok := field.Value.(int8)
			if !ok || temperature == basetype.Sint8Invalid {
				continue
			}
			m["temperature"] = temperature
		}
	}

	return m
}

func title(s string) string {
	s = strings.Map(func(r rune) rune {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ' '
		}
		return r
	}, s)
	s = cases.Title(language.English).String(s)
	return s
}
