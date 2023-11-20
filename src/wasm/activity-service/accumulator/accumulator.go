package accumulator

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Accumulator[T constraints.Integer | constraints.Float] struct {
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
