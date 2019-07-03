package errors

// Err defines the error interface.
type Err interface {
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
	Unwrap() Err
}

// ex is a thing
type ex struct {
	caller Caller
	err    error
	prev   error
}

// Caller implements Error.
func (e ex) Caller() Caller {
	return e.caller
}

// Error implements Error.
func (e ex) Error() string {
	return e.err.Error()
}

// Has implements Error.
func (e ex) Has(test error) bool {
	return Has(e, test)
}

// Is implements Error.
func (e ex) Is(test error) bool {
	return Is(e, test)
}

// Error implements Error.
func (e ex) Unwrap() Err {
	return Unwrap(e)
}

// list will convert the error stack into a simple array.
func list(e error) []error {
	ret := []error{}

	if nil != e {
		if tmp, ok := e.(ex); ok {
			ret = append(ret, e)
			ret = append(ret, list(tmp.prev)...)
		}
	}

	return ret
}
