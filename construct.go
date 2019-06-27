package errors

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// newCaller returns a new caller instance containing data for the current
// call stack.
func newCaller() caller {
	trace := []Caller{}
	clr := caller{}
	a := 0
	for {
		traceCaller := caller{}
		if traceCaller.pc, traceCaller.file, traceCaller.line, traceCaller.ok = runtime.Caller(a); traceCaller.ok {
			trace = append(trace, traceCaller)
			if !clr.ok &&
				(!strings.Contains(strings.ToLower(traceCaller.file), "github.com/bdlm/errors") ||
					strings.HasSuffix(strings.ToLower(traceCaller.file), "_test.go")) {
				clr.pc = traceCaller.pc
				clr.file = traceCaller.file
				clr.line = traceCaller.line
				clr.ok = traceCaller.ok
			}
		} else {
			break
		}
		a++
	}
	clr.trace = trace
	return clr
}

// newErr returns a new err instance.
func newErr(e error) err {
	return err{
		e:      e,
		caller: newCaller(),
	}
}

// newStack returns a new error stack.
func newStack(msg string, data ...interface{}) stack {
	return newStackFromErr(fmt.Errorf(msg, data...))
}

// newStackFromError returns a new error stack.
func newStackFromErr(e error) stack {
	return stack{
		stack: []err{newErr(e)},
		mux:   &sync.Mutex{},
	}
}
