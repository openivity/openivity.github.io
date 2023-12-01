package activity

import (
	"time"

	"github.com/muktihari/openactivity-fit/kit"
)

const Unknown = "Unknown"

type Creator struct {
	Name         string
	Manufacturer *uint16
	Product      *uint16
	TimeCreated  time.Time
}

func (c *Creator) ToMap() map[string]any {
	m := map[string]any{}

	if c.Name == "" {
		c.Name = Unknown
	}
	m["name"] = c.Name
	if c.Manufacturer != nil {
		m["manufacturer"] = *c.Manufacturer
	}
	if c.Product != nil {
		m["product"] = *c.Product
	}
	if !c.TimeCreated.IsZero() {
		m["timeCreated"] = c.TimeCreated.Format(time.RFC3339)
	}

	return m
}

func (c *Creator) Clone() *Creator {
	cre := &Creator{
		Name:        c.Name,
		TimeCreated: c.TimeCreated,
	}

	if c.Manufacturer != nil {
		cre.Manufacturer = kit.Ptr(*c.Manufacturer)
	}
	if c.Product != nil {
		cre.Product = kit.Ptr(*c.Product)
	}

	return cre
}
