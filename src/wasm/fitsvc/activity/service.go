package activity

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"
	"time"

	"github.com/muktihari/fit/decoder"
	"golang.org/x/exp/slices"
)

type Result struct {
	Err           string
	Took          int64
	DecodeResults []any
}

func (m Result) ToMap() map[string]any {
	return map[string]any{
		"err":           m.Err,
		"took":          m.Took,
		"decodeResults": m.DecodeResults,
	}
}

type DecodeResult struct {
	Err          string        `json:"err"`
	Feature      *Feature      `json:"feature"`
	ActivityFile *ActivityFile `json:"activityFile"`
}

func (m DecodeResult) ToMap() map[string]any {
	return map[string]any{
		"err":          m.Err,
		"feature":      m.Feature.ToMap(),
		"activityFile": m.ActivityFile.ToMap(),
	}
}

func decodeWorker(rc <-chan io.Reader, resc chan<- DecodeResult) {
	r := <-rc

	lis := NewListener()
	dec := decoder.New(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
		decoder.WithNoComponentExpansion(),
	)

	_, err := dec.Decode()
	if err != nil {
		fmt.Printf("could not decode: %s\n", err)
		resc <- DecodeResult{Err: err.Error()}
		return
	}

	lis.Wait()

	resc <- DecodeResult{
		Feature:      lis.Feature(),
		ActivityFile: lis.ActivityFile(),
	}
}

func Decode() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		begin := time.Now()

		input := args[0] // input is an Array<Uint8Array>
		if input.Length() == 0 {
			return []DecodeResult{{Err: "zero input"}}
		}

		rc := make(chan io.Reader, input.Length())
		resc := make(chan DecodeResult, input.Length())

		for i := 0; i < input.Length(); i++ {
			go decodeWorker(rc, resc)

			b := make([]byte, input.Index(i).Length())
			js.CopyBytesToGo(b, input.Index(i))
			rc <- bytes.NewReader(b)
		}

		decodeResults := make([]DecodeResult, input.Length())
		for i := 0; i < input.Length(); i++ {
			decodeResults[i] = <-resc

			if decodeResults[i].Err != "" {
				return Result{Err: decodeResults[i].Err}.ToMap()
			}
		}

		slices.SortStableFunc(decodeResults, func(a, b DecodeResult) int {
			if a.ActivityFile.FileId.TimeCreated.Before(b.ActivityFile.FileId.TimeCreated) {
				return -1
			}
			return 1
		})

		results := make([]any, 0, len(decodeResults))
		for i := range decodeResults {
			results = append(results, decodeResults[i].ToMap())
		}

		return Result{DecodeResults: results, Took: time.Since(begin).Milliseconds()}.ToMap()
	})
}
