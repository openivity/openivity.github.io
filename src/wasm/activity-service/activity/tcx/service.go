package tcx

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/tcx/schema"
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

	var tcx schema.TCX
	if err := tcx.UnmarshalXML(dec, *start); err != nil {
		return nil, err
	}

	act := new(activity.Activity)
	if tcx.Author != nil {
		act.Creator.Name = tcx.Author.Name
	}
	if len(tcx.Activities) > 0 && tcx.Activities[0].Activity != nil {
		act.Creator.TimeCreated = tcx.Activities[0].Activity.ID
	}

	sessions := make([]*activity.Session, 0, len(tcx.Activities))

	for i := range tcx.Activities {
		a := tcx.Activities[i]
		if a.Activity == nil {
			continue
		}

		sport := activity.FormatTitle(a.Activity.Sport)
		if sport == "" || sport == "Other" {
			sport = activity.SportGeneric
		}

		laps := make([]*activity.Lap, 0, len(a.Activity.Laps))

		var recordCount int
		for j := range a.Activity.Laps {
			for k := range a.Activity.Laps[j].Tracks {
				recordCount += len(a.Activity.Laps[j].Tracks[k].Trackpoints)
			}
		}

		records := make([]*activity.Record, 0, recordCount)

		for j := range a.Activity.Laps {
			activityLap := a.Activity.Laps[j]

			var lapRecordCount int
			for k := range activityLap.Tracks {
				lapRecordCount += len(activityLap.Tracks[k].Trackpoints)
			}
			lapRecords := make([]*activity.Record, 0, lapRecordCount)

			var prevRec *activity.Record
			for k := range activityLap.Tracks { // flattening tracks-trackpoints
				for l := range activityLap.Tracks[k].Trackpoints {
					tp := activityLap.Tracks[k].Trackpoints[l]

					rec := NewRecord(&tp, prevRec)
					prevRec = rec

					lapRecords = append(lapRecords, rec)
				}
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

			records = append(records, lapRecords...)

			lap := activity.NewLapFromRecords(lapRecords, sport)
			lap.StartTime = activityLap.StartTime
			lap.TotalDistance = activityLap.DistanceMeters
			lap.TotalCalories = activityLap.Calories
			lap.TotalElapsedTime = activityLap.TotalTimeSeconds

			if activityLap.AverageHeartRateBpm != nil {
				lap.AvgHeartRate = activityLap.AverageHeartRateBpm
			}
			if activityLap.MaximumHeartRateBpm != nil {
				lap.MaxHeartRate = activityLap.MaximumHeartRateBpm
			}

			laps = append(laps, lap)
		}

		if len(laps) == 0 {
			continue
		}

		session := activity.NewSessionFromLaps(laps, sport)
		session.StartTime = a.Activity.ID

		session.Laps = laps
		session.Records = records

		sessions = append(sessions, session)

		if act.Creator.TimeCreated.IsZero() {
			act.Creator.TimeCreated = session.StartTime
			act.Creator.Name = a.Activity.Creator.Name
			act.Creator.Product = &a.Activity.Creator.ProductID
		}
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("tcx has no activity data")
	}

	act.Sessions = sessions

	return []activity.Activity{*act}, nil
}
