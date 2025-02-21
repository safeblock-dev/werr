package werr

import (
	"fmt"
)

// Wrap takes an error and returns a new wrapped error.
// If the input error (err) is nil, the function returns nil.
// Otherwise, it creates a new wrapped error using the input error
// and an empty message text.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	return newError(err, "")
}

// Wrapf takes an error, a format string, and optional arguments, and returns a new wrapped error.
// If the input error (err) is nil, the function returns nil.
// Otherwise, it creates a new wrapped error using the input error
// and formats the message text based on the provided format and arguments.
func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	return newError(err, fmt.Sprintf(format, a...))
}

// Wrapt takes a value, an error and returns a new wrapped error.
// If the input error (err) is nil, the function returns nil.
// Otherwise, it creates a new wrapped error using the input error
// and an empty message text.
func Wrapt[T interface{}](val T, err error) (T, error) { //nolint: ireturn
	if err == nil {
		return val, nil
	}

	return val, newError(err, "")
}

// Unwrap retrieves the underlying error wrapped by the provided error.
// If the error is of type Error, it returns the original error stored within it.
// Otherwise, it returns nil.
func Unwrap(err error) error {
	if e, ok := err.(Error); ok { //nolint: errorlint
		return e.Unwrap()
	}

	return nil
}

// Cause returns the root cause of the error if it is an Error.
func Cause(err error) error {
	if e, ok := err.(Error); ok { //nolint: errorlint
		return e.Cause()
	}

	return err
}

// IsWrap checks if the provided error is of type Error.
func IsWrap(err error) bool {
	_, ok := err.(Error) //nolint: errorlint

	return ok
}
