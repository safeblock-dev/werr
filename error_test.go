package werr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCause(t *testing.T) {
	t.Parallel()

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3 := newError(err2, "wrap level 1")
		err4 := newError(err3, "wrap level 2")

		require.Equal(t, err2, Cause(err4))
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, Cause(nil))
	})

	t.Run("when without wrap", func(t *testing.T) {
		t.Parallel()

		errWithoutWrap := errors.New("error without wrap")
		require.Equal(t, errWithoutWrap, Cause(errWithoutWrap))
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

		require.Equal(t, err3, Unwrap(err4))
		require.Equal(t, err2, Unwrap(err3))
		require.Equal(t, err2, Unwrap(err2))
		require.Equal(t, err1, Unwrap(err1))
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, Unwrap(nil))
	})
}

func TestUnwrapAll(t *testing.T) {
	t.Parallel()

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3 := newError(err2, "wrap level 1")
		err4 := newError(err3, "wrap level 2")

		require.Equal(t, err1, UnwrapAll(err4))
		require.Equal(t, err1, UnwrapAll(err3))
		require.Equal(t, err1, UnwrapAll(err2))
		require.Equal(t, err1, UnwrapAll(err1))
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, UnwrapAll(nil))
	})
}

func TestAsWrap(t *testing.T) {
	t.Parallel()

	// Test case when the error is not a wrappedError
	t.Run("not wrappedError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("not a wrapped error")
		_, ok := AsWrap(err)
		require.False(t, ok)
	})

	// Test case when the error is a wrappedError
	t.Run("is wrappedError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("original error")
		wrappedErr := newError(err, "additional message")
		wErr, ok := AsWrap(wrappedErr)
		require.True(t, ok)
		require.Equal(t, wrappedErr, wErr)
	})

	// Test case when the error is a nested wrappedError
	t.Run("nested wrappedError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("original error")
		wrappedErr1 := newError(err, "wrap level 1")
		wrappedErr2 := newError(wrappedErr1, "wrap level 2")
		wErr, ok := AsWrap(wrappedErr2)
		require.True(t, ok)
		require.Equal(t, wrappedErr2, wErr)
	})

	// Test case when the error is nil
	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		_, ok := AsWrap(nil)
		require.False(t, ok)
	})
}

func TestIsWrap(t *testing.T) {
	t.Parallel()

	t.Run("WrappedError", func(t *testing.T) {
		t.Parallel()

		require.True(t, IsWrap(Wrap(errors.New("wrapped error"))))
	})

	t.Run("CustomError", func(t *testing.T) {
		t.Parallel()

		require.False(t, IsWrap(errors.New("custom error")))
	})

	t.Run("NilError", func(t *testing.T) {
		t.Parallel()

		require.False(t, IsWrap(nil))
	})
}

//
// WrappedError
//

func TestWrappedError_Error(t *testing.T) {
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

func TestWrappedError_Cause(t *testing.T) {
	t.Parallel()

	// Mock errors for testing
	err1 := errors.New("original error")
	err2 := fmt.Errorf("fmt wrap: %w", err1)
	err3 := newError(err2, "wrap level 1")
	err4 := newError(err3, "wrap level 2")

	t.Run("single wrapped error", func(t *testing.T) {
		t.Parallel()

		// Create a wrapped error instance
		wrappedErr := wrappedError{
			file:     "main.go",
			funcName: "main.main",
			line:     42,
			err:      err1,
			msg:      "additional message",
		}

		// Calling Cause() should return the original error
		cause := wrappedErr.Cause()
		require.Equal(t, err1, cause)
	})

	t.Run("chain of wrapped errors", func(t *testing.T) {
		t.Parallel()

		// Create a nested chain of wrapped errors
		wrappedErr4 := wrappedError{
			file:     "main.go",
			funcName: "main.main",
			line:     42,
			err:      err4,
			msg:      "level 2 message",
		}

		// Calling Cause() on the top-level wrapped error should return the original error in the chain
		cause := wrappedErr4.Cause()
		require.Equal(t, err2, cause)
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		// Create a wrapped error instance with nil error
		wrappedErr := wrappedError{
			file:     "main.go",
			funcName: "main.main",
			line:     42,
			msg:      "nil error wrap",
		}

		// Calling Cause() on an instance with nil error should return nil
		cause := wrappedErr.Cause()
		require.NoError(t, cause)
	})
}

//
// Tests for errors package
//

func TestPkgErrorsUnwrap(t *testing.T) {
	t.Parallel()

	err1 := errors.New("original error")
	err2 := newError(err1, "wrap level 1")
	err3 := newError(err2, "wrap level 2")

	t.Run("unwrap level 1", func(t *testing.T) {
		t.Parallel()
		require.EqualError(t, errors.Unwrap(err3), err2.Error())
	})

	t.Run("unwrap level 2", func(t *testing.T) {
		t.Parallel()
		require.EqualError(t, errors.Unwrap(err2), err1.Error())
	})

	t.Run("unwrap original", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, errors.Unwrap(err1))
	})
}

func TestPkgErrorsIs(t *testing.T) {
	t.Parallel()

	err1 := errors.New("original error")
	err2 := fmt.Errorf("fmt wrap: %w", err1)
	err3 := newError(err2, "wrap level 1")
	err4 := newError(err3, "wrap level 2")

	t.Run("is original error", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t, err4, err1)
	})

	t.Run("is fmt wrapped error", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t, err4, err2)
	})

	t.Run("is wrapped error 1", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t, err4, err3)
	})

	t.Run("is different error", func(t *testing.T) {
		t.Parallel()
		otherErr := errors.New("different error")
		require.NotErrorIs(t, err4, otherErr)
	})
}

func TestPkgErrorsAs(t *testing.T) {
	t.Parallel()

	err1 := errors.New("original error")
	err2 := newError(err1, "wrap level 1")
	err3 := newError(err2, "wrap level 2")

	t.Run("as wrappedError level 2", func(t *testing.T) {
		t.Parallel()
		var wErr wrappedError
		require.ErrorAs(t, err3, &wErr)
		require.Equal(t, wErr.Error(), err3.Error())
	})

	t.Run("as wrappedError level 1", func(t *testing.T) {
		t.Parallel()
		var wErr wrappedError
		require.ErrorAs(t, err2, &wErr)
		require.Equal(t, wErr.Error(), err2.Error())
	})

	t.Run("as different error type", func(t *testing.T) {
		t.Parallel()
		var otherErr *wrappedError
		otherErrorInstance := errors.New("different error")
		require.False(t, errors.As(otherErrorInstance, &otherErr))
	})
}
