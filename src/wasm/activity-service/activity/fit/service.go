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

package fit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/mem"
	"golang.org/x/exp/slices"
)

var _ activity.Service = &service{}

type service struct {
	preprocessor *activity.Preprocessor
}

// NewService creates new FIT service.
func NewService(preproc *activity.Preprocessor) activity.Service {
	return &service{
		preprocessor: preproc,
	}
}

func (s *service) Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	lis := filedef.NewListener()
	defer lis.Close()

	dec := decoder.New(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
		decoder.WithIgnoreChecksum(),
	)

	activities := make([]activity.Activity, 0, 1) // In most cases, 1 fit == 1 activity
	for dec.Next() {
		fileId, err := dec.PeekFileId()
		if err != nil {
			return nil, err
		}

		if fileId.Type != typedef.FileActivity {
			if err = dec.Discard(); err != nil {
				return nil, err
			}
			continue
		}

		_, err = dec.DecodeWithContext(ctx)
		if err != nil {
			return nil, err
		}

		activityFile := lis.File().(*filedef.Activity)
		if len(activityFile.Records) == 0 {
			continue
		}

		activities = append(activities, s.convertToActivity(activityFile))
	}

	if len(activities) == 0 {
		return nil, fmt.Errorf("fit: %w", activity.ErrNoActivity)
	}

	// Sumarizing...
	for i := range activities {
		act := &activities[i]
		for j := range act.Sessions {
			ses := &act.Sessions[j]

			s.preprocessor.CalculateDistanceAndSpeed(ses.Records)
			s.preprocessor.SmoothingElevation(ses.Records)
			s.preprocessor.CalculateGrade(ses.Records)
			if activity.HasPace(ses.Sport) {
				s.preprocessor.CalculatePace(ses.Sport, ses.Records)
			}

			s.recalculateSummary(ses)
		}
	}

	return activities, nil
}

func (s *service) convertToActivity(activityFile *filedef.Activity) activity.Activity {
	var timezone int8
	if activityFile.Activity != nil {
		localTimestamp := activityFile.Activity.LocalTimestamp
		timestamp := activityFile.Activity.Timestamp
		if !localTimestamp.IsZero() && !timestamp.IsZero() {
			timezone = int8(datetime.TzOffsetHours(localTimestamp, timestamp))
		}
	}

	fileId := activityFile.FileId
	act := activity.Activity{
		Creator:           activity.CreateCreator(&fileId),
		Timezone:          timezone,
		UnrelatedMessages: s.handleUnrelatedMessages(activityFile),
	}

	// Convert Records, Laps and Sessions to activity's structs
	records := make([]activity.Record, len(activityFile.Records))
	for i := range activityFile.Records {
		records[i] = activity.CreateRecord(activityFile.Records[i])
	}

	laps := make([]activity.Lap, len(activityFile.Laps))
	for i := range activityFile.Laps {
		laps[i] = activity.CreateLap(activityFile.Laps[i])
	}

	sessions := make([]activity.Session, len(activityFile.Sessions))
	for i := range activityFile.Sessions {
		sessions[i] = activity.CreateSession(activityFile.Sessions[i])
	}

	records = s.preprocessor.AggregateByTimestamp(records)

	// Create Sessions and Laps if not exist. This could happen only if:
	//  - FIT file is truncated, so only some Records could be retrieved.
	//  - Some devices may not create Lap even though it's required for an Activity File.
	//    ref: https://developer.garmin.com/fit/file-types/activity
	if len(sessions) == 0 {
		if len(laps) == 0 {
			lap := activity.NewLapFromRecords(records, typedef.SportGeneric)
			laps = append(laps, lap)
		}

		ses := activity.NewSessionFromLaps(laps)
		ses.Records = records
		ses.Laps = laps
		act.Sessions = []activity.Session{ses}

		return act
	} else if len(laps) == 0 {
		// Some devices may only create sessions without laps.
		// Let's create laps from sessions, 1 session should at least have 1 lap.
		laps = make([]activity.Lap, len(sessions))
		for i := range sessions {
			laps[i] = activity.NewLapFromSession(&sessions[i])
		}
	}

	const anomalyThreshold = 10
	act.Sessions = make([]activity.Session, 0, len(sessions))
	for i := range sessions {
		ses := sessions[i]

		laps = ses.PutLaps(laps...)
		records = ses.PutRecords(records...)

		if len(ses.Laps) == 0 {
			ses.Laps = append(ses.Laps, activity.NewLapFromSession(&sessions[i]))
		}

		// Include anomaly records below threshold to the last session.
		if i == len(sessions)-1 && len(records) < anomalyThreshold {
			ses.Records = append(ses.Records, records...)
			records = records[:0]
		}
		act.Sessions = append(act.Sessions, ses)
	}

	if len(records) == 0 {
		return act
	}

	// Handle remaining laps and records that don't belong any session
	// This could happen when file is truncated or file is not properly encoded.
	if len(laps) != 0 {
		ses := activity.NewSessionFromLaps(laps)
		ses.Laps = laps
		records = ses.PutRecords(records...)
		act.Sessions = append(act.Sessions, ses)
	}

	// Handle remaining records that don't belong anywhere.
	if len(records) != 0 {
		lap := activity.NewLapFromRecords(records, typedef.SportGeneric)
		laps = []activity.Lap{lap}
		ses := activity.NewSessionFromLaps(laps)
		ses.Laps = laps
		ses.Records = records
		act.Sessions = append(act.Sessions, ses)
	}

	return act
}

