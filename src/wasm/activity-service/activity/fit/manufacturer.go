package fit

type Manufacturer struct {
	ID       uint16
	Name     string
	Products []ManufacturerProduct
}

func (m *Manufacturer) ToMap() map[string]any {
	products := make([]any, len(m.Products))
	for i := range m.Products {
		products[i] = m.Products[i].ToMap()
	}
	return map[string]any{
		"id":       uint16(m.ID),
		"name":     m.Name,
		"products": products,
	}
}

type ManufacturerProduct struct {
	ID   uint16
	Name string
}

func (p *ManufacturerProduct) ToMap() map[string]any {
	return map[string]any{
		"id":   p.ID,
		"name": p.Name,
	}
}
