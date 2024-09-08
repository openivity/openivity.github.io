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

package tcx

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/muktihari/fit/kit/scaleoffset"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/activity/tcx/schema"
	"github.com/openivity/activity-service/aggregator"
	"github.com/openivity/activity-service/mem"
	"github.com/openivity/activity-service/service"
	"github.com/openivity/activity-service/strutils"
	"github.com/openivity/activity-service/xmlutils"
	"golang.org/x/exp/slices"
)

const (
	applicationName = "openitivy.github.io"
)

var _ service.DecodeEncoder = (*DecodeEncoder)(nil)

type DecodeEncoder struct {
	preprocessor *activity.Preprocessor
}

// NewDecodeEncoder creates new TCX decode-encoder.
func NewDecodeEncoder(preproc *activity.Preprocessor) *DecodeEncoder {
	return &DecodeEncoder{preprocessor: preproc}
}

func (s *DecodeEncoder) Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	tok := xmltokenizer.New(r)

	var tcx schema.TCX
loop:
	for {
		token, err := tok.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch string(token.Name.Local) {
		case "TrainingCenterDatabase":
			se := xmltokenizer.GetToken().Copy(token)
			err = tcx.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return nil, err
			}
			break loop
		}
	}

	act := activity.CreateActivity()
	if len(tcx.Activities) > 0 {
		act.Creator.TimeCreated = tcx.Activities[0].Activity.ID
	}

	sessions := make([]activity.Session, 0, len(tcx.Activities))

	for i := range tcx.Activities {
		a := tcx.Activities[i]

		if act.Creator.Name == "" && a.Activity.Creator != nil {
			act.Creator.Name = a.Activity.Creator.Name
		}

		sport := typedef.SportFromString(strutils.ToLowerSnakeCase(a.Activity.Sport))
		if sport == typedef.SportInvalid {
			sport = typedef.SportGeneric
		}

		var recordCount int
		for j := range a.Activity.Laps {
			for k := range a.Activity.Laps[j].Tracks {
				recordCount += len(a.Activity.Laps[j].Tracks[k].Trackpoints)
			}
		}

		laps := make([]activity.Lap, 0, len(a.Activity.Laps))
		records := make([]activity.Record, 0, recordCount)
		recordsByLap := make([][]activity.Record, 0, len(a.Activity.Laps))
		for j := range a.Activity.Laps {
			activityLap := &a.Activity.Laps[j]

			var lapRecordCount int
			for k := range activityLap.Tracks {
				lapRecordCount += len(activityLap.Tracks[k].Trackpoints)
			}
			lapRecords := make([]activity.Record, 0, lapRecordCount)

			for k := range activityLap.Tracks { // flattening tracks-trackpoints
				for l := range activityLap.Tracks[k].Trackpoints {
					trackpoint := &activityLap.Tracks[k].Trackpoints[l]
					lapRecords = append(lapRecords, trackpoint.ToRecord())
				}
			}

			if len(lapRecords) == 0 {
				continue
			}

			records = append(records, lapRecords...)
			recordsByLap = append(recordsByLap, lapRecords)

			lap := activity.CreateLap(nil)
			lap.StartTime = activityLap.StartTime
			if !math.IsNaN(activityLap.DistanceMeters) {
				lap.TotalDistance = uint32(scaleoffset.Discard(activityLap.DistanceMeters, 100, 0))
			}
			lap.TotalCalories = activityLap.Calories
			if !math.IsNaN(activityLap.DistanceMeters) {
				lap.TotalElapsedTime = uint32(scaleoffset.Discard(activityLap.TotalTimeSeconds, 1000, 0))
			}
			lap.AvgHeartRate = activityLap.AverageHeartRateBpm
			lap.MaxHeartRate = activityLap.MaximumHeartRateBpm

			laps = append(laps, lap)
		}

		// Preprocessing...
		s.preprocessor.CalculateDistanceAndSpeed(records)
		if activity.HasPace(sport) {
			s.preprocessor.CalculatePace(sport, records)
		}

		s.preprocessor.SmoothingElevation(records)
		s.preprocessor.CalculateGrade(records)

		// We can only calculate laps' summary after preprocessing
		for i := range laps {
			lap := &laps[i]
			lapFromRecords := activity.NewLapFromRecords(recordsByLap[i], sport)
			aggregator.Fill(lap.Lap, lapFromRecords.Lap)
		}

		if len(laps) == 0 {
			continue
		}

		session := activity.NewSessionFromLaps(laps)
		if !a.Activity.ID.IsZero() {
			session.StartTime = a.Activity.ID
		}

		session.Laps = laps
		session.Records = records
		session.Summarize()

		sessions = append(sessions, session)

		if act.Creator.TimeCreated.IsZero() {
			act.Creator.TimeCreated = session.StartTime
			act.Creator.Name = a.Activity.Creator.Name
			act.Creator.Product = a.Activity.Creator.ProductID
		}
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("tcx: %w", activity.ErrNoActivity)
	}

	act.Sessions = sessions

	return []activity.Activity{act}, nil
}

