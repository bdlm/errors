package errors

import (
	"fmt"
)

// Errorf formats according to a format specifier and returns an error that
// contains caller data.
func Errorf(msg string, data ...interface{}) Err {
	return New(fmt.Sprintf(msg, data...))
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

// New returns an error that contains caller data.
func New(msg string) Err {
	return ex{
		caller: NewCaller(),
		err:    fmt.Errorf(msg),
	}
}

// Trace adds an additional caller line to the error trace trace on an error
// to aid in debugging and forensic analysis.
func Trace(e error) Err {
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

// Track updates the error stack with additional caller data.
func Track(e error) Err {
	err, ok := e.(ex)
	if !ok {
		err = ex{
			caller: NewCaller(),
			err:    e,
		}
	}

	return ex{
		caller: err.caller,
		err:    err.err,
		prev: ex{
			caller: NewCaller(),
			err:    fmt.Errorf("%s (tracked)", e),
			prev:   err.prev,
		},
	}
}

// Unwrap returns the previous error.
func Unwrap(e error) Err {
	if tmp, ok := e.(ex); ok {
		if nil != tmp.prev {
			if tmp2, ok := tmp.prev.(ex); ok {
				return tmp2
			}
			return ex{
				caller: NewCaller(),
				err:    tmp,
			}
		}
	}
	return nil
}

// Wrap returns a new error that wraps the provided error.
func Wrap(e error, msg string, data ...interface{}) Err {
	return WrapE(e, fmt.Errorf(msg, data...))
}

// WrapE returns a new error that wraps the provided error.
func WrapE(e, err error) Err {
	return ex{
		caller: NewCaller(),
		err:    err,
		prev:   e,
	}
}
