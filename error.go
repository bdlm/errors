package errors

import (
	std_errors "errors"

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
//
// The wrapped error's message is included, "outer: inner", as fmt.Errorf("%w") does. Returning only
// this frame's own message discarded everything underneath it wherever an error is rendered by its
// Error method rather than by a %+v format verb -- which is most places that matter, including an
// HTTP or gRPC response body. A caller was told "unable to resolve the signing key" with no hint of
// the cause that was already recorded one link down.
func (e *E) Error() string {
	if nil == e {
		return ""
	}
	// A frame can legitimately carry no message of its own: Trace adds caller data without
	// annotating, and WrapE accepts a nil annotation. Such a frame must be transparent rather
	// than swallowing the chain or emitting a leading ": ".
	if nil == e.err {
		if nil == e.prev {
			return ""
		}
		return e.prev.Error()
	}
	if nil == e.prev {
		return e.err.Error()
	}
	return e.err.Error() + ": " + e.prev.Error()
}

// message is THIS frame's own message, without the wrapped chain.
//
// Error() deliberately includes the cause, which is what a caller rendering a single string wants.
// The trace formats want the opposite: they already print one line per frame, so a line carrying its
// own downstream chain repeats everything below it. This is the accessor those formats use.
func (e *E) message() string {
	if nil == e || nil == e.err {
		return ""
	}
	return e.err.Error()
}

// Is implements std_error.Error.
func (e *E) Is(test error) bool {
	if nil == e || nil == test {
		return false
	}

	if comparableErrors(e, test) && error(e) == test {
		return true
	}
	if comparableErrors(e.err, test) && e.err == test {
		return true
	}

	if testE, ok := test.(*E); ok {
		if comparableErrors(e, testE) && error(e) == error(testE) {
			return true
		}
		if comparableErrors(e.err, testE.err) && e.err == testE.err {
			return true
		}
	}

	// This method searches the CHAIN, not just this frame. That is this package's own documented
	// behavior (E.Is is part of std_error.Error and is used directly, not only as the errors.Is
	// hook), and it is safe to keep: the standard library's errors.Is calls the hook as
	// `ok && x.Is(target)`, so a false answer here does not end its walk -- it unwraps and asks the
	// next link. The defect was never here. It was the PACKAGE-level Is in export.go, which returned
	// this method's answer directly and so did end the walk, and Unwrap, which never yielded a
	// wrapped error of a foreign type for anything to be asked about.
	if err := e.Unwrap(); nil != err {
		if err, ok := err.(interface{ Is(error) bool }); ok {
			return err.Is(test)
		}
		return Is(err, test)
	}

	return false
}

// As implements the standard library's As hook, interface{ As(interface{}) bool }, which
// errors.As consults on each link before unwrapping.
//
// It exists because WrapE puts its annotation in e.err, which is NOT on the Unwrap chain -- Unwrap
// yields e.prev. Is already special-cases e.err, so without this hook Is could find an error that
// As could not, and a caller had no way to know which of the two would work. For an error package
// the two must agree about what is in the chain.
//
// Delegating to the standard library rather than comparing types by hand also searches the
// annotation's OWN chain, which is what errors.As would have done had the value been reachable.
func (e *E) As(target interface{}) bool {
	if nil == e || nil == e.err || nil == target {
		return false
	}
	return std_errors.As(e.err, target)
}

// Unwrap implements std_error.Wrapper.
//
// It returns the wrapped error ITSELF, whatever its concrete type. Previously a wrapped error that
// did not implement std_error.Error was re-boxed as &E{err: prev}, and because that box carried no
// prev of its own the chain ended there: the original error value was never handed to the caller, so
// its own Unwrap, Is and As were never consulted.
//
// That broke every consumer of the standard errors package. errors.As could only ever match *E,
// since no other concrete type was ever yielded -- so, for example, grpc's status.FromError, whose
// whole job is errors.As(err, &interface{ GRPCStatus() *Status }), could not find a status error even
// when one was wrapped one level down. errors.Is could only reach sentinels that happened to be
// created by this package, because those alone satisfied std_error.Error and were passed through.
//
// Returning e.prev is what the errors.Unwrap contract asks for: "the result of calling Unwrap is the
// underlying error", not a copy of it.
func (e *E) Unwrap() error {
	if nil == e {
		return nil
	}
	return e.prev
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
