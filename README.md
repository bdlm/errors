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

`bdlm/errors` provides simple, concise, useful error handling and annotation.

`go get github.com/bdlm/errors`

One of the biggest frustrations with Go error handling is the lack of forensic and meta information errors can provide. Out of the box errors are just a string and possibly a type. They can't tell you where they occurred or the path through the call stack they followed. The error implementation in Go is robust enough to control program flow but it's not very efficient for troubleshooting or analasys.

Since the idom in Go is that we pass the error back up the stack anyway:
```go
if nil != err {
	return err
}
```
it's trivial to make errors much more informative with a simple error package. `bdlm/errors` makes this easy and supports tracing the call stack and the error callers with relative ease. Custom error types are also fully compatible with this package and can be used freely.

## Quick start

See the [Godoc](https://godoc.org/github.com/bdlm/errors) for more examples.

Create an error:
```go
var MyError = errors.New("My error")
```

Create an error using formatting verbs:
```go
var MyError = errors.Errorf("My error #%d", 1)
```

Wrap an error:
```go
if nil != err {
	return errors.Wrap(err, "the operation failed")
}
```

Wrap an error with another error:
```go
err := try1()
if nil != err {
	err2 := try2()
	if nil != err2 {
		return errors.Wrap(err, err2)
	}
	return err
}
```

Get the previous error, if any:
```go
err := doWork()
if prevErr := errors.Unwrap(err); nil != prevErr {
	...
}
```

Test for a specific error type:
```go
var MyError = errors.New("My error")
func main() {
	err := doWork()
	if errors.Is(err, MyError) {
		...
	}
}
```

Test to see if a specific error type exists anywhere in an error stack:
```go
var MyError = errors.New("My error")
func main() {
	err := doWork()
	if errors.Has(err, MyError) {
		...
	}
}
```



## The `Error` interface

The exported package methods return an interface that exposes additional metadata and troubleshooting information:

```go
// Error defines the error interface.
type Error interface {
	// Caller returns the associated Caller instance.
	Caller() Caller

	// Error implements error.
	Error() string

	// Has tests to see if the test error exists anywhere in the error
	// stack.
	Has(test error) bool

	// Is tests to see if the test error matches most recent error in the
	// stack.
	Is(test error) bool

	// Unwrap returns the next error, if any.
	Unwrap() Error
}

// Caller holds runtime.Caller data.
type Caller interface {
	// File returns the file in which the call occurred.
	File() string

	// Func returns the name of the function in which the call occurred.
	Func() string

	// Line returns the line number in the file in which the call occurred.
	Line() int

	// Pc returns the program counter.
	Pc() uintptr

	// Trace returns the call stack.
	Trace() []Caller
}
```



## Error stacks

An error stack is an array of errors.

### Create a new stack

```go
if !someWork() {
    err := errs.New("validation failed")
}
```

### Base a new stack off any error

```go
if err := someWork(); nil != err {
	return errs.Wrap(err, "could not read configuration")
}
```

## Define a new error with an error code

Creating a new error defines the root of a backtrace.
```go
_, err := ioutil.ReadAll(r)
if err != nil {
	return errs.New("read failed")
}
```

## Adding context to an error

The errors.Wrap function returns a new error stack, adding context as the top error in the stack:

```go
data, err := ioutil.ReadAll(r)
if err != nil {
	return errs.Wrap(err, "read failed")
}
```

In this case, if the original `err` is not an instance of `Err`, that error becomes the root of the error stack.

## Building an error stack

Most cases will build a stack trace off a series of errors returned from the call stack:

```go
import (
	"fmt"

	errs "github.com/bdlm/errors"
)

const (
	// Error codes below 1000 are reserved future use by the
	// "github.com/bdlm/errors" package.
	ConfigurationNotValid errs.Code = iota + 1000
)

func loadConfig() error {
	err := decodeConfig()
	return errs.Wrap(err, ConfigurationNotValid, "service configuration could not be loaded")
}

func decodeConfig() error {
	err := readConfig()
	return errs.Wrap(err, errs.ErrInvalidJSON, "could not decode configuration data")
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return errs.Wrap(err, errs.ErrEOF, "could not read configuration file")
}

func someWork() error {
	return fmt.Errorf("failed to do work")
}
```

But for cases where a set of errors need to be captured from a single procedure, the `Add()` call can be used. The with call adds an error to the stack behind the lead error:

```go
import (
	errs "github.com/bdlm/errors"
)

func doSteps() error {
	var errStack errs.Err

	err := doStep1()
	if nil != err {
		errStack.Add(err, "step 1 failed")
	}


	err = doStep2()
	if nil != err {
		errStack.Add(err, "step 2 failed")
	}

	err = doStep3()
	if nil != err {
		errStack.Add(err, "step 3 failed")
	}

	return errStack
}
```

## Root cause of an error stack

Retrieving the root cause of an error stack is straightforward:

```go
log.Println(err.(errs.Error).Cause())
```

You can easily switch on the type of any error in the stack (including the causer) as usual:

```go
switch err.(errs.Error).Cause().(type) {
case MyError:
        // handle specifically
default:
        // handle generically
}
```

## Iterating the error stack

Becase an error stack is just an array of errors iterating through it is trivial:

```go
for _, e := range err.(errs.Error).Stack() {
	fmt.Println(e.Code())
	fmt.Println(e.Error())
	fmt.Println(e.Msg())  // In the case of Wrap(), it is possible to suppliment
	                      // an error with additional information, which is
	                      // returned by Msg(). Otherwise, Msg() returns the same
	                      // string as Error().
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