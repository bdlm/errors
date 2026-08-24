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
//
//	%s      Returns the error string of the last error added
//	%v      Alias for %s
//
// Flags:
//
//	#      JSON formatted output, useful for logging
//	-      Output caller details, useful for troubleshooting
//	+      Output full error stack details, useful for debugging
//	' '    (space) Add whitespace formatting for readability, useful for development
//
// Examples:
//
//	%s:    An error occurred
//	%v:    An error occurred
//	%-v:   #0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors) - An error occurred
//	%+v:   #0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors) - An error occurred #1 stack_test.go:39 (github.com/bdlm/error_test.TestErrors) - An error occurred
//	%#v:   {"error":"An error occurred"}
//	%#-v:  {"caller":"#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"}
//	%#+v:  [{"caller":"#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"},{"caller":"#0 stack_test.go:39 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"}]
func (e *E) Format(state fmt.State, verb rune) {
	str := bytes.NewBuffer([]byte{})

	switch verb {
	default:
		fmt.Fprint(str, e.Error())

	case 'v':
		var (
			flagDetail bool
			flagFormat bool
			flagTrace  bool
			modeJSON   bool
			lastE      error
			key        int
			nextE      error
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

		for key, nextE = range list(e) {
			sp, jsonData, str = format(key, nextE, sp, jsonData, str, flagDetail, flagFormat, flagTrace, modeJSON)
			if !flagTrace {
				break
			}

			if !flagDetail &&
				!flagFormat &&
				!flagTrace &&
				!modeJSON {
				break
			}

			if err, ok := nextE.(*E); ok {
				lastE = err.prev
			}
		}
		if nil != lastE {
			_, jsonData, str = format(key+1, lastE, sp, jsonData, str, flagDetail, flagFormat, flagTrace, modeJSON)
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

func format(key int, nextE error, sp string, jsonData []map[string]interface{}, str *bytes.Buffer, flagDetail bool, flagFormat bool, flagTrace bool, modeJSON bool) (string, []map[string]interface{}, *bytes.Buffer) {
	err, ok := nextE.(*E)

	if modeJSON {
		data := map[string]interface{}{}
		if flagDetail || flagTrace {
			if ok && nil != err.Caller() {
				data["caller"] = fmt.Sprintf("#%d %s:%d (%s)",
					key,
					path.Base(err.Caller().File()),
					err.Caller().Line(),
					runtime.FuncForPC(err.Caller().Pc()).Name(),
				)
			} else {
				data["caller"] = fmt.Sprintf("#%d n/a",
					key,
				)
			}
		}
		if message := messageFor(nextE, flagTrace); "" != message {
			data["error"] = message
		}
		jsonData = append(jsonData, data)

	} else {
		message := messageFor(nextE, flagTrace)
		if "" != message {
			fmt.Fprintf(str, "%s%s", sp, message)
		}

		if flagDetail || flagTrace {
			if "" != message {
				fmt.Fprintf(str, " - ")
			}
			if ok && nil != err.Caller() {
				fmt.Fprintf(str, "#%d %s:%d (%s);",
					key,
					path.Base(err.Caller().File()),
					err.Caller().Line(),
					runtime.FuncForPC(err.Caller().Pc()).Name(),
				)
			} else {
				fmt.Fprintf(str, "#%d n/a",
					key,
				)
			}
		}

		if flagFormat {
			str = bytes.NewBuffer([]byte(strings.Trim(str.String(), " ")))
			fmt.Fprintf(str, "\n")
		} else if flagTrace {
			sp = " "
		}
	}
	return sp, jsonData, str
}

// messageFor chooses between the whole chain and this frame alone.
//
// The trace verbs (%+v) print one line per frame, so each line must carry only its own message or
// every line repeats the tail below it. Every OTHER form renders a single string, and that string
// has to be the full message -- the same thing Error() and %s produce.
//
// %v used to take the per-frame path and stop at the first frame, so it printed only the outermost
// message and silently discarded the cause. It is the most-used verb for an error in Go and the one
// most loggers call, so the effect was that a wrapped cause was recorded everywhere except where
// anyone looked: err.Error() had it, %s had it, %v did not.
func messageFor(err error, perFrame bool) string {
	if perFrame {
		return frameMessage(err)
	}
	if nil == err {
		return ""
	}
	return err.Error()
}

// frameMessage is one link's own message. For an *E that is its frame message without the wrapped
// chain -- Error() now includes the cause, and a trace that prints one line per frame would otherwise
// repeat the whole tail on every line. Any other error type has only its own message to give.
func frameMessage(err error) string {
	if e, ok := err.(*E); ok {
		return e.message()
	}
	if nil == err {
		return ""
	}
	return err.Error()
}
