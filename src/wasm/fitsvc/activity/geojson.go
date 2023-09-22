package activity

type Feature struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties,omitempty"`
	Geometry   Geometry       `json:"geometry"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}
