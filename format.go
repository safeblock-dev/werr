package werr

import (
	"path"
	"strconv"
	"strings"
)

type FormatFn func(file string, line int, funcName string, err error, msg string) string

var _defaultFormatter FormatFn = defaultFormatter

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

func SetFormatter(fn FormatFn) {
	_defaultFormatter = fn
}
