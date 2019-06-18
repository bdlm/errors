package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Caller holds runtime.Caller data.
type Caller struct {
	loaded bool
	file   string
	line   int
	ok     bool
	pc     uintptr
	trace  []Caller
}

// File returns the caller file name.
func (caller Caller) File() string {
	return caller.file
}

// Line returns the caller line number.
func (caller Caller) Line() int {
	return caller.line
}

// Ok returns whether the caller data was successfully recovered.
func (caller Caller) Ok() bool {
	return caller.ok
}

// Pc returns the caller program counter.
func (caller Caller) Pc() uintptr {
	return caller.pc
}

// String implements the Stringer interface
func (caller Caller) String() string {
	return fmt.Sprintf(
		"%s:%d %s",
		caller.file,
		caller.line,
		runtime.FuncForPC(caller.pc).Name(),
	)
}

// Trace returns the call stack leading to this caller.
func (caller Caller) Trace() []Caller {
	return caller.trace
}

func getCaller() Caller {
	var traceCaller Caller
	var caller Caller
	trace := []Caller{}
	a := 0
	for {
		if traceCaller.pc, traceCaller.file, traceCaller.line, traceCaller.ok = runtime.Caller(a); traceCaller.ok {
			trace = append(trace, traceCaller)
			if !strings.Contains(strings.ToLower(caller.file), "github.com/bdlm/errors") ||
				strings.HasSuffix(strings.ToLower(caller.file), "_test.go") {
				caller = traceCaller
			}
		} else {
			break
		}
		a++
	}
	caller.trace = trace
	return caller
}
