package werr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapError_Error(t *testing.T) {
	t.Parallel()

	err := errors.New("original error")

	t.Run("with message", func(t *testing.T) {
		t.Parallel()
		wrappedErr := wrapError{
			caller:   "caller",
			err:      err,
			funcName: "function",
			msg:      "additional message",
		}

		const expected = "caller\tfunction\tadditional message\noriginal error"
		require.Equal(t, expected, wrappedErr.Error())
	})

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()
		subWrappedErr := wrapError{
			caller:   "subCaller",
			err:      err,
			funcName: "subFunction",
			msg:      "",
		}
		wrappedErr := wrapError{
			caller:   "caller",
			err:      subWrappedErr,
			funcName: "function",
			msg:      "additional message",
		}

		const expected = "caller\tfunction\tadditional message\nsubCaller\tsubFunction\noriginal error"
		require.Equal(t, expected, wrappedErr.Error())
	})

	t.Run("without message", func(t *testing.T) {
		t.Parallel()
		wrappedErr := wrapError{
			caller:   "caller",
			err:      err,
			funcName: "function",
		}

		const expected = "caller\tfunction\noriginal error"
		require.Equal(t, expected, wrappedErr.Error())
	})
}

func TestCause(t *testing.T) {
	t.Parallel()

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3 := newError(err2, "wrap level 1")
		err4 := newError(err3, "wrap level 2")

		require.EqualError(t, Cause(err4), err1.Error())
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, Cause(nil))
	})

	t.Run("when without wrap", func(t *testing.T) {
		t.Parallel()

		errWithoutWrap := errors.New("error without wrap")
		require.ErrorIs(t, Cause(errWithoutWrap), errWithoutWrap)
	})
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3 := newError(err2, "wrap level 1")
		err4 := newError(err3, "wrap level 2")

		require.EqualError(t, Unwrap(err4), err3.Error())
		require.EqualError(t, Unwrap(err3), err2.Error())
		require.EqualError(t, Unwrap(err2), err1.Error())
		require.NoError(t, Unwrap(err1))
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, Unwrap(nil))
	})

	t.Run("when without wrap", func(t *testing.T) {
		t.Parallel()

		errWithoutWrap := errors.New("error without wrap")
		require.NoError(t, Unwrap(errWithoutWrap))
	})
}
