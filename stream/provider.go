package stream

import (
	"github.com/go-thic/generic/optional"
)

func WithValues[T any](s ...T) func() optional.Optional[T] {
	a := make([]T, len(s))
	copy(a, s)

	return func() optional.Optional[T] {
		if len(a) > 0 {
			var next T
			next, a = a[0], a[1:]
			return optional.New(next, true)
		}
		return optional.None[T]()
	}
}

func WithGenerator[T any](next func() T) func() optional.Optional[T] {
	return func() optional.Optional[T] {
		return optional.New[T](next(), true)
	}
}

type Countable interface {
	int64 | int32 | int | int16 | int8 | float64 | float32
}

func StartCountingFrom[T Countable](start T) func() optional.Optional[T] {
	next := start
	return func() optional.Optional[T] {
		ret := next
		next++
		return optional.New[T](ret, true)
	}
}
