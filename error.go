package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

// Caller holds runtime.Caller data.
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

// ex is a thing
type ex struct {
	caller Caller
	err    error
	prev   error
}

// Caller does things.
func (e ex) Caller() Caller {
	return e.caller
}

// Error implements error.
func (e ex) Error() string {
	return e.err.Error()
}

// Has implements error.
func (e ex) Has(test error) bool {
	return Has(e, test)
}

// Is implements error.
func (e ex) Is(test error) bool {
	return Is(e, test)
}

// Error implements error.
func (e ex) Unwrap() error {
	return Unwrap(e)
}

func list(e error) []error {
	ret := []error{}
	if tmp, ok := e.(ex); ok {
		ret = append(ret, e)
		ret = append(ret, list(tmp.prev)...)
	}

	return ret
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
func (e ex) Format(state fmt.State, verb rune) {
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
			flagTrace = true
		}

		jsonData := []map[string]interface{}{}

		for a, b := range list(e) {
			err, ok := b.(ex)
			if !ok {
				break
			}

			if modeJSON {
				data := map[string]interface{}{}
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
				if flagTrace {
					trace := []string{}
					for b, caller := range err.Caller().Trace() {
						trace = append(trace, fmt.Sprintf("#%d %s:%d (%s)",
							b,
							path.Base(caller.File()),
							caller.Line(),
							runtime.FuncForPC(caller.Pc()).Name(),
						))
					}
					data["trace"] = trace
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

				if flagTrace {
					start := ""
					end := ""
					if flagFormat {
						start = "\t"
						end = "\n"
					}
					for b, caller := range err.Caller().Trace() {
						fmt.Fprintf(str, "%s#%d %s:%d (%s)%s",
							start,
							b,
							path.Base(caller.File()),
							caller.Line(),
							runtime.FuncForPC(caller.Pc()).Name(),
							end,
						)
					}
				}
			}
		}

		if modeJSON {
			var byts []byte
			if flagFormat {
				byts, _ = json.MarshalIndent(jsonData, "", "    ")
			} else {
				byts, _ = json.Marshal(jsonData)
			}

			str.Write(byts)
		}
	}

	fmt.Fprintf(state, "%s", str.String())
}
