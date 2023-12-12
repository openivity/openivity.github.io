package tcx

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/tcx/schema"
	"github.com/muktihari/openactivity-fit/kit"
	kxml "github.com/muktihari/openactivity-fit/kit/xml"
	"github.com/muktihari/openactivity-fit/preprocessor"
)

const (
	applicationName = "openitivy.github.io"
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

		sport := kit.FormatTitle(a.Activity.Sport)
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

		recordsByLap := make([][]*activity.Record, 0, len(a.Activity.Laps))
		for j := range a.Activity.Laps {
			activityLap := a.Activity.Laps[j]

			var lapRecordCount int
			for k := range activityLap.Tracks {
				lapRecordCount += len(activityLap.Tracks[k].Trackpoints)
			}
			lapRecords := make([]*activity.Record, 0, lapRecordCount)

			for k := range activityLap.Tracks { // flattening tracks-trackpoints
				for l := range activityLap.Tracks[k].Trackpoints {
					tp := &activityLap.Tracks[k].Trackpoints[l]
					rec := NewRecord(tp)
					lapRecords = append(lapRecords, rec)
				}
			}

			if len(lapRecords) == 0 {
				continue
			}

			records = append(records, lapRecords...)

			recordsByLap = append(recordsByLap, lapRecords)

			lap := &activity.Lap{
				StartTime:        activityLap.StartTime,
				TotalDistance:    activityLap.DistanceMeters,
				TotalCalories:    activityLap.Calories,
				TotalElapsedTime: activityLap.TotalTimeSeconds,
				AvgHeartRate:     activityLap.AverageHeartRateBpm,
				MaxHeartRate:     activityLap.MaximumHeartRateBpm,
			}

			laps = append(laps, lap)
		}

		// Preprocessing...
		s.preprocessor.CalculateDistanceAndSpeed(records)
		if activity.HasPace(sport) {
			s.preprocessor.CalculatePace(sport, records)
		}

		s.preprocessor.SmoothingElev(records)
		s.preprocessor.CalculateGrade(records)

		// We can only calculate laps' summary after preprocessing
		for i := range laps {
			lap := laps[i]
			lapFromRecords := activity.NewLapFromRecords(recordsByLap[i], sport)

			activity.CombineLap(lap, lapFromRecords)
		}

		if len(laps) == 0 {
			continue
		}

		session := activity.NewSessionFromLaps(laps, sport)
		if !a.Activity.ID.IsZero() {
			session.StartTime = a.Activity.ID
		}

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
		return nil, fmt.Errorf("tcx: %w", activity.ErrNoActivity)
	}

	act.Sessions = sessions

	s.preprocessor.SetSessionsWorkoutType(act.Sessions...)

	return []activity.Activity{*act}, nil
}

func (s *service) Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error) {
	bs := make([][]byte, len(activities))

	for i := range activities {
		tcx := s.convertActivityToTCX(&activities[i])
		b, err := kxml.Marshal(tcx)
		if err != nil {
			return nil, fmt.Errorf("could not marshal tcx: %w", err)
		}
		bs[i] = b
	}

	return bs, nil
}

func (s *service) convertActivityToTCX(act *activity.Activity) *schema.TCX {
	tcx := &schema.TCX{
		Author: &schema.Application{
			Name: applicationName,
		},
		Activities: make([]schema.ActivityList, 0, len(act.Sessions)),
	}

	for i := range act.Sessions {
		ses := act.Sessions[i]

		activityList := schema.ActivityList{
			Activity: &schema.Activity{
				ID:    ses.Timestamp,
				Sport: ses.Sport,
				Creator: &schema.Device{
					Name: act.Creator.Name,
				},
			},
		}

		if activityList.Activity.ID.IsZero() {
			activityList.Activity.ID = ses.StartTime
		}

		sesRecords := ses.Records
		for j := range ses.Laps {
			lap := ses.Laps[j]

			activityLap := schema.ActivityLap{
				StartTime:           lap.StartTime,
				TotalTimeSeconds:    lap.TotalElapsedTime,
				DistanceMeters:      lap.TotalDistance,
				MaximumSpeed:        lap.MaxSpeed,
				Calories:            lap.TotalCalories,
				AverageHeartRateBpm: lap.AvgHeartRate,
				MaximumHeartRateBpm: lap.MaxHeartRate,
				Cadence:             lap.AvgCadence,
			}

			track := schema.Track{}
			remainingRecords := make([]*activity.Record, 0)
			for k := range sesRecords {
				rec := sesRecords[k]

				if lap.IsBelongToThisLap(rec.Timestamp) {
					trackpoint := schema.Trackpoint{
						Time:           rec.Timestamp,
						AltitudeMeters: rec.Altitude,
						DistanceMeters: rec.Distance,
						HeartRateBpm:   rec.HeartRate,
						Cadence:        rec.Cadence,
					}

					if rec.PositionLat != nil && rec.PositionLong != nil {
						trackpoint.Position = &schema.Position{
							LatitudeDegrees:  *rec.PositionLat,
							LongitudeDegrees: *rec.PositionLong,
						}
					}

					if rec.Speed != nil {
						trackpoint.Extensions = &schema.TrackpointExtension{
							Speed: rec.Speed,
						}
					}
					track.Trackpoints = append(track.Trackpoints, trackpoint)
				} else {
					remainingRecords = append(remainingRecords, rec)
				}
			}
			sesRecords = remainingRecords

			activityLap.Tracks = []schema.Track{track}
			activityList.Activity.Laps = append(activityList.Activity.Laps, activityLap)
		}

		tcx.Activities = append(tcx.Activities, activityList)
	}

	return tcx
}
