package errors

import (
	"fmt"
	"runtime"
	"strings"

	std "github.com/bdlm/std/error"
)

/*
Call implements bdlm/std/error.Caller, holding runtime.Caller data.
*/
type Call struct {
	loaded bool
	file   string
	line   int
	ok     bool
	pc     uintptr
}

/*
File implements bdlm/std/error.Caller, returning the caller file name.
*/
func (call Call) File() string {
	return call.file
}

/*
Line implements bdlm/std/error.Caller, returning the caller line number.
*/
func (call Call) Line() int {
	return call.line
}

/*
Ok implements bdlm/std/error.Caller, returning whether the caller data was
successfully recovered.
*/
func (call Call) Ok() bool {
	return call.ok
}

/*
Pc implements bdlm/std/error.Caller, returning the caller program counter.
*/
func (call Call) Pc() uintptr {
	return call.pc
}

/*
String implements the Stringer interface
*/
func (call Call) String() string {
	return fmt.Sprintf(
		"%s:%d %s",
		call.file,
		call.line,
		runtime.FuncForPC(call.pc).Name(),
	)
}

func getCaller() std.Caller {
	var caller Call
	a := 0
	for {
		if caller.pc, caller.file, caller.line, caller.ok = runtime.Caller(a); caller.ok {
			if !strings.Contains(strings.ToLower(caller.file), "github.com/bdlm/errors") ||
				strings.HasSuffix(strings.ToLower(caller.file), "_test.go") {
				break
			}
		} else {
			break
		}
		a++
	}
	return caller
}

func getTrace() std.Trace {
	var trace std.Trace
	var caller Call
	a := 0
	for {
		if caller.pc, caller.file, caller.line, caller.ok = runtime.Caller(a); caller.ok {
			trace = append(trace, caller)
		} else {
			break
		}
		a++
	}
	return trace
}
