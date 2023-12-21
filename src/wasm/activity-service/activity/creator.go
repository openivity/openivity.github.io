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
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('{')

	buf.WriteString("\"name\":\"")
	buf.WriteString(c.Name)
	buf.WriteString("\",")

	if c.Manufacturer != nil {
		buf.WriteString("\"manufacturer\":")
		buf.WriteString(strconv.FormatUint(uint64(*c.Manufacturer), 10))
		buf.WriteByte(',')
	}
	if c.Product != nil {
		buf.WriteString("\"product\":")
		buf.WriteString(strconv.FormatUint(uint64(*c.Product), 10))
		buf.WriteByte(',')
	}
	if !c.TimeCreated.IsZero() {
		buf.WriteString("\"timeCreated\":\"")
		buf.WriteString(c.TimeCreated.Format(time.RFC3339))
		buf.WriteString("\"")
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
