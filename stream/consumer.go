package stream

import (
	"github.com/go-thic/generic/optional"
	"log"
)

func NewConsumer[T any](s *Stream, consume func(elem T) (optional.Optional[T], bool)) *Stream {
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
			if val, ok := v.(T); ok {
				consumed, stopConsuming := consume(val)
				if consumed.IsSome() {
					consumerStream.Write(consumed.Val())
				}
				if stopConsuming {
					break
				}
			}
		}
	}()

	return consumerStream
}
