package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Caller holds runtime.Caller data.
type Caller struct {
	ok    bool
	pc    uintptr
	trace []Caller

	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

// String implements the Stringer interface
func (caller Caller) String() string {
	return fmt.Sprintf(
		"%s:%d %s",
		caller.File,
		caller.Line,
		runtime.FuncForPC(caller.pc).Name(),
	)
}

// getCaller returns the caller and backtrace of this error.
func getCaller() Caller {
	caller := Caller{}
	trace := []Caller{}
	a := 0
	for {
		traceCaller := Caller{}
		if traceCaller.pc, traceCaller.File, traceCaller.Line, traceCaller.ok = runtime.Caller(a); traceCaller.ok {
			trace = append(trace, traceCaller)
			if !caller.ok &&
				(!strings.Contains(strings.ToLower(traceCaller.File), "github.com/bdlm/errors") ||
					strings.HasSuffix(strings.ToLower(traceCaller.File), "_test.go")) {
				caller.pc = traceCaller.pc
				caller.File = traceCaller.File
				caller.Line = traceCaller.Line
				caller.ok = traceCaller.ok
			}
		} else {
			break
		}
		a++
	}
	caller.trace = trace
	return caller
}
