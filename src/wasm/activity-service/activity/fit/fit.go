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
	"sync"
	"time"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/kit/datetime"
	"github.com/muktihari/fit/profile/basetype"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/proto"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/aggregator"
	"github.com/openivity/activity-service/mem"
	"github.com/openivity/activity-service/service"
	"golang.org/x/exp/slices"
)

var _ service.DecodeEncoder = (*DecodeEncoder)(nil)

var decoderPool = sync.Pool{New: func() any { return decoder.New(nil) }}
var encoderPool = sync.Pool{New: func() any { return encoder.New(nil) }}

type DecodeEncoder struct {
	preprocessor *activity.Preprocessor
}

// NewDecodeEncoder creates new FIT decode-encoder.
func NewDecodeEncoder(preproc *activity.Preprocessor) *DecodeEncoder {
	return &DecodeEncoder{
		preprocessor: preproc,
	}
}

func (s *DecodeEncoder) Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	lis := filedef.NewListener()
	defer lis.Close()

	lis.Reset(filedef.WithFileFunc(typedef.FileActivity,
		func() filedef.File { return &wrapActivity{activity: filedef.NewActivity()} }))

	dec := decoderPool.Get().(*decoder.Decoder)
	defer decoderPool.Put(dec)

	dec.Reset(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
		decoder.WithIgnoreChecksum(),
	)

	activities := make([]activity.Activity, 0, 1) // In most cases, 1 fit == 1 activity
	for dec.Next() {
		fileId, err := dec.PeekFileId()
		if err != nil {
			return nil, fmt.Errorf("could not peek: %w", err)
		}

		if fileId.Type != typedef.FileActivity {
			if err = dec.Discard(); err != nil {
				return nil, fmt.Errorf("could not discard: %w", err)
			}
			continue
		}

		_, err = dec.DecodeWithContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not decode: %w", err)
		}

		wa := lis.File().(*wrapActivity)
		if len(wa.activity.Records) == 0 {
			continue
		}

		activities = append(activities, s.convertToActivity(wa.activity))
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

func (s *DecodeEncoder) convertToActivity(activityFile *filedef.Activity) activity.Activity {
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
		Sports:            activityFile.Sports,
		SplitSummaries:    activityFile.SplitSummaries,
		Activity:          activityFile.Activity,
		UnrelatedMessages: activityFile.UnrelatedMessages,
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
func (s *DecodeEncoder) recalculateSummary(ses *activity.Session) {
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
		aggregator.Fill(lap.Lap, lapFromRecords.Lap)
		records = records[pos:]
	}
	sesFromLaps := activity.NewSessionFromLaps(ses.Laps)
	aggregator.Fill(ses, sesFromLaps.Session)
	ses.Summarize()
}

func (s *DecodeEncoder) Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error) {
	buf := mem.GetBuffer()
	defer mem.PutBuffer(buf)

	bufAt := &bytesBufferAt{buf}

	enc := encoderPool.Get().(*encoder.Encoder)
	defer encoderPool.Put(enc)

	bs := make([][]byte, len(activities))
	for i := range activities {
		a := &activities[i]
		s.makeLastSummary(a)

		wa := wrapActivity{activity: filedef.NewActivity()}
		wa.activity.FileId = *a.Creator.FileId
		for j := range a.Sessions {
			ses := &a.Sessions[j]
			for k := range ses.Laps {
				wa.activity.Laps = append(wa.activity.Laps, ses.Laps[k].Lap)
			}
			for k := range ses.Records {
				wa.activity.Records = append(wa.activity.Records, ses.Records[k].Record)
			}
			wa.activity.Sessions = append(wa.activity.Sessions, ses.Session)
		}
		wa.activity.SplitSummaries = a.SplitSummaries
		wa.activity.Activity = a.Activity
		wa.activity.UnrelatedMessages = a.UnrelatedMessages

		fit := wa.ToFIT(nil)

		enc.Reset(bufAt,
			encoder.WithProtocolVersion(proto.V2),
			encoder.WithHeaderOption(encoder.HeaderOptionNormal, 15),
		)
		if err := enc.EncodeWithContext(ctx, &fit); err != nil {
			return nil, fmt.Errorf("could not encode: %w", err)
		}
		bs[i] = slices.Clone(bufAt.Buffer.Bytes())
		bufAt.Buffer.Reset()
	}

	return bs, nil
}

func (s *DecodeEncoder) makeLastSummary(a *activity.Activity) {
	var lastTimestamp time.Time
	var totalTimerTime uint32
	for i := len(a.Sessions) - 1; i >= 0; i-- {
		ses := a.Sessions[i]
		totalTimerTime += ses.TotalTimerTime

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

	// Ensure we got the latest timestamp across all messages.
	lastTimestampUint32 := datetime.ToUint32(lastTimestamp)
	for i := range a.UnrelatedMessages {
		timestamp := a.UnrelatedMessages[i].FieldValueByNum(proto.FieldNumTimestamp).Uint32()
		if timestamp == basetype.Uint32Invalid {
			continue
		}
		if timestamp < lastTimestampUint32 {
			break
		}
		lastTimestamp = datetime.ToTime(timestamp) // We get latest timestamp
	}

	for i := range a.Sessions {
		a.Sessions[i].Timestamp = lastTimestamp
	}

	if a.Activity == nil {
		a.Activity = mesgdef.NewActivity(nil)
	}

	a.Activity.Timestamp = lastTimestamp
	a.Activity.LocalTimestamp = lastTimestamp.Add(time.Duration(a.Timezone) * time.Hour)
	a.Activity.TotalTimerTime = totalTimerTime
	a.Activity.Type = typedef.ActivityAutoMultiSport
	a.Activity.NumSessions = uint16(len(a.Sessions))
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
