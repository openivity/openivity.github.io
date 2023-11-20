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

	activities := make([]activity.Activity, 0, 1) // In most of cases, 1 fit == 1 activity
	for dec.Next() {
		_, err := dec.DecodeWithContext(ctx)
		if err != nil {
			return nil, err
		}

		act := lis.Activity()

		sessions := make([]*activity.Session, 0, len(act.Sessions))
		for i := range act.Sessions {
			ses := act.Sessions[i]

			if len(ses.Records) == 0 {
				continue
			}

			if activity.HasPace(ses.Sport) {
				s.preprocessor.CalculatePace(ses.Sport, ses.Records)
			}

			s.preprocessor.SmoothingElev(ses.Records)
			s.preprocessor.CalculateGrade(ses.Records)

			sessions = append(sessions, ses)
		}

		act.Sessions = sessions

		if len(act.Sessions) == 0 {
			continue
		}

		activities = append(activities, *act)
	}

	if len(activities) == 0 {
		return nil, fmt.Errorf("fit has no activity data")
	}

	return activities, nil
}
