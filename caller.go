package werr

import (
	"runtime"
)

const defaultCallerSkip = 3

func caller(skip int) (funcName string, file string, line int) {
	// skip current func call.
	if skip < 1 {
		skip = 1
	}

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", "", 0
	}

	funcName = runtime.FuncForPC(pc).Name()

	return funcName, file, line
}
