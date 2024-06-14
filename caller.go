package werr

import (
	"runtime"
)

const defaultCallerSkip = 3

// caller returns the name of the calling function, its file path,
// and line number after skipping `skip` levels in the call stack.
//
// Example output:
//
//	funcName: "main.main"
//	file: "/path/to/your/file/main.go"
//	line: 42
func caller(skip int) (string, string, int) {
	// skip current func call.
	if skip < 1 {
		skip = 1
	}

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", "", 0
	}

	funcName := runtime.FuncForPC(pc).Name()

	return funcName, file, line
}
