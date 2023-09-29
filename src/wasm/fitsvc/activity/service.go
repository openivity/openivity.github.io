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
	Err          string        `json:"err"`
	Feature      *Feature      `json:"feature"`
	ActivityFile *ActivityFile `json:"activityFile"`
}

func (m Result) ToMap() map[string]any {
	defer func(begin time.Time) {
		fmt.Printf("wasm: serialization took: %s\n", time.Since(begin))
	}(time.Now())

	return map[string]any{
		"err":          m.Err,
		"feature":      m.Feature.ToMap(),
		"activityFile": m.ActivityFile.ToMap(),
	}
}

type Service interface {
	DecodeFit(rs []io.Reader) []Result
}

type service struct{}

func NewService() Service {
	s := &service{}

	return s
}

func (s *service) decodeWorker(rc <-chan io.Reader, resc chan<- Result) {
	r := <-rc

	lis := NewListener()
	dec := decoder.New(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
	)

	_, err := dec.Decode()
	if err != nil {
		fmt.Printf("could not decode: %s\n", err)
		return
	}

	lis.Wait()

	resc <- Result{
		Feature:      lis.Feature(),
		ActivityFile: lis.ActivityFile(),
	}
}

func (s *service) DecodeFit(rs []io.Reader) []Result {
	rc := make(chan io.Reader, len(rs))
	resc := make(chan Result, len(rs))

	for i := range rs {
		go s.decodeWorker(rc, resc)
		rc <- rs[i]
	}

	results := make([]Result, len(rs))
	for i := range rs {
		results[i] = <-resc
	}

	return results
}

var svc = NewLogging(NewService())

func Decode(this js.Value, args []js.Value) any {
	data := args[0]
	b := make([]byte, data.Length())
	js.CopyBytesToGo(b, data)

	results := svc.DecodeFit([]io.Reader{bytes.NewBuffer(b)})

	slices.SortStableFunc(results, func(a, b Result) int {
		if a.ActivityFile.FileId.TimeCreated.Before(b.ActivityFile.FileId.TimeCreated) {
			return -1
		} else if a.ActivityFile.FileId.TimeCreated.Equal(b.ActivityFile.FileId.TimeCreated) {
			return 0
		}
		return 1
	})

	result := results[0].ToMap()

	return result
}
