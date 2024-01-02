package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/fit"
	"github.com/muktihari/openactivity-fit/kit"
	"github.com/muktihari/openactivity-fit/service/result"
	"github.com/muktihari/openactivity-fit/service/spec"
	"golang.org/x/exp/slices"
)

var (
	ErrFileTypeUnsupported = errors.New("file type is unsupported")
)

type Service interface {
	Decode(ctx context.Context, rs []io.Reader) result.Decode
	Encode(ctx context.Context, encodeSpec spec.Encode) result.Encode
	ManufacturerList() result.ManufacturerList
	SportList() result.SportList
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
		if a.Creator.TimeCreated.After(b.Creator.TimeCreated) {
			return 1
		}
		// TimeCreated is equal, compare by first record's timestamp of each activity.
		firstTimestampA := firstNonZeroTimestamp(&a)
		firstTimestampB := firstNonZeroTimestamp(&b)
		if firstTimestampA.Before(firstTimestampB) {
			return -1
		}
		if firstTimestampA.After(firstTimestampB) {
			return 1
		}
		return 0
	})

	return result.Decode{
		DecodeTook: time.Since(begin),
		Activities: activities,
	}
}

func firstNonZeroTimestamp(act *activity.Activity) time.Time {
	for i := range act.Sessions {
		for j := range act.Sessions[i].Records {
			if !act.Sessions[i].Records[j].Timestamp.IsZero() {
				return act.Sessions[i].Records[j].Timestamp
			}
		}
	}
	return time.Time{}
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
	case spec.FileTypeFIT:
		return s.fitService.Decode(ctx, r)
	case spec.FileTypeGPX:
		return s.gpxService.Decode(ctx, r)
	case spec.FileTypeTCX:
		return s.tcxService.Decode(ctx, r)
	default:
		return nil, ErrFileTypeUnsupported
	}
}

func (s *service) readType(r io.Reader) (spec.FileType, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return spec.FileTypeUnsupported, err
	}
	return spec.FileType(b[0]), nil
}

func (s *service) Encode(ctx context.Context, encodeSpec spec.Encode) result.Encode {
	begin := time.Now()

	activities, err := s.preprocessEncode(ctx, encodeSpec)
	if err != nil {
		return result.Encode{Err: err}
	}

	var bs [][]byte
	switch encodeSpec.TargetFileType {
	case spec.FileTypeFIT:
		bs, err = s.fitService.Encode(ctx, activities)
	case spec.FileTypeGPX:
		bs, err = s.gpxService.Encode(ctx, activities)
	case spec.FileTypeTCX:
		bs, err = s.tcxService.Encode(ctx, activities)
	default:
		return result.Encode{Err: fmt.Errorf("encode: invalid filetype")}
	}

	return result.Encode{
		FileName:   fmt.Sprintf("openivity-%d-%s", begin.Unix(), encodeSpec.ToolMode),
		FileType:   encodeSpec.TargetFileType.String(),
		FilesBytes: bs,
		Err:        err,
		EncodeTook: time.Since(begin),
	}
}

