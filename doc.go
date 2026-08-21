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

# Quick Start

All package methods work with any `error` type as well as `nil` values -- including a nil *E
receiver -- and error instances implement the Unwrap, Is, As, Marshaler, and Formatter interfaces as
well as the github.com/bdlm/std/errors interfaces.

Two guarantees are worth stating explicitly, because both are load-bearing and both are covered by
tests: nothing in this package panics, for any combination of nil arguments and nil receivers; and
Is and As agree with the standard library's Is and As on every chain shape, so it never matters
which of the two a caller reaches for.

New and Wrap take a MESSAGE, stored verbatim. Only Errorf, and Wrap when given arguments, treat it
as a format string -- so a message containing a percent sign survives intact.

Create an error:

	var MyError = errors.New("My error")

Create an error using formatting verbs:

	var MyError = errors.Errorf("My error #%d", 1)

Wrap an error:

	if nil != err {
		return errors.Wrap(err, "the operation failed")
	}

Wrap an error with another error. Note this is WrapE: Wrap takes a message string, so passing an
error to it does not compile.

	err := try1()
	if nil != err {
		err2 := try2()
		if nil != err2 {
			return errors.WrapE(err, err2)
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

Find a specific error type anywhere in a chain:

	var target *MyErrorType
	if errors.As(err, &target) {
		...
	}

As takes the same arguments as the standard library's, target being a pointer to the type or
interface being looked for, and returns whether one was found. It panics if target is not a non-nil
pointer to either a type implementing error or to any interface type -- an unusable target is a
mistake at the call site, not a condition to report.

Is searches the whole chain, so there is no separate Has: the test above answers "does this
error, or anything it wraps, match" -- at any depth, through fmt.Errorf("%w") wrappers, foreign
error types and joined errors alike.

Iterate through an error stack:

	err := doWork()
	for nil != err {
		fmt.Println(err)
		err = errors.Unwrap(err)
	}
*/
package errors
