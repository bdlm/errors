# errors

<a href="https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#mature"><img src="https://img.shields.io/badge/stability-mature-008000.svg" alt="Mature"></a> This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). This package is considered mature, you should expect package stability in <strong>Minor</strong> and <strong>Patch</strong> version releases

- **Major**: backwards incompatible package updates
- **Minor**: feature additions
- **Patch**: bug fixes, API route/ingress updates, DNS updates

**[CHANGELOG](CHANGELOG.md)**<br>

<p align="center">
	<a href="https://github.com/bdlm/errors/blob/master/LICENSE"><img src="https://img.shields.io/github/license/bdlm/errors.svg" alt="BSD-3-Clause"></a>
	<a href="https://travis-ci.org/bdlm/errors"><img src="https://travis-ci.org/bdlm/errors.svg?branch=master" alt="Build status"></a>
	<a href="https://codecov.io/gh/bdlm/errors"><img src="https://img.shields.io/codecov/c/github/bdlm/errors/master.svg" alt="Coverage status"></a>
	<a href="https://goreportcard.com/report/github.com/bdlm/errors"><img src="https://goreportcard.com/badge/github.com/bdlm/errors" alt="Go Report Card"></a>
	<a href="https://github.com/bdlm/errors/issues"><img src="https://img.shields.io/github/issues-raw/bdlm/errors.svg" alt="Github issues"></a>
	<a href="https://github.com/bdlm/errors/pulls"><img src="https://img.shields.io/github/issues-pr/bdlm/errors.svg" alt="Github pull requests"></a>
	<a href="https://godoc.org/github.com/bdlm/errors"><img src="https://godoc.org/github.com/bdlm/errors?status.svg" alt="GoDoc"></a>
</p>

`bdlm/errors` provides simple, concise, useful error handling and annotation.

One of the biggest frustrations with Go error handling is the lack of forensic and meta information errors should provide. By default errors are just a string and possibly a type. They can't tell you where they occurred or the path through the call stack they followed. The error implementation in Go is robust enough to control program flow but it's not very efficient for troubleshooting or analysis.

Since the idom in Go is that we pass the error back up the stack anyway:
```go
if nil != err {
	return err
}
```
it's trivial to make errors much more informative with a simple error package. `bdlm/errors` makes this easy and supports tracing the call stack and the error callers with relative ease. Custom error types are also fully compatible with this package and can be used freely.

## Changelog

All notable changes to this project are documented in the [`CHANGELOG`](CHANGELOG.md). The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Quick start

See the [documentation](https://godoc.org/github.com/bdlm/errors#pkg-examples) for more examples.

```
go get github.com/bdlm/errors
```

#### Create an error
```go
var MyError = errors.New("My error")
```

#### Create an error using formatting verbs
```go
var MyError = errors.Errorf("My error #%d", 1)
```

#### Wrap an error
```go
if nil != err {
	return errors.Wrap(err, "the operation failed")
}
```

#### Wrap an error with another error
```go
err := try1()
if nil != err {
	err2 := try2()
	if nil != err2 {
		return errors.WrapE(err, err2)
	}
	return err
}
```

#### Get the previous error, if any
```go
err := doWork()
if prevErr := errors.Unwrap(err); nil != prevErr {
	...
}
```

#### Test for a specific error type
```go
var MyError = errors.New("My error")
func main() {
	err := doWork()
	if errors.Is(err, MyError) {
		...
	}
}
```

#### Test to see if a specific error type exists anywhere in an error stack
```go
var MyError = errors.New("My error")
func main() {
	err := doWork()
	if errors.Has(err, MyError) {
		...
	}
}
```

#### Iterate through an error stack
```go
err := doWork()
for nil != err {
	fmt.Println(err)
	err = errors.Unwrap(err)
}
```

#

See the [documentation](https://godoc.org/github.com/bdlm/errors#pkg-examples) for more examples.


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
