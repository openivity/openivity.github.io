package fit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/proto"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/preprocessor"
	"golang.org/x/exp/slices"
)

var _ activity.Service = &service{}

type service struct {
	preprocessor  *preprocessor.Preprocessor
	manufacturers map[uint16]Manufacturer
}

func NewService(preproc *preprocessor.Preprocessor, manufacturers map[uint16]Manufacturer) activity.Service {
	return &service{
		preprocessor:  preproc,
		manufacturers: manufacturers,
	}
}

func (s *service) Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	lis := NewListener()
	defer lis.WaitAndClose()

	dec := decoder.New(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
		decoder.WithIgnoreChecksum(),
	)

	activities := make([]activity.Activity, 0, 1) // In most cases, 1 fit == 1 activity
	for dec.Next() {
		_, err := dec.DecodeWithContext(ctx)
		if err != nil {
			return nil, err
		}

		res := lis.Result()
		act := s.convertListenerResultToActivity(res)
		if act == nil {
			continue
		}

		activities = append(activities, *act)
	}

	if len(activities) == 0 {
		return nil, fmt.Errorf("fit: %w", activity.ErrNoActivity)
	}

	return activities, nil
}

func (s *service) convertListenerResultToActivity(result *ListenerResult) *activity.Activity {
	if len(result.Records) == 0 {
		return nil
	}

	s.sanitize(result)

	creator := result.Creator
	if creator.Manufacturer != nil && creator.Product != nil {
		creator.Name = s.creatorName(*creator.Manufacturer, *creator.Product)
	}

	act := &activity.Activity{
		Creator:  *result.Creator,
		Timezone: result.Timezone,
	}

	result.Records = s.preprocessor.AggregateByTimestamp(result.Records)

	// Create Sessions and Laps if not exist. This could happen only if:
	//  - Fit file is truncated, so only some Records could be retrieved.
	//  - Some devices may not create Lap even though it's actually required for an Activity File.
	//    ref: https://developer.garmin.com/fit/file-types/activity
	if len(result.Sessions) == 0 {
		if len(result.Laps) == 0 {
			lap := activity.NewLapFromRecords(result.Records, activity.SportGeneric)
			result.Laps = append(result.Laps, lap)
		}

		ses := activity.NewSessionFromLaps(result.Laps, activity.SportGeneric)
		ses.Records = result.Records
		ses.Laps = result.Laps
		s.preprocessingRecords(ses.Records, ses.Sport)
		act.Sessions = []*activity.Session{ses}

		return act
	} else if len(result.Laps) == 0 {
		// Some devices may only create sessions but no laps, to ensure we don't lose any summary data in session
		// since there are information that only available in session but not in records.
		// Let's create laps from sessions, 1 session should at least have 1 lap.
		result.Laps = make([]*activity.Lap, len(result.Sessions))
		for i := range result.Sessions {
			result.Laps[i] = activity.NewLapFromSession(result.Sessions[i])
		}
	}

	for i := range result.Sessions {
		ses := result.Sessions[i]

		result.Laps = ses.PutLaps(result.Laps...)
		result.Records = ses.PutRecords(result.Records...)

		if len(ses.Records) == 0 {
			continue
		}

		s.preprocessingRecords(ses.Records, ses.Sport)
		s.finalizeSession(ses)

		act.Sessions = append(act.Sessions, ses)
	}

	if len(result.Records) != 0 && len(act.Sessions) != 0 {
		s.handleAnomalyRecords(result, act.Sessions[len(act.Sessions)-1])
	}

	if len(result.Records) == 0 {
		return act
	}

	// Handle remaining laps and records that don't belong any session
	// This could happen when file is truncated or file is not properly encoded.
	sport := activity.SportGeneric // Mark as Generic
	s.preprocessingRecords(result.Records, sport)

	if len(result.Laps) != 0 {
		ses := activity.NewSessionFromLaps(result.Laps, sport)
		ses.Laps = result.Laps

		result.Records = ses.PutRecords(result.Records...)

		if len(ses.Records) != 0 {
			s.finalizeSession(ses)
			act.Sessions = append(act.Sessions, ses)
		}
	}

	// Handle remaining records that don't belong anywhere.
	if len(result.Records) != 0 {
		lap := activity.NewLapFromRecords(result.Records, sport)
		laps := []*activity.Lap{lap}

		ses := activity.NewSessionFromLaps(laps, sport)
		ses.Laps = laps
		ses.Records = result.Records

		act.Sessions = append(act.Sessions, ses)
	}

	return act
}

