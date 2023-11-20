package tcx

import (
	"context"
	"fmt"
	"io"

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
	return nil, fmt.Errorf("TODO: not yet implemented")
}
