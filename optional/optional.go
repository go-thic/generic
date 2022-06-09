package optional

type Optional[T any] interface {
	Option
	Val(orElse ...func() T) T
}

type Option interface {
	IsSome() bool
	IsNone() bool
}

func New[V any](v V, isSome bool) Optional[V] {
	return Value[V]{isSome: isSome, val: v}
}

func None[V any]() Optional[V] {
	return Value[V]{
		val:    getZero[V](),
		isSome: false,
	}
}

type Value[T any] struct {
	isSome bool
	val    T
}

func (v Value[T]) IsNone() bool {
	return !v.isSome
}

func (v Value[T]) IsSome() bool {
	return v.isSome
}

func (v Value[T]) Val(orElse ...func() T) T {
	if v.IsNone() {
		if orElse != nil {
			return orElse[0]()
		}
		panic("trying to call Val() on None")
	}

	return v.val
}

func getZero[T any]() T {
	var result T
	return result
}

func OrElse[T any](v T) func() T {
	return func() T {
		return v
	}
}

func OrZero[T any]() func() T {
	return func() T {
		var result T
		return result
	}
}
