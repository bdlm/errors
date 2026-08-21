package errors

import (
	std_errors "errors"
	"fmt"
	"reflect"

	std_caller "github.com/bdlm/std/v2/caller"
	std_error "github.com/bdlm/std/v2/errors"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// comparableErrors reports whether == is defined for both operands.
//
// reflect.TypeOf(nil) returns a NIL reflect.Type, and calling Comparable on that panics rather than
// answering false. A nil operand is ordinary here -- a frame can carry no annotation of its own --
// so every comparability check has to screen for it first.
func comparableErrors(a, b error) bool {
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	return nil != ta && nil != tb && ta.Comparable() && tb.Comparable()
}

// As searches the error stack for an error that can be cast to the test
// argument, which must be a pointer. If it succeeds it performs the
// As finds the first error in err's chain that matches target, and if one is found, sets target to
// that error value and returns true. Otherwise it returns false.
//
// The chain consists of err itself followed by the errors obtained by repeatedly unwrapping it. An
// error matches target if the error's concrete type is assignable to the value pointed to by
// target, or if the error has a method As(interface{}) bool such that As(target) returns true.
//
// As panics if target is not a non-nil pointer to either a type that implements error, or to any
// interface type. That is deliberate and matches the standard library: an unusable target is a
// programming error at the call site, not a runtime condition to be reported, and returning false
// for it would silently hide the mistake in exactly the code paths that are hardest to observe.
//
// SIGNATURE CHANGE. This previously read As(err, test error) error, returning the matched error or
// nil. That forced the target to satisfy the error interface, so a caller could not ask for a bare
// interface type and had to choose between a value and a pointer receiver purely to make the call
// compile. It also meant this package's central helper did not match the standard library it aims
// to implement, so code moving between the two had to change shape.
func As(err error, target interface{}) bool {
	if nil == err {
		return false
	}
	if nil == target {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if reflect.Ptr != typ.Kind() || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	targetType := typ.Elem()
	if reflect.Interface != targetType.Kind() && !targetType.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}

	for {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		// `ok && x.As(target)`, not a bare return: a hook answering false means THIS link does not
		// match, not that the search is over. (*E) implements this hook, and it is what reaches the
		// annotation slot -- which is not on the Unwrap chain -- so Is and As agree about the chain.
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		// An error may wrap SEVERAL causes -- errors.Join, or fmt.Errorf with more than one %w --
		// and Unwrap() []error cannot be expressed by the single-error Unwrap below, so without this
		// branch the walk stops at the join and every cause underneath is unreachable.
		if multi, ok := err.(interface{ Unwrap() []error }); ok {
			for _, branch := range multi.Unwrap() {
				if nil != branch && As(branch, target) {
					return true
				}
			}
			return false
		}
		if err = Unwrap(err); nil == err {
			return false
		}
	}
}

// Caller returns the Caller associated with an error, if any.
func Caller(err error) std_caller.Caller {
	if e, ok := err.(std_error.Caller); ok {
		return e.Caller()
	}
	return nil
}

// Errorf formats according to a format specifier and returns an error that
// contains caller data.
func Errorf(msg string, data ...interface{}) *E {
	return New(fmt.Sprintf(msg, data...))
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

	if comparableErrors(err, test) && err == test {
		return true
	}

	// The annotation slot, which is not on the Unwrap chain -- see (*E).As. Note e.err can be nil
	// (Trace annotates nothing), which is why this goes through comparableErrors rather than
	// calling reflect.TypeOf(...).Comparable() directly: reflect.TypeOf(nil) is a NIL Type and
	// calling any method on it panics.
	if e, ok := err.(*E); ok {
		if comparableErrors(e.err, test) && e.err == test {
			return true
		}
	}

	// A custom Is answering "no" means THIS error does not match -- not that the search is over. The
	// previous code returned that answer directly, which ended the walk at the first link with an Is
	// method, and every *E has one. So a sentinel two links down was unreachable.
	if e, ok := err.(interface{ Is(error) bool }); ok && e.Is(test) {
		return true
	}

	// An error may wrap SEVERAL causes -- errors.Join, or fmt.Errorf with more than one %w. Those
	// implement Unwrap() []error, which the single-error Unwrap below cannot express: it returns nil
	// for them, so without this branch the walk stops dead at the join and every cause underneath is
	// unreachable. That made Is answer false for a joined sentinel even when the join was the error
	// passed in, unwrapped -- golang-jwt joins its sentinels, so no jwt error was ever matchable.
	//
	// The whole tree is searched, not just the first branch, which is what the standard library does.
	if multi, ok := err.(interface{ Unwrap() []error }); ok {
		for _, branch := range multi.Unwrap() {
			if Is(branch, test) {
				return true
			}
		}
		return false
	}

	if err = Unwrap(err); err == nil {
		return false
	}

	return Is(err, test)
}

// New returns an error that contains caller data.
//
// msg is a MESSAGE, not a format string: it is stored verbatim. Passing it through fmt.Errorf
// corrupted any message containing a percent sign -- New("100% complete") produced
// "100%!c(MISSING)omplete", and "disk usage at 95%" produced "...95%!(NOVERB)". Callers that want
// formatting have Errorf.
func New(msg string) *E {
	return &E{
		caller: NewCaller(),
		err:    std_errors.New(msg),
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
	if stdClr, ok := e.(std_error.Caller); ok {
		clr.trace = append(clr.trace, stdClr.Caller().Trace()...)
	}

	// prev, NOT err. Trace adds a caller line; it does not annotate. Holding the wrapped error in
	// the message slot left prev nil, so Unwrap returned nothing and the whole chain below the
	// trace became unreachable -- errors.Is could no longer find a sentinel through it. Error()
	// treats a frame with no message of its own as transparent, so the rendered text is unchanged.
	return &E{
		caller: clr,
		prev:   e,
	}
}

// Track updates the error stack with additional caller data.
//
// The tracked error stays ON THE UNWRAP CHAIN. It was previously held in the annotation slot with
// prev taken from a synthetic box, so for anything other than an *E -- a fmt.Errorf("%w") wrapper, a
// joined error, any foreign type that wraps -- prev was nil and everything the error itself wrapped
// became unreachable: errors.Is could no longer find a sentinel through it. A decorator that adds
// caller data must not cost the caller their chain.
//
// The inserted frame's message is the marker alone rather than the error's text repeated, since the
// error it wraps renders immediately after it.
func Track(e error) *E {
	if nil == e {
		return nil
	}

	// Adopt the error's own caller when it has one, so tracking does not overwrite the origin.
	clr := NewCaller()
	if stdClr, ok := e.(std_error.Caller); ok && nil != stdClr.Caller() {
		clr = stdClr.Caller()
	}

	return &E{
		caller: clr,
		prev: &E{
			caller: NewCaller(),
			err:    std_errors.New("(tracked)"),
			prev:   e,
		},
	}
}

// Unwrap returns the previous error.
//
// A multi-error -- one implementing Unwrap() []error -- has no single previous error, so this
// returns nil for it, exactly as the standard library's errors.Unwrap does. Callers that walk a
// chain must branch on Unwrap() []error themselves; Is and As do.
func Unwrap(err error) error {
	if e, ok := err.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}
	return nil
}

// Wrap returns a new error that wraps the provided error.
//
// msg is treated as a format string ONLY when data is supplied, matching New/Errorf: with no
// arguments there is nothing to interpolate, and interpreting it anyway corrupts any message
// containing a percent sign.
func Wrap(e error, msg string, data ...interface{}) *E {
	if 0 == len(data) {
		return WrapE(e, std_errors.New(msg))
	}
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
