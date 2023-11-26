package kit

import (
	"strings"
	"unicode"

	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
