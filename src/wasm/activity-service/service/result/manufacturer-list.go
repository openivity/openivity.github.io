package result

import "github.com/muktihari/openactivity-fit/activity/fit"

type ManufacturerList struct {
	Manufacturers []fit.Manufacturer
}

func (m ManufacturerList) ToMap() map[string]any {
	manufacturers := make([]any, len(m.Manufacturers))
	for i := range m.Manufacturers {
		manufacturers[i] = m.Manufacturers[i].ToMap()
	}
	return map[string]any{
		"manufacturers": manufacturers,
	}
}
