package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"sync"
)

// stack represents an error stack.
type stack struct {
	stack []err
	mux   *sync.Mutex
}

// newStackFromError returns a new error stack.
func newEmptyStack() stack {
	return stack{
		stack: make([]err, 1),
		mux:   &sync.Mutex{},
	}
}

// newStack returns a new error stack.
func newStack(msg string, data ...interface{}) stack {
	return newStackFromErr(fmt.Errorf(msg, data...))
}

// newStackFromError returns a new error stack.
func newStackFromErr(e error) stack {
	s := newEmptyStack()
	s.stack[0] = newErr(e)
	return s
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
// Verbs:
//     %s      Returns the error string of the last error added
//     %v      Alias for %s
//
//  Flags:
//      #      JSON formatted output, useful for logging
//      -      Output caller details, useful for troubleshooting
//      +      Output full error stack details, useful for debugging
//      ' '    Add whitespace for readability, useful for development
//
// Examples:
//      %s:    An error occurred
//      %v:    An error occurred
//      %-v:   #0 stack_test.go:40 (github.com/bdlm/errors_test.TestErrors) - An error occurred
//      %+v:   #0 stack_test.go:40 (github.com/bdlm/errors_test.TestErrors) - An error occurred #1 stack_test.go:39 (github.com/bdlm/errors_test.TestErrors) - An error occurred
//      %#v:   {"error":"An error occurred"}
//      %#-v:  {"caller":"#0 stack_test.go:40 (github.com/bdlm/errors_test.TestErrors)","error":"An error occurred"}
//      %#+v:  [{"caller":"#0 stack_test.go:40 (github.com/bdlm/errors_test.TestErrors)","error":"An error occurred"},{"caller":"#0 stack_test.go:39 (github.com/bdlm/errors_test.TestErrors)","error":"An error occurred"}]
//      %# v:  {
//                 "error":"An error occurred"
//             }
//      %# -v: {
//                 "caller":"#0 stack_test.go:40 (github.com/bdlm/errors_test.TestErrors)",
//                 "error":"An error occurred"
//             }
//      %# +v: [
//                 {
//                     "caller":"#0 stack_test.go:40 (github.com/bdlm/errors_test.TestErrors)",
//                     "error":"An error occurred"
//                 },
//                 {
//                     "caller":"#0 stack_test.go:39 (github.com/bdlm/errors_test.TestErrors)",
//                     "error":"An error occurred"
//                 }
//             ]
func (e stack) Format(state fmt.State, verb rune) {
	str := bytes.NewBuffer([]byte{})

	switch verb {
	default:
		fmt.Fprintf(str, e.Error())

	case 'v':
		var (
			flagDetail bool
			flagFormat bool
			flagTrace  bool
			modeJSON   bool
		)

		if state.Flag('#') {
			modeJSON = true
		}
		if state.Flag(' ') {
			flagFormat = true
		}
		if state.Flag('-') {
			flagDetail = true
		}
		if state.Flag('+') {
			flagDetail = true
			flagTrace = true
		}

		jsonData := []map[string]string{}
		for a, err := range e.stack {
			if modeJSON {
				data := map[string]string{}
				if flagDetail {
					data["caller"] = fmt.Sprintf("#%d %s:%d (%s)",
						a,
						path.Base(err.Caller().File()),
						err.Caller().Line(),
						runtime.FuncForPC(err.Caller().Pc()).Name(),
					)
				}
				if "" != err.Error() {
					data["error"] = err.Error()
				}
				jsonData = append(jsonData, data)

			} else {
				if flagDetail {
					fmt.Fprintf(str, "#%d %s:%d (%s) ",
						a,
						path.Base(err.Caller().File()),
						err.Caller().Line(),
						runtime.FuncForPC(err.Caller().Pc()).Name(),
					)
					if "" != err.Error() {
						fmt.Fprintf(str, "- ")
					}
				}

				if "" != err.Error() {
					fmt.Fprintf(str, "%s ", err.Error())
				}

				if flagFormat {
					fmt.Fprintf(str, "\n")
				}

			}

			if !flagTrace {
				break
			}
		}

		if modeJSON {
			var data interface{}
			if !flagTrace {
				data = jsonData[0]
			} else {
				data = jsonData
			}

			var byts []byte
			if flagFormat {
				byts, _ = json.MarshalIndent(data, "", "    ")
			} else {
				byts, _ = json.Marshal(data)
			}

			str.Write(byts)
		}
	}

	fmt.Fprintf(state, "%s", str.String())
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
