// Copyright (C) 2023 Openivity

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

package accumulator

import (
	"math"

	"github.com/muktihari/openactivity-fit/kit"
)

type Accumulator[T kit.Numeric] struct {
	min     *T
	max     *T
	sum     float64 // we don't need that precision
	counter int
}

func (a *Accumulator[T]) Collect(v *T) {
	if v == nil {
		return
	}

	if math.IsInf(float64(*v), 0) {
		return
	}

	if a.min == nil {
		a.min = v
	}
	if a.max == nil {
		a.max = v
	}

	if *v < *a.min {
		a.min = v
	}
	if *v > *a.max {
		a.max = v
	}

	a.sum += float64(*v)
	a.counter++
}

func (a *Accumulator[T]) Min() *T {
	return a.min
}

func (a *Accumulator[T]) Max() *T {
	return a.max
}

func (a *Accumulator[T]) Avg() *T {
	if a.counter == 0 {
		return nil
	}

	avg := T(a.sum / float64(a.counter))

	return &avg
}

func (a *Accumulator[T]) Sum() *T {
	if a.counter == 0 {
		return nil
	}
	sum := T(a.sum)
	return &sum
}
