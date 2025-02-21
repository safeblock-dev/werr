package werr

// ArgsFormat is a format string used to represent function arguments in error messages.
// When used with Wrapf, it formats the provided arguments using the %+v verb, e.g.:
// "args=[foo bar]".
const ArgsFormat = "args=%+v"

// UnwrapAll recursively unwraps an error that implements the Unwrap() method,
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

// Args converts a list of variadic arguments into a slice.
// It is a generic helper function that works with any type.
// Example: werr.Wrapf(errors.New("error"), werr.ArgsFormat, werr.Args("arg", 1)).
func Args[T any](args ...T) []T {
	return args
}
