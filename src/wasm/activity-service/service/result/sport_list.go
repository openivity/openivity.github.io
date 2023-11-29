package result

type SportList struct {
	Sports []Sport
}

type Sport struct {
	ID   uint8
	Name string
}

func (s SportList) ToMap() map[string]any {
	sports := make([]any, len(s.Sports))
	for i := range s.Sports {
		sports[i] = map[string]any{
			"id":   uint8(s.Sports[i].ID),
			"name": s.Sports[i].Name,
		}
	}
	return map[string]any{
		"sports": sports,
	}
}
