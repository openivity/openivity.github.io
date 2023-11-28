package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/fit"
	"github.com/muktihari/openactivity-fit/service/result"
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
	Decode(ctx context.Context, rs []io.Reader) result.Decode
	Encode(ctx context.Context, activities []activity.Activity) result.Encode
	ManufacturerList() result.ManufacturerList
}

type service struct {
	fitService    activity.Service
	gpxService    activity.Service
	tcxService    activity.Service
	manufacturers map[uint16]fit.Manufacturer
}

func New(fitService, gpxService, tcxService activity.Service, manufacturers map[uint16]fit.Manufacturer) Service {
	return &service{
		fitService:    fitService,
		gpxService:    gpxService,
		tcxService:    tcxService,
		manufacturers: manufacturers,
	}
}

func (s *service) Decode(ctx context.Context, rs []io.Reader) result.Decode {
	begin := time.Now()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(rs))
	rc := make(chan io.Reader, len(rs))
	resc := make(chan result.DecodeWorker, len(rs))

	for i := range rs {
		i := i
		go s.decodeWorker(ctx, rc, resc, &wg, i)
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
		return result.Decode{Err: err}
	}

	slices.SortStableFunc(activities, func(a, b activity.Activity) int {
		if a.Creator.TimeCreated.Before(b.Creator.TimeCreated) {
			return -1
		}
		return 1
	})

	return result.Decode{
		DecodeTook: time.Since(begin),
		Activities: activities,
	}
}

func (s *service) decodeWorker(ctx context.Context, rc <-chan io.Reader, resc chan<- result.DecodeWorker, wg *sync.WaitGroup, index int) {
	defer wg.Done()

	activities, err := s.decode(ctx, <-rc)
	if err != nil {
		resc <- result.DecodeWorker{Err: err, Index: index}
		return
	}

	for i := range activities {
		resc <- result.DecodeWorker{Activity: &activities[i], Index: index}
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

func (s *service) Encode(ctx context.Context, activities []activity.Activity) result.Encode {
	return result.Encode{Err: fmt.Errorf("encode: Not yet implemented")}
}

func (s *service) ManufacturerList() result.ManufacturerList {
	manufacturers := make([]fit.Manufacturer, 0, len(s.manufacturers))
	for _, v := range s.manufacturers {
		manufacturers = append(manufacturers, v)
	}
	slices.SortFunc(manufacturers, func(a, b fit.Manufacturer) int {
		if a.Name < b.Name {
			return -1
		}
		return 1
	})
	return result.ManufacturerList{Manufacturers: manufacturers}
}
