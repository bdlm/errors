package errors

import (
	"reflect"

	std_caller "github.com/bdlm/std/v2/caller"
	std_error "github.com/bdlm/std/v2/errors"
)

// E is a github.com/bdlm/std.Error interface implementation and simply wraps
// the exported package methods as a convenience.
type E struct {
	caller std_caller.Caller
	err    error
	prev   error
}

// Caller implements std_error.Caller.
func (e *E) Caller() std_caller.Caller {
	if nil == e {
		return nil
	}
	return e.caller
}

// Error implements std_error.Error.
func (e *E) Error() string {
	return e.err.Error()
}

// Is implements std_error.Error.
func (e *E) Is(test error) bool {
	if nil == test {
		return false
	}

	isComparable := reflect.TypeOf(e).Comparable() && reflect.TypeOf(test).Comparable()
	if isComparable && e == test {
		return true
	}
	isComparable = reflect.TypeOf(e.err).Comparable() && reflect.TypeOf(test).Comparable()
	if isComparable && e.err == test {
		return true
	}

	if testE, ok := test.(*E); ok {
		isComparable = reflect.TypeOf(e).Comparable() && reflect.TypeOf(testE).Comparable()
		if isComparable && e == testE {
			return true
		}
		isComparable = reflect.TypeOf(e.err).Comparable() && reflect.TypeOf(testE.err).Comparable()
		if isComparable && e.err == testE.err {
			return true
		}
	}

	if err := e.Unwrap(); nil != err {
		if err, ok := err.(interface{ Is(error) bool }); ok {
			return err.Is(test)
		}
	}

	return false
}

// Unwrap implements std_error.Wrapper.
func (e *E) Unwrap() error {
	if nil == e {
		return nil
	}
	if nil == e.prev {
		return nil
	}
	if prev, ok := e.prev.(std_error.Error); ok {
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
		if std, ok := e.(std_error.Wrapper); ok {
			ret = append(ret, e)
			ret = append(ret, list(std.Unwrap())...)
		}
	}

	return ret
}
