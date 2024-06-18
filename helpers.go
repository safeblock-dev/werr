package werr

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
