package kit

// Ptr returns pointer of v
func Ptr[T any](v T) *T { return &v }
