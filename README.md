# go-errors

<p align="center">
	<a href="https://github.com/mkenney/go-errors/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mkenney/go-errors.svg" alt="MIT License"></a>
	<a href="https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#beta"><img src="https://img.shields.io/badge/stability-beta-33bbff.svg" alt="Beta"></a>
	<a href="https://travis-ci.org/mkenney/go-errors"><img src="https://travis-ci.org/mkenney/go-errors.svg?branch=master" alt="Build status"></a>
	<a href="https://codecov.io/gh/mkenney/go-errors"><img src="https://img.shields.io/codecov/c/github/mkenney/go-errors/master.svg" alt="Coverage status"></a>
	<a href="https://goreportcard.com/report/github.com/mkenney/go-errors"><img src="https://goreportcard.com/badge/github.com/mkenney/go-errors" alt="Go Report Card"></a>
	<a href="https://github.com/mkenney/go-errors/issues"><img src="https://img.shields.io/github/issues-raw/mkenney/go-errors.svg" alt="Github issues"></a>
	<a href="https://github.com/mkenney/go-errors/pulls"><img src="https://img.shields.io/github/issues-pr/mkenney/go-errors.svg" alt="Github pull requests"></a>
	<a href="https://godoc.org/github.com/mkenney/go-errors"><img src="https://godoc.org/github.com/mkenney/go-errors?status.svg" alt="GoDoc"></a>
</p>


```go
import (
	errs "github.com/mkenney/go-errors"
)
```

Go errors is inspired by [`pkg/errors`](https://github.com/pkg/errors) and uses a similar API but adds support for error codes. Error codes are always optional.

## Error stacks

An error stack is an array of errors.

### Create a new stack

```go
if !decodeSomeJSON() {
    err := errs.New("validation failed")
}
```

### Base a new stack off any error

```go
err := decodeSomeJSON()
err = errs.Wrap(err, "could not read configuration")
```

## Define error codes

Adding support for error codes is the primary motivation behind this project. See [`codes.go`](https://github.com/mkenney/go-errors/blob/master/codes.go). `HTTPStatus` is optional and a convenience property that allows automation of HTTP status responses based on internal error codes. The `Code` definition associated with error at the top of the stack (most recent error) should be used for HTTP status output.

```go
import (
	errs "github.com/mkenney/go-errors"
)

const (
	// Error codes below 1000 are reserved future use by the errors
	// package.
	UserError errs.Code = iota + 1000
	InternalError
)

func init() {
	errs.Codes[UserError] = errs.Metadata{
		Internal:   "bad user input",
		External:   "A user error occurred",
		HTTPStatus: 400,
	}
	errs.Codes[InternalError] = errs.Metadata{
		Internal:   "could not save data",
		External:   "An internal server occurred",
		HTTPStatus: 500,
	}
}

func SomeFunc() error {
	return errs.New("SomeFunc failed because of things", InternalError)
}
```

## Define a new error with an error code

Creating a new error defines the root of a backtrace.
```go
_, err := ioutil.ReadAll(r)
if err != nil {
	return errs.New("read failed", errs.ErrUnknown)
}
```

## Adding context to an error

The errors.Wrap function returns a new error that adds context to the original error and starts an error stack:
```go
_, err := ioutil.ReadAll(r)
if err != nil {
	return errs.Wrap(err, "read failed", errs.ErrUnknown)
}
```

In this case, if the original `err` is not an instance of `Stack`, that error becomes the root of the error stack.

## Building an error stack

Most cases will build a stack trace off a series of errors returned from the call stack:

```go
import (
	"fmt"
	errs "github.com/mkenney/go-errors"
)

func main() {
	err := loadConfig()
	fmt.Printf("%#v", err)
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return errs.Wrap(err, "could not read configuration file", errs.ErrEOF)
}

func decodeConfig() error {
	err := readConfig()
	return errs.Wrap(err, "could not decode configuration data", errs.ErrInvalidJSON)
}

func loadConfig() error {
	err := decodeConfig()
	return errs.Wrap(err, "service configuration could not be loaded", errs.ErrFatal)
}
```

But for cases where a set of errors need to be captured from a single procedure, the `With()` call can be used. The with call adds an error to the stack behind the leading error:

```go
import (
	errs "github.com/mkenney/go-errors"
)

func doSteps() error {
	var errStack errs.Err

	err := doStep1()
	if nil != err {
		errStack.With(err, "step 1 failed")
	}


	err = doStep2()
	if nil != err {
		errStack.With(err, "step 2 failed")
	}

	err = doStep3()
	if nil != err {
		errStack.With(err, "step 3 failed")
	}

	return errStack
}
```

## Root cause of an error stack

Retrieving the root cause of an error stack is straightforward:

```go
log.Println(err.(errs.Stack).Cause())
```

Similar to `pkg/errors`, you can easily switch on the type of any error in the stack (including the causer):

```go
switch err.(errs.Err).Cause().err.(type) {
case *MyError:
        // handle specifically
default:
        // unknown error
}
```

## Output formats

The Formatter interface has been implemented to provide access to a stack trace with the `%v` verb.

Standard error output, use with error codes to ensure appropriate user-facing messages `%s`:
```
0002: Internal Server Error
```

Single-line stack trace, useful for logging `%v`:
```
#0 - "service configuration could not be loaded" example_test.go:22 `github.com/mkenney/go-errors_test.loadConfig` {0002: a fatal error occurred} #1 - "could
not decode configuration data" example_test.go:17 `github.com/mkenney/go-errors_test.decodeConfig` {0200: invalid JSON data could not be decoded} #2 - "could
not read configuration file" example_test.go:12 `github.com/mkenney/go-errors_test.readConfig` {0100: unexpected EOF}
```

Multi-line condensed stack trace `%#v`:
```
#0 - "service configuration could not be loaded" example_test.go:22 `github.com/mkenney/go-errors_test.loadConfig` {0002: a fatal error occurred}
#1 - "could not decode configuration data" example_test.go:17 `github.com/mkenney/go-errors_test.decodeConfig` {0200: invalid JSON data could not be decoded}
#2 - "could not read configuration file" example_test.go:12 `github.com/mkenney/go-errors_test.readConfig` {0100: unexpected EOF}
```

Multi-line detailed stack trace `%+v`:
```
#0: `github.com/mkenney/go-errors_test.loadConfig`
        error:   service configuration could not be loaded
        line:    example_test.go:22
        code:    2 - a fatal error occurred
        entry:   17741072
        message: Internal Server Error

#1: `github.com/mkenney/go-errors_test.decodeConfig`
        error:   could not decode configuration data
        line:    example_test.go:17
        code:    200 - invalid JSON data could not be decoded
        entry:   17740848
        message: Invalid JSON Data

#2: `github.com/mkenney/go-errors_test.readConfig`
        error:   could not read configuration file
        line:    example_test.go:12
        code:    100 - unexpected EOF
        entry:   17740576
        message: End of input
```
