package werr

import (
	"fmt"
)

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

// Arg syntactic sugar. Returns a message in the format "arg=?".
//
//	Example 0:
//	func Calculate(x, y int) {
//		// calculate...
//		return werr.Wrapf(err, werr.Arg(x,y))
//	}
//	// Output:	"arg=[100, 50]"
//
//	Example 1:
//	type CalculateParams struct {
//		X, Y int
//	}
//	func Calculate(arg CalculateParams) {
//		// calculate...
//		return werr.Wrapf(err, werr.Arg(arg))
//	}
//	// Output:	"arg={X:100 Y:50}"
func Arg(arg ...interface{}) string {
	switch len(arg) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("arg=%+v", arg[0])
	default:
		return fmt.Sprintf("arg=%+v", arg)
	}
}
