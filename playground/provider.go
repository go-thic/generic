package playground

import (
	"github.com/go-thic/generic/optional"
	"log"
)

type ProviderFunc[V VAL] func() optional.Optional[V]

func NewProvider[V VAL](generate func() optional.Optional[V]) *Stream[V] {
	valueChan := make(chan V)
	stopChan := make(chan empty)

	providerStream := New(valueChan, stopChan)

	go func() {
		defer func() {
			close(valueChan)
			close(stopChan)
			if r := recover(); r != nil {
				log.Printf("panic: %q", r)
			}
		}()

	generating:
		for {
			select {
			case <-stopChan:
				break generating
			default:
				v := generate()
				if v.IsSome() {
					providerStream.Write(v.Val())
				} else {
					break generating
				}
			}
		}
	}()

	return providerStream
}

func WithValues[T VAL](s ...T) ProviderFunc[T] {
	a := make([]T, len(s))
	copy(a, s)

	return func() optional.Optional[T] {
		var next T
		if len(a) > 0 {
			next, a = a[0], a[1:]
			return optional.New(next, true)
		}
		return optional.None[T]()
	}
}
