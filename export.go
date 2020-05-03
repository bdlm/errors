package errors

import (
	"fmt"

	std_err "github.com/bdlm/std/v2/errors"
)

// Caller returns the Caller associated with an Error, if any.
func Caller(err error) std_err.Caller {
	if e, ok := err.(std_err.Error); ok {
		return e.Caller()
	}
	return nil
}

// Errorf formats according to a format specifier and returns an error that
// contains caller data.
func Errorf(msg string, data ...interface{}) *E {
	return New(fmt.Sprintf(msg, data...))
}

// Has returns whether an error or an error stack stack is or contains the
// referenced error type.
func Has(err, test error) bool {
	if nil == err || nil == test {
		return false
	}
	if std, ok := err.(std_err.Error); ok {
		return std.Has(test)
	}
	return Is(err, test)
}

// Is returns whether an error is the referenced error type.
func Is(err, test error) bool {
	if nil == err || nil == test {
		return false
	}
	if std, ok := err.(std_err.Error); ok {
		return std.Is(test)
	}
	return err == test && err.Error() == test.Error()
}

// New returns an error that contains caller data.
func New(msg string) *E {
	return &E{
		caller: NewCaller(),
		err:    fmt.Errorf(msg),
	}
}

// Trace adds an additional caller line to the error trace trace on an error
// to aid in debugging and forensic analysis.
func Trace(e error) *E {
	if nil == e {
		return nil
	}

	clr := NewCaller().(*caller)
	if std, ok := e.(std_err.Error); ok {
		clr.trace = std_err.Trace{clr.trace[0]}
		clr.trace = append(clr.trace, std.Caller().Trace()...)
	}

	return &E{
		caller: clr,
		err:    e,
	}
}

// Track updates the error stack with additional caller data.
func Track(e error) *E {
	var stdE std_err.Error
	if nil == e {
		return nil
	}

	stdE, ok := e.(std_err.Error)
	if !ok {
		stdE = &E{
			caller: NewCaller(),
			err:    e,
		}
	}

	return &E{
		caller: stdE.Caller(),
		err:    e,
		prev: &E{
			caller: NewCaller(),
			err:    fmt.Errorf("%s (tracked)", e),
			prev:   stdE.Unwrap(),
		},
	}
}

// Unwrap returns the previous error.
func Unwrap(e error) std_err.Error {
	if std, ok := e.(std_err.Error); ok {
		return std.Unwrap()
	}
	return nil
}

// Wrap returns a new error that wraps the provided error.
func Wrap(e error, msg string, data ...interface{}) *E {
	return WrapE(e, fmt.Errorf(msg, data...))
}

// WrapE returns a new error that wraps the provided error.
func WrapE(e, err error) *E {
	return &E{
		caller: NewCaller(),
		err:    err,
		prev:   e,
	}
}
