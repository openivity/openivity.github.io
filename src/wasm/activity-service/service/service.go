// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/aggregator"
	"github.com/openivity/activity-service/service/result"
	"github.com/openivity/activity-service/service/spec"
	"github.com/openivity/activity-service/strutils"
	"golang.org/x/exp/slices"
)

var ErrFileTypeUnsupported = errors.New("file type is unsupported")

// DecodeEncoder is a contract that any types implement these methods can be used by the Service.
type DecodeEncoder interface {
	// Decode decodes the given r into activities and returns any encountered errors.
	Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error)
	// Encode encodes the given activities into a slice of bytes and returns any encountered errors.
	Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error)
}

// Service is an activity service. It handle decoding and encoding these following file formats: FIT, GPX and TCX.
type Service struct {
	fit           DecodeEncoder
	gpx           DecodeEncoder
	tcx           DecodeEncoder
	manufacturers map[typedef.Manufacturer]activity.Manufacturer
}

// New creates new activity service to handle decoding and encoding these following file formats: FIT, GPX and TCX.
func New(fit, gpx, tcx DecodeEncoder, manufacturers map[typedef.Manufacturer]activity.Manufacturer) *Service {
	return &Service{
		fit:           fit,
		gpx:           gpx,
		tcx:           tcx,
		manufacturers: manufacturers,
	}
}

func (s *Service) Decode(ctx context.Context, rs []io.Reader) result.Decode {
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

func (s *Service) decodeWorker(ctx context.Context, rc <-chan io.Reader, resc chan<- result.DecodeWorker, wg *sync.WaitGroup, index int) {
	defer wg.Done()

	activities, err := s.decode(ctx, <-rc)
	if err != nil {
		resc <- result.DecodeWorker{Err: err, Index: index}
		return
	}

	for i := range activities {
		if activities[i].Creator.Name == "" {
			activities[i].Creator.Name = s.creatorName(activities[i].Creator.Manufacturer, activities[i].Creator.Product)
		}
		resc <- result.DecodeWorker{Activity: &activities[i], Index: index}
	}
}

func (s *Service) decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	fileType, err := s.readType(r)
	if err != nil {
		return nil, err
	}
	switch fileType {
	case spec.FileTypeFIT:
		return s.fit.Decode(ctx, r)
	case spec.FileTypeGPX:
		return s.gpx.Decode(ctx, r)
	case spec.FileTypeTCX:
		return s.tcx.Decode(ctx, r)
	default:
		return nil, ErrFileTypeUnsupported
	}
}

func (s *Service) readType(r io.Reader) (spec.FileType, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return spec.FileTypeUnsupported, err
	}
	return spec.FileType(b[0]), nil
}

func (s *Service) creatorName(manufacturerID typedef.Manufacturer, productID uint16) string {
	manufacturer, ok := s.manufacturers[manufacturerID]
	if !ok {
		return activity.Unknown
	}

	var productName string
	for i := range manufacturer.Products {
		product := manufacturer.Products[i]
		if product.ID == productID {
			productName = product.Name
			break
		}
	}
	if productName == "" {
		productName = "(" + strconv.FormatUint(uint64(productID), 10) + ")"
	}

	return manufacturer.Name + " " + productName
}

