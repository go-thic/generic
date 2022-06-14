package stream

import (
	"github.com/go-thic/generic/optional"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("WithValues", func(t *testing.T) {
		generator := WithValues()
		assert.NotPanics(t, func() {
			_ = generator()
		})
		assert.False(t, generator().IsSome())
		assert.Panics(t, func() {
			generator().Val()
		})
	})

	t.Run("WithValues generating int Values", func(t *testing.T) {
		generator := WithValues(1, 2, 3)
		a := generator()
		assert.True(t, a.IsSome())
		assert.Equal(t, 1, a.Val())
		assert.Equal(t, 2, generator().Val())
		assert.Equal(t, 3, generator().Val())
		assert.Equal(t, optional.None[VAL](), generator())
		assert.Panics(t, func() {
			generator().Val()
		})
	})

	t.Run("WithValues generating string values", func(t *testing.T) {
		generator := WithValues("a", "b", "c")
		a := generator()

		assert.True(t, a.IsSome())
		assert.Equal(t, "a", a.Val())
		assert.Equal(t, "b", generator().Val())
		assert.Equal(t, "c", generator().Val())
		assert.Equal(t, optional.None[VAL](), generator())
		assert.Panics(t, func() {
			generator().Val()
		})
	})
}

func TestNewProvider(t *testing.T) {
	t.Run("Test if values arrive", func(t *testing.T) {
		s := NewProvider(WithValues(1, 2, 3))

		assert.Equal(t, 1, <-s.values)
		assert.Equal(t, 2, <-s.values)
		assert.Equal(t, 3, <-s.values)
		assert.NotPanics(t, func() {
			_ = <-s.values
		})
	})
}

func TestStream_Limit(t *testing.T) {
	s := NewProvider(WithValues("a", "b", "c")).Limit(Count(3))

	assert.NotNil(t, s)
	assert.Equal(t, "a", <-s.values)
	assert.Equal(t, "b", <-s.values)
	assert.Equal(t, "c", <-s.values)
	assert.NotPanics(t, func() {
		<-s.values
	})
}

func TestNewMapper(t *testing.T) {
	x := Transpose(ToString[int])(1)
	assert.Equal(t, "1", x.Val())

	m := NewMapper(NewProvider(WithValues(1, 2, 3, 4)), Transpose(ToString[int])).Limit(Count(4))
	assert.NotNil(t, m)
	assert.Equal(t, []DST{"1", "2", "3", "4"}, m)
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
	t.Run("Do ints to strings", func(t *testing.T) {
		var strSlice []string
		NewProvider(WithValues(1, 2, 3, 4, 5)).Do(Map(FizzBuzz)).Finally(Do(func(s string) {
			strSlice = append(strSlice, s)
		}))
		assert.NotNil(t, strSlice)
		assert.Equal(t, []string{"1", "2", "Fizz", "4", "Buzz"}, strSlice)
	})
}

func TestImpl_Do(t *testing.T) {
	t.Run("Filter and finally sum up", func(t *testing.T) {
		x := 0
		unevenNumbers := func(n int) bool { return n%2 != 0 }
		sumUp := func(elem int) {
			x += elem
		}
		NewProvider(WithValues(1, 2, 3, 4, 5)).Do(Filter(unevenNumbers)).Finally(Do(sumUp))

		assert.Equal(t, 6, x)
	})

	t.Run("Filter and collect in a slice of strings", func(t *testing.T) {
		var result []string

		NewProvider(WithValues("Hello", "sad, sad", "World")).Do(Filter(Equals("sad, sad"))).Finally(Do(func(s string) {
			result = append(result, s)
		}))

		assert.Equal(t, []string{"Hello", "World"}, result)
	})
}
