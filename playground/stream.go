package playground

import (
	"github.com/go-thic/generic/optional"
	"sync"
)

type empty interface{}

type VAL interface {
	any
}

type Stream[V VAL] struct {
	sync.Mutex

	values chan V
	stop   chan<- empty

	running bool
}

func New[V VAL](values chan V, stop chan<- empty) *Stream[V] {
	return &Stream[V]{
		values:  values,
		stop:    stop,
		running: true,
	}
}

func (s *Stream[V]) Stop() {
	s.Lock()
	defer s.Unlock()

	if s.running {
		s.stop <- nil
		s.running = false
	}
}

func (s *Stream[V]) Write(val V) {
	if s.running {
		s.values <- val
	}
}

func (s *Stream[V]) Limit(limit func(elem VAL) bool) *Stream[VAL] {
	return NewConsumer[V](s, WithLimit[V](limit))
}

func (s *Stream[V]) Filter(filter func(VAL) bool) *Stream[VAL] {
	return NewConsumer(s, func(elem V) (optional.Optional[V], bool) {
		return optional.New[V](elem, !filter(elem)), false
	})
}

func (s *Stream[V]) Map(doMap func(elem V) (VAL, bool)) *Stream[VAL] {
	return NewMapper(s, func(elem V) (optional.Optional[VAL], bool) {
		var mapped, stopMapping = doMap(elem)
		return optional.New[VAL](mapped, !stopMapping), stopMapping
	})
}

func (s *Stream[V]) Each(fun func(VAL) bool) *Stream[VAL] {
	return NewConsumer(s, func(elem V) (optional.Optional[V], bool) {
		return optional.New[V](elem, true), fun(elem)
	})
}

func (s *Stream[V]) Collect() []V {
	var coll []V
	for v := range s.values {
		coll = append(coll, v)
	}
	return coll
}
