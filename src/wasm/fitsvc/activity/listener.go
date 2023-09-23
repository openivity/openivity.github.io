package activity

import (
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
)

func NewListener() *Listener {
	l := &Listener{
		mesgc:        make(chan proto.Message, 1000),
		done:         make(chan struct{}),
		coordinates:  make([]any, 0),
		activityFile: &ActivityFile{},
	}
	go l.loop()
	return l
}

type Listener struct {
	mesgc        chan proto.Message
	done         chan struct{}
	coordinates  []any
	activityFile *ActivityFile
}

func (l *Listener) loop() {
	for mesg := range l.mesgc {
		switch mesg.Num {
		case mesgnum.FileId:
			l.activityFile.FileId = NewFileId(mesg)
		case mesgnum.Session:
			l.activityFile.Sessions = append(l.activityFile.Sessions, NewSession(mesg))
		// case mesgnum.Lap:
		// l.activityFile.Laps = append(l.activityFile.Laps, NewLap(mesg))
		case mesgnum.Record:
			record := NewRecord(mesg)
			l.activityFile.Records = append(l.activityFile.Records, record)
			if record["positionLong"] == nil || record["positionLat"] == nil {
				continue
			}
			l.coordinates = append(l.coordinates, []any{
				record["positionLong"],
				record["positionLat"],
			})

			// long, ok := mesg.FieldByNum(fieldnum.RecordPositionLong)
			// if !ok {
			// 	continue
			// }

			// lat, ok := mesg.FieldByNum(fieldnum.RecordPositionLat)
			// if !ok {
			// 	continue
			// }

			// l.coordinates = append(l.coordinates, []any{
			// 	semicircles.ToDegrees(typeconv.ToSint32[int32](long.Value)),
			// 	semicircles.ToDegrees(typeconv.ToSint32[int32](lat.Value)),
			// })

		}
	}
	close(l.done)
}

func (l *Listener) Wait() {
	close(l.mesgc)
	<-l.done
}

func (l *Listener) OnMesg(mesg proto.Message) { l.mesgc <- mesg }

func (l *Listener) Feature() *Feature {
	return &Feature{
		Type: "Feature",
		Geometry: Geometry{
			Type:        "LineString",
			Coordinates: l.coordinates,
		},
	}
}

func (l *Listener) ActivityFile() *ActivityFile {
	return l.activityFile
}
