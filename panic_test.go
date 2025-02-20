package werr_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/werr"
)

func TestPanicToError(t *testing.T) {
	t.Parallel()

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		err := werr.PanicToError(nil)
		require.NoError(t, err)
	})

	t.Run("when error", func(t *testing.T) {
		t.Parallel()

		inputErr := errors.New("some error")
		err := werr.PanicToError(inputErr)
		require.Error(t, err)
		require.ErrorIs(t, err, inputErr)
	})

	t.Run("when string", func(t *testing.T) {
		t.Parallel()

		err := werr.PanicToError("some string panic")
		require.Error(t, err)
		require.Contains(t, err.Error(), "panic recovered")
		require.Contains(t, err.Error(), "some string panic")
	})

	t.Run("other type", func(t *testing.T) {
		t.Parallel()

		err := werr.PanicToError(123)
		require.Error(t, err)
		require.Contains(t, err.Error(), "panic recovered")
		require.Contains(t, err.Error(), "123")
	})

	t.Run("panic example", func(t *testing.T) {
		t.Parallel()

		var v any

		func() {
			defer func() { v = recover() }()

			_ = []string{}[1] // panic
		}()

		err := werr.PanicToError(v)
		require.Error(t, err)

		require.Contains(t, err.Error(), "runtime/debug.Stack()")
		require.Contains(t, err.Error(), "runtime error: index out of range [1] with length 0")
	})
}
