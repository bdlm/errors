package errors

import (
	std_err "github.com/bdlm/std/v2/errors"
)

// E is a github.com/bdlm/std.Error interface implementation and simply wraps
// the exported package methods as a convenience.
type E struct {
	caller std_err.Caller
	err    error
	prev   error
}

// Caller implements std.Error.
func (e *E) Caller() std_err.Caller {
	if nil == e {
		return nil
	}
	return e.caller
}

// Error implements std.Error.
func (e *E) Error() string {
	return e.err.Error()
}

// Has implements std.Error.
func (e *E) Has(test error) bool {
	if e.Is(test) {
		return true
	}
	if prev := e.Unwrap(); nil != prev {
		return prev.Has(test)
	}
	return false
}

// Is implements std.Error.
func (e *E) Is(test error) bool {
	if nil == test {
		return false
	}
	if e.err == test && e.Error() == test.Error() {
		return true
	}
	if std, ok := test.(std_err.Error); ok {
		return func(e1, e2 std_err.Error) bool {
			return e1 == e2 && e1.Error() == e2.Error()
		}(e, std)
	}
	return false
}

// Unwrap implements std.Error.
func (e *E) Unwrap() error {
	if nil == e {
		return nil
	}
	if nil == e.prev {
		return nil
	}
	if prev, ok := e.prev.(std_err.Error); ok {
		return prev
	}
	return &E{
		err: e.prev,
	}
}

// list will convert the error stack into a simple array.
func list(e error) []error {
	ret := []error{}

	if nil != e {
		if std, ok := e.(std_err.Error); ok {
			ret = append(ret, e)
			ret = append(ret, list(std.Unwrap())...)
		}
	}

	return ret
}
