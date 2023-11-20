package service

import (
	"time"

	"github.com/muktihari/openactivity-fit/activity"
)

type Result struct {
	Err        error
	Took       time.Duration
	Activities []activity.Activity
}

func (r Result) ToMap() map[string]any {
	if r.Err != nil {
		return map[string]any{"err": r.Err.Error()}
	}

	activities := make([]any, len(r.Activities))
	for i := range r.Activities {
		activities[i] = r.Activities[i].ToMap()
	}

	m := map[string]any{
		"err":        nil,
		"took":       r.Took.Milliseconds(),
		"activities": activities,
	}

	return m
}

type DecodeResult struct {
	Err      error
	Index    int
	Activity *activity.Activity
}
