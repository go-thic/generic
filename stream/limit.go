package stream

import "github.com/go-thic/generic/optional"

func withLimitFunc[V any](limit func(elem V) bool) func(elem V) (optional.Optional[V], bool) {
	return func(elem V) (optional.Optional[V], bool) {
		isLimit := limit(elem)
		stopConsuming := isLimit
		return optional.New(elem, !isLimit), stopConsuming
	}
}

func Count(i int) func(elem VAL) bool {
	return func(_ VAL) bool {
		i--
		return i < 0
	}
}
