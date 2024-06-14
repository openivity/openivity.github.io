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

package gpx

import (
	"context"
	"fmt"
	"io"

	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/xmltokenizer"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/activity/gpx/schema"
	"github.com/openivity/activity-service/mem"
	"github.com/openivity/activity-service/strutils"
	"github.com/openivity/activity-service/xmlutils"
	"golang.org/x/exp/slices"
)

const (
	metadataDesc = "The GPX file is created by openivity.github.io"
	metadataLink = "https://openivity.github.io"
)

var _ activity.Service = (*service)(nil)

type service struct {
	preprocessor *activity.Preprocessor
}

// NewService creates new GPX service.
func NewService(preproc *activity.Preprocessor) activity.Service {
	return &service{preprocessor: preproc}
}

func (s *service) Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	tok := xmltokenizer.New(r)

	var gpx schema.GPX
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
		case "gpx":
			se := xmltokenizer.GetToken().Copy(token)
			err = gpx.UnmarshalToken(tok, se)
			xmltokenizer.PutToken(se)
			if err != nil {
				return nil, err
			}
			break loop
		}
	}

	act := activity.CreateActivity()
	act.Creator.Name = gpx.Creator
	act.Creator.TimeCreated = gpx.Metadata.Time

	sessions := make([]activity.Session, 0, len(gpx.Tracks))

	for i := range gpx.Tracks { // Sessions
		trk := gpx.Tracks[i]

		sport := typedef.SportFromString(strutils.ToLowerSnakeCase(trk.Type))
		if sport == typedef.SportInvalid {
			sport = typedef.SportGeneric
		}

		var recordCount int
		for i := range trk.TrackSegments {
			recordCount += len(trk.TrackSegments[i].Trackpoints)
		}

		laps := make([]activity.Lap, 0, len(trk.TrackSegments))
		records := make([]activity.Record, 0, recordCount)
		recordsByLap := make([][]activity.Record, 0, len(trk.TrackSegments))
		for j := range trk.TrackSegments { // Laps
			trkseg := trk.TrackSegments[j]

			lapRecords := make([]activity.Record, 0, len(trkseg.Trackpoints))
			for k := range trkseg.Trackpoints { // Records
				trkpt := &trkseg.Trackpoints[k]
				lapRecords = append(lapRecords, trkpt.ToRecord())
			}

			if len(lapRecords) == 0 {
				continue
			}

			recordsByLap = append(recordsByLap, lapRecords)
			records = append(records, lapRecords...)
		}

		// Preprocessing...
		s.preprocessor.CalculateDistanceAndSpeed(records)
		if activity.HasPace(sport) {
			s.preprocessor.CalculatePace(sport, records)
		}
		s.preprocessor.SmoothingElevation(records)
		s.preprocessor.CalculateGrade(records)

		// We can only calculate laps' summary after preprocessing.
		for i := range recordsByLap {
			lap := activity.NewLapFromRecords(recordsByLap[i], sport)
			laps = append(laps, lap)
		}

		if len(laps) == 0 {
			continue
		}

		session := activity.NewSessionFromLaps(laps)
		session.Records = records
		session.Laps = laps
		session.Summarize()

		sessions = append(sessions, session)

		if act.Creator.TimeCreated.IsZero() {
			act.Creator.TimeCreated = session.StartTime
		}
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("gpx: %w", activity.ErrNoActivity)
	}

	act.Sessions = sessions

	return []activity.Activity{act}, nil
}

func (s *service) Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error) {
	bs := make([][]byte, len(activities))

	buf := mem.GetBuffer()
	defer mem.PutBuffer(buf)

	for i := range activities {
		gpx := s.convertActivityToGPX(&activities[i])
		if err := gpx.Validate(); err != nil {
			return nil, fmt.Errorf("invalid gpx: %w", err)
		}
		buf.Reset()
		if err := xmlutils.MarshalWrite(buf, &gpx); err != nil {
			return nil, fmt.Errorf("could not marshal gpx[%d]: %w", i, err)
		}
		bs[i] = slices.Clone(buf.Bytes())
	}

	return bs, nil
}

func (s *service) convertActivityToGPX(act *activity.Activity) schema.GPX {
	gpx := schema.GPX{
		Creator: act.Creator.Name,
		Metadata: schema.Metadata{
			Time: act.Creator.TimeCreated,
			Desc: metadataDesc,
			Link: &schema.Link{Href: metadataLink},
		},
		Tracks: make([]schema.Track, 0, len(act.Sessions)),
	}

	for i := range act.Sessions {
		ses := &act.Sessions[i]
		track := schema.Track{
			Name:          strutils.ToTitle(ses.Sport.String()),
			Type:          strutils.ToTitle(ses.Sport.String()),
			TrackSegments: make([]schema.TrackSegment, 0, len(ses.Laps)),
		}

		sesRecords := ses.Records
		for j := range ses.Laps {
			lap := &ses.Laps[j]
			trackSegment := schema.TrackSegment{}

			remainingRecords := make([]activity.Record, 0)
			for k := range sesRecords {
				rec := sesRecords[k]

				if lap.IsBelongToThisLap(rec.Timestamp) {
					waypoint := schema.Waypoint{
						Time: rec.Timestamp,
						Lat:  rec.PositionLatDegrees(),
						Lon:  rec.PositionLongDegrees(),
						Ele:  rec.AltitudeScaled(),
						TrackPointExtension: schema.TrackPointExtension{
							Cadence:     rec.Cadence,
							Distance:    rec.DistanceScaled(),
							HeartRate:   rec.HeartRate,
							Temperature: rec.Temperature,
							Power:       rec.Power,
						},
					}
					trackSegment.Trackpoints = append(trackSegment.Trackpoints, waypoint)
				} else {
					remainingRecords = append(remainingRecords, rec)
				}
			}
			sesRecords = remainingRecords

			track.TrackSegments = append(track.TrackSegments, trackSegment)
		}

		gpx.Tracks = append(gpx.Tracks, track)
	}

	return gpx
}
