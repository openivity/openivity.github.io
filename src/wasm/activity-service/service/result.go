package service

import (
	"time"

	"github.com/muktihari/openactivity-fit/activity"
)

type Result struct {
	Err               error
	DecodeTook        time.Duration
	SerializationTook time.Duration
	TotalElapsed      time.Duration
	Activities        []activity.Activity
}

func (r Result) ToMap() map[string]any {
	if r.Err != nil {
		return map[string]any{"err": r.Err.Error()}
	}

	begin := time.Now()

	activities := make([]any, len(r.Activities))
	for i := range r.Activities {
		activities[i] = r.Activities[i].ToMap()
	}

	r.SerializationTook = time.Since(begin)
	r.TotalElapsed = r.DecodeTook + r.SerializationTook

	m := map[string]any{
		"err":               nil,
		"activities":        activities,
		"decodeTook":        r.DecodeTook.Milliseconds(),
		"serializationTook": r.SerializationTook.Milliseconds(),
		"totalElapsed":      r.TotalElapsed.Milliseconds(),
	}

	return m
}

type DecodeResult struct {
	Err      error
	Index    int
	Activity *activity.Activity
}
