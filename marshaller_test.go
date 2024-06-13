package werr_test

import (
	"errors"
	"testing"

	"github.com/safeblock-dev/werr"
	"github.com/stretchr/testify/require"
)

func TestDefaultErrorStackMarshaler(t *testing.T) {
	t.Parallel()

	file, line, funcName := "main.go", 42, "main.TestFunction"
	err, msg := errors.New("example error"), "custom message"

	format := werr.DefaultErrorStackMarshaler(file, line, funcName, err, msg)

	exp := "main/main.go:42\tTestFunction()\tcustom message\nexample error"
	require.Equal(t, exp, format)
}
