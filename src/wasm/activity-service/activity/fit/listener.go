package fit

import (
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/kit"
)

const bufferSize = 1000

type Listener struct {
	poolc  chan proto.Message // pool of reusable objects to minimalize slice allocations. do not close this channel.
	mesgc  chan proto.Message // queue messages to be processed concurrently.
	done   chan struct{}
	active bool

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
	l := &Listener{active: true}
	l.reset()

	l.poolc = make(chan proto.Message, bufferSize)
	for i := 0; i < bufferSize; i++ {
		l.poolc <- proto.Message{} // fill pool with empty message and alloc its slices as needed.
	}

	go l.loop()
	return l
}

func (l *Listener) loop() {
	for mesg := range l.mesgc {
		l.processMesg(mesg)
		l.poolc <- mesg // put the message back to the pool to be recycled.
	}
	close(l.done)
}

func (l *Listener) processMesg(mesg proto.Message) {
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

func (l *Listener) OnMesg(mesg proto.Message) {
	if !l.active {
		l.reset()
		go l.loop()
		l.active = true
	}
	l.mesgc <- l.prep(mesg)
}

func (l *Listener) prep(mesg proto.Message) proto.Message {
	m := <-l.poolc

	if cap(m.Fields) < len(mesg.Fields) {
		m.Fields = make([]proto.Field, len(mesg.Fields))
	}
	copy(m.Fields, mesg.Fields)
	mesg.Fields = m.Fields[:len(mesg.Fields)]

	if mesg.DeveloperFields == nil {
		return mesg
	}

	if cap(m.DeveloperFields) < len(mesg.DeveloperFields) {
		m.DeveloperFields = make([]proto.DeveloperField, len(mesg.DeveloperFields))
	}
	copy(m.DeveloperFields, mesg.DeveloperFields)
	mesg.DeveloperFields = m.DeveloperFields[:len(mesg.DeveloperFields)]

	return mesg
}

func (l *Listener) Result() *ListenerResult {
	l.WaitAndClose()

	return &ListenerResult{
		Creator:  l.creator,
		Timezone: l.timezone,
		Records:  l.records,
		Laps:     l.laps,
		Sessions: l.sessions,
	}
}

func (l *Listener) reset() {
	l.mesgc = make(chan proto.Message, bufferSize)
	l.done = make(chan struct{})

	l.creator = nil
	l.timezone = 0
	l.records = nil
	l.laps = nil
	l.sessions = nil
}

func (l *Listener) WaitAndClose() {
	if !l.active {
		return
	}
	close(l.mesgc)
	<-l.done
	l.active = false
}
