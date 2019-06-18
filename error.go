package errors

// Error defines a single error stack entry.
type Error struct {
	err    error
	caller Caller
}

func newError(e error) Error {
	return Error{
		err:    e,
		caller: getCaller(),
	}
}

// Caller implements Err.
func (err Error) Caller() Caller {
	return err.caller
}

// Error implements error.
func (err Error) Error() string {
	return err.String()
}

// String implements Stringer.
func (err Error) String() string {
	return err.err.Error()
}

// Trace implements Err.
func (err Error) Trace() []Caller {
	return err.Caller().Trace()
}
