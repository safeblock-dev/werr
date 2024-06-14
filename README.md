## Highlights

The **werr** library provides efficient error wrapping capabilities for Go. It is designed with a focus on performance and enables the recording of functions where an error occurred. Here are some key features of the library:

* Recording information about functions that triggered the error for easy debugging.
* Support for custom error messages to make errors more informative.
* Very high performance compared to other error-wrapping libraries.

## Introduction

When handling errors in Go, simplicity and performance matter. **werr** offers a streamlined approach:

```go
return werr.Wrap(err)
```

This captures the context of where the error originated, aiding in quick debugging.

# werr

To use **werr** in your Go application, simply import it:

```go
import "github.com/safeblock-dev/werr"
```

## Features

* **Error Creation**: Create errors using `errors.New("error message")` as usual.
* **Error Wrapping**: Enhance errors with context using `werr.Wrap(err)`.
* **Custom Messages**: Add custom messages to errors with `werr.Wrapf(err, "custom error message")`.
* **Error Unwrapping**: Retrieve the original error with `werr.Unwrap(err)` for seamless error propagation.
* **Full Unwrapping**: Get the root cause of wrapped errors with `werr.UnwrapAll(err)`.
* **Direct Cause**: Identify the immediate cause of an error with `werr.Cause(err)`.
* **Typed Assertion**: Check error types with `werr.AsWrap(err)` for precise error handling.

## Example

```go
package main

import (
	"errors"
	"fmt"

	"github.com/safeblock-dev/werr"
)

var errExample = errors.New("find me")

func main() {
	err := example()
	if errors.Is(err, errExample) {
		fmt.Printf("trace: \n%v\n", err)
		fmt.Printf("\nunwrap: %v\n", werr.Unwrap(err))
	}
}

func example() error {
	return werr.Wrap(example2())
}

func example2() error {
	return werr.Wrapf(example3(), "without if")
}

func example3() error {
	if err := newError(); err != nil {
		return werr.Wrapf(err, "wow error!")
	}

	return nil
}

func newError() error {
	return errExample
}
```

#### Result

```
trace: 
main/main.go:21 example()
main/main.go:25 example2()      without if
main/main.go:30 example3()      wow error!
find me

unwrap: find me
```

## Stack Traces Benchmark

Performance benchmarks showcase **werr**'s efficiency in error handling:

Result sample, MacBook Air M1 @ 3.2GHz:

| name                           |     runs | ns/op | note                                         |
|--------------------------------|---------:|------:|----------------------------------------------|
| BenchmarkSimpleError10         | 39252192 | 28.67 | simple error, 10 frames deep                 |
| BenchmarkWrapError10           |  2190848 | 543.0 | with wrap error                              |
| BenchmarkWrapMsgError10        |  1881180 | 588.4 | with message                                 |
| BenchmarkErrorxError10         |   969931 | 1365  | errorx library                               |
| BenchmarkGoErrorsError10       |  1000000 | 1171  | go-errors library                            |
|                                |          |       |                                              |
| BenchmarkSimpleError100        |  1911195 | 633.1 | simple error, 100 frames deep                |
| BenchmarkWrapError100          |  1000000 | 1276  | with wrap error                              |
| BenchmarkWrapMsgError100       |   978968 | 1250  | with message                                 |
| BenchmarkErrorxError100        |   312757 | 3890  | errorx library                               |
| BenchmarkGoErrorsError100      |   519740 | 2407  | go-errors library                            |
|                                |          |       |                                              |
| BenchmarkSimpleErrorPrint100   |  1721454 | 695.1 | simple error, 100 frames deep, format output |
| BenchmarkWrapErrorPrint100     |   856266 | 1384  | with wrap error                              |
| BenchmarkWrapMsgErrorPrint100  |   834480 | 1441  | with message                                 |
| BenchmarkErrorxErrorPrint100   |    38940 | 30826 | errorx library, format output                |
| BenchmarkGoErrorsErrorPrint100 |   505380 | 2376  | go-errors library, format output             |

Key takeaways:

* **werr** provides efficient error creation and wrapping with minimal overhead.
* Use `werr.Unwrap(err)` to retrieve the original error for seamless propagation.
* Enhance error context with custom messages using `werr.Wrapf(err, "custom message")`.
* Identify the immediate cause of an error with `werr.Cause(err)` for precise error handling.
* Verify error types with `werr.AsWrap(err)` to handle errors based on specific types.

## More

Portions of the description and benchmark were adapted from the project [errorx](https://github.com/joomcode/errorx).

---

This update integrates the new functionality descriptions and benchmarks, ensuring clarity and interest in error handling with the **werr** library.