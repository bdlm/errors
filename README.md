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

Iterate through an error stack:
```go
err := doWork()
for nil != err {
	fmt.Println(err)
	err = errors.Unwrap(err)
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