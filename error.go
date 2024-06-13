package werr

type wrapError struct {
	caller   string
	err      error
	funcName string
	msg      string
}

func newError(err error, msg string) error {
	sourceCaller, funcName := caller(defaultCallerSkip)

	return wrapError{
		caller:   sourceCaller,
		err:      err,
		msg:      msg,
		funcName: funcName,
	}
}

// Unwrap returns the result of calling the Unwrap method on err,
// if err's type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	u, ok := err.(interface{ Unwrap() error })
	if !ok {
		return nil
	}

	return u.Unwrap()
}

// Cause recursively traverses the wrapped errors and returns the innermost non-wrapped error.
// It checks for both Cause() and Unwrap() methods.
func Cause(err error) error {
	for {
		// Check if the error implements the causer interface
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

func (e wrapError) Error() string {
	return errorStackMarshaler(e.caller, e.err, e.funcName, e.msg)
}

func (e wrapError) ErrorFormat(fn ErrorStackMarshalerFn) string {
	return fn(e.caller, e.err, e.funcName, e.msg)
}

func (e wrapError) Unwrap() error {
	return e.err
}

func (e wrapError) Cause() error {
	return e.err
}

func (e wrapError) Caller() string {
	return e.caller
}

func (e wrapError) FuncName() string {
	return e.funcName
}

func (e wrapError) Msg() string {
	return e.msg
}
