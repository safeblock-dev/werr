package werr

// Error represents an error with additional context such as funcName, file, line, and msg.
type Error struct {
	funcName string // funcName represents the fully qualified function name ("<pkg>.<name>").
	file     string // file is the file name where the error occurred.
	line     int    // line is the line number in the file where the error occurred.
	err      error  // err is the original error that was wrapped.
	msg      string // msg is an optional message to provide additional context for the error.
}

// newError creates a new wrapped error with caller information and an optional additional message.
func newError(err error, msg string) error {
	funcName, file, line := caller(defaultCallerSkip)

	return Error{
		file:     file,
		funcName: funcName,
		line:     line,
		err:      err,
		msg:      msg,
	}
}

// Error returns a string representation of the wrapped error.
func (e Error) Error() string {
	return _defaultFormatter(e.file, e.line, e.funcName, e.err, e.msg)
}

// Format returns a custom formatted string representation of the wrapped error using a provided formatter function.
func (e Error) Format(fn FormatFn) string {
	return fn(e.file, e.line, e.funcName, e.err, e.msg)
}

// Unwrap returns the underlying error wrapped by this structure.
func (e Error) Unwrap() error {
	return e.err
}

// Cause returns the root cause of the wrapped error by recursively unwrapping it if it is also an Error.
func (e Error) Cause() error {
	for {
		wrapped, ok := e.err.(Error) //nolint: errorlint
		if !ok {
			break
		}

		e.err = wrapped.Unwrap()
	}

	return e.err
}

// File returns the file path associated with the wrapped error.
func (e Error) File() string {
	return e.file
}

// Line returns the line number associated with the wrapped error.
func (e Error) Line() int {
	return e.line
}

// FuncName returns the fully qualified function name associated with the wrapped error.
// Example: "main.main".
func (e Error) FuncName() string {
	return e.funcName
}

// Message returns the additional message associated with the wrapped error.
func (e Error) Message() string {
	return e.msg
}
