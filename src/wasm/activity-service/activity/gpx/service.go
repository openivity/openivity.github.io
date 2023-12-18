package gpx

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/gpx/schema"
	"github.com/muktihari/openactivity-fit/kit"
	kxml "github.com/muktihari/openactivity-fit/kit/xml"
	"github.com/muktihari/openactivity-fit/preprocessor"
)

const (
	metadataDesc = "The GPX file is created by openivity.github.io"
	metadataLink = "https://openivity.github.io"
)

var _ activity.Service = &service{}

type service struct {
	preprocessor *preprocessor.Preprocessor
}

func NewService(preproc *preprocessor.Preprocessor) activity.Service {
	return &service{preprocessor: preproc}
}

func (s *service) Decode(ctx context.Context, r io.Reader) ([]activity.Activity, error) {
	dec := xml.NewDecoder(r)

	// NOTE: We manually define xml.Unmarshaler to void the reflection used
	//       in default decoder, which is particularly slow expecially in WASM.

	// Find start element
	var start *xml.StartElement
	if start == nil {
		for {
			tok, err := dec.Token()
			if err != nil {
				return nil, err
			}
			if t, ok := tok.(xml.StartElement); ok {
				start = &t
				break
			}
		}
	}

	if start == nil { // In case we got invalid xml, avoid panic.
		return nil, fmt.Errorf("not a valid xml")
	}

	var gpx schema.GPX
	if err := gpx.UnmarshalXML(dec, *start); err != nil {
		return nil, err
	}

	act := new(activity.Activity)
	act.Creator.Name = gpx.Creator
	act.Creator.TimeCreated = gpx.Metadata.Time

	sessions := make([]*activity.Session, 0, len(gpx.Tracks))

	for i := range gpx.Tracks { // Sessions
		trk := gpx.Tracks[i]

		sport := kit.FormatTitle(trk.Type)
		if sport == "" || sport == "Other" {
			sport = activity.SportGeneric
		}

		laps := make([]*activity.Lap, 0, len(trk.TrackSegments))

		var recordCount int
		for i := range trk.TrackSegments {
			recordCount += len(trk.TrackSegments[i].Trackpoints)
		}
		records := make([]*activity.Record, 0, recordCount)

		recordsByLap := make([][]*activity.Record, 0, len(trk.TrackSegments))
		for j := range trk.TrackSegments { // Laps
			trkseg := trk.TrackSegments[j]

			lapRecords := make([]*activity.Record, 0, len(trkseg.Trackpoints))
			for k := range trkseg.Trackpoints { // Records
				trkpt := &trkseg.Trackpoints[k]
				rec := NewRecord(trkpt)
				lapRecords = append(lapRecords, rec)
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
		s.preprocessor.SmoothingElev(records)
		s.preprocessor.CalculateGrade(records)

		// We can only calculate laps' summary after preprocessing.
		for i := range recordsByLap {
			lap := activity.NewLapFromRecords(recordsByLap[i], sport)
			laps = append(laps, lap)
		}

		if len(laps) == 0 {
			continue
		}

		session := activity.NewSessionFromLaps(laps, sport)
		session.Records = records
		session.Laps = laps

		sessions = append(sessions, session)

		if act.Creator.TimeCreated.IsZero() {
			act.Creator.TimeCreated = session.StartTime
		}
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("gpx: %w", activity.ErrNoActivity)
	}

	act.Sessions = sessions

	s.preprocessor.SetSessionsWorkoutType(act.Sessions...)

	return []activity.Activity{*act}, nil
}

func (s *service) Encode(ctx context.Context, activities []activity.Activity) ([][]byte, error) {
	bs := make([][]byte, len(activities))

	for i := range activities {
		gpx := s.convertActivityToGPX(&activities[i])
		if err := gpx.Validate(); err != nil {
			return nil, fmt.Errorf("invalid gpx: %w", err)
		}
		b, err := kxml.Marshal(gpx)
		if err != nil {
			return nil, fmt.Errorf("could not marshal gpx[%d]: %w", i, err)
		}
		bs[i] = b
	}

	return bs, nil
}

func (s *service) convertActivityToGPX(act *activity.Activity) *schema.GPX {
	gpx := &schema.GPX{
		Creator: act.Creator.Name,
		Metadata: schema.Metadata{
			Time: act.Creator.TimeCreated,
			Desc: metadataDesc,
			Link: &schema.Link{Href: metadataLink},
		},
		Tracks: make([]schema.Track, 0, len(act.Sessions)),
	}

	for i := range act.Sessions {
		ses := act.Sessions[i]

		track := schema.Track{
			Name:          ses.Sport,
			Type:          ses.Sport,
			TrackSegments: make([]schema.TrackSegment, 0, len(ses.Laps)),
		}

		sesRecords := ses.Records
		for j := range ses.Laps {
			lap := ses.Laps[j]
			trackSegment := schema.TrackSegment{}

			remainingRecords := make([]*activity.Record, 0)
			for k := range sesRecords {
				rec := sesRecords[k]

				if lap.IsBelongToThisLap(rec.Timestamp) {
					waypoint := schema.Waypoint{
						Time: rec.Timestamp,
						Lat:  rec.PositionLat,
						Lon:  rec.PositionLong,
						Ele:  rec.Altitude,
						TrackPointExtension: &schema.TrackPointExtension{
							Cadence:     rec.Cadence,
							Distance:    rec.Distance,
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
