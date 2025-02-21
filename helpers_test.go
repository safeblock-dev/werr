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