func (s *Service) Encode(ctx context.Context, encodeSpec spec.Encode) result.Encode {
	begin := time.Now()

	activities, err := s.preprocessEncode(encodeSpec)
	if err != nil {
		return result.Encode{Err: err}
	}

	var bs [][]byte
	switch encodeSpec.TargetFileType {
	case spec.FileTypeFIT:
		bs, err = s.fit.Encode(ctx, activities)
	case spec.FileTypeGPX:
		bs, err = s.gpx.Encode(ctx, activities)
	case spec.FileTypeTCX:
		bs, err = s.tcx.Encode(ctx, activities)
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

func (s *Service) preprocessEncode(encodeSpec spec.Encode) ([]activity.Activity, error) {
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
	var validActivityCount int
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

		validActivityCount++
		s.changeSport(activities, encodeSpec.Sports)
		s.removeFields(activity, removeFields)
	}

	if validActivityCount == 0 {
		return nil, fmt.Errorf("no activity data after processed")
	}

	var newActivities []activity.Activity
	switch encodeSpec.ToolMode {
	case spec.ToolModeEdit:
		for i := range activities {
			activities[i].Creator.Manufacturer = encodeSpec.ManufacturerID
			activities[i].Creator.Product = encodeSpec.ProductID
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

func (s *Service) combineActivity(activities []activity.Activity, manufacturer typedef.Manufacturer, product uint16) activity.Activity {
	creator := activity.CreateCreator(nil)
	creator.FileId.
		SetType(typedef.FileActivity).
		SetManufacturer(manufacturer).
		SetProduct(product).
		SetTimeCreated(activities[0].Creator.TimeCreated)

	newActivity := activity.Activity{
		Creator:           creator,
		Timezone:          activities[0].Timezone,
		Sessions:          activities[0].Sessions,
		SplitSummaries:    activities[0].SplitSummaries,
		UnrelatedMessages: activities[0].UnrelatedMessages,
	}

	lastDistance := getLastDistanceOfRecords(newActivity.Sessions[0].Records)
	lastAccumulatedPower := getLastAccumulatedPower(newActivity.Sessions[0].Records)

	for i := 1; i < len(activities); i++ {
		cur := &activities[i]

		if newActivity.Creator.Name == "" {
			newActivity.Creator.Name = cur.Creator.Name
		}

		newActLastSes := &newActivity.Sessions[len(newActivity.Sessions)-1]
		curActFirstSes := cur.Sessions[0]

		// Adjust distance before combine
		for j := range cur.Sessions {
			ses := &cur.Sessions[j]
			for k := range ses.Records {
				rec := &ses.Records[k]
				if rec.Distance != basetype.Uint32Invalid {
					rec.Distance += lastDistance
				}
				if rec.AccumulatedPower != basetype.Uint32Invalid {
					rec.AccumulatedPower += lastAccumulatedPower
				}
			}
			lastDistance = getLastDistanceOfRecords(ses.Records)
			lastAccumulatedPower = getLastAccumulatedPower(ses.Records)
		}

		if newActLastSes.Sport == curActFirstSes.Sport {
			// Combine records and laps to newActivity's last session
			newActLastSes.Records = append(newActLastSes.Records, curActFirstSes.Records...)
			newActLastSes.Laps = append(newActLastSes.Laps, curActFirstSes.Laps...)

			// Update summary
			gap := (curActFirstSes.StartTime.Sub(newActLastSes.EndTime()).Seconds() * 1000)
			newActLastSes.TotalElapsedTime += uint32(gap)
			newActLastSes.TotalTimerTime += uint32(gap)
			aggregator.Aggregate(newActLastSes.Session, curActFirstSes.Session)
			newActLastSes.Summarize()

			cur.Sessions = cur.Sessions[1:]
		}

		newActivity.Sessions = append(newActivity.Sessions, cur.Sessions...)

		for _, m := range cur.SplitSummaries {
			var ok bool
			for _, v := range newActivity.SplitSummaries {
				if v.SplitType == m.SplitType {
					aggregator.Aggregate(v, m)
					ok = true
					break
				}
			}
			if !ok {
				newActivity.SplitSummaries = append(newActivity.SplitSummaries, m)
			}
			continue
		}

		for j := 0; j < len(cur.UnrelatedMessages); j++ {
			if cur.UnrelatedMessages[j].Num == mesgnum.FileCreator {
				continue
			}
			newActivity.UnrelatedMessages = append(newActivity.UnrelatedMessages, cur.UnrelatedMessages[j])
		}
	}

	return newActivity
}

func getLastDistanceOfRecords(records []activity.Record) uint32 {
	for i := len(records) - 1; i >= 0; i-- {
		rec := records[i]
		if rec.Distance != basetype.Uint32Invalid {
			return rec.Distance
		}
	}
	return 0
}

func getLastAccumulatedPower(records []activity.Record) uint32 {
	for i := len(records) - 1; i >= 0; i-- {
		rec := records[i]
		if rec.AccumulatedPower != basetype.Uint32Invalid {
			return rec.AccumulatedPower
		}
	}
	return 0
}

func (s *Service) splitActivityPerSession(activities []activity.Activity, manufacturer typedef.Manufacturer, product uint16) []activity.Activity {
	newActivities := make([]activity.Activity, 0)

	for i := range activities {
		act := &activities[i]

		for j := range act.Sessions {
			ses := act.Sessions[j]

			creator := activity.CreateCreator(nil)
			creator.FileId.
				SetType(typedef.FileActivity).
				SetManufacturer(manufacturer).
				SetProduct(product).
				SetTimeCreated(activities[0].Creator.TimeCreated)

			newActivity := activity.Activity{
				Creator:  creator,
				Timezone: act.Timezone,
			}

			if j != 0 {
				// We should remove accumulated distance and accumulated power from previous session.
				// To achieve that, we must get the previous values.
				prevSes := act.Sessions[j-1]
				var prevDistance = basetype.Uint32Invalid
				for k := len(prevSes.Records) - 1; k >= 0; k-- {
					if prevSes.Records[k].Distance != basetype.Uint32Invalid {
						prevDistance = prevSes.Records[k].Distance
						break
					}
				}
				var prevAccumulatedPower = basetype.Uint32Invalid
				for k := len(prevSes.Records) - 1; k >= 0; k-- {
					if prevSes.Records[k].AccumulatedPower != basetype.Uint32Invalid {
						prevAccumulatedPower = prevSes.Records[k].AccumulatedPower
						break
					}
				}
				// Subtract current values with previous values.
				for k := range ses.Records {
					rec := &ses.Records[k]
					if prevDistance != basetype.Uint32Invalid && rec.Distance >= prevDistance {
						rec.Distance -= prevDistance
					}
					if prevAccumulatedPower != basetype.Uint32Invalid && rec.AccumulatedPower >= prevAccumulatedPower {
						rec.AccumulatedPower -= prevAccumulatedPower
					}
				}
			}

			newActivity.Sessions = []activity.Session{ses}
			newActivities = append(newActivities, newActivity)
		}
	}
	return newActivities
}

func (s *Service) changeSport(activities []activity.Activity, sports []string) {
	var cur int
	for i := range activities {
		act := &activities[i]
		for j := range act.Sessions {
			ses := &act.Sessions[j]
			ses.Sport = typedef.SportFromString(strutils.ToLowerSnakeCase(sports[cur]))
			cur++
		}
	}
}

// trimRecords trims Records based on the markers (1 marker correspond to 1 session) and recalculate the summary since the records is trimmed.
func (s *Service) trimRecords(a *activity.Activity, markers []spec.EncodeMarker) error {
	if len(markers) < len(a.Sessions) {
		return fmt.Errorf("trim: marker size is less than sessions size")
	}

	for i := range a.Sessions {
		ses := &a.Sessions[i]
		marker := markers[i]

		if marker.StartN == 0 && marker.EndN == len(ses.Records)-1 {
			continue
		}

		if marker.StartN == marker.EndN {
			ses.Records = nil
			continue
		}

		// Adjust distance since ses.Records[marker.StartN] will be the beginning of record, its distance should be zero.
		// Find the exact or nearest distance as the substraction number.
		var distanceAdjustment uint32
		for i := marker.StartN; i >= 0; i-- {
			rec := &ses.Records[i]
			if rec.Distance != basetype.Uint32Invalid {
				distanceAdjustment = rec.Distance
				break
			}
		}

		ses.Records = ses.Records[marker.StartN : marker.EndN+1]
		for i := range ses.Records {
			rec := &ses.Records[i]
			if rec.Distance != basetype.Uint32Invalid {
				rec.Distance -= distanceAdjustment
			}
		}

		if len(ses.Records) == 0 {
			continue
		}

		records := slices.Clone(ses.Records)

		// Recalculate Lap and Session Summary
		newLaps := make([]activity.Lap, 0)
		for i := range ses.Laps {
			lap := &ses.Laps[i]

			var pos int
			for j := range records {
				rec := &records[j]
				if lap.IsBelongToThisLap(rec.Timestamp) {
					records[j], records[pos] = records[pos], records[j]
					pos++
				}
			}
			lapRecords := records[:pos]

			if len(lapRecords) != 0 {
				lapFromRecords := activity.NewLapFromRecords(lapRecords, ses.Sport)
				newLaps = append(newLaps, lapFromRecords)
			}
			records = records[pos:]
		}

		newSes := activity.NewSessionFromLaps(newLaps)
		newSes.Laps = newLaps
		newSes.Records = ses.Records
		newSes.Summarize()
		*ses = newSes
	}

	// Validate Records in Sessions
	validSessions := make([]activity.Session, 0, len(a.Sessions))
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
func (s *Service) concealGPSPositions(a *activity.Activity, markers []spec.EncodeMarker) error {
	if len(markers) < len(a.Sessions) {
		return fmt.Errorf("conceal: marker size is less than sessions size")
	}

	for i := range a.Sessions {
		ses := a.Sessions[i]
		marker := markers[i]

		if marker.StartN == 0 && marker.EndN == len(ses.Records)-1 {
			continue
		}

		if marker.StartN == marker.EndN {
			marker.StartN++
			marker.EndN--
			continue
		}

		for j := 0; j < marker.StartN; j++ {
			ses.Records[j].PositionLat = basetype.Sint32Invalid
			ses.Records[j].PositionLong = basetype.Sint32Invalid
		}

		for j := marker.EndN + 1; j < len(ses.Records); j++ {
			ses.Records[j].PositionLat = basetype.Sint32Invalid
			ses.Records[j].PositionLong = basetype.Sint32Invalid
		}
	}

	return nil
}

// removeFields removes field from the entire records as well as the summary of it.
func (s *Service) removeFields(a *activity.Activity, fields map[string]struct{}) {
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
			ses.AvgAltitude = basetype.Uint16Invalid
			ses.MaxAltitude = basetype.Uint16Invalid
		}
		if _, ok := fields["heartRate"]; ok {
			ses.AvgHeartRate = basetype.Uint8Invalid
			ses.MaxHeartRate = basetype.Uint8Invalid
		}
		if _, ok := fields["cadence"]; ok {
			ses.AvgCadence = basetype.Uint8Invalid
			ses.MaxCadence = basetype.Uint8Invalid
		}
		if _, ok := fields["speed"]; ok {
			ses.AvgSpeed = basetype.Uint16Invalid
			ses.MaxSpeed = basetype.Uint16Invalid
		}
		if _, ok := fields["power"]; ok {
			ses.AvgPower = basetype.Uint16Invalid
			ses.MaxPower = basetype.Uint16Invalid
		}
		if _, ok := fields["temperature"]; ok {
			ses.AvgTemperature = basetype.Sint8Invalid
			ses.MaxTemperature = basetype.Sint8Invalid
		}

		for j := range ses.Laps {
			lap := ses.Laps[j]

			if _, ok := fields["distance"]; ok {
				lap.TotalDistance = 0
			}
			if _, ok := fields["altitude"]; ok {
				lap.TotalAscent = 0
				lap.TotalDescent = 0
				lap.AvgAltitude = basetype.Uint16Invalid
				lap.MaxAltitude = basetype.Uint16Invalid
			}
			if _, ok := fields["heartRate"]; ok {
				lap.AvgHeartRate = basetype.Uint8Invalid
				lap.MaxHeartRate = basetype.Uint8Invalid
			}
			if _, ok := fields["cadence"]; ok {
				lap.AvgCadence = basetype.Uint8Invalid
				lap.MaxCadence = basetype.Uint8Invalid
			}
			if _, ok := fields["speed"]; ok {
				lap.AvgSpeed = basetype.Uint16Invalid
				lap.MaxSpeed = basetype.Uint16Invalid
			}
			if _, ok := fields["power"]; ok {
				lap.AvgPower = basetype.Uint16Invalid
				lap.MaxPower = basetype.Uint16Invalid
			}
			if _, ok := fields["temperature"]; ok {
				lap.AvgTemperature = basetype.Sint8Invalid
				lap.MaxTemperature = basetype.Sint8Invalid
			}
		}

		for j := range ses.Records {
			rec := ses.Records[j]

			if _, ok := fields["positionLat"]; ok {
				rec.PositionLat = basetype.Sint32Invalid
			}
			if _, ok := fields["positionLong"]; ok {
				rec.PositionLong = basetype.Sint32Invalid
			}
			if _, ok := fields["distance"]; ok {
				rec.Distance = basetype.Uint32Invalid
			}
			if _, ok := fields["altitude"]; ok {
				rec.Altitude = basetype.Uint16Invalid
			}
			if _, ok := fields["heartRate"]; ok {
				rec.HeartRate = basetype.Uint8Invalid
			}
			if _, ok := fields["cadence"]; ok {
				rec.Cadence = basetype.Uint8Invalid
			}
			if _, ok := fields["speed"]; ok {
				rec.Speed = basetype.Uint16Invalid
			}
			if _, ok := fields["power"]; ok {
				rec.Power = basetype.Uint16Invalid
			}
			if _, ok := fields["temperature"]; ok {
				rec.Temperature = basetype.Sint8Invalid
			}
		}
	}
}

func (s *Service) ManufacturerList() result.ManufacturerList {
	manufacturers := make([]activity.Manufacturer, 0, len(s.manufacturers))
	for _, v := range s.manufacturers {
		manufacturers = append(manufacturers, v)
	}
	slices.SortFunc(manufacturers, func(a, b activity.Manufacturer) int {
		if strings.ToLower(a.Name) < strings.ToLower(b.Name) {
			return -1
		}
		return 1
	})
	return result.ManufacturerList{Manufacturers: manufacturers}
}

func (s *Service) SportList() result.SportList {
	sportList := typedef.ListSport()
	sports := make([]activity.Sport, 0, len(sportList))

	for _, v := range sportList {
		sport := activity.Sport{
			ID:   v,
			Name: strutils.ToTitle(v.String()),
		}
		sport.ToleranceMovingSpeed = activity.ToleranceMovingSpeed(sport.ID)
		sports = append(sports, sport)
	}

	slices.SortFunc(sports, func(a, b activity.Sport) int {
		if a.Name < b.Name {
			return -1
		}
		return 0
	})

	return result.SportList{Sports: sports}
}
