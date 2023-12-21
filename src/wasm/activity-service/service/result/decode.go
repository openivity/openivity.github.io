package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/muktihari/openactivity-fit/activity"
)

type Decode struct {
	Err               error
	DecodeTook        time.Duration
	SerializationTook time.Duration
	TotalElapsed      time.Duration
	Activities        []activity.Activity
}

var _ json.Marshaler = &Decode{}

func (d *Decode) MarshalJSON() ([]byte, error) {
	if d.Err != nil {
		return []byte(fmt.Sprintf("{%q:%q}", "err", d.Err)), nil
	}

	begin := time.Now()
	buf := new(bytes.Buffer)
	buf.WriteByte('{')

	buf.WriteString("\"err\":null,")

	if len(d.Activities) != 0 {
		buf.WriteString("\"activities\":[")
		for i := range d.Activities {
			b, _ := d.Activities[i].MarshalJSON()
			buf.Write(b)
			if i != len(d.Activities)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteString("],")
	}

	d.SerializationTook = time.Since(begin)
	d.TotalElapsed = d.DecodeTook + d.SerializationTook

	buf.WriteString("\"decodeTook\":")
	buf.WriteString(strconv.FormatInt(d.DecodeTook.Milliseconds(), 10))
	buf.WriteByte(',')

	buf.WriteString("\"serializationTook\":")
	buf.WriteString(strconv.FormatInt(d.SerializationTook.Milliseconds(), 10))
	buf.WriteByte(',')

	buf.WriteString("\"totalElapsed\":")
	buf.WriteString(strconv.FormatInt(d.TotalElapsed.Milliseconds(), 10))

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

type DecodeWorker struct {
	Err      error
	Index    int
	Activity *activity.Activity
}
