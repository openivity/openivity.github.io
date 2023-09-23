package activity

type Feature struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties,omitempty"`
	Geometry   Geometry       `json:"geometry"`
}

func (m Feature) ToMap() map[string]any {
	return map[string]any{
		"type":       m.Type,
		"properties": m.Properties,
		"geometry":   m.Geometry.ToMap(),
	}
}

type Geometry struct {
	Type        string `json:"type"`
	Coordinates []any  `json:"coordinates"`
}

func (m Geometry) ToMap() map[string]any {
	return map[string]any{
		"type":        m.Type,
		"coordinates": m.Coordinates,
	}
}
