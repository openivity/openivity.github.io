package fit

import (
	"context"
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

	activities := make([]activity.Activity, 0, 1) // In most of cases, 1 fit == 1 activity
	for dec.Next() {
		_, err := dec.DecodeWithContext(ctx)
		if err != nil {
			return nil, err
		}
		act := lis.Activity()
		if len(act.Laps) == 0 {
			// Some devices may not write a Lap even thought it's required for an Activity File
			// (ref: https://developer.garmin.com/fit/file-types/activity).
			// Let's treat a Session as a Lap.
			act.Laps = make([]*activity.Lap, len(act.Sessions))
			for i := range act.Sessions {
				act.Laps[i] = NewLapFromSession(act.Sessions[i])
			}
		}
		s.preprocessor.SmoothingElev(act.Records)
		s.preprocessor.CalculateGrade(act.Records)
		activities = append(activities, *act)
	}

	lis.WaitAndClose()

	return activities, nil
}
