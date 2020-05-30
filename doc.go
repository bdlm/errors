/*
Package errors provides simple, concise, useful error handling and annotation. This package
aims to implement the Error Inspection and Error Values Go2 draft designs.

https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md
https://go.googlesource.com/proposal/+/master/design/go2draft.md

	import (
		"github.com/bdlm/errors/v2"
	)

One of the biggest frustrations with Go error handling is the lack of forensic and meta
information errors should provide. By default errors are just a string and possibly a type.
They can't tell you where they occurred or the path through the call stack they followed.
The error implementation in Go is robust enough to control program flow but it's not very
efficient for troubleshooting or analysis.

Since the idom in Go is that we pass the error back up the stack anyway:

	if nil != err {
		return err
	}

it's trivial to make errors much more informative with a simple error package. `bdlm/errors`
makes this easy and supports tracing the call stack and the error callers with relative ease.
Custom error types are also fully compatible with this package and can be used freely.

Install

	go get github.com/bdlm/errors/v2

Quick Start

All package methods work with any `error` type as well as `nil` values, and error instances
implement the Unwrap, Is, Marshaler, and Formatter interfaces as well as the github.com/bdlm/std/errors
interfaces.

Create an error:

	var MyError = errors.New("My error")

Create an error using formatting verbs:

	var MyError = errors.Errorf("My error #%d", 1)


Wrap an error:

	if nil != err {
		return errors.Wrap(err, "the operation failed")
	}

Wrap an error with another error:

	err := try1()
	if nil != err {
		err2 := try2()
		if nil != err2 {
			return errors.Wrap(err, err2)
		}
		return err
	}

Get the previous error, if any:

	err := doWork()
	if prevErr := errors.Unwrap(err); nil != prevErr {
		...
	}

Test for a specific error type:

	var MyError = errors.New("My error")
	func main() {
		err := doWork()
		if errors.Is(err, MyError) {
			...
		}
	}

Test to see if a specific error type exists anywhere in an error stack:

	var MyError = errors.New("My error")
	func main() {
		err := doWork()
		if errors.Has(err, MyError) {
			...
		}
	}

Iterate through an error stack:

	err := doWork()
	for nil != err {
		fmt.Println(err)
		err = errors.Unwrap(err)
	}
*/
package errors
