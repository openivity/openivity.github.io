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

	sessions := make([]*activity.Session, 0, len(gpx.Tracks))

	for i := range gpx.Tracks { // Sessions
		trk := gpx.Tracks[i]

		sport := activity.SportUnknown
		if trk.Type != "" {
			sport = activity.FormatTitle(trk.Type)
		}

		laps := make([]*activity.Lap, 0, len(trk.TrackSegments))

		var recordCount int
		for i := range trk.TrackSegments {
			recordCount += len(trk.TrackSegments[i].Trackpoints)
		}
		records := make([]*activity.Record, 0, recordCount)

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

			if len(lapRecords) == 0 {
				continue
			}

			// Preprocessing...
			if activity.HasPace(sport) {
				s.preprocessor.CalculatePace(sport, lapRecords)
			}
			s.preprocessor.SmoothingElev(lapRecords)
			s.preprocessor.CalculateGrade(lapRecords)

			lap := activity.NewLapFromRecords(lapRecords, sport)
			laps = append(laps, lap)

			records = append(records, lapRecords...)
		}

		if len(laps) == 0 {
			continue
		}

		session := activity.NewSessionFromLaps(laps, sport)
		session.Records = records
		session.Laps = laps

		sessions = append(sessions, session)

		if act.Creator.TimeCreated.IsZero() {
			act.Creator.TimeCreated = session.StartTime
		}
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("gpx has no activity data")
	}

	act.Sessions = sessions

	return []activity.Activity{*act}, nil
}
