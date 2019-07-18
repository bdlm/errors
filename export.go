package errors

import (
	"fmt"
)

// Errorf formats according to a format specifier and returns an error that
// contains caller data.
func Errorf(msg string, data ...interface{}) Error {
	return New(fmt.Sprintf(msg, data...))
}

// GetCaller returns the Caller associated with an Error, if any.
func GetCaller(err error) Caller {
	if e, ok := err.(Error); ok {
		return e.Caller()
	}
	return nil
}

// Has returns whether an error or an error stack stack is or contains the
// referenced error type.
func Has(err, test error) bool {
	if nil == err || nil == test {
		return false
	}
	if tmp, ok := err.(E); ok {
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
	if tmp, ok := err.(E); ok {
		return tmp.err == test
	}
	return err == test
}

// New returns an error that contains caller data.
func New(msg string) Error {
	return E{
		caller: NewCaller(),
		err:    fmt.Errorf(msg),
	}
}

// Trace adds an additional caller line to the error trace trace on an error
// to aid in debugging and forensic analysis.
func Trace(e error) Error {
	if nil == e {
		return nil
	}

	clr := NewCaller().(caller)
	if tmp, ok := e.(E); ok {
		clr.trace = []Caller{clr.trace[0]}
		clr.trace = append(clr.trace, tmp.caller.Trace()...)
		tmp.caller = clr
		return tmp
	}

	return E{
		caller: clr,
		err:    e,
	}
}

// Track updates the error stack with additional caller data.
func Track(e error) Error {
	err, ok := e.(E)
	if !ok {
		err = E{
			caller: NewCaller(),
			err:    e,
		}
	}

	return E{
		caller: err.caller,
		err:    err.err,
		prev: E{
			caller: NewCaller(),
			err:    fmt.Errorf("%s (tracked)", e),
			prev:   err.prev,
		},
	}
}

// Unwrap returns the previous error.
func Unwrap(e error) Error {
	if tmp, ok := e.(E); ok {
		if nil == tmp.prev {
			return nil
		}
		if tmp2, ok := tmp.prev.(E); ok {
			return tmp2
		}
		return E{
			caller: NewCaller(),
			err:    tmp.prev,
		}
	}
	return nil
}

// Wrap returns a new error that wraps the provided error.
func Wrap(e error, msg string, data ...interface{}) Error {
	return WrapE(e, fmt.Errorf(msg, data...))
}

// WrapE returns a new error that wraps the provided error.
func WrapE(e, err error) Error {
	return E{
		caller: NewCaller(),
		err:    err,
		prev:   e,
	}
}
