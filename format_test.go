package werr

import (
	"errors"
	"strconv"
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

// nolint: paralleltest
func TestSetFormatter(t *testing.T) {
	testCases := []struct {
		name      string
		formatter FormatFn
		file      string
		line      int
		funcName  string
		err       error
		msg       string
		expected  string
	}{
		{
			name: "CustomFormatter",
			formatter: func(_ string, line int, funcName string, err error, msg string) string {
				return funcName + "#" + strconv.Itoa(line) + " - " + msg + ": " + err.Error()
			},
			file:     "file.go",
			line:     42,
			funcName: "main.main",
			err:      errors.New("an error occurred"),
			msg:      "something went wrong",
			expected: "main.main#42 - something went wrong: an error occurred",
		},
		{
			name:      "DefaultFormatter",
			formatter: defaultFormatter,
			file:      "file.go",
			line:      42,
			funcName:  "main.main",
			err:       errors.New("an error occurred"),
			msg:       "something went wrong",
			expected:  "main/file.go:42\tmain()\tsomething went wrong\nan error occurred",
		},
	}

	for _, testCase := range testCases {
		tt := testCase
		t.Run(tt.name, func(t *testing.T) {
			SetFormatter(tt.formatter)
			result := _defaultFormatter(tt.file, tt.line, tt.funcName, tt.err, tt.msg)
			require.Equal(t, tt.expected, result)
		})
	}
}
