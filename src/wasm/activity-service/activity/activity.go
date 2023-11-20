package activity

type Activity struct {
	Creator  Creator
	Timezone int
	Sessions []*Session
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

	return m
}
