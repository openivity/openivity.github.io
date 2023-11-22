package fit

import (
	"time"

	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/kit"
)

type Listener struct {
	mesgc chan proto.Message
	done  chan struct{}

	creator  *activity.Creator
	timezone int
	records  []*activity.Record
	laps     []*activity.Lap
	sessions []*activity.Session
}

func NewListener() *Listener {
	l := &Listener{
		mesgc: make(chan proto.Message, 1000),
		done:  make(chan struct{}),
	}
	go l.loop()
	return l
}

func (l *Listener) loop() {
	for mesg := range l.mesgc {
		switch mesg.Num {
		case mesgnum.FileId:
			l.creator = kit.Ptr(NewCreator(mesg))
		case mesgnum.Activity:
			l.timezone = CreateTimezone(mesg)
		case mesgnum.Session:
			ses := NewSession(mesg)
			if len(l.sessions) > 0 {
				l.sessions[len(l.sessions)-1].EndTime = ses.StartTime
			}
			l.sessions = append(l.sessions, ses)
		case mesgnum.Lap:
			l.laps = append(l.laps, NewLap(mesg))
		case mesgnum.Record:
			record := NewRecord(mesg)
			l.records = append(l.records, record)
		}
	}
	close(l.done)
}

func (l *Listener) OnMesg(mesg proto.Message) { l.mesgc <- mesg }

func (l *Listener) Activity() *activity.Activity {
	l.WaitAndClose()

	act := &activity.Activity{
		Creator:  *l.creator,
		Timezone: l.timezone,
	}

	// Create Sessions and Laps if not exist. This could happen only if:
	//  - Fit file is truncated, so only some Records could be retrieved.
	//  - Some devices may not create Lap even though it's actually required for an Activity File.
	//    ref: https://developer.garmin.com/fit/file-types/activity
	if len(l.sessions) == 0 {
		if len(l.laps) == 0 {
			lap := activity.NewLapFromRecords(l.records, activity.SportUnknown)
			if lap != nil {
				l.laps = append(l.laps, lap)
			}
			ses := activity.NewSessionFromLaps(l.laps, activity.SportUnknown)
			if ses != nil {
				l.sessions = append(l.sessions, ses)
			}
		} else {
			ses := activity.NewSessionFromLaps(l.laps, activity.SportUnknown)
			if ses != nil {
				l.sessions = append(l.sessions, ses)
			}
		}
	} else {
		if len(l.laps) == 0 {
			for i := range l.sessions { // 1 Session = 1 Lap
				lap := activity.NewLapFromSession(l.sessions[i])
				if lap != nil {
					l.laps = append(l.laps, lap)
				}
			}
		}
	}

	for i := range l.sessions {
		ses := l.sessions[i]

		for j := range l.laps {
			lap := l.laps[j]
			if isBelongToSession(lap.Timestamp, ses.StartTime, ses.EndTime) {
				ses.Laps = append(ses.Laps, lap)
			}
		}

		var lastRecordTimestamp time.Time
		for j := range l.records {
			rec := l.records[j]
			if isBelongToSession(rec.Timestamp, ses.StartTime, ses.EndTime) {
				ses.Records = append(ses.Records, rec)
				lastRecordTimestamp = rec.Timestamp
			}
		}

		l.sessions[i].EndTime = lastRecordTimestamp
	}

	act.Sessions = l.sessions

	l.reset()

	go l.loop()

	return act
}

func (l *Listener) reset() {
	l.mesgc = make(chan proto.Message, 1000)
	l.done = make(chan struct{})

	l.creator = nil
	l.timezone = 0
	l.records = nil
	l.laps = nil
	l.sessions = nil
}

func (l *Listener) WaitAndClose() {
	close(l.mesgc)
	<-l.done
}

func isBelongToSession(timestamp, sessionStartTime, sessionEndTime time.Time) bool {
	if timestamp.Equal(sessionStartTime) {
		return true
	}
	if sessionEndTime.IsZero() && timestamp.After(sessionStartTime) { // Last Session has no EndTime
		return true
	} else if timestamp.After(sessionStartTime) && timestamp.Before(sessionEndTime) {
		return true
	}
	return false
}
