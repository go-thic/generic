package stream

import (
	"github.com/go-thic/generic/optional"
	"log"
)

type ConsumerFunc[V VAL] func(elem V) (optional.Optional[V], bool)

func NewConsumer(s *Stream, consume ConsumerFunc[VAL]) *Stream {
	valueChan := make(chan any)

	consumerStream := newImpl(valueChan)

	go func() {
		defer func() {
			close(valueChan)
			if r := recover(); r != nil {
				log.Printf("panic: %q", r)
			}
		}()

		for v := range s.values {
			consumed, stopConsuming := consume(v)
			if consumed.IsSome() {
				consumerStream.Write(consumed.Val())
			}
			if stopConsuming {
				break
			}
		}
	}()

	return consumerStream
}
