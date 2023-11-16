package activity

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
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
	ActivityFile *ActivityFile `json:"activityFile"`
}

func (m DecodeResult) ToMap() map[string]any {
	return map[string]any{
		"err":          m.Err,
		"activityFile": m.ActivityFile.ToMap(),
	}
}

func decodeWorker(ctx context.Context, rc <-chan io.Reader, resc chan<- DecodeResult, wg *sync.WaitGroup) {
	defer wg.Done()

	r := <-rc

	lis := NewListener()
	dec := decoder.New(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
		decoder.WithNoComponentExpansion(),
		decoder.WithIgnoreChecksum(),
	)

	for {
		_, err := dec.DecodeWithContext(ctx)
		if err != nil {
			fmt.Printf("could not decode: %s\n", err)
			resc <- DecodeResult{Err: err.Error()}
			return
		}

		resc <- DecodeResult{
			ActivityFile: lis.ActivityFile(),
		}

		if !dec.Next() {
			break
		}
	}

	lis.WaitAndClose()
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

		var wg sync.WaitGroup
		wg.Add(input.Length())

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for i := 0; i < input.Length(); i++ {
			go decodeWorker(ctx, rc, resc, &wg)

			b := make([]byte, input.Index(i).Length())
			js.CopyBytesToGo(b, input.Index(i))
			rc <- bytes.NewReader(b)
		}

		decodeResults := make([]DecodeResult, 0, input.Length())
		var err string
		done := make(chan struct{})
		go func() {
			for decodeResult := range resc {
				if decodeResult.Err != "" {
					err = decodeResult.Err
					cancel()
					break
				}
				decodeResults = append(decodeResults, decodeResult)
			}
			close(done)
		}()

		wg.Wait()

		close(rc)
		close(resc)

		<-done

		if err != "" {
			return Result{Err: err}.ToMap()
		}

		slices.SortStableFunc(decodeResults, func(a, b DecodeResult) int {
			if a.ActivityFile.Creator.TimeCreated.Before(b.ActivityFile.Creator.TimeCreated) {
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
