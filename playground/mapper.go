package playground

import (
	"fmt"
	"log"

	"github.com/go-thic/generic/optional"
)

func NewMapper[S, D VAL](s *Stream[S], doMap func(elem S) (optional.Optional[D], bool)) *Stream[D] {
	valueChan := make(chan D)
	stopChan := make(chan empty)

	mapperStream := New[D](valueChan, stopChan)

	go func() {
		defer func() {
			close(valueChan)
			close(stopChan)
			if r := recover(); r != nil {
				log.Printf("panic: %q", r)
			}
		}()

		for v := range s.values {
			mapped, stopMapping := doMap(v)
			if mapped.IsSome() {
				mapperStream.Write(mapped.Val())
			}
			if stopMapping {
				break
			}
		}
	}()

	return mapperStream
}

func Map[S, D VAL](mapFunc func(elem S) (D, bool)) func(elem S) (optional.Optional[D], bool) {
	return func(elem S) (optional.Optional[D], bool) {
		mappedVal, stopMapping := mapFunc(elem)
		return optional.New(mappedVal, !stopMapping), stopMapping
	}
}

func ToString[T VAL](s T) (string, bool) {
	if v, ok := VAL(s).(fmt.Stringer); ok {
		return v.String(), false
	}
	return fmt.Sprint(s), false
}