// recalculateSummary recalculates values based on Laps and Records.
func (s *service) recalculateSummary(ses *activity.Session) {
	records := slices.Clone(ses.Records)
	if len(ses.Laps) == 1 { // Ensure lap's time windows match with session, FIT produces by Strava contains wrong time.
		ses.Laps[0].StartTime = ses.StartTime
		ses.Laps[0].TotalElapsedTime = ses.TotalElapsedTime
		ses.Laps[0].TotalTimerTime = ses.TotalTimerTime
	}
	for i := range ses.Laps {
		lap := &ses.Laps[i]
		var pos int

		for j := range records {
			if lap.IsBelongToThisLap(records[j].Timestamp) {
				records[j], records[pos] = records[pos], records[j]
				pos++
			}
		}
		if len(records) == 0 {
			continue
		}
		lapFromRecords := activity.NewLapFromRecords(records[:pos], ses.Sport)
		lap.ReplaceValues(&lapFromRecords)
		records = records[pos:]
	}
	sesFromLaps := activity.NewSessionFromLaps(ses.Laps)
	ses.ReplaceValues(&sesFromLaps)
	ses.Summarize()
}

func (s *service) handleUnrelatedMessages(activityFile *filedef.Activity) []proto.Message {
	size := len(activityFile.DeveloperDataIds) +
		len(activityFile.FieldDescriptions) +
		len(activityFile.DeviceInfos) +
		len(activityFile.Events) +
		len(activityFile.Lengths) +
		len(activityFile.SegmentLap) +
		len(activityFile.ZonesTargets) +
		len(activityFile.Workouts) +
		len(activityFile.WorkoutSteps) +
		len(activityFile.HRs) +
		len(activityFile.HRVs) +
		len(activityFile.UnrelatedMessages)

	if activityFile.UserProfile != nil {
		size += 1
	}

	unrelatedMessages := make([]proto.Message, 0, size)

	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.DeveloperDataId, activityFile.DeveloperDataIds)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.DeveloperDataId, activityFile.FieldDescriptions)

	if activityFile.UserProfile != nil {
		unrelatedMessages = append(unrelatedMessages, activityFile.UserProfile.ToMesg(nil))
	}

	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.DeviceInfo, activityFile.DeviceInfos)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.Event, activityFile.Events)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.Length, activityFile.Lengths)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.SegmentLap, activityFile.SegmentLap)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.ZonesTarget, activityFile.ZonesTargets)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.Workout, activityFile.Workouts)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.WorkoutStep, activityFile.WorkoutSteps)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.Hr, activityFile.HRs)
	filedef.ToMesgs(&unrelatedMessages, nil, mesgnum.Hrv, activityFile.HRVs)

	unrelatedMessages = append(unrelatedMessages, activityFile.UnrelatedMessages...)

	return unrelatedMessages
}

func (s *service) Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error) {
	buf := mem.GetBuffer()
	defer mem.PutBuffer(buf)

	bufAt := &bytesBufferAt{buf}

	opts := []encoder.Option{
		encoder.WithProtocolVersion(proto.V2),
		encoder.WithNormalHeader(15),
	}
	enc := encoder.New(bufAt, opts...)

	bs := make([][]byte, len(activities))
	for i := range activities {
		s.makeLastSummary(&activities[i])
		fit := activities[i].ToFIT(nil)

		if err := enc.EncodeWithContext(ctx, &fit); err != nil {
			return nil, fmt.Errorf("could not encode: %w", err)
		}

		bs[i] = slices.Clone(bufAt.Buffer.Bytes())
		bufAt.Buffer.Reset()
		enc.Reset(bufAt, opts...)
	}

	return bs, nil
}

func (s *service) makeLastSummary(a *activity.Activity) {
	var lastTimestamp time.Time
	for i := len(a.Sessions) - 1; i >= 0; i-- {
		ses := a.Sessions[i]

		for j := len(ses.Records) - 1; j >= 0; j-- {
			rec := ses.Records[j]
			if !rec.Timestamp.IsZero() {
				lastTimestamp = rec.Timestamp
				break
			}
		}

		for j := len(ses.Laps) - 1; j >= 0; j-- {
			lap := ses.Laps[j]
			if !lap.Timestamp.IsZero() && lap.Timestamp.After(lastTimestamp) {
				lastTimestamp = lap.Timestamp
				break
			}
		}

		if !lastTimestamp.IsZero() {
			break
		}
	}

	for i := range a.Sessions {
		a.Sessions[i].Timestamp = lastTimestamp
	}
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
