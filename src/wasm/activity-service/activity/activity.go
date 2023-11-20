package activity

type Activity struct {
	Creator  Creator
	Timezone int
	Sessions []*Session
	Laps     []*Lap
	Records  []*Record
}

func (a *Activity) ToMap() map[string]any {
	m := map[string]any{
		"timezone": a.Timezone,
	}

	m["creator"] = a.Creator.ToMap()

	sessions := make([]any, 0, len(a.Sessions))
	for i := range a.Sessions {
		sessions = append(sessions, a.Sessions[i].ToMap())
	}
	m["sessions"] = sessions

	laps := make([]any, 0, len(a.Laps))
	for i := range a.Laps {
		laps = append(laps, a.Laps[i].ToMap())
	}
	m["laps"] = laps

	records := make([]any, 0, len(a.Records))
	for i := range a.Records {
		records = append(records, a.Records[i].ToMap())
	}
	m["records"] = records

	return m
}
