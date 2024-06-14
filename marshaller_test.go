package werr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultFormatter(t *testing.T) {
	t.Parallel()

	file, line, funcName := "main.go", 42, "main.TestFunction"
	err, msg := errors.New("example error"), "custom message"

	format := defaultFormatter(file, line, funcName, err, msg)

	exp := "main/main.go:42\tTestFunction()\tcustom message\nexample error"
	require.Equal(t, exp, format)
}
