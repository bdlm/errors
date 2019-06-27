package errors

// Caller holds runtime.Caller data.
type Caller interface {
	// File returns the file in which the call occurred.
	File() string

	// Func returns the name of the function in which the call occurred.
	Func() string

	// Line returns the line number in the file in which the call occurred.
	Line() int
}

// Error defines the error interface.
type Error interface {
	// Caller returns the most recent error caller.
	Caller() Caller

	// Cause returns the root cause of an error stack.
	Cause() Error

	// Error returns the most recent error message.
	Error() string

	// String implements the stringer interface.
	String() string

	// Trace returns the call stack.
	Trace() []Caller
}
