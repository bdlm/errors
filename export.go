package errors

import (
	"fmt"
	"reflect"

	std_caller "github.com/bdlm/std/v2/caller"
	std_error "github.com/bdlm/std/v2/errors"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// As searches the error stack for an error that can be cast to the test
// argument, which must be a pointer. If it succeeds it performs the
// assignment and returns the result, otherwise it returns nil.
func As(err, test error) error {
	if test == nil {
		return nil
	}

	val := reflect.ValueOf(test)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return nil
	}
	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
		return nil
	}
	testType := typ.Elem()
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(testType) {
			val.Elem().Set(reflect.ValueOf(err))
			return err
		}
		if e, ok := err.(interface{ As(error) error }); ok {
			return e.As(test)
		}
		err = Unwrap(err)
	}
	return nil
}

// Caller returns the Caller associated with an error, if any.
func Caller(err error) std_caller.Caller {
	if e, ok := err.(interface{ Caller() std_caller.Caller }); ok {
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
	if e, ok := err.(interface{ Has(error) bool }); ok {
		return e.Has(test)

	}
	return Is(err, test)
}

// Is reports whether any error in err's chain matches test.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a test if it is equal to that test or if
// it implements a method Is(error) bool such that Is(test) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(test error) bool { return test == os.ErrExist }
//
// then Is(MyError{}, os.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library.
func Is(err, test error) bool {
	if nil == err || nil == test {
		return false
	}

	isComparable := reflect.TypeOf(err).Comparable() && reflect.TypeOf(test).Comparable()
	if isComparable && err == test {
		return true
	}

	if e, ok := err.(*E); ok {
		isComparable := reflect.TypeOf(e.err).Comparable() && reflect.TypeOf(test).Comparable()
		if isComparable && e.err == test {
			return true
		}
	}

	if e, ok := err.(interface{ Is(error) bool }); ok {
		return e.Is(test)
	}

	if err = Unwrap(err); err == nil {
		return false
	}

	return Is(err, test)
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
	clr.trace = std_caller.Trace{clr.trace[0]}
	if stdClr, ok := e.(std_error.ErrorCaller); ok {
		clr.trace = append(clr.trace, stdClr.Caller().Trace()...)
	}

	return &E{
		caller: clr,
		err:    e,
	}
}

// Track updates the error stack with additional caller data.
func Track(e error) *E {
	var stdE *E
	if nil == e {
		return nil
	}

	stdE, ok := e.(*E)
	if !ok {
		stdE = &E{
			err: e,
		}
		if clr, ok := e.(std_error.ErrorCaller); ok {
			stdE.caller = clr.Caller()
		} else {
			stdE.caller = NewCaller()
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
func Unwrap(err error) error {
	if e, ok := err.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
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
