// Copyright (C) 2024 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package strutils_test

import (
	"testing"

	"github.com/openivity/activity-service/strutils"
)

func TestToTitle(t *testing.T) {
	tt := []struct {
		in  string
		out string
	}{
		{in: "snow_boarding", out: "Snow Boarding"},
		{in: "snow boarding", out: "Snow Boarding"},
		{in: "snow_boardinG", out: "Snow Boarding"},
	}

	for _, tc := range tt {
		s := strutils.ToTitle(tc.in)
		if s != tc.out {
			t.Fatalf("expected: %q, got: %q", tc.out, s)
		}
	}
}

func TestToLowerSnakeCase(t *testing.T) {
	tt := []struct {
		in  string
		out string
	}{
		{in: "Snow Boarding", out: "snow_boarding"},
		{in: "Snow_Boarding", out: "snow_boarding"},
		{in: "snow boarding", out: "snow_boarding"},
	}

	for _, tc := range tt {
		s := strutils.ToLowerSnakeCase(tc.in)
		if s != tc.out {
			t.Fatalf("expected: %q, got: %q", tc.out, s)
		}
	}
}