func (s *DecodeEncoder) Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error) {
	bs := make([][]byte, len(activities))

	buf := mem.GetBuffer()
	defer mem.PutBuffer(buf)

	for i := range activities {
		tcx := s.convertActivityToTCX(&activities[i])
		buf.Reset()
		if err := xmlutils.MarshalWrite(buf, &tcx); err != nil {
			return nil, fmt.Errorf("could not marshal tcx: %w", err)
		}
		bs[i] = slices.Clone(buf.Bytes())
	}

	return bs, nil
}

func (s *DecodeEncoder) convertActivityToTCX(act *activity.Activity) schema.TCX {
	tcx := schema.TCX{
		Author: &schema.Application{
			Name: applicationName,
		},
		Activities: make([]schema.ActivityList, 0, len(act.Sessions)),
	}

	for i := range act.Sessions {
		ses := act.Sessions[i]

		activityList := schema.ActivityList{
			Activity: schema.Activity{
				ID:    ses.Timestamp,
				Sport: strutils.ToTitle(ses.Sport.String()),
				Creator: &schema.Device{
					Name: act.Creator.Name,
				},
			},
		}

		if activityList.Activity.ID.IsZero() {
			activityList.Activity.ID = ses.StartTime
		}

		sesRecords := ses.Records
		for j := range ses.Laps {
			lap := ses.Laps[j]

			activityLap := schema.ActivityLap{
				StartTime:           lap.StartTime,
				TotalTimeSeconds:    lap.TotalElapsedTimeScaled(),
				DistanceMeters:      lap.TotalDistanceScaled(),
				MaximumSpeed:        lap.MaxSpeedScaled(),
				Calories:            lap.TotalCalories,
				AverageHeartRateBpm: lap.AvgHeartRate,
				MaximumHeartRateBpm: lap.MaxHeartRate,
				Cadence:             lap.AvgCadence,
			}

			track := schema.Track{}
			remainingRecords := make([]activity.Record, 0)
			for k := range sesRecords {
				rec := sesRecords[k]

				if lap.IsBelongToThisLap(rec.Timestamp) {
					trackpoint := schema.Trackpoint{
						Time: rec.Timestamp,
						Position: schema.Position{
							LatitudeDegrees:  rec.PositionLatDegrees(),
							LongitudeDegrees: rec.PositionLongDegrees(),
						},
						AltitudeMeters: rec.AltitudeScaled(),
						DistanceMeters: rec.DistanceScaled(),
						HeartRateBpm:   rec.HeartRate,
						Cadence:        rec.Cadence,
						Extensions: schema.TrackpointExtension{
							Speed: rec.SpeedScaled(),
						},
					}
					track.Trackpoints = append(track.Trackpoints, trackpoint)
				} else {
					remainingRecords = append(remainingRecords, rec)
				}
			}
			sesRecords = remainingRecords

			activityLap.Tracks = []schema.Track{track}
			activityList.Activity.Laps = append(activityList.Activity.Laps, activityLap)
		}

		tcx.Activities = append(tcx.Activities, activityList)
	}

	return tcx
}
