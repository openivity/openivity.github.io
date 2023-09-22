package activity

import (
	"time"

	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/kit/semicircles"
	"github.com/muktihari/fit/kit/typeconv"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	"github.com/muktihari/fit/proto"
)

type ActivityFile struct {
	FileId   FileId    `json:"fileId"`
	Activity Activity  `json:"activity"`
	Sessions []Session `json:"sessions,omitempty"`
	Laps     []Lap     `json:"laps,omitempty"`
	Records  []Record  `json:"records,omitempty"`
}

type FileId struct {
	Manufacturer string    `json:"manufacturer"`
	Product      uint16    `json:"product"`
	TimeCreated  time.Time `json:"timeCreated"`
}

func NewFileId(mesg proto.Message) FileId {
	fileId := FileId{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.FileIdManufacturer:
			fileId.Manufacturer = typeconv.ToUint16[typedef.Manufacturer](field.Value).String()
		case fieldnum.FileIdProduct:
			fileId.Product, _ = field.Value.(uint16)
		case fieldnum.FileIdTimeCreated:
			fileId.TimeCreated = datetime.ToTime(field.Value)
		}
	}

	return fileId
}

type Activity struct {
	Timestamp      time.Time `json:"timestamp"`
	TotalTimerTime uint32    `json:"totalTimerTime"`
	NumSessions    uint16    `json:"numSessions"`
	Type           string    `json:"type"`
	Event          string    `json:"event"`
	EventType      string    `json:"eventType"`
}

func NewActivity(mesg proto.Message) Activity {
	activity := Activity{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.ActivityTimestamp:
			activity.Timestamp = datetime.ToTime(field.Value)
		case fieldnum.ActivityTotalTimerTime:
			activity.TotalTimerTime, _ = field.Value.(uint32)
		case fieldnum.ActivityNumSessions:
			activity.NumSessions, _ = field.Value.(uint16)
		case fieldnum.ActivityType:
			activityType := typeconv.ToEnum[typedef.ActivityType](field.Value)
			activity.Type = activityType.String()
		case fieldnum.ActivityEvent:
			activityEvent := typeconv.ToEnum[typedef.Event](field.Value)
			activity.Event = activityEvent.String()
		case fieldnum.ActivityEventType:
			activityEventType := typeconv.ToEnum[typedef.EventType](field.Value)
			activity.EventType = activityEventType.String()
		}
	}

	return activity
}

type Session struct {
	Sport            *string  `json:"sport,omitempty"`
	SubSport         *string  `json:"subSport,omitempty"`
	TotalMovingTime  *float64 `json:"totalMovingTime,omitempty"`  // units: s;
	TotalElapsedTime *float64 `json:"totalElapsedTime,omitempty"` // Units: s;
	TotalTimerTime   *float64 `json:"totalTimerTime,omitempty"`   // Units: s;
	TotalDistance    *float64 `json:"totalDistance,omitempty"`    // Units: m;
	TotalAscent      *uint16  `json:"totalAscent,omitempty"`      // Units: m;
	TotalDescent     *uint16  `json:"totalDescent,omitempty"`     // Units: m;
	TotalCycles      *uint32  `json:"totalCycles"`                // Units: cycles;
	TotalCalories    *uint16  `json:"totalCalories"`              // Units: kcal;
	AvgSpeed         *float64 `json:"avgSpeed"`                   // Scale: 1000; Units: m/s; total_distance / total_timer_time
	MaxSpeed         *float64 `json:"maxSpeed"`                   // Scale: 1000; Units: m/s;
	AvgHeartRate     *uint8   `json:"avgHeartRate"`               // Units: bpm; average heart rate (excludes pause time)
	MaxHeartRate     *uint8   `json:"maxHeartRate"`               // Units: bpm;
	AvgCadence       *uint8   `json:"avgCadence"`                 // Units: rpm; total_cycles / total_timer_time if non_zero_avg_cadence otherwise total_cycles / total_elapsed_time
	MaxCadence       *uint8   `json:"maxCadence"`                 // Units: rpm;
	AvgPower         *uint16  `json:"avgPower"`                   // Units: watts; total_power / total_timer_time if non_zero_avg_power otherwise total_power / total_elapsed_time
	MaxPower         *uint16  `json:"maxPower"`                   // Units: watts;
	AvgTemperature   *int8    `json:"avgTemperature"`             // Units: C;
	MaxTemperature   *int8    `json:"maxTemperature"`             // Units: C;
	AvgAltitude      *float64 `json:"avgAltitude"`                // Scale: 5; Offset: 500; Units: m;
	MaxAltitude      *float64 `json:"maxAltitude"`                // Scale: 5; Offset: 500; Units: m;
}

