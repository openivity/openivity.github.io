package activity

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func IsConsideredMoving(sport string, speed *float64) bool {
	if speed == nil {
		return false
	}

	switch sport {
	case "Hiking", "Running", "Walking", "Swimming": // slow moving exercise
		return *speed > 0.2778 // 1km/h ≈ 0.2778 m/s
	default:
		return *speed > 1.6
	}
}

func HasPace(sport string) bool {
	switch sport {
	case "Hiking", "Walking", "Running":
		return true
	default:
		return false
	}
}

func FormatSport(s string) string {
	s = strings.Map(func(r rune) rune {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ' '
		}
		return r
	}, s)
	s = cases.Title(language.English).String(s)
	return s
}
