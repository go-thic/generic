package stream

import (
	"fmt"
	"log"

	"github.com/go-thic/generic/optional"
)

type SRC interface {
	any
}

type DST interface {
	any
}

func NewMapper[S, D any](s *Stream, doMap func(elem S) optional.Optional[D]) *Stream {
	valueChan := make(chan any)

	mapperStream := newImpl(valueChan)

	go func() {
		defer func() {
			close(valueChan)
			if r := recover(); r != nil {
				log.Printf("panic: %q", r)
			}
		}()

		for v := range s.values {
			if val, ok := v.(S); ok {
				mapped := doMap(val)
				if mapped.IsSome() {
					mapperStream.Write(mapped.Val())
				}
			}
		}
	}()

	return mapperStream
}

func transpose[S, D any](mapFunc func(elem S) (D, bool)) func(elem S) optional.Optional[D] {
	return func(elem S) optional.Optional[D] {
		mappedVal, isSome := mapFunc(elem)
		return optional.New(mappedVal, !isSome)
	}
}

func Map[S, D any](mapper func(S) D) func(elem SRC) (DST, bool) {
	return func(elem SRC) (DST, bool) {
		if s, ok := elem.(S); ok {
			return mapper(s), false
		}
		var zero DST
		return zero, true
	}
}

func Filter[S any](filter func(elem S) bool) func(elem SRC) (DST, bool) {
	return func(elem SRC) (DST, bool) {
		if s, ok := elem.(S); ok {
			return s, filter(s)
		}
		var zero DST
		return zero, false
	}
}

func ToString[T any](s T) (string, bool) {
	if v, ok := any(s).(fmt.Stringer); ok {
		return v.String(), false
	}
	return fmt.Sprint(s), false
}
