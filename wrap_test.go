package werr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/werr"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("original error")
		wrappedErr := werr.Wrap(originalErr)

		// Ensure that the wrapped error is of type error
		require.IsType(t, werr.Error{}, wrappedErr)

		// Ensure that the wrapped error contains the original error
		require.ErrorIs(t, wrappedErr, originalErr)
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()
		// Ensure that wrapping a nil error results in nil
		require.NoError(t, werr.Wrap(nil))
	})

	t.Run("check skip count is sufficient", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("original error")
		wrappedErr := werr.Wrap(originalErr)

		// Ensure the skip count 3 is enough
		const exp = "github.com/safeblock-dev/werr_test.TestWrap/wrap_test.go:39\tfunc3()\noriginal error"

		require.Equal(t, exp, wrappedErr.Error())
	})
}

func TestWrapf(t *testing.T) {
	t.Parallel()

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("original error")
		wrappedErr := werr.Wrapf(originalErr, "additional message: %s", "some details")

		// Ensure that the wrapped error is of type *wrapError
		require.IsType(t, werr.Error{}, wrappedErr)

		// Ensure that the wrapped error contains the original error
		require.ErrorIs(t, wrappedErr, originalErr)

		// Ensure that the wrapped error contains the additional message
		require.Contains(t, wrappedErr.Error(), "additional message: some details")
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		wrappedErr := werr.Wrapf(nil, "additional message: %s", "some details")

		// Ensure that wrapping a nil error results in nil
		require.NoError(t, wrappedErr, "Expected wrapped error to be nil")
	})
}

func TestWrapfWithArgs(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")

	t.Run("when arg is equal primitive", func(t *testing.T) {
		t.Parallel()

		args := werr.Args("test")
		act := werr.Wrapf(originalErr, werr.ArgsFormat, args).(werr.Error).Message()
		require.Equal(t, "args=[test]", act)
	})

	t.Run("when arg is equal map", func(t *testing.T) {
		t.Parallel()

		args := werr.Args(map[string]any{
			"foo":  "bar",
			"bar":  1,
			"buzz": false,
		})
		act := werr.Wrapf(originalErr, werr.ArgsFormat, args).(werr.Error).Message()
		require.Equal(t, "args=[map[bar:1 buzz:false foo:bar]]", act)
	})

	t.Run("when arg is equal slice", func(t *testing.T) {
		t.Parallel()

		args := werr.Args("foo", 1, false)
		act := werr.Wrapf(originalErr, werr.ArgsFormat, args).(werr.Error).Message()
		require.Equal(t, "args=[foo 1 false]", act)
	})
}

func TestWrapt(t *testing.T) {
	t.Parallel()

	fn := func(err error) (int, error) {
		return 100, err
	}

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("original error")
		_, wrappedErr := werr.Wrapt(fn(originalErr))

		// Ensure that the wrapped error is of type error
		require.IsType(t, werr.Error{}, wrappedErr)

		// Ensure that the wrapped error contains the original error
		require.ErrorIs(t, wrappedErr, originalErr)
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()
		// Ensure that wrapping a nil error results in nil
		v, err := werr.Wrapt(fn(nil))
		require.NoError(t, err)
		require.NotEmpty(t, v)
	})

	t.Run("check skip count is sufficient", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("original error")
		_, wrappedErr := werr.Wrapt(fn(originalErr))

		// Ensure the skip count 3 is enough
		const exp = "github.com/safeblock-dev/werr_test.TestWrapt/wrap_test.go:143\tfunc4()\noriginal error"

		require.Equal(t, exp, wrappedErr.Error())
	})
}

func TestCause(t *testing.T) {
	t.Parallel()

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3 := werr.Wrapf(err2, "wrap level 1")
		err4 := werr.Wrapf(err3, "wrap level 2")

		require.Equal(t, err2, werr.Cause(err4))
	})

	t.Run("when join", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3p0 := werr.Wrapf(err2, "wrap level 1")
		err3p1 := werr.Wrapf(err2, "wrap level 2")
		err3 := errors.Join(err3p0, err3p1)

		require.Equal(t, err3, werr.Cause(err3))
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, werr.Cause(nil))
	})

	t.Run("when without wrap", func(t *testing.T) {
		t.Parallel()

		errWithoutWrap := errors.New("error without wrap")
		require.Equal(t, errWithoutWrap, werr.Cause(errWithoutWrap))
	})
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3 := werr.Wrapf(err2, "wrap level 1")
		err4 := werr.Wrapf(err3, "wrap level 2")

		require.Equal(t, err3, werr.Unwrap(err4))
		require.Equal(t, err2, werr.Unwrap(err3))

		require.NoError(t, werr.Unwrap(err2))
		require.Equal(t, err1, errors.Unwrap(err2))

		require.NoError(t, werr.Unwrap(err1))
		require.NoError(t, errors.Unwrap(err1))
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, werr.Unwrap(nil))
	})
}

func TestIsWrap(t *testing.T) {
	t.Parallel()

	t.Run("wrapped error", func(t *testing.T) {
		t.Parallel()

		require.True(t, werr.IsWrap(werr.Wrap(errors.New("wrapped error"))))
	})

	t.Run("custom error", func(t *testing.T) {
		t.Parallel()

		require.False(t, werr.IsWrap(errors.New("custom error")))
	})

	t.Run("join error", func(t *testing.T) {
		t.Parallel()

		err0 := errors.New("custom error 0")
		err1 := errors.New("custom error 1")
		err := errors.Join(err0, err1)

		require.False(t, werr.IsWrap(err))
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		require.False(t, werr.IsWrap(nil))
	})
}
