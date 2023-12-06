package fit

import (
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

type ListenerResult struct {
	Creator  *activity.Creator
	Timezone int
	Records  []*activity.Record
	Laps     []*activity.Lap
	Sessions []*activity.Session
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
			if len(l.sessions) > 0 && l.sessions[len(l.sessions)-1].EndTime.IsZero() {
				l.sessions[len(l.sessions)-1].EndTime = ses.StartTime
			}
			l.sessions = append(l.sessions, ses)
		case mesgnum.Lap:
			lap := NewLap(mesg)
			if len(l.laps) > 0 && l.laps[len(l.laps)-1].EndTime.IsZero() {
				l.laps[len(l.laps)-1].EndTime = lap.StartTime
			}
			l.laps = append(l.laps, lap)
		case mesgnum.Record:
			record := NewRecord(mesg)
			l.records = append(l.records, record)
		}
	}
	close(l.done)
}

func (l *Listener) OnMesg(mesg proto.Message) { l.mesgc <- mesg }

func (l *Listener) Result() *ListenerResult {
	l.WaitAndClose()

	r := &ListenerResult{
		Creator:  l.creator,
		Timezone: l.timezone,
		Records:  l.records,
		Laps:     l.laps,
		Sessions: l.sessions,
	}

	l.reset()

	go l.loop()

	return r
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
