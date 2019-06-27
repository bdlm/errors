package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
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
			byts, _ := json.Marshal(e)
			fmt.Fprintf(str, string(byts))

		} else {
			for a := len(e.stack) - 1; a >= 0; a-- {
				err := e.stack[a]

				switch {
				// Extended stack trace
				//case state.Flag('+'):
				//	if nil != err.err {
				//		fmt.Fprintf(str, "%s - %s:%d (%s) \n",
				//			//n,
				//			err.Error(),
				//			path.Base(err.Caller().File),
				//			err.Caller().Line,
				//			runtime.FuncForPC(err.Caller().pc).Name(),
				//		)
				//	} else {
				//		fmt.Fprintf(str, "%s:%d (%s) \n",
				//			//a,
				//			path.Base(err.Caller().File),
				//			err.Caller().Line,
				//			runtime.FuncForPC(err.Caller().pc).Name(),
				//		)
				//	}
				//	//fmt.Fprintf(str, "#%d: `%s`\n", a, runtime.FuncForPC(err.Caller().pc).Name())
				//	//fmt.Fprintf(str, "\terror:   %s\n", err.Error())
				//	//fmt.Fprintf(str, "\tline:    %s:%d\n", path.Base(err.Caller().File), err.Caller().Line)

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
	jsonData := []struct {
		Err    string `json:"error,omitempty"`
		Caller string `json:"caller,omitempty"`
	}{}
	e.mux.Lock()
	if len(e.stack) > 1 {
		for a := len(e.stack) - 1; a >= 0; a-- {
			err := e.stack[a]
			jsonData = append(jsonData, struct {
				Err    string `json:"error,omitempty"`
				Caller string `json:"caller,omitempty"`
			}{
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
	e.mux.Lock()
	e.stack = append(e.stack, errors...)
	e.mux.Unlock()
	return e
}

// last returns the last error appended to the stack.
func (e stack) last() err {
	var err err
	e.mux.Lock()
	err = e.stack[len(e.stack)-1]
	e.mux.Unlock()
	return err
}
