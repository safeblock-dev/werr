package werr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/werr"
)

func TestUnwrapAll(t *testing.T) {
	t.Parallel()

	t.Run("when wrap chain", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("original error")
		err2 := fmt.Errorf("fmt wrap: %w", err1)
		err3 := werr.Wrapf(err2, "wrap level 1")
		err4 := werr.Wrapf(err3, "wrap level 2")

		require.Equal(t, err1, werr.UnwrapAll(err4))
		require.Equal(t, err1, werr.UnwrapAll(err3))
		require.Equal(t, err1, werr.UnwrapAll(err2))
		require.Equal(t, err1, werr.UnwrapAll(err1))
	})

	t.Run("when nil", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, werr.UnwrapAll(nil))
	})
}

func TestArg(t *testing.T) {
	t.Parallel()

	t.Run("when arg is equal primitive", func(t *testing.T) {
		t.Parallel()

		arg := "test"
		act := werr.Arg(arg)
		require.Equal(t, fmt.Sprintf("arg=%+v", arg), act)
	})

	t.Run("when arg is equal map", func(t *testing.T) {
		t.Parallel()

		arg := map[string]interface{}{
			"foo":  "bar",
			"bar":  1,
			"buzz": false,
		}
		act := werr.Arg(arg)
		require.Equal(t, fmt.Sprintf("arg=%+v", arg), act)
	})

	t.Run("when arg is equal slice", func(t *testing.T) {
		t.Parallel()

		arg := []interface{}{"foo", 1, false}
		act := werr.Arg(arg)
		require.Equal(t, fmt.Sprintf("arg=%+v", arg), act)
	})
}

func TestArgs(t *testing.T) {
	t.Parallel()

	t.Run("when arg is equal primitive", func(t *testing.T) {
		t.Parallel()

		arg := "test"
		act := werr.Args(arg)
		require.Equal(t, fmt.Sprintf("args=[%+v]", arg), act)
	})

	t.Run("when arg is equal map", func(t *testing.T) {
		t.Parallel()

		arg := map[string]interface{}{
			"foo":  "bar",
			"bar":  1,
			"buzz": false,
		}
		act := werr.Args(arg)
		require.Equal(t, fmt.Sprintf("args=[%+v]", arg), act)
	})

	t.Run("when arg is equal slice", func(t *testing.T) {
		t.Parallel()

		arg := []interface{}{"foo", 1, false}
		act := werr.Args(arg)
		require.Equal(t, fmt.Sprintf("args=[%+v]", arg), act)
	})
}
