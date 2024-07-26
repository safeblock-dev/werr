package werr

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// PanicToError converts a recovered panic to an error.
func PanicToError(p any) error {
	msg := "panic recovered\n"
	msg += string(debug.Stack())

	switch v := p.(type) {
	case nil:
		return nil
	case error:
		return newError(v, msg)
	case string:
		return newError(errors.New(v), msg)
	default:
		return newError(fmt.Errorf("%#v", v), msg)
	}
}
