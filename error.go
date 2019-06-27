package errors

// err defines a single error stack entry.
type err struct {
	e      error
	caller caller
}

// Caller implements Error.
func (err err) Caller() Caller {
	return err.caller
}

// Cause implements Error.
func (err err) Cause() Error {
	return err
}

// Error implements error.
func (err err) Error() string {
	return err.String()
}

// String implements Stringer.
func (err err) String() string {
	if nil == err.e {
		return ""
	}
	return err.e.Error()
}

// Trace implements Error.
func (err err) Trace() []Caller {
	return err.Caller().(caller).trace
}
