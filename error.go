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

// Caller implements Error.
func (e E) Caller() std_err.Caller {
	return e.caller
}

// Error implements Error.
func (e E) Error() string {
	return e.err.Error()
}

// Has implements Error.
func (e E) Has(test error) bool {
	return Has(e, test)
}

// Is implements Error.
func (e E) Is(test error) bool {
	return Is(e, test)
}

// Unwrap implements Error.
func (e E) Unwrap() std_err.Error {
	return Unwrap(e)
}

// list will convert the error stack into a simple array.
func list(e error) []error {
	ret := []error{}

	if nil != e {
		if tmp, ok := e.(E); ok {
			ret = append(ret, e)
			ret = append(ret, list(tmp.prev)...)
		}
	}

	return ret
}