func NewSession(mesg proto.Message) Session {
	session := Session{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.SessionSport:
			sport, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			sSport := typedef.Sport(sport).String()
			session.Sport = &sSport
		case fieldnum.SessionSubSport:
			subSport, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			sSubSport := typedef.SubSport(subSport).String()
			session.SubSport = &sSubSport
		case fieldnum.SessionTotalMovingTime:
			totalMovingTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fTotalMovingTime := float64(totalMovingTime) / 1000
			session.TotalMovingTime = &fTotalMovingTime
		case fieldnum.SessionTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fTotalElapsedTime := float64(totalElapsedTime) / 1000
			session.TotalElapsedTime = &fTotalElapsedTime
		case fieldnum.SessionTotalTimerTime:
			totalTimerTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fTotalTimerTime := float64(totalTimerTime) / 1000
			session.TotalTimerTime = &fTotalTimerTime
		case fieldnum.SessionTotalDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fTotalDistance := float64(totalDistance) / 100
			session.TotalDistance = &fTotalDistance
		case fieldnum.SessionTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			session.TotalAscent = &totalAscent
		case fieldnum.SessionTotalDescent:
			totalDescent, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			session.TotalDescent = &totalDescent
		case fieldnum.SessionTotalCycles:
			totalCycles, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			session.TotalCycles = &totalCycles
		case fieldnum.SessionTotalCalories:
			totalCalories, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			session.TotalCalories = &totalCalories
		case fieldnum.SessionAvgSpeed:
			avgSpeed, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			fAvgSpeed := float64(avgSpeed) / 1000
			session.AvgSpeed = &fAvgSpeed
		case fieldnum.SessionMaxSpeed:
			maxSpeed, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			fMaxSpeed := float64(maxSpeed) / 1000
			session.MaxSpeed = &fMaxSpeed
		case fieldnum.SessionAvgHeartRate:
			avgHeartRate, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			session.AvgHeartRate = &avgHeartRate
		case fieldnum.SessionMaxHeartRate:
			maxHeartRate, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			session.MaxHeartRate = &maxHeartRate
		case fieldnum.SessionAvgCadence:
			avgCadence, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			session.AvgCadence = &avgCadence
		case fieldnum.SessionMaxCadence:
			maxCadence, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			session.MaxCadence = &maxCadence
		case fieldnum.SessionAvgPower:
			avgPower, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			session.AvgPower = &avgPower
		case fieldnum.SessionMaxPower:
			maxPower, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			session.MaxPower = &maxPower
		case fieldnum.SessionAvgTemperature:
			avgTemperature, ok := field.Value.(int8)
			if !ok {
				continue
			}
			session.AvgTemperature = &avgTemperature
		case fieldnum.SessionMaxTemperature:
			maxTemperature, ok := field.Value.(int8)
			if !ok {
				continue
			}
			session.MaxTemperature = &maxTemperature
		case fieldnum.SessionAvgAltitude:
			avgAltitude, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			fAvgAltitude := (float64(avgAltitude) / 5) - 500
			session.AvgAltitude = &fAvgAltitude
		case fieldnum.SessionMaxAltitude:
			maxAltitude, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			fMaxAltitude := (float64(maxAltitude) / 5) - 500
			session.MaxAltitude = &fMaxAltitude
		}
	}

	return session
}

