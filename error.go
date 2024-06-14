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
	// funcName function name in the format "<pkg>.<name>", e.g. "main.main".
	funcName string

	// file and line are the file name and line number of the
	// location in this frame. For non-leaf frames, this will be
	// the location of a call. These may be the empty string and
	// zero, respectively, if not known.
	file string
	line int

	// err error for which you want to save its funcName, file, line and msg.
	err error

	// msg optional message to help analyze the error.
	msg string
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

// Cause returns the cause of the error if it is a wrappedError.
// If the error is not a wrappedError, it returns the original error.
func Cause(err error) error {
	var wErr wrappedError
	if errors.As(err, &wErr) {
		return wErr.Cause()
	}

	return err
}

// Error returns a string representation of the wrapped error.
func (e wrappedError) Error() string {
	return DefaultErrorStackMarshaler(e.file, e.line, e.funcName, e.err, e.msg)
}

// Format returns a custom formatted string representation of the wrapped error.
func (e wrappedError) Format(fn ErrorStackMarshalerFn) string {
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
// Example: "main.main"
func (e wrappedError) FuncName() string {
	return e.funcName
}

// Msg returns the additional message associated with the wrapped error.
func (e wrappedError) Msg() string {
	return e.msg
}
