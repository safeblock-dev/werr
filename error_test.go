package werr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//
// Error
//

func TestError_Error(t *testing.T) {
	t.Parallel()

	err := errors.New("original error")

	t.Run("with message", func(t *testing.T) {
		t.Parallel()
		wrappedErr := Error{
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
		subWrappedErr := Error{
			file:     "main.go",
			funcName: "main.main2",
			line:     84,
			err:      err,
			msg:      "",
		}
		wrappedErr := Error{
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
		wrappedErr := Error{
			file:     "main.go",
			funcName: "main.main",
			line:     42,
			err:      err,
		}

		exp := "main/main.go:42\tmain()\noriginal error"
		require.Equal(t, exp, wrappedErr.Error())
	})
}

func TestError_Cause(t *testing.T) {
	t.Parallel()

	// Mock errors for testing
	err1 := errors.New("original error")
	err2 := fmt.Errorf("fmt wrap: %w", err1)
	err3 := newError(err2, "wrap level 1")
	err4 := newError(err3, "wrap level 2")

	t.Run("single wrapped error", func(t *testing.T) {
		t.Parallel()

		// Create a wrapped error instance
		wrappedErr := Error{
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
		wrappedErr4 := Error{
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
		wrappedErr := Error{
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

	t.Run("as Error level 2", func(t *testing.T) {
		t.Parallel()
		var wErr Error
		require.ErrorAs(t, err3, &wErr)
		require.Equal(t, wErr.Error(), err3.Error())
	})

	t.Run("as Error level 1", func(t *testing.T) {
		t.Parallel()
		var wErr Error
		require.ErrorAs(t, err2, &wErr)
		require.Equal(t, wErr.Error(), err2.Error())
	})

	t.Run("as different error type", func(t *testing.T) {
		t.Parallel()
		var otherErr *Error
		otherErrorInstance := errors.New("different error")
		require.False(t, errors.As(otherErrorInstance, &otherErr))
	})
}
