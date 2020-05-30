package errors

import (
	"fmt"
	"runtime"
	"strings"

	std_caller "github.com/bdlm/std/v2/caller"
)

// caller is a github.com/bdlm/std.Caller interface implementation and holds
// runtime.Caller data.
type caller struct {
	file  string
	line  int
	ok    bool
	pc    uintptr
	trace std_caller.Trace
}

// NewCaller returns a new Caller containing data for the current call stack.
func NewCaller() std_caller.Caller {
	trace := std_caller.Trace{}
	clr := &caller{}
	a := 0
	for {
		traceCaller := &caller{}
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
func (caller *caller) File() string {
	return caller.file
}

// Func implements Caller.
func (caller *caller) Func() string {
	return runtime.FuncForPC(caller.pc).Name()
}

// Line implements Caller.
func (caller *caller) Line() int {
	return caller.line
}

// Pc implements Caller.
func (caller *caller) Pc() uintptr {
	return caller.pc
}

// String implements fmt.Stringer.
func (caller *caller) String() string {
	return fmt.Sprintf(
		"%s:%d",
		runtime.FuncForPC(caller.pc).Name(),
		caller.line,
	)
}

// Trace implements Caller.
func (caller *caller) Trace() std_caller.Trace {
	return caller.trace
}
