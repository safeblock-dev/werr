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
		wrappedErr := wrappedError{
			file:     "main.go",
			funcName: "main.main",
			line:     42,
			err:      err,
			msg:      "additional message",
		}

		exp := "main/main.go:42\tmain()\tadditional message\noriginal error"
		require.Equal(t, exp, wrappedErr.Error())
	})

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()
		subWrappedErr := wrappedError{
			file:     "main.go",
			funcName: "main.main2",
			line:     84,
			err:      err,
			msg:      "",
		}
		wrappedErr := wrappedError{
			file:     "main.go",
			funcName: "main.main",
			line:     42,
			err:      subWrappedErr,
			msg:      "additional message",
		}

		exp := "main/main.go:42\tmain()\tadditional message\nmain/main.go:84\tmain2()\noriginal error"
		require.Equal(t, exp, wrappedErr.Error())
	})

	t.Run("without message", func(t *testing.T) {
		t.Parallel()
		wrappedErr := wrappedError{
			file:     "main.go",
			funcName: "main.main",
			line:     42,
			err:      err,
		}

		exp := "main/main.go:42\tmain()\noriginal error"
		require.Equal(t, exp, wrappedErr.Error())
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

		require.EqualError(t, Cause(err4), err2.Error())
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