func (s *service) preprocessEncode(ctx context.Context, encodeSpec spec.Encode) ([]activity.Activity, error) {
	if encodeSpec.ToolMode == spec.ToolModeUnknown {
		return nil, fmt.Errorf("encode mode '%v' not recognized", encodeSpec.ToolMode)
	}

	activities := encodeSpec.Activities
	if len(activities) == 0 {
		return nil, fmt.Errorf("no activity is retrieved")
	}

	if encodeSpec.TargetFileType == spec.FileTypeFIT {
		if _, ok := s.manufacturers[encodeSpec.ManufacturerID]; !ok {
			return nil, fmt.Errorf("manufacturer %d does not exist", encodeSpec.ManufacturerID)
		}
	}

	removeFields := make(map[string]struct{})
	for _, v := range encodeSpec.RemoveFields {
		removeFields[v] = struct{}{}
	}

	// Preprocess data before encoding
	var validActivityCounter int
	for i := range activities {
		activity := &activities[i]
		n := len(activities[i].Sessions) + i // markers is based on session across activities.

		if err := s.concealGPSPositions(activity, encodeSpec.ConcealMarkers[i:n]); err != nil {
			return nil, err
		}
		if err := s.trimRecords(activity, encodeSpec.TrimMarkers[i:n]); err != nil {
			return nil, err
		}

		if len(activity.Sessions) == 0 {
			continue
		}

		validActivityCounter++
		s.changeSport(activities, encodeSpec.Sports)
		s.RemoveFields(activity, removeFields)
	}

	if validActivityCounter == 0 {
		return nil, fmt.Errorf("no activity data after processed")
	}

	var newActivities []activity.Activity
	switch encodeSpec.ToolMode {
	case spec.ToolModeEdit:
		for i := range activities {
			activities[i].Creator.Manufacturer = &encodeSpec.ManufacturerID
			activities[i].Creator.Product = &encodeSpec.ProductID
			activities[i].Creator.Name = encodeSpec.DeviceName
		}
		newActivities = activities
	case spec.ToolModeCombine:
		newActivity := s.combineActivity(encodeSpec.Activities, encodeSpec.ManufacturerID, encodeSpec.ProductID)
		newActivities = []activity.Activity{newActivity}
	case spec.ToolModeSplitPerSession:
		newActivities = s.splitActivityPerSession(activities, encodeSpec.ManufacturerID, encodeSpec.ProductID)
	}

	return newActivities, nil
}

func (s *service) combineActivity(activities []activity.Activity, manufacturer, product uint16) activity.Activity {
	newActivity := activity.Activity{
		Creator: activity.Creator{
			Manufacturer: &manufacturer,
			Product:      &product,
			TimeCreated:  activities[0].Creator.TimeCreated,
		},
		Timezone: activities[0].Timezone,
		Sessions: activities[0].Sessions,
	}

	lastDistance := getLastDistanceOfRecords(newActivity.Sessions[0].Records)

	for i := 1; i < len(activities); i++ {
		cur := &activities[i]

		newActLastSes := newActivity.Sessions[len(newActivity.Sessions)-1]
		curActFirstSes := cur.Sessions[0]

		if newActLastSes.Sport != curActFirstSes.Sport { // Sport is not match, append as it is
			newActivity.Sessions = append(newActivity.Sessions, cur.Sessions...)
			continue
		}

		// Adjust distance before combine
		for j := range curActFirstSes.Records {
			rec := curActFirstSes.Records[j]
			if rec.Distance != nil {
				*rec.Distance += lastDistance
			}
		}

		// Combine records and laps to newActivity's last session
		newActLastSes.Records = append(newActLastSes.Records, curActFirstSes.Records...)
		newActLastSes.Laps = append(newActLastSes.Laps, curActFirstSes.Laps...)

		lastDistance = getLastDistanceOfRecords(newActLastSes.Records)

		if len(cur.Sessions) > 1 {
			newActivity.Sessions = append(newActivity.Sessions, cur.Sessions[1:]...)
		}

		// Update summary
		activity.AccumulateSession(newActLastSes, curActFirstSes)
	}

	return newActivity
}

func getLastDistanceOfRecords(records []*activity.Record) float64 {
	for i := len(records) - 1; i >= 0; i-- {
		rec := records[i]
		if rec.Distance != nil {
			return *rec.Distance
		}
	}
	return 0
}

func (s *service) splitActivityPerSession(activities []activity.Activity, manufacturer, product uint16) []activity.Activity {
	newActivities := make([]activity.Activity, 0)

	for i := range activities {
		act := &activities[i]

		for j := range act.Sessions {
			ses := act.Sessions[j]

			newActivity := activity.Activity{
				Creator: activity.Creator{
					Manufacturer: &manufacturer,
					Product:      &product,
					TimeCreated:  act.Creator.TimeCreated,
				},
				Timezone: act.Timezone,
			}
			newActivity.Sessions = []*activity.Session{ses}
			newActivities = append(newActivities, newActivity)
		}
	}
	return newActivities
}

func (s *service) changeSport(activities []activity.Activity, sports []string) {
	var cur int
	for i := range activities {
		act := &activities[i]
		for j := range act.Sessions {
			ses := act.Sessions[j]
			ses.Sport = sports[cur]
			cur++
		}
	}
}

