package activity

import (
	"bytes"
	"strconv"
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

func (c *Creator) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('{')

	buf.WriteString("\"name\":\"" + c.Name + "\",")
	if c.Manufacturer != nil {
		buf.WriteString("\"manufacturer\":" + strconv.FormatUint(uint64(*c.Manufacturer), 10) + ",")
	}
	if c.Product != nil {
		buf.WriteString("\"product\":" + strconv.FormatUint(uint64(*c.Product), 10) + ",")
	}
	if !c.TimeCreated.IsZero() {
		buf.WriteString("\"timeCreated\":\"" + c.TimeCreated.Format(time.RFC3339) + "\"")
	}

	b := buf.Bytes()
	if b[len(b)-1] == ',' {
		b[len(b)-1] = '}'
		return b, nil
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
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
