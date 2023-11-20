package activity

import (
	"time"
)

type Creator struct {
	Name         string
	Manufacturer *uint16
	Product      *uint16
	TimeCreated  time.Time
}

func (c *Creator) ToMap() map[string]any {
	m := map[string]any{}

	if c.Name == "" {
		c.Name = "Unknown"
	}
	m["name"] = c.Name
	if c.Manufacturer != nil {
		m["manufacturer"] = *c.Manufacturer
	}
	if c.Product != nil {
		m["product"] = *c.Product
	}
	if c.TimeCreated != (time.Time{}) {
		m["timeCreated"] = c.TimeCreated.Format(time.RFC3339)
	}

	return m
}
