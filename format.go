package werr

import (
	"path"
	"strconv"
	"strings"
)

// FormatFn defines a function signature for custom error formatting.
type FormatFn func(file string, line int, funcName string, err error, msg string) string

var _defaultFormatter FormatFn = defaultFormatter //nolint: gochecknoglobals

// defaultFormatter provides a default formatting style for error messages.
func defaultFormatter(file string, line int, funcName string, err error, msg string) string {
	idx := strings.LastIndex(funcName, ".")
	pkg := funcName[:idx]

	var fn string
	if len(funcName) > idx {
		fn = funcName[idx+1:] + "()"
	}

	source := pkg + "/" + path.Base(file) + ":" + strconv.Itoa(line)

	if msg != "" {
		msg = "\t" + msg
	}

	return source + "\t" + fn + msg + "\n" + err.Error()
}

// SetFormatter allows setting a custom error formatting function.
func SetFormatter(fn FormatFn) {
	_defaultFormatter = fn
}