// handleAnomalyRecords handles anomaly leftover records caused by session's EndTime miscalculation.
//
// Study case: FIT Files retrieved from Strava.
func (s *service) handleAnomalyRecords(result *ListenerResult, lastSession *activity.Session) {
	// In our study case, the number of anomaly record is 1, but just in case, we add a tolerance threshold to 10.
	// There should be no session with less than 10 records, if it's caused by truncated file, let's just add to the last session.
	// Having less than 10 records in a session doesn't add any value for analytic.
	const threshold = 10
	if len(result.Records) < threshold {
		ses := result.Sessions[len(result.Sessions)-1]
		ses.Records = append(ses.Records, result.Records...)
		for i := len(result.Records) - 1; i >= 0; i-- {
			rec := result.Records[i]
			if rec.Timestamp.After(ses.EndTime) {
				ses.EndTime = rec.Timestamp
				break
			}
		}
		result.Records = result.Records[:0]
	}
}

// sanitize removes any invalid item from given result.
func (s *service) sanitize(result *ListenerResult) {
	if len(result.Records) == 0 {
		return
	}

	validLaps := make([]*activity.Lap, 0)

	for i := range result.Laps {
		lap := result.Laps[i]

		// Timestamp, Start Time, and Total Elapsed Time are required fields for all summary messages.
		// If any of this field is missing, let's mark as invalid and use our own lap calculation later.
		if lap.Timestamp.IsZero() {
			continue
		}
		if lap.StartTime.IsZero() {
			continue
		}
		if lap.TotalElapsedTime == 0 {
			continue
		}

		// Most activity has 1 Lap and first Lap's StartTime should match first record's timestamp.
		// We should not try to accommodate all bad encoding practices and this guard should be sufficient for most cases.
		if i == 0 && !lap.StartTime.Equal(result.Records[0].Timestamp) {
			continue
		}

		validLaps = append(validLaps, lap)
	}

	result.Laps = validLaps
}

// preprocessingRecords pre-processes records per session since 1 session corresponds to 1 sport.
// We should not process different sports as one.
func (s *service) preprocessingRecords(records []*activity.Record, sport string) {
	s.preprocessor.CalculateDistanceAndSpeed(records)
	s.preprocessor.SmoothingElev(records)
	s.preprocessor.CalculateGrade(records)
	if activity.HasPace(sport) {
		s.preprocessor.CalculatePace(sport, records)
	}
}

// finalizeSession finalises session by creating lap if not exist and calculating its summary as well as calculating session's summary.
func (s *service) finalizeSession(ses *activity.Session) {
	if len(ses.Laps) == 0 {
		lap := activity.NewLapFromRecords(ses.Records, ses.Sport)
		ses.Laps = append(ses.Laps, lap)
		sesFromLaps := activity.NewSessionFromLaps(ses.Laps, ses.Sport)
		activity.CombineSession(ses, sesFromLaps)
		return
	}

	remainingRecords := ses.Records
	for j := range ses.Laps {
		lap := ses.Laps[j]

		lapRecords := make([]*activity.Record, 0)
		remainingLapRecords := make([]*activity.Record, 0)
		for k := range remainingRecords {
			rec := remainingRecords[k]

			if lap.IsBelongToThisLap(rec.Timestamp) {
				lapRecords = append(lapRecords, rec)
			} else {
				remainingLapRecords = append(remainingLapRecords, rec)
			}
		}
		remainingRecords = remainingLapRecords

		lapFromRecords := activity.NewLapFromRecords(lapRecords, ses.Sport)
		activity.CombineLap(lap, lapFromRecords)
	}

	// Handle remaining records that don't belong to any lap.
	if len(remainingRecords) != 0 {
		lap := activity.NewLapFromRecords(remainingRecords, ses.Sport)
		ses.Laps = append(ses.Laps, lap)
	}

	slices.SortFunc(ses.Laps, func(l1, l2 *activity.Lap) int {
		if l1.StartTime.Equal(l2.StartTime) {
			return 0
		}
		if l1.StartTime.Before(l2.StartTime) {
			return -1
		}
		return 1
	})

	sesFromLaps := activity.NewSessionFromLaps(ses.Laps, ses.Sport)
	activity.CombineSession(ses, sesFromLaps)

	s.preprocessor.SetSessionsWorkoutType(ses)
}

