package activity

import "time"

type Event struct {
	Timestamp time.Time
	Event     uint8
	EventType uint8
}
