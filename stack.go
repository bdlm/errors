package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
)

// stack represents an error stack.
type stack struct {
	stack []err
	mux   *sync.Mutex
}

// Caller returns the most recent error caller.
func (e stack) Caller() Caller {
	return e.last().Caller()
}

// Cause returns the root cause of an error stack.
func (e stack) Cause() Error {
	return e.stack[0]
}

// Error returns the most recent error message.
func (e stack) Error() string {
	return fmt.Sprintf("%v", e)
}

// Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing
//
// Format formats the stack trace output. Several verbs are supported:
//  %s  - Returns the error string of the last error added
//
//  %v  - Alias for %s
//
//  %-v - Returns the full call stack trace in a single line, useful for
//        logging. Same as %#v with the newlines escaped.
//
//  %+v - Returns a multi-line call stack trace including the full trace of
//        each addition to the call stack. Useful for development.
//
//  %#v - Returns a full call stack trace as a JSON object, useful for
//        logging.
func (e stack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})

		//e.mux.Lock()
		//defer e.mux.Unlock()

		// JSON format
		if state.Flag('#') {
			var byts []byte
			if state.Flag('+') {
				byts, _ = json.MarshalIndent(e, "", "    ")
			} else {
				byts, _ = json.Marshal(e)
			}
			fmt.Fprintf(str, string(byts))

		} else {
			for a, err := range e.stack { // a := len(e.stack) - 1; a >= 0; a--

				switch {
				// Extended stack trace
				case state.Flag('+'):
					if "" != err.Error() {
						fmt.Fprintf(str, "#%d %s:%d (%s) - %s\n",
							a,
							path.Base(err.Caller().File()),
							err.Caller().Line(),
							runtime.FuncForPC(err.Caller().(caller).pc).Name(),
							err.Error(),
						)
					} else {
						fmt.Fprintf(str, "#%d %s:%d (%s) \n",
							a,
							path.Base(err.Caller().File()),
							err.Caller().Line(),
							runtime.FuncForPC(err.Caller().(caller).pc).Name(),
						)
					}

				// Inline stack trace
				case state.Flag('-'):
					if nil != err.e {
						fmt.Fprintf(str, "#%d %s - %s:%d (%s) ",
							a,
							err.Error(),
							path.Base(err.Caller().File()),
							err.Caller().Line(),
							err.Caller().Func(),
						)
					} else {
						fmt.Fprintf(str, "#%d - %s:%d (%s) ",
							a,
							path.Base(err.Caller().File()),
							err.Caller().Line(),
							err.Caller().Func(),
						)
					}

				// Default output
				default:
					fmt.Fprintf(state, err.Error())
					return
				}
			}
		}
		fmt.Fprintf(state, "%s", strings.Trim(str.String(), " \n\t"))
	default:
		// Default output
		fmt.Fprintf(state, e.Error())
	}
}

// MarshalJSON implements the json.Marshaller interface.
func (e stack) MarshalJSON() ([]byte, error) {
	type data struct {
		Err    string `json:"error,omitempty"`
		Caller string `json:"caller,omitempty"`
	}

	jsonData := []data{}

	e.mux.Lock()
	if len(e.stack) > 1 {
		for _, err := range e.stack {
			jsonData = append(jsonData, data{
				Err: err.Error(),
				Caller: fmt.Sprintf("%s:%d (%s)",
					path.Base(err.Caller().File()),
					err.Caller().Line(),
					err.Caller().Func(),
				),
			})
		}
	}
	e.mux.Unlock()

	return json.Marshal(jsonData)
}

// String implements the stringer interface.
func (e stack) String() string {
	return e.Error()
}

// Trace returns the call stack.
func (e stack) Trace() []Caller {
	var callers []Caller

	for _, caller := range e.stack {
		callers = append(callers, caller.Caller())
	}

	return callers
}

// append appends an error to the stack.
func (e stack) append(errors ...err) stack {
	ret := newEmptyStack()
	ret.stack = make([]err, e.len()+len(errors))

	var (
		a int
		b err
	)

	e.mux.Lock()
	for a, b = range e.stack {
		ret.stack[a] = b
	}
	e.mux.Unlock()

	for _, b = range errors {
		a++
		ret.stack[a] = b
	}

	return ret
}

// clone returns a clone of the stack.
func (e stack) clone(errors ...err) stack {
	ret := newEmptyStack()
	ret.stack = make([]err, len(e.stack))

	e.mux.Lock()
	for a, b := range e.stack {
		ret.stack[a] = b
	}
	e.mux.Unlock()

	return ret
}

// first returns the first error added to the stack.
func (e stack) first() err {
	var ret err

	e.mux.Lock()
	ret = e.stack[len(e.stack)-1]
	e.mux.Unlock()

	return ret
}

// last returns the last error added to the stack.
func (e stack) last() err {
	var ret err

	e.mux.Lock()
	ret = e.stack[0]
	e.mux.Unlock()

	return ret
}

// len returns the current length of the stack.
func (e stack) len() int {
	var ret int

	e.mux.Lock()
	ret = len(e.stack)
	e.mux.Unlock()

	return ret
}

// prepend prepends an error to the stack.
func (e stack) prepend(errors ...err) stack {
	ret := newEmptyStack()
	ret.stack = make([]err, e.len()+len(errors))

	var (
		a int
		b err
	)

	for a, b = range errors {
		ret.stack[a] = b
	}

	e.mux.Lock()
	for _, b = range e.stack {
		a++
		ret.stack[a] = b
	}
	e.mux.Unlock()

	return ret
}
