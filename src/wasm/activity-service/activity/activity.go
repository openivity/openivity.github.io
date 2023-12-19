package activity

import (
	"bytes"
	"strconv"
)

type Activity struct {
	Creator  Creator
	Timezone int
	Sessions []*Session
}

func (a *Activity) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteByte('{')

	buf.WriteString("\"creator\":")
	b, _ := a.Creator.MarshalJSON()
	buf.Write(b)
	buf.WriteByte(',')

	buf.WriteString("\"timezone\":" + strconv.FormatInt(int64(a.Timezone), 10) + ",")

	if len(a.Sessions) != 0 {
		buf.WriteString("\"sessions\":[")
		for i := range a.Sessions {
			b, _ = a.Sessions[i].MarshalJSON()
			buf.Write(b)
			if i != len(a.Sessions)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte(']')
	}

	b = buf.Bytes()
	if b[len(b)-1] == ',' {
		b[len(b)-1] = '}'
		return b, nil
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (a *Activity) Clone() *Activity {
	act := &Activity{
		Creator:  *a.Creator.Clone(),
		Timezone: a.Timezone,
	}

	act.Sessions = make([]*Session, len(a.Sessions))
	for i := range a.Sessions {
		act.Sessions[i] = a.Sessions[i].Clone()
	}

	return act
}
