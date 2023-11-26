package result

import (
	"time"

	"github.com/muktihari/openactivity-fit/activity"
)

type Decode struct {
	Err               error
	DecodeTook        time.Duration
	SerializationTook time.Duration
	TotalElapsed      time.Duration
	Activities        []activity.Activity
}

func (d Decode) ToMap() map[string]any {
	if d.Err != nil {
		return map[string]any{"err": d.Err.Error()}
	}

	begin := time.Now()

	activities := make([]any, len(d.Activities))
	for i := range d.Activities {
		activities[i] = d.Activities[i].ToMap()
	}

	d.SerializationTook = time.Since(begin)
	d.TotalElapsed = d.DecodeTook + d.SerializationTook

	m := map[string]any{
		"err":               nil,
		"activities":        activities,
		"decodeTook":        d.DecodeTook.Milliseconds(),
		"serializationTook": d.SerializationTook.Milliseconds(),
		"totalElapsed":      d.TotalElapsed.Milliseconds(),
	}

	return m
}

type DecodeWorker struct {
	Err      error
	Index    int
	Activity *activity.Activity
}
