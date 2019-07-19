package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strings"
)

// Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing
//
// Verbs:
//     %s      Returns the error string of the last error added
//     %v      Alias for %s
//
// Flags:
//      #      JSON formatted output, useful for logging
//      -      Output caller details, useful for troubleshooting
//      +      Output full error stack details, useful for debugging
//      ' '    (space) Add whitespace formatting for readability, useful for development
//
// Examples:
//      %s:    An error occurred
//      %v:    An error occurred
//      %-v:   #0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors) - An error occurred
//      %+v:   #0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors) - An error occurred #1 stack_test.go:39 (github.com/bdlm/error_test.TestErrors) - An error occurred
//      %#v:   {"error":"An error occurred"}
//      %#-v:  {"caller":"#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"}
//      %#+v:  [{"caller":"#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"},{"caller":"#0 stack_test.go:39 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"}]
func (e E) Format(state fmt.State, verb rune) {
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
		sp := ""

		for a, b := range list(e) {
			err, ok := b.(E)
			if !ok {
				break
			}

			if modeJSON {
				data := map[string]interface{}{}
				if flagDetail || flagTrace {
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
				if "" != err.Error() {
					fmt.Fprintf(str, "%s%s", sp, err.Error())
				}

				if flagDetail || flagTrace {
					if "" != err.Error() {
						fmt.Fprintf(str, " - ")
					}
					fmt.Fprintf(str, "#%d %s:%d (%s);",
						a,
						path.Base(err.Caller().File()),
						err.Caller().Line(),
						runtime.FuncForPC(err.Caller().Pc()).Name(),
					)
				}

				if flagFormat {
					str = bytes.NewBuffer([]byte(strings.Trim(str.String(), " ")))
					fmt.Fprintf(str, "\n")
				} else if flagTrace {
					sp = " "
				}
			}

			if !flagTrace {
				break
			}

			if !flagDetail &&
				!flagFormat &&
				!flagTrace &&
				!modeJSON {
				break
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

	fmt.Fprintf(state, "%s", strings.Trim(str.String(), "\r\n\t"))
}
