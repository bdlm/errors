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
// assignment and returns the result, otherwise it returns nil.
func As(err, test error) error {
	if nil == err || nil == test {
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
		// An *E's annotation (e.err) is not on the Unwrap chain -- Unwrap yields e.prev -- so it
		// has to be searched explicitly, exactly as Is does. Without this the package's own As
		// disagreed with its Is about what the chain contains.
		if e, ok := err.(*E); ok && nil != e.err {
			if found := As(e.err, test); nil != found {
				return found
			}
		}
		// Same multi-error branch as Is, for the same reason: Unwrap cannot express Unwrap() []error,
		// so the walk would terminate at a join and miss every cause below it.
		if multi, ok := err.(interface{ Unwrap() []error }); ok {
			for _, branch := range multi.Unwrap() {
				if found := As(branch, test); nil != found {
					return found
				}
			}
			return nil
		}
		err = Unwrap(err)
	}

	return nil
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
