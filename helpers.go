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

// Arg is syntactic sugar. Returns a message in the format "arg=?".
//
//	Code:
//
//	type CalculateParams struct {
//	  X, Y int
//	}
//
//	func Calculate(arg CalculateParams) (int, error) {
//	  if arg.X == 0 && arg.Y == 0 {
//	    return 0, werr.Wrapf(errors.New("zero arguments"), werr.Arg(arg))
//	  }
//
//	  return arg.X+arg.Y, nil
//	}
//
//	Output:	"arg={X:100 Y:50}"
func Arg(arg interface{}) string {
	return fmt.Sprintf("arg=%+v", arg)
}

// Args is syntactic sugar. Returns a message in the format "args=?".
//
//	Code:
//
//	func Calculate(x, y int) (int, error) {
//	  if x == 0 && y == 0 {
//	    return 0, werr.Wrapf(errors.New("zero arguments"), werr.Args(x, y))
//	  }
//
//	  return x+y, nil
//	}
//
//	Output:	"arg=[100, 50]"
func Args(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	return fmt.Sprintf("args=%+v", args)
}
