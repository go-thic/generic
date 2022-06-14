package playground

import (
	"github.com/go-thic/generic/optional"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("WithValues", func(t *testing.T) {
		generator := WithValues[int]()
		assert.NotPanics(t, func() {
			_ = generator()
		})
		assert.False(t, generator().IsSome())
		assert.Panics(t, func() {
			generator().Val()
		})
	})

	t.Run("WithValues generating Values", func(t *testing.T) {
		generator := WithValues(1, 2, 3)
		a := generator()
		assert.True(t, a.IsSome())
		assert.Equal(t, 1, a.Val())
		assert.Equal(t, 2, generator().Val())
		assert.Equal(t, 3, generator().Val())
		assert.Equal(t, optional.None[int](), generator())
		assert.Panics(t, func() {
			generator().Val()
		})
	})
}

func TestNewProvider(t *testing.T) {
	s := NewProvider(WithValues(1, 2, 3))

	assert.Equal(t, 1, <-s.values)
	assert.Equal(t, 2, <-s.values)
	assert.Equal(t, 3, <-s.values)
	assert.NotPanics(t, func() {
		_ = <-s.values
	})

	a := NewProvider(WithValues("a", "b", "c")).Collect()
	assert.Equal(t, []string{"a", "b", "c"}, a)
}

func TestStream_Filter(t *testing.T) {
	s := NewProvider(WithValues(1, 2, 3)).Filter(func(elem VAL) bool {
		return elem == 2
	})

	assert.NotNil(t, s)
	assert.Equal(t, 1, <-s.values)
	assert.Equal(t, 3, <-s.values)
	assert.NotPanics(t, func() {
		<-s.values
	})

	x := NewProvider(WithValues("a", "b", "c")).Filter(Equals("b")).Collect()
	assert.Equal(t, []VAL{"a", "c"}, x)
}

func TestStream_Limit(t *testing.T) {
	s := NewProvider(WithValues("a", "b", "c")).Limit(Count[VAL](3))

	assert.NotNil(t, s)
	assert.Equal(t, "a", <-s.values)
	assert.Equal(t, "b", <-s.values)
	assert.Equal(t, "c", <-s.values)
	assert.NotPanics(t, func() {
		<-s.values
	})
}

func TestNewMapper(t *testing.T) {
	x, b := Map(ToString[int])(1)
	assert.Equal(t, "1", x.Val())
	assert.False(t, b)

	m := NewMapper(NewProvider(WithValues(1, 2, 3, 4)), Map(ToString[int])).Limit(Count[VAL](4)).Collect()
	assert.NotNil(t, m)
	assert.Equal(t, []VAL{"1", "2", "3", "4"}, m)
}

// FizzBuzz just for fun!
func FizzBuzz(number int) string {
	theString := func(s string) string { return s }

	makeMapper := func(forNumber int, toString string) func(func(string) string) func(string) string {
		return func(fun func(string) string) func(string) string {
			if number%forNumber == 0 {
				return func(_ string) string {
					return toString + fun("")
				}
			}
			return fun
		}
	}

	fizz := makeMapper(3, "Fizz")
	buzz := makeMapper(5, "Buzz")

	return fizz(buzz(theString))(strconv.Itoa(number))
}

func TestStream_Map(t *testing.T) {
	str := NewProvider(WithValues(1, 2, 3, 4, 5)).Map(func(in int) (VAL, bool) {
		return FizzBuzz(in), false
	}).Collect()

	assert.NotNil(t, str)
	assert.Equal(t, []VAL{"1", "2", "Fizz", "4", "Buzz"}, str)
}
