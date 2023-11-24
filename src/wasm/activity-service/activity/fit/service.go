package fit

import (
	"context"
	"fmt"
	"io"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/openactivity-fit/activity"
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
	lis := NewListener()
	dec := decoder.New(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
		decoder.WithNoComponentExpansion(),
		decoder.WithIgnoreChecksum(),
	)

	defer lis.WaitAndClose()

	activities := make([]activity.Activity, 0, 1) // In most cases, 1 fit == 1 activity
	for dec.Next() {
		_, err := dec.DecodeWithContext(ctx)
		if err != nil {
			return nil, err
		}

		res := lis.Result()
		act := s.convertListenerResultToActivity(res)
		if act == nil {
			continue
		}
		activities = append(activities, *act)
	}

	if len(activities) == 0 {
		return nil, fmt.Errorf("fit: %w", activity.ErrNoActivity)
	}

	return activities, nil
}

func (s *service) convertListenerResultToActivity(result *ListenerResult) *activity.Activity {
	if len(result.Records) == 0 {
		return nil
	}

	s.sanitize(result)

	act := &activity.Activity{
		Creator:  *result.Creator,
		Timezone: result.Timezone,
	}

	result.Records = s.preprocessor.AggregateByTimestamp(result.Records)

	// Create Sessions and Laps if not exist. This could happen only if:
	//  - Fit file is truncated, so only some Records could be retrieved.
	//  - Some devices may not create Lap even though it's actually required for an Activity File.
	//    ref: https://developer.garmin.com/fit/file-types/activity
	if len(result.Sessions) == 0 {
		if len(result.Laps) == 0 {
			lap := activity.NewLapFromRecords(result.Records, activity.SportGeneric)
			result.Laps = append(result.Laps, lap)
		}

		ses := activity.NewSessionFromLaps(result.Laps, activity.SportGeneric)
		ses.Records = result.Records
		ses.Laps = result.Laps
		s.preprocessingRecords(ses.Records)
		act.Sessions = []*activity.Session{ses}

		return act
	}

	for i := range result.Sessions {
		ses := result.Sessions[i]

		result.Laps = ses.PutLaps(result.Laps...)
		result.Records = ses.PutRecords(result.Records...)

		if len(ses.Records) == 0 {
			continue
		}

		s.preprocessingRecords(ses.Records)
		s.finalizeSession(ses)

		act.Sessions = append(act.Sessions, ses)
	}

	if len(result.Records) == 0 {
		return act
	}

	// Handle remaining laps and records that don't belong any session
	// This could happen when file is truncated or file is not properly encoded.
	sport := activity.SportGeneric // Mark as Generic
	s.preprocessingRecords(result.Records)

	if len(result.Laps) != 0 {
		ses := activity.NewSessionFromLaps(result.Laps, sport)
		ses.Laps = result.Laps

		result.Records = ses.PutRecords(result.Records...)

		if len(ses.Records) != 0 {
			s.finalizeSession(ses)
			act.Sessions = append(act.Sessions, ses)
		}
	}

	// Handle remaining records that don't belong anywhere.
	if len(result.Records) != 0 {
		lap := activity.NewLapFromRecords(result.Records, sport)
		laps := []*activity.Lap{lap}

		ses := activity.NewSessionFromLaps(laps, sport)
		ses.Laps = laps

		act.Sessions = append(act.Sessions, ses)
	}

	return act
}

// sanitize removes any invalid item from given result.
func (s *service) sanitize(result *ListenerResult) {
	if len(result.Records) == 0 {
		return
	}

	validLaps := make([]*activity.Lap, 0)

	for i := range result.Laps {
		lap := result.Laps[i]

		// Timestamp, Start Time, and Total Elapsed Time are required fields for all summary messages.
		// If any of this field is missing, let's mark as invalid and use our own lap calculation later.
		if lap.Timestamp.IsZero() {
			continue
		}
		if lap.StartTime.IsZero() {
			continue
		}
		if lap.TotalElapsedTime == 0 {
			continue
		}

		// Most activity has 1 Lap and first Lap's StartTime should match first record's timestamp.
		// We should not try to accommodate all bad encoding practices and this guard should be sufficient for most cases.
		if i == 0 && !lap.StartTime.Equal(result.Records[0].Timestamp) {
			continue
		}

		validLaps = append(validLaps, lap)
	}

	result.Laps = validLaps
}

// preprocessingRecords pre-processes records per session since 1 session corresponds to 1 sport.
// We should not process different sports as one.
func (s *service) preprocessingRecords(records []*activity.Record) {
	s.preprocessor.CalculateDistanceAndSpeed(records)
	s.preprocessor.SmoothingElev(records)
	s.preprocessor.CalculateGrade(records)
	if activity.HasPace(activity.SportGeneric) {
		s.preprocessor.CalculatePace(activity.SportGeneric, records)
	}
}

// finalizeSession finalises session by creating lap if not exist and calculating its summary as well as calculating session's summary.
func (s *service) finalizeSession(ses *activity.Session) {
	if len(ses.Laps) == 0 {
		lap := activity.NewLapFromRecords(ses.Records, ses.Sport)
		ses.Laps = append(ses.Laps, lap)
		sesFromLaps := activity.NewSessionFromLaps(ses.Laps, ses.Sport)
		activity.CombineSession(ses, sesFromLaps)
		return
	}

	remainingRecords := ses.Records
	for j := range ses.Laps {
		lap := ses.Laps[j]

		lapRecords := make([]*activity.Record, 0)
		remainingLapRecords := make([]*activity.Record, 0)
		for k := range remainingRecords {
			rec := remainingRecords[k]

			if lap.IsBelongToThisLap(rec.Timestamp) {
				lapRecords = append(lapRecords, rec)
			} else {
				remainingLapRecords = append(remainingLapRecords, rec)
			}
		}
		remainingRecords = remainingLapRecords

		lapFromRecords := activity.NewLapFromRecords(lapRecords, ses.Sport)
		activity.CombineLap(lap, lapFromRecords)
	}

	// Handle remaining records that don't belong to any lap.
	if len(remainingRecords) != 0 {
		lap := activity.NewLapFromRecords(remainingRecords, ses.Sport)
		ses.Laps = append(ses.Laps, lap)
	}

	sesFromLaps := activity.NewSessionFromLaps(ses.Laps, ses.Sport)
	activity.CombineSession(ses, sesFromLaps)
}
