package playground

import (
	"github.com/go-thic/generic/optional"
	"log"
)

type ConsumerFunc[V VAL] func(elem V) (optional.Optional[V], bool)

func NewConsumer[V VAL](s *Stream[V], consume ConsumerFunc[V]) *Stream[VAL] {
	valueChan := make(chan VAL)
	stopChan := make(chan empty)

	consumerStream := New(valueChan, stopChan)

	go func() {
		defer func() {
			close(valueChan)
			close(stopChan)
			if r := recover(); r != nil {
				log.Printf("panic: %q", r)
			}
		}()

		for v := range s.values {
			res, stopConsuming := consume(v)
			if res.IsSome() {
				consumerStream.Write(res.Val())
			}
			if stopConsuming {
				break
			}
		}
	}()

	return consumerStream
}
