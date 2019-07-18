package errors

// Error defines the error interface.
type Error interface {
	// Caller returns the associated Caller instance.
	Caller() Caller

	// Error implements error.
	Error() string

	// Has tests to see if the test error exists anywhere in the error
	// stack.
	Has(test error) bool

	// Is tests to see if the test error matches most recent error in the
	// stack.
	Is(test error) bool

	// Unwrap returns the next error, if any.
	Unwrap() Error
}

// E is a thing
type E struct {
	caller Caller
	err    error
	prev   error
}

// Caller implements Error.
func (e E) Caller() Caller {
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
func (e E) Unwrap() Error {
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