// trimRecords trims Records based on the markers (1 marker correspond to 1 session) and recalculate the summary since the records is trimmed.
func (s *service) trimRecords(a *activity.Activity, markers []spec.EncodeMarker) error {
	if len(markers) == 0 {
		return nil
	}
	if len(markers) < len(a.Sessions) {
		return fmt.Errorf("trim: marker size is less than sessions size")
	}

	for i := range a.Sessions {
		ses := a.Sessions[i]
		marker := markers[i]

		if marker.StartN == 0 && marker.EndN == len(ses.Records)-1 { // no data to be trimmed
			continue
		}

		if marker.StartN >= len(ses.Records)-1 {
			ses.Records = nil
			continue
		}

		if marker.EndN == 0 {
			ses.Records = nil
			continue
		}

		// Adjust distance since ses.Records[marker.StartN] will be the beginning of record, its distance should be zero.
		// Find the exact or nearest distance as the substraction number.
		var distanceAdjustment float64
		for i := marker.StartN; i >= 0; i-- {
			rec := ses.Records[i]
			if rec.Distance != nil {
				distanceAdjustment = *rec.Distance
				break
			}
		}

		records := make([]*activity.Record, 0, marker.EndN-marker.StartN)
		for i := marker.StartN; i <= marker.EndN; i++ {
			rec := ses.Records[i]
			if rec.Distance != nil {
				*rec.Distance -= distanceAdjustment
			}
			records = append(records, rec)
		}

		if len(records) == 0 {
			continue
		}

		ses.Records = records

		// Recalculate Lap and Session Summary
		remainingRecords := make([]*activity.Record, 0)
		sesRecords := ses.Records
		newLaps := make([]*activity.Lap, 0)
		for j := range ses.Laps {
			lap := ses.Laps[j]
			lapRecords := make([]*activity.Record, 0)

			for k := range sesRecords {
				rec := sesRecords[k]
				if lap.IsBelongToThisLap(rec.Timestamp) {
					lapRecords = append(lapRecords, rec)
				} else {
					remainingRecords = append(remainingRecords, rec)
				}
			}

			lapFromRecords := activity.NewLapFromRecords(lapRecords, ses.Sport)
			if lapFromRecords != nil {
				newLaps = append(newLaps, lapFromRecords)
			}
			sesRecords = remainingRecords
		}

		newSes := activity.NewSessionFromLaps(newLaps, ses.Sport)
		newSes.Laps = newLaps
		newSes.Records = ses.Records
		*ses = *newSes
	}

	// Validate Records in Sessions
	validSessions := make([]*activity.Session, 0, len(a.Sessions))
	for i := range a.Sessions {
		if len(a.Sessions[i].Records) == 0 {
			continue
		}
		validSessions = append(validSessions, a.Sessions[i])
	}

	a.Sessions = validSessions

	return nil
}

// concealGPSPositions conceal positions from the records by removing PositionLat and PositionLong.
func (s *service) concealGPSPositions(a *activity.Activity, markers []spec.EncodeMarker) error {
	if len(markers) == 0 {
		return nil
	}
	if len(markers) < len(a.Sessions) {
		return fmt.Errorf("conceal: marker size is less than sessions size")
	}

	for i := range a.Sessions {
		ses := a.Sessions[i]
		marker := markers[i]

		if marker.StartN == 0 && marker.EndN == 0 {
			continue
		}

		if marker.StartN >= len(ses.Records)-1 {
			marker.StartN = len(ses.Records) - 1
		}

		if marker.EndN == 0 {
			marker.EndN = len(ses.Records) - 1
		}

		for j := 0; j < marker.StartN+1; j++ {
			ses.Records[j].PositionLat = nil
			ses.Records[j].PositionLong = nil
		}

		for j := marker.EndN + 1; j < len(ses.Records); j++ {
			ses.Records[j].PositionLat = nil
			ses.Records[j].PositionLong = nil
		}
	}

	return nil
}

