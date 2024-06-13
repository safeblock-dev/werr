package werr

import (
	"errors"
)

// wrapError represents an error with additional context such as caller, function name, and message.
type wrapError struct {
	caller   string
	err      error
	funcName string
	msg      string
}

// newError creates a new wrapped error with caller information and additional message.
func newError(err error, msg string) error {
	sourceCaller, funcName := caller(defaultCallerSkip)

	return wrapError{
		caller:   sourceCaller,
		err:      err,
		msg:      msg,
		funcName: funcName,
	}
}

// Unwrap calls errors.Unwrap to retrieve the underlying error.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Cause returns the cause of the error if it is a wrapError.
// If the error is not a wrapError, it returns the original error.
func Cause(err error) error {
	var wErr wrapError
	if errors.As(err, &wErr) {
		return wErr.Cause()
	}

	return err
}

// UnwrapAll recursively traverses the wrapped errors and returns the innermost non-wrapped error.
// It checks for both Cause() and Unwrap() methods.
func UnwrapAll(err error) error {
	for {
		// Check if the error implements the Cause() method
		if e, ok := err.(interface{ Cause() error }); ok {
			err = e.Cause()

			continue
		}

		// Check if the error implements the Unwrap() method
		if e, ok := err.(interface{ Unwrap() error }); ok {
			err = e.Unwrap()

			continue
		}

		// If neither interface is implemented, return the error
		return err
	}
}

// Error returns a string representation of the wrapped error.
func (e wrapError) Error() string {
	return errorStackMarshaler(e.caller, e.err, e.funcName, e.msg)
}

// ErrorFormat returns a custom formatted string representation of the wrapped error.
func (e wrapError) ErrorFormat(fn ErrorStackMarshalerFn) string {
	return fn(e.caller, e.err, e.funcName, e.msg)
}

// Unwrap returns the underlying error wrapped by this structure.
func (e wrapError) Unwrap() error {
	return e.err
}

// Cause returns the cause of the wrapped error by recursively checking for wrapError type
// and calling Unwrap() on it.
func (e wrapError) Cause() error {
	var wErr wrapError
	for errors.As(e.err, &wErr) {
		e.err = wErr.Unwrap()
	}

	return e.err
}

// Caller returns the caller information associated with the wrapped error.
func (e wrapError) Caller() string {
	return e.caller
}

// FuncName returns the function name associated with the wrapped error.
func (e wrapError) FuncName() string {
	return e.funcName
}

// Msg returns the additional message associated with the wrapped error.
func (e wrapError) Msg() string {
	return e.msg
}