func (s *service) creatorName(manufacturerID, productID uint16) string {
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

func (s *service) Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error) {
	bs := make([][]byte, len(activities))
	bufAt := &bytesBufferAt{bytes.NewBuffer(make([]byte, 0, 1_000<<10))} // Buffer 1 MB

	for i := range activities {
		fit := s.convertActivityToFit(&activities[i])

		enc := encoder.New(bufAt, encoder.WithProtocolVersion(proto.V2))
		if err := enc.EncodeWithContext(ctx, fit); err != nil {
			return nil, fmt.Errorf("could not encode: %w", err)
		}

		bs[i] = slices.Clone(bufAt.Buffer.Bytes())
		bufAt.Buffer.Reset()
	}

	return bs, nil
}

// bytesBufferAt wraps bytes.Buffer to implement io.WriterAt enabling fast encoding.
type bytesBufferAt struct {
	*bytes.Buffer
}

func (b *bytesBufferAt) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return n, fmt.Errorf("negative offset")
	}
	l := off + int64(len(p))
	if l > int64(b.Len()) {
		return n, fmt.Errorf("offset > len")
	}
	n = copy(b.Bytes()[off:l], p)
	return
}

func (s *service) convertActivityToFit(act *activity.Activity) *proto.Fit {
	var lapCount, recordCount int
	sessionCount := len(act.Sessions)

	for i := range act.Sessions {
		lapCount += len(act.Sessions[i].Laps)
		recordCount += len(act.Sessions[i].Records)
	}

	fit := new(proto.Fit)
	fit.Messages = make([]proto.Message, 0, sessionCount+lapCount+recordCount+2) // +2 for FileId and Activity messages

	filedIdMesg := convertCreatorToMesg(&act.Creator) // Must be first the message
	fit.Messages = append(fit.Messages, filedIdMesg)

	eventStart := &activity.Event{
		Timestamp: act.Sessions[0].Records[0].Timestamp,
		Event:     uint8(typedef.EventTimer),
		EventType: uint8(typedef.EventTypeStart),
	}

	fit.Messages = append(fit.Messages, convertEventToMesg(eventStart)) // add event start

	for i := range act.Sessions {
		ses := act.Sessions[i]

		for j := range ses.Records {
			recMesg := convertRecordToMesg(ses.Records[j])
			fit.Messages = append(fit.Messages, recMesg)
		}

		if i == len(act.Sessions)-1 { // before last session add event stop all
			eventStopAll := &activity.Event{
				Timestamp: ses.Records[len(ses.Records)-1].Timestamp,
				Event:     uint8(typedef.EventTimer),
				EventType: uint8(typedef.EventTypeStopAll),
			}
			fit.Messages = append(fit.Messages, convertEventToMesg(eventStopAll))
		}

		for j := range ses.Laps {
			lapMesg := convertLapToMesg(ses.Laps[j])
			fit.Messages = append(fit.Messages, lapMesg)
		}

		sessionMesg := convertSessionToMesg(ses)
		fit.Messages = append(fit.Messages, sessionMesg)
	}

	activityMesg := createActivityMesg(act.Creator.TimeCreated, act.Timezone, uint16(sessionCount))
	fit.Messages = append(fit.Messages, activityMesg)

	return fit
}
