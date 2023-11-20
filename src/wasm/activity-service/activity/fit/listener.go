package fit

import (
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
)

type Listener struct {
	mesgc        chan proto.Message
	done         chan struct{}
	activityFile *activity.Activity
}

func NewListener() *Listener {
	l := &Listener{
		mesgc:        make(chan proto.Message, 1000),
		done:         make(chan struct{}),
		activityFile: &activity.Activity{},
	}
	go l.loop()
	return l
}

func (l *Listener) loop() {
	for mesg := range l.mesgc {
		switch mesg.Num {
		case mesgnum.FileId:
			l.activityFile.Creator = NewCreator(mesg)
		case mesgnum.Activity:
			l.activityFile.Timezone = CreateTimezone(mesg)
		case mesgnum.Session:
			l.activityFile.Sessions = append(l.activityFile.Sessions, NewSession(mesg))
		case mesgnum.Lap:
			l.activityFile.Laps = append(l.activityFile.Laps, NewLap(mesg))
		case mesgnum.Record:
			record := NewRecord(mesg)
			l.activityFile.Records = append(l.activityFile.Records, record)
			if record.PositionLat == nil || record.PositionLong == nil {
				continue
			}
		}
	}
	close(l.done)
}

func (l *Listener) OnMesg(mesg proto.Message) { l.mesgc <- mesg }

func (l *Listener) Activity() *activity.Activity {
	l.WaitAndClose()

	activityFile := *l.activityFile

	l.mesgc = make(chan proto.Message, 1000)
	l.done = make(chan struct{})
	l.activityFile = &activity.Activity{}

	go l.loop()

	return &activityFile
}

func (l *Listener) WaitAndClose() {
	close(l.mesgc)
	<-l.done
}
