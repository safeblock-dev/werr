package werr

import (
	"path"
	"strconv"
	"strings"
)

//nolint:gochecknoglobals // for custom setting
var (
	// errorStackMarshaler extract the stack from err.
	errorStackMarshaler ErrorStackMarshalerFn = DefaultErrorStackMarshaler
)

type ErrorStackMarshalerFn func(file string, line int, funcName string, err error, msg string) string

func DefaultErrorStackMarshaler(file string, line int, funcName string, err error, msg string) string {
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
