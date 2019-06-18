package errors

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
)

// newStack returns a new error stack.
func newStack(msg string, data ...interface{}) Stack {
	return Stack{
		stack: []Error{
			Error{
				err:    fmt.Errorf(msg, data...),
				caller: getCaller(),
			},
		},
		mux: &sync.Mutex{},
	}
}

// newStackFromError returns a new error stack.
func newStackFromErr(err error) Stack {
	return Stack{
		stack: []Error{
			Error{
				err:    err,
				caller: getCaller(),
			},
		},
		mux: &sync.Mutex{},
	}
}

type Stack struct {
	stack []Error
	mux   *sync.Mutex
}

// last returns the last error appended to the stack.
func (err Stack) last() Error {
	var e Error
	err.mux.Lock()
	e = err.stack[len(err.stack)-1]
	err.mux.Unlock()
	return e
}

// append appends an Error to the stack.
func (err Stack) append(e ...Error) Stack {
	err.mux.Lock()
	err.stack = append(err.stack, e...)
	err.mux.Unlock()
	return err
}

// Caller returns the most recent error caller.
func (err Stack) Caller() Caller {
	return err.last().Caller()
}

// Cause returns the root cause of an error stack.
func (err Stack) Cause() Error {
	var e Error
	err.mux.Lock()
	e = err.stack[0]
	err.mux.Unlock()
	return e
}

// Error returns the most recent error message.
func (err Stack) Error() string {
	return fmt.Sprintf("%v", err)
}

// Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing
//
// Format formats the stack trace output. Several verbs are supported:
//  %s  - Returns the error string of the last error added
//
//  %v  - Alias for %s
//
//  %#v - Returns the full call stack trace in a single line, useful for
//        logging. Same as %#v with the newlines escaped.
//
//  %-v - Returns a multi-line call stack trace, one column-delimited line
//        per error. Useful for development.
//
//  %+v - Returns a multi-line call stack trace including the full trace of
//        each addition to the call stack. Useful for development.
func (err Stack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})

		err.mux.Lock()
		defer err.mux.Unlock()

		for a := len(err.stack) - 1; a >= 0; a-- {
			e := err.stack[a]

			switch {
			case state.Flag('+'):
				// Extended stack trace
				fmt.Fprintf(str, "#%d: `%s`\n", a, runtime.FuncForPC(e.Caller().Pc()).Name())
				fmt.Fprintf(str, "\terror:   %s\n", e.Error())
				fmt.Fprintf(str, "\tline:    %s:%d\n", path.Base(e.Caller().File()), e.Caller().Line())

			case state.Flag('#'):
				// Condensed stack trace
				fmt.Fprintf(str, "#%d - \"%s\" %s:%d (%s)\n",
					a,
					e.Error(),
					path.Base(e.Caller().File()),
					e.Caller().Line(),
					runtime.FuncForPC(e.Caller().Pc()).Name(),
				)

			case state.Flag('-'):
				// Inline stack trace
				fmt.Fprintf(str, "#%d - \"%s\" %s:%d (%s) ",
					a,
					e.Error(),
					path.Base(e.Caller().File()),
					e.Caller().Line(),
					runtime.FuncForPC(e.Caller().Pc()).Name(),
				)

			default:
				// Default output
				fmt.Fprintf(state, e.Error())
				return
			}
		}
		fmt.Fprintf(state, "%s", strings.Trim(str.String(), " \n\t"))
	default:
		// Externally-safe error message
		fmt.Fprintf(state, err.Error())
	}
}

// String implements the stringer interface.
func (err Stack) String() string {
	return err.Error()
}

// Trace returns the call stack.
func (err Stack) Trace() []Caller {
	var callers []Caller
	for _, caller := range err.stack {
		callers = append(callers, caller.Caller())
	}
	return callers
}
