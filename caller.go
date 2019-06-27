package errors

import (
	"fmt"
	"runtime"
)

// caller holds runtime.Caller data.
type caller struct {
	file  string
	line  int
	ok    bool
	pc    uintptr
	trace []Caller
}

// File implements Caller.
func (caller caller) File() string {
	return caller.file
}

// Func implements Caller.
func (caller caller) Func() string {
	return runtime.FuncForPC(caller.pc).Name()
}

// Line implements Caller.
func (caller caller) Line() int {
	return caller.line
}

// String implements Stringer.
func (caller caller) String() string {
	return fmt.Sprintf(
		"%s:%d %s",
		caller.file,
		caller.line,
		runtime.FuncForPC(caller.pc).Name(),
	)
}
