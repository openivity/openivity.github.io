package kit

import (
	"strings"
	"unicode"

	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Numeric interface {
	constraints.Integer | constraints.Float
}

// Ptr returns pointer of v
func Ptr[T any](v T) *T { return &v }

// PickNonZeroValue returns x if x != 0, otherwise return y.
func PickNonZeroValue[T Numeric](x, y T) T {
	if x == 0 {
		return y
	}
	return x
}

// PickNonZeroValuePtr same as PickNonZeroValue but for pointers and return a pointer.
func PickNonZeroValuePtr[T Numeric](x, y *T) *T {
	if x == nil {
		return y
	}
	if *x == 0 {
		return y
	}
	return x
}

// FormatTitle returns init capital for every word. "snow boarding", "snow_boarding", "SNOW_boardinG" -> "Show Boarding".
func FormatTitle(s string) string {
	s = strings.Map(func(r rune) rune {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ' '
		}
		return r
	}, s)
	s = cases.Title(language.English).String(s)
	return s
}

// FormatToLowerSnakeCase returns all lower case with word separated by underscore `_`. "Snow Boarding", "snow boarding", "SNOW boardinG" -> "show_boarding".
func FormatToLowerSnakeCase(s string) string {
	s = strings.Map(func(r rune) rune {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return '_'
		}
		return r
	}, s)
	s = cases.Lower(language.English).String(s)
	return s
}

// Avg returns average of two non-nil values. Otherwise, return any non-nil value if possible.
func Avg[T Numeric](a, b *T) *T {
	if a == nil {
		return b
	}
	if b != nil {
		c := T((float64(*a) + float64(*b)) / 2)
		return &c
	}
	return a
}

// Max returns max of two non-nil values. Otherwise, return any non-nil value if possible.
func Max[T Numeric](x, y *T) *T {
	if x == nil {
		return y
	}
	if y != nil && *x > *y {
		return x
	}
	return y
}
