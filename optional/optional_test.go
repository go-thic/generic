package optional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNone(t *testing.T) {
	t.Run("None is not Some", func(t *testing.T) {
		n := None[string]()
		if got := n.IsSome(); got != false {
			t.Errorf("IsSome() = %v, want %v", got, false)
		}
	})

	t.Run("None is None", func(t *testing.T) {
		n := None[int]()
		if got := n.IsNone(); got != true {
			t.Errorf("IsNone() = %v, want %v", got, true)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("New of None returns None that is some", func(t *testing.T) {
		none := None[int]()
		got := New(none, true)

		assert.True(t, got.IsSome())
		assert.False(t, got.Val().IsSome())
	})
}

func TestNone_Val(t *testing.T) {
	t.Run("Val of None", func(t *testing.T) {
		n := None[string]()

		assert.Panics(t, func() { _ = n.Val() })
		assert.Equal(t, "hello", n.Val(OrElse("hello")))
		assert.Equal(t, "", n.Val(OrZero[string]()))
	})
}
