package stream

import "log"

type VAL interface {
	any
}

type Stream struct {
	values chan any
}

func New(generate ProviderFunc) *Stream {
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

func newImpl(values chan any) *Stream {
	return &Stream{
		values: values,
	}
}

func (s *Stream) Write(val any) {
	s.values <- val
}

func (s *Stream) Limit(limit func(elem VAL) bool) *Stream {
	return NewConsumer(s, withLimitFunc[VAL](limit))
}

func (s *Stream) Do(theNeedful func(elem SRC) (DST, bool)) *Stream {
	return NewMapper(s, transpose(theNeedful))
}

func (s *Stream) Finally(fun func(elem VAL)) {
	for v := range s.values {
		fun(v)
	}
}

func Do[S VAL](process func(S)) func(VAL) {
	return func(val VAL) {
		if v, ok := val.(S); ok {
			process(v)
		}
	}
}
