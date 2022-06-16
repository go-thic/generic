package stream

import "github.com/go-thic/generic/optional"

type LimitFunc[V VAL] func(elem V) bool

func withLimitFunc[V VAL](limit LimitFunc[VAL]) func(elem V) (optional.Optional[V], bool) {
	return func(elem V) (optional.Optional[V], bool) {
		isLimit := limit(elem)
		stopConsuming := isLimit
		return optional.New(elem, !isLimit), stopConsuming
	}
}

func Count(i int) func(_ VAL) bool {
	return func(_ VAL) bool {
		i--
		return i < 0
	}
}
