package werr

import (
	"errors"
)

type WrappedError interface {
	Error
	Wrapped
}

type Wrapped interface {
	Cause() error
	Unwrap() error
}

type Error interface {
	File() string
	FuncName() string
	Line() int
	Msg() string
}

var _ WrappedError = (*wrappedError)(nil)

// wrappedError represents an error with additional context such as funcName, file, line and msg.
type wrappedError struct {
	funcName string // funcName (function name) in the format "<pkg>.<name>", e.g. "main.main".
	file     string // file name of the location in this frame.
	line     int    // line number of the location in this frame.
	err      error  // err is original error.
	msg      string // msg is optional message to help analyze the error.
}

// newError creates a new wrapped error with caller information and additional message.
func newError(err error, msg string) error {
	funcName, file, line := caller(defaultCallerSkip)

	return wrappedError{
		file:     file,
		funcName: funcName,
		line:     line,

		err: err,
		msg: msg,
	}
}

// Unwrap calls errors.Unwrap to retrieve the underlying error.
func Unwrap(err error) error {
	var wErr wrappedError
	if errors.As(err, &wErr) {
		return wErr.Unwrap()
	}

	return err
}

// UnwrapAll function recursively unwraps an error that implements the Unwrap() error method,
// effectively retrieving the root cause error from a chain of wrapped errors.
func UnwrapAll(err error) error {
	for {
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}

		err = u.Unwrap()
	}
}

// Cause returns the cause of the error if it is a wrappedError.
func Cause(err error) error {
	var wErr wrappedError
	if errors.As(err, &wErr) {
		return wErr.Cause()
	}

	return err
}

// AsWrap checks if the error is a wrappedError and returns it.
func AsWrap(err error) (WrappedError, bool) {
	var wErr wrappedError
	if errors.As(err, &wErr) {
		return wErr, true
	}

	return nil, false
}

// Error returns a string representation of the wrapped error.
func (e wrappedError) Error() string {
	return _defaultFormatter(e.file, e.line, e.funcName, e.err, e.msg)
}

// Format returns a custom formatted string representation of the wrapped error.
func (e wrappedError) Format(fn FormatFn) string {
	return fn(e.file, e.line, e.funcName, e.err, e.msg)
}

// Unwrap returns the underlying error wrapped by this structure.
func (e wrappedError) Unwrap() error {
	return e.err
}

// Cause returns the cause of the wrapped error by recursively checking for wrappedError type
// and calling Unwrap() on it.
func (e wrappedError) Cause() error {
	var wErr wrappedError
	for errors.As(e.err, &wErr) {
		e.err = wErr.Unwrap()
	}

	return e.err
}

// File returns the file information associated with the wrapped error.
func (e wrappedError) File() string {
	return e.file
}

// Line returns the line information associated with the wrapped error.
func (e wrappedError) Line() int {
	return e.line
}

// FuncName returns the function name associated with the wrapped error.
//
// Example: "main.main".
func (e wrappedError) FuncName() string {
	return e.funcName
}

// Msg returns the additional message associated with the wrapped error.
func (e wrappedError) Msg() string {
	return e.msg
}