// RemoveFields removes field from the entire records as well as the summary of it.
func (s *service) RemoveFields(a *activity.Activity, fields map[string]struct{}) {
	if len(fields) == 0 {
		return
	}

	for i := range a.Sessions {
		ses := a.Sessions[i]

		if _, ok := fields["distance"]; ok {
			ses.TotalDistance = 0
		}
		if _, ok := fields["altitude"]; ok {
			ses.TotalAscent = 0
			ses.TotalDescent = 0
			ses.AvgAltitude = nil
			ses.MaxAltitude = nil
		}
		if _, ok := fields["heartRate"]; ok {
			ses.AvgHeartRate = nil
			ses.MaxHeartRate = nil
		}
		if _, ok := fields["cadence"]; ok {
			ses.AvgCadence = nil
			ses.MaxCadence = nil
		}
		if _, ok := fields["speed"]; ok {
			ses.AvgSpeed = nil
			ses.MaxSpeed = nil
		}
		if _, ok := fields["power"]; ok {
			ses.AvgPower = nil
			ses.MaxPower = nil
		}
		if _, ok := fields["temperature"]; ok {
			ses.AvgTemperature = nil
			ses.MaxTemperature = nil
		}

		for j := range ses.Laps {
			lap := ses.Laps[j]

			if _, ok := fields["distance"]; ok {
				lap.TotalDistance = 0
			}
			if _, ok := fields["altitude"]; ok {
				lap.TotalAscent = 0
				lap.TotalDescent = 0
				lap.AvgAltitude = nil
				lap.MaxAltitude = nil
			}
			if _, ok := fields["heartRate"]; ok {
				lap.AvgHeartRate = nil
				lap.MaxHeartRate = nil
			}
			if _, ok := fields["cadence"]; ok {
				lap.AvgCadence = nil
				lap.MaxCadence = nil
			}
			if _, ok := fields["speed"]; ok {
				lap.AvgSpeed = nil
				lap.MaxSpeed = nil
			}
			if _, ok := fields["power"]; ok {
				lap.AvgPower = nil
				lap.MaxPower = nil
			}
			if _, ok := fields["temperature"]; ok {
				lap.AvgTemperature = nil
				lap.MaxTemperature = nil
			}
		}

		for j := range ses.Records {
			rec := ses.Records[j]

			if _, ok := fields["positionLat"]; ok {
				rec.PositionLat = nil
			}
			if _, ok := fields["positionLong"]; ok {
				rec.PositionLong = nil
			}
			if _, ok := fields["distance"]; ok {
				rec.Distance = nil
			}
			if _, ok := fields["altitude"]; ok {
				rec.Altitude = nil
			}
			if _, ok := fields["heartRate"]; ok {
				rec.HeartRate = nil
			}
			if _, ok := fields["cadence"]; ok {
				rec.Cadence = nil
			}
			if _, ok := fields["speed"]; ok {
				rec.Speed = nil
			}
			if _, ok := fields["power"]; ok {
				rec.Power = nil
			}
			if _, ok := fields["temperature"]; ok {
				rec.Temperature = nil
			}
		}
	}
}

func (s *service) ManufacturerList() result.ManufacturerList {
	manufacturers := make([]fit.Manufacturer, 0, len(s.manufacturers))
	for _, v := range s.manufacturers {
		manufacturers = append(manufacturers, v)
	}
	slices.SortFunc(manufacturers, func(a, b fit.Manufacturer) int {
		if strings.ToLower(a.Name) < strings.ToLower(b.Name) {
			return -1
		}
		return 1
	})
	return result.ManufacturerList{Manufacturers: manufacturers}
}

func (s *service) SportList() result.SportList {
	rawSports := typedef.ListSport()
	sports := make([]result.Sport, 0, len(rawSports)-1)

	for i := range rawSports {
		rs := rawSports[i]
		if rs == typedef.SportInvalid {
			continue
		}

		sport := result.Sport{
			ID:   uint8(rs),
			Name: kit.FormatTitle(rs.String()),
		}

		sport.ToleranceMovingSpeed = activity.ToleranceMovingSpeed(sport.Name)

		sports = append(sports, sport)
	}

	slices.SortFunc(sports, func(a, b result.Sport) int {
		if a.Name < b.Name {
			return -1
		}
		return 1
	})

	return result.SportList{Sports: sports}
}
