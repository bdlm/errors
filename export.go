package errors

import (
	"fmt"
)

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(msg string, data ...interface{}) error {
	return New(msg, data...)
}

// Has returns whether an error or an error stack stack is or contains the
// referenced error type.
func Has(err, test error) bool {
	if nil == err || nil == test {
		return false
	}
	if tmp, ok := err.(ex); ok {
		if tmp.err == test {
			return true
		}
		return Has(tmp.prev, test)
	}
	return err == test

}

// Is returns whether an error or an error stack stack is the referenced
// error type.
func Is(err, test error) bool {
	if tmp, ok := err.(ex); ok {
		return tmp.err == test
	}
	return err == test
}

// New formats according to a format specifier and returns the string
// as a value that satisfies error.
func New(msg string, data ...interface{}) error {
	return ex{
		caller: NewCaller(),
		err:    fmt.Errorf(msg, data...),
	}
}

// Track updates caller metadata on an error as it's passed back up the
// stack.
func Track(e error) error {
	if nil == e {
		return nil
	}

	a := NewCaller()
	clr := a.(caller)

	if tmp, ok := e.(ex); ok {
		clr.trace = []Caller{clr.trace[0]}
		clr.trace = append(clr.trace, tmp.caller.Trace()...)
		tmp.caller = clr
		return tmp
	}

	return ex{
		caller: clr,
		err:    e,
	}
}

// Unwrap returns the previous error.
func Unwrap(e error) error {
	if tmp, ok := e.(ex); ok {
		return tmp.prev
	}
	return nil
}

// Wrap returns a new error that wrapps the provided error.
func Wrap(e error, msg string, data ...interface{}) error {
	return WrapE(e, fmt.Errorf(msg, data...))
}

// WrapE returns a new error that wrapps the provided error.
func WrapE(e, err error) error {
	return ex{
		caller: NewCaller(),
		err:    err,
		prev:   e,
	}
}
