package playground

import "github.com/go-thic/generic/optional"

type LimitFunc[V VAL] func(elem V) bool

func WithLimit[V VAL](limit LimitFunc[VAL]) func(elem V) (optional.Optional[V], bool) {
	return func(elem V) (optional.Optional[V], bool) {
		isLimit := limit(elem)
		stopConsuming := isLimit
		return optional.New(elem, !isLimit), stopConsuming
	}
}

func Count[V any](i int) func(_ V) bool {
	return func(_ V) bool {
		i--
		return i < 0
	}
}
