package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/muktihari/openactivity-fit/activity"
	"golang.org/x/exp/slices"
)

var (
	ErrFileTypeUnsupported = errors.New("file type is unsupported")
)

type FileType byte

const (
	FileTypeUnsupported FileType = iota
	FileTypeFIT
	FileTypeGPX
	FileTypeTCX
)

type Service interface {
	Decode(ctx context.Context, rs []io.Reader) Result
}

type service struct {
	fitService activity.Service
	gpxService activity.Service
	tcxService activity.Service
}

func New(fitService, gpxService, tcxService activity.Service) Service {
	return &service{
		fitService: fitService,
		gpxService: gpxService,
		tcxService: tcxService,
	}
}

func (s *service) Decode(ctx context.Context, rs []io.Reader) Result {
	begin := time.Now()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(rs))
	rc := make(chan io.Reader, len(rs))
	resc := make(chan DecodeResult, len(rs))

	for i := range rs {
		i := i
		go s.worker(ctx, rc, resc, &wg, i)
	}

	for i := range rs {
		rc <- rs[i]
	}

	activities := make([]activity.Activity, 0, len(rs))
	var err error
	done := make(chan struct{})
	go func() {
		for decodeResult := range resc {
			if decodeResult.Err != nil {
				err = fmt.Errorf("[%d]: %w", decodeResult.Index, decodeResult.Err)
				cancel()
				break
			}
			activities = append(activities, *decodeResult.Activity)
		}
		close(done)
	}()

	wg.Wait()

	close(rc)
	close(resc)

	<-done

	if err != nil {
		return Result{Err: err}
	}

	slices.SortStableFunc(activities, func(a, b activity.Activity) int {
		if a.Creator.TimeCreated.Before(b.Creator.TimeCreated) {
			return -1
		}
		return 1
	})

	return Result{
		DecodeTook: time.Since(begin),
		Activities: activities,
	}
}

func (s *service) worker(ctx context.Context, rc <-chan io.Reader, resc chan<- DecodeResult, wg *sync.WaitGroup, index int) {
	defer wg.Done()

	activities, err := s.decode(ctx, <-rc)
	if err != nil {
		resc <- DecodeResult{Err: err, Index: index}
		return
	}

	for i := range activities {
		resc <- DecodeResult{Activity: &activities[i], Index: index}
	}
}

func (s *service) decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	fileType, err := s.readType(r)
	if err != nil {
		return nil, err
	}

	switch fileType {
	case FileTypeFIT:
		return s.fitService.Decode(ctx, r)
	case FileTypeGPX:
		return s.gpxService.Decode(ctx, r)
	case FileTypeTCX:
		return s.tcxService.Decode(ctx, r)
	default:
		return nil, ErrFileTypeUnsupported
	}
}

func (s *service) readType(r io.Reader) (FileType, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return FileTypeUnsupported, err
	}
	return FileType(b[0]), nil
}