type Lap struct {
	Timestamp        time.Time `json:"timestamp,omitempty"`
	TotalElapsedTime *float64  `json:"totalElapsedTime,omitempty"` // Units: s;
	TotalTimerTime   *float64  `json:"totalTimerTime,omitempty"`   // Units: s;
	TotalDistance    *float64  `json:"totalDistance,omitempty"`    // Units: m;
	TotalAscent      *uint16   `json:"totalAscent,omitempty"`      // Units: m;
}

func NewLap(mesg proto.Message) Lap {
	lap := Lap{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.LapTimestamp:
			lap.Timestamp = datetime.ToTime(field.Value)
		case fieldnum.LapTotalElapsedTime:
			totalElapsedTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fTotalElapsedTime := float64(totalElapsedTime) / 1000
			lap.TotalElapsedTime = &fTotalElapsedTime
		case fieldnum.LapTotalTimerTime:
			totalTimerTime, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fTotalTimerTime := float64(totalTimerTime) / 1000
			lap.TotalTimerTime = &fTotalTimerTime
		case fieldnum.TotalsDistance:
			totalDistance, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fTotalDistance := float64(totalDistance) / 100
			lap.TotalDistance = &fTotalDistance
		case fieldnum.LapTotalAscent:
			totalAscent, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			lap.TotalAscent = &totalAscent
		}
	}

	return lap
}

type Record struct {
	Timestamp    time.Time `json:"timestamp,omitempty"`
	PositionLat  *float64  `json:"positionLat,omitempty"`  // Units: degrees;
	PositionLong *float64  `json:"positionLong,omitempty"` // Units: degrees;
	Altitude     *float64  `json:"altitude,omitempty"`     // Units: m;
	HeartRate    *uint8    `json:"heartRate,omitempty"`    // Units: bpm;
	Cadence      *uint8    `json:"cadence,omitempty"`      // Units: rpm;
	Distance     *float64  `json:"distance,omitempty"`     // Units: m;
	Speed        *float64  `json:"speed,omitempty"`        // Units: m/s;
	Power        *uint16   `json:"power,omitempty"`        // Units: watts;
}

func NewRecord(mesg proto.Message) Record {
	record := Record{}
	for i := range mesg.Fields {
		field := &mesg.Fields[i]

		switch field.Num {
		case fieldnum.RecordTimestamp:
			record.Timestamp = datetime.ToTime(field.Value)
		case fieldnum.RecordPositionLat:
			lat, ok := field.Value.(int32)
			if !ok {
				continue
			}
			latDeg := semicircles.ToDegrees(lat)
			record.PositionLat = &latDeg
		case fieldnum.RecordPositionLong:
			long, ok := field.Value.(int32)
			if !ok {
				continue
			}
			longDeg := semicircles.ToDegrees(long)
			record.PositionLong = &longDeg
		case fieldnum.RecordAltitude:
			altitude, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			faltitude := (float64(altitude) / 5) - 500
			record.Altitude = &faltitude
		case fieldnum.RecordHeartRate:
			heartRate, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			record.HeartRate = &heartRate
		case fieldnum.RecordCadence:
			cadence, ok := field.Value.(uint8)
			if !ok {
				continue
			}
			record.Cadence = &cadence
		case fieldnum.RecordDistance:
			distance, ok := field.Value.(uint32)
			if !ok {
				continue
			}
			fDistance := float64(distance) / 100
			record.Distance = &fDistance
		case fieldnum.RecordSpeed:
			speed, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			fSpeed := float64(speed) / 1000
			record.Speed = &fSpeed
		case fieldnum.RecordPower:
			power, ok := field.Value.(uint16)
			if !ok {
				continue
			}
			record.Power = &power
		}
	}

	return record
}
