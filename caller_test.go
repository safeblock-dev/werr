package werr

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func c() (string, string, int) { return caller(0) }
func b() (string, string, int) { return c() }
func a() (string, string, int) { return b() }

var varFuncName, varFile, varLine = caller(0) //nolint: gochecknoglobals

func TestCaller(t *testing.T) {
	t.Parallel()

	t.Run("when called in a function", func(t *testing.T) {
		t.Parallel()

		funcName, file, line := a()
		require.Equal(t, "caller_test.go:11", filepath.Base(file)+":"+strconv.Itoa(line))
		require.Equal(t, "github.com/safeblock-dev/werr.c", funcName)
	})

	t.Run("when called from outside", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, "caller_test.go:15", filepath.Base(varFile)+":"+strconv.Itoa(varLine))
		require.Equal(t, "github.com/safeblock-dev/werr.init", varFuncName)
	})
}
