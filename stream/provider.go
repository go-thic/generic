package stream

import (
	"github.com/go-thic/generic/optional"
	"log"
)

type ProviderFunc func() optional.Optional[VAL]

func NewProvider(generate ProviderFunc) *Stream {
	valueChan := make(chan any)

	providerStream := newImpl(valueChan)

	go func() {
		defer func() {
			close(valueChan)
			if r := recover(); r != nil {
				log.Printf("panic: %q", r)
			}
		}()

	generating:
		for {
			v := generate()
			if v.IsSome() {
				providerStream.Write(v.Val())
			} else {
				break generating
			}
		}
	}()

	return providerStream
}

func WithValues(s ...VAL) ProviderFunc {
	a := make([]VAL, len(s))
	copy(a, s)

	return func() optional.Optional[VAL] {
		if len(a) > 0 {
			var next VAL
			next, a = a[0], a[1:]
			return optional.New(next, true)
		}
		return optional.None[VAL]()
	}
}

type Countable interface {
	int64 | int32 | int | int16 | int8 | float64 | float32
}

func StartCountingFrom[T Countable](start T) ProviderFunc {
	next := start
	return func() optional.Optional[VAL] {
		ret := next
		next++
		return optional.New[VAL](ret, true)
	}
}
