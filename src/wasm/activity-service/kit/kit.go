package kit

import "golang.org/x/exp/constraints"

// Ptr returns pointer of v
func Ptr[T any](v T) *T { return &v }

// PickNonZeroValue returns x if x != 0, otherwise return y.
func PickNonZeroValue[T constraints.Integer | constraints.Float](x, y T) T {
	if x == 0 {
		return y
	}
	return x
}

// PickNonZeroValuePtr same as PickNonZeroValue but for pointers and return a pointer.
func PickNonZeroValuePtr[T constraints.Integer | constraints.Float](x, y *T) *T {
	if x == nil {
		return y
	}
	if *x == 0 {
		return y
	}
	return x
}
