# errors

<p align="center">
	<a href="https://github.com/bdlm/errors/blob/master/LICENSE"><img src="https://img.shields.io/github/license/bdlm/errors.svg" alt="BSD-2-Clause"></a>
	<a href="https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#beta"><img src="https://img.shields.io/badge/stability-beta-33bbff.svg" alt="Beta"></a>
	<a href="https://travis-ci.org/bdlm/errors"><img src="https://travis-ci.org/bdlm/errors.svg?branch=master" alt="Build status"></a>
	<a href="https://codecov.io/gh/bdlm/errors"><img src="https://img.shields.io/codecov/c/github/bdlm/errors/master.svg" alt="Coverage status"></a>
	<a href="https://goreportcard.com/report/github.com/bdlm/errors"><img src="https://goreportcard.com/badge/github.com/bdlm/errors" alt="Go Report Card"></a>
	<a href="https://github.com/bdlm/errors/issues"><img src="https://img.shields.io/github/issues-raw/bdlm/errors.svg" alt="Github issues"></a>
	<a href="https://github.com/bdlm/errors/pulls"><img src="https://img.shields.io/github/issues-pr/bdlm/errors.svg" alt="Github pull requests"></a>
	<a href="https://godoc.org/github.com/bdlm/errors"><img src="https://godoc.org/github.com/bdlm/errors?status.svg" alt="GoDoc"></a>
</p>


```go
import (
	errs "github.com/bdlm/errors"
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

Adding support for error codes is the primary motivation behind this project. See [`codes.go`](https://github.com/bdlm/errors/blob/master/codes.go). `HTTPStatus` is optional and a convenience property that allows automation of HTTP status responses based on internal error codes. The `Code` definition associated with error at the top of the stack (most recent error) should be used for HTTP status output.

```go
import (
	errs "github.com/bdlm/errors"
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
	errs "github.com/bdlm/errors"
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
	errs "github.com/bdlm/errors"
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
switch err.(errs.Err).Cause().(type) {
case *MyError:
        // handle specifically
default:
        // unknown error
}
```

## Output formats

The Formatter interface has been implemented to provide access to a stack trace with the `%v` verb.

Standard error output, use with error codes to ensure appropriate user-facing messages `%v`:
```
0000: failed to load configuration
```

Single-line stack trace, useful for logging `%-v`:
```
#4 - "failed to load configuration" examples_test.go:36 `github.com/bdlm/errors_test.ExampleWrap_backtrace` {0000: unknown error} #3 - "service configuration could not be loaded" mocks_test.go:16 `github.com/bdlm/errors_test.loadConfig` {0001: fatal error} #2 - "could not decode configuration data" mocks_test.go:21 `github.com/bdlm/errors_test.decodeConfig` {0200: invalid JSON data could not be decoded} #1 - "could not read configuration file" mocks_test.go:26 `github.com/bdlm/errors_test.readConfig` {0100: unexpected EOF} #0 - "read: end of input" mocks_test.go:26 `github.com/bdlm/errors_test.readConfig` {0000: unknown error}
```

Multi-line condensed stack trace `%#v`:
```
#4 - "failed to load configuration" examples_test.go:36 `github.com/bdlm/errors_test.ExampleWrap_backtrace` {0000: unknown error}
#3 - "service configuration could not be loaded" mocks_test.go:16 `github.com/bdlm/errors_test.loadConfig` {0001: fatal error}
#2 - "could not decode configuration data" mocks_test.go:21 `github.com/bdlm/errors_test.decodeConfig` {0200: invalid JSON data could not be decoded}
#1 - "could not read configuration file" mocks_test.go:26 `github.com/bdlm/errors_test.readConfig` {0100: unexpected EOF}
#0 - "read: end of input" mocks_test.go:26 `github.com/bdlm/errors_test.readConfig` {0000: unknown error}
```

Multi-line detailed stack trace `%+v`:
```
#4: `github.com/bdlm/errors_test.ExampleWrap_backtrace`
	error:   failed to load configuration
	line:    examples_test.go:36
	code:    0000: unknown error
	message: 0000: failed to load configuration
#3: `github.com/bdlm/errors_test.loadConfig`
	error:   service configuration could not be loaded
	line:    mocks_test.go:16
	code:    0001: fatal error
	message: 0001: Internal Server Error
#2: `github.com/bdlm/errors_test.decodeConfig`
	error:   could not decode configuration data
	line:    mocks_test.go:21
	code:    0200: invalid JSON data could not be decoded
	message: 0200: Invalid JSON Data
#1: `github.com/bdlm/errors_test.readConfig`
	error:   could not read configuration file
	line:    mocks_test.go:26
	code:    0100: unexpected EOF
	message: 0100: End of input
#0: `github.com/bdlm/errors_test.readConfig`
	error:   read: end of input
	line:    mocks_test.go:26
	code:    0000: unknown error
	message: 0000: read: end of input
```
