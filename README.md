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

`bdlm/errors` is inspired by [`pkg/errors`](https://github.com/pkg/errors) and uses a similar API.

```go
import (
	errs "github.com/bdlm/errors"
)
```

One of the most common frustrations with go error handling is the lack of exceptions in the language. You can use the `panic`/`recover` method to simulate the behavior but that's akward, clunky, requires many lines of code, and is hard to follow. It is also very common practice to handle an exception and run additional code or simply use an exception to inform higher-level code about internal behavior from deeper in the call stack.

As Pike says, errors are values, and when `panic`/`recover` isn't a reasonable solution you have to handle passing that information up the stack yourself. Which kind of sucks and leaves us with the `if err != nil` idiom which is fairly useless without a solid pattern behind it.

Since the idom is that we handle the error all the way up the stack anyway, it's trivial to make errors much, much more useful with a good error package. This package makes this simple and supports tracing the call stack and the error callers.

It is still possible to define custom errors. Custom error types are fully compatible with this package as well and can be used freely.

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