package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Caller holds runtime.Caller data for an error.
type Caller interface {
	// File returns the file in which the call occurred.
	File() string

	// Func returns the name of the function in which the call occurred.
	Func() string

	// Line returns the line number in the file in which the call occurred.
	Line() int

	// Pc returns the program counter.
	Pc() uintptr

	// Trace returns the call stack.
	Trace() []Caller
}

// caller holds runtime.Caller data.
type caller struct {
	file  string
	line  int
	ok    bool
	pc    uintptr
	trace []Caller
}

// NewCaller returns a new caller instance containing data for the current
// call stack.
func NewCaller() Caller {
	trace := []Caller{}
	clr := caller{}
	a := 0
	for {
		traceCaller := caller{}
		if traceCaller.pc, traceCaller.file, traceCaller.line, traceCaller.ok = runtime.Caller(a); traceCaller.ok {
			if !strings.Contains(strings.ToLower(traceCaller.file), "github.com/bdlm/errors") ||
				strings.HasSuffix(strings.ToLower(traceCaller.file), "_test.go") {
				trace = append(trace, traceCaller)
				if !clr.ok {
					clr.pc = traceCaller.pc
					clr.file = traceCaller.file
					clr.line = traceCaller.line
					clr.ok = traceCaller.ok
				}
			}
		} else {
			break
		}
		a++
	}
	clr.trace = trace
	return clr
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

// Pc implements Caller.
func (caller caller) Pc() uintptr {
	return caller.pc
}

// String implements fmt.Stringer.
func (caller caller) String() string {
	return fmt.Sprintf(
		"%s:%d",
		runtime.FuncForPC(caller.pc).Name(),
		caller.line,
	)
}

// Trace implements Caller.
func (caller caller) Trace() []Caller {
	return caller.trace
}
