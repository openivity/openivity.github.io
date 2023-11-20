package gpx

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/gpx/schema"
	"github.com/muktihari/openactivity-fit/preprocessor"
)

var _ activity.Service = &service{}

type service struct {
	preprocessor *preprocessor.Preprocessor
}

func NewService(preproc *preprocessor.Preprocessor) activity.Service {
	return &service{preprocessor: preproc}
}

func (s *service) Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	dec := xml.NewDecoder(r)

	// NOTE: We manually define xml.Unmarshaler to void the reflection used
	//       in default decoder, which is particularly slow expecially in WASM.

	// Find start element
	var start *xml.StartElement
	if start == nil {
		for {
			tok, err := dec.Token()
			if err != nil {
				return nil, err
			}
			if t, ok := tok.(xml.StartElement); ok {
				start = &t
				break
			}
		}
	}

	if start == nil { // In case we got invalid xml, avoid panic.
		return nil, fmt.Errorf("not a valid xml")
	}

	var gpx schema.GPX
	if err := gpx.UnmarshalXML(dec, *start); err != nil {
		return nil, err
	}

	act := new(activity.Activity)
	act.Creator.Name = gpx.Creator
	act.Creator.TimeCreated = gpx.Metadata.Time

	lapCount := 0
	recordCount := 0
	for i := range gpx.Tracks {
		trk := gpx.Tracks[i]
		lapCount += len(trk.TrackSegments)

		for j := range trk.TrackSegments {
			trkseg := trk.TrackSegments[j]
			recordCount += len(trkseg.Trackpoints)
		}
	}

	records := make([]*activity.Record, 0, recordCount)
	laps := make([]*activity.Lap, 0, lapCount)
	sessions := make([]*activity.Session, 0, len(gpx.Tracks))

	for i := range gpx.Tracks { // Sessions
		trk := gpx.Tracks[i]
		sport := activity.FormatSport(trk.Type)

		sessionLaps := make([]*activity.Lap, 0, len(trk.TrackSegments))
		for j := range trk.TrackSegments { // Laps
			trkseg := trk.TrackSegments[j]

			var prevRec *activity.Record
			lapRecords := make([]*activity.Record, 0, len(trkseg.Trackpoints))
			for k := range trkseg.Trackpoints { // Records
				trkpt := trkseg.Trackpoints[k]

				rec := NewRecord(&trkpt, prevRec)
				prevRec = rec

				lapRecords = append(lapRecords, rec)
			}

			// Preprocessing...
			if activity.HasPace(sport) {
				s.preprocessor.CalculatePace(sport, lapRecords)
			}
			s.preprocessor.SmoothingElev(lapRecords)
			s.preprocessor.CalculateGrade(lapRecords)

			lap := NewLap(lapRecords, sport)
			records = append(records, lapRecords...)
			sessionLaps = append(sessionLaps, lap)
		}

		laps = append(laps, sessionLaps...)
		session := NewSession(sessionLaps, sport)
		sessions = append(sessions, session)
	}

	act.Sessions = sessions
	act.Laps = laps
	act.Records = records

	return []activity.Activity{*act}, nil
}
