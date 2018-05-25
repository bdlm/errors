package errors

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"text/tabwriter"
)

/*
Err defines an error heap
*/
type Err []Msg

/*
New returns an error with caller information for debugging. `code` is
optional. Although you can pass multiple codes, only the first is
accepted.
*/
func New(msg string, code ...Code) Err {
	var errCode Code
	if len(code) > 0 {
		errCode = code[0]
	}
	return Err{Msg{
		err:    errors.New(msg),
		caller: getCaller(),
		code:   errCode,
		msg:    msg,
	}}
}

/*
Callers returns the caller stack.
*/
func (errs Err) Callers() []Caller {
	callers := []Caller{}
	for _, msg := range errs {
		if msg.caller.Ok {
			callers = append(callers, msg.caller)
		}
	}
	return callers
}

/*
Cause returns the root cause of an error stack.
*/
func (errs Err) Cause() Msg {
	if len(errs) > 0 {
		return errs[0]
	}
	return Msg{}
}

/*
Code returns the most recent error code.
*/
func (errs Err) Code() Code {
	code := ErrUnspecified
	if len(errs) > 0 {
		code = errs[len(errs)-1].code
	}
	return code
}

/*
Error implements the error interface.
*/
func (errs Err) Error() string {
	meta, ok := Codes[errs.Code()]
	if !ok {
		meta = Codes[ErrUnspecified]
	}
	return meta.External
}

/*
Format implements fmt.Formatter.

https://golang.org/pkg/fmt/#hdr-Printing

Format formats the stack trace output. Several verbs are supported:
	%s  - Returns the user-safe error string mapped to the error code or
		  "Internal Server Error" if none is specified.

	%v  - Returns the full stack trace in a single line, useful for
	      logging. Same as %#v with the newlines escaped.

	%#v - Returns a multi-line stack trace, one column-delimited line
		  per error.

	%+v - Returns a multi-line detailed stack trace with multiple lines
	      per error. Only useful for human consumption.
*/
func (errs Err) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		stackPos := 0
		for k := len(errs) - 1; k >= 0; k-- {
			err := errs[k]
			msg, ok := Codes[err.code]
			if !ok {
				msg = Codes[ErrUnspecified]
			}
			switch {
			case s.Flag('+'):
				w := tabwriter.NewWriter(s, 0, 0, 4, ' ', 0)
				fmt.Fprintf(w, "%2d: %s\n", stackPos, runtime.FuncForPC(err.caller.Pc).Name())
				fmt.Fprintf(w, "\t\tline: %s: %d\n", path.Base(err.caller.File), err.caller.Line)
				fmt.Fprintf(w, "\t\tcode: %d: %s\n", err.code, msg.Internal)
				fmt.Fprintf(w, "\t\tmesg: %s\n\n", err.msg)
				w.Flush()

			case s.Flag('#'):
				// Condensed stack trace
				w := tabwriter.NewWriter(s, 5, 1, 4, ' ', 0)
				fmt.Fprintf(w, "%2d - %s:%d\t%s\t%d:%s\t%s\t\n",
					stackPos,
					path.Base(err.caller.File),
					err.caller.Line,
					runtime.FuncForPC(err.caller.Pc).Name(),
					err.code,
					msg.Internal,
					err.msg,
				)
				w.Flush()

			default:
				// Condensed stack trace
				w := tabwriter.NewWriter(s, 5, 1, 4, ' ', 0)
				str := fmt.Sprintf(
					"%2d - %s:%d\t%s\t%d:%s\t%s\t",
					stackPos,
					path.Base(err.caller.File),
					err.caller.Line,
					runtime.FuncForPC(err.caller.Pc).Name(),
					err.code,
					msg.Internal,
					err.msg,
				)

				// Condensed stack trace
				if s.Flag('#') {
					str += "\n"

					// Inline stack trace
				} else {
					str += "\\n"
				}

				fmt.Fprint(w, str)
				w.Flush()
			}
			stackPos++
		}
	default:
		// Simple error messages
		fmt.Fprintf(s, "%s", errs.Error())
	}
}

/*
With adds a new error to the stack
*/
func (errs Err) With(err error) Err {
	if msg, ok := err.(Err); ok {
		errs = append(errs, msg...)
	} else if msg, ok := err.(Msg); ok {
		errs = append(errs, msg)
	} else {
		errs = append(errs, Msg{
			err:    err,
			caller: getCaller(),
			code:   0,
			msg:    "",
		})
	}
	return errs
}

/*
Wrap wraps an error into the stack.
*/
func Wrap(err error, msg string, code Code) Err {
	// Can't wrap a nil...
	if nil == err {
		return nil
	}
	var errs Err
	var ok bool
	if errs, ok = err.(Err); ok {
		errs = errs.With(Msg{
			err:    err,
			caller: getCaller(),
			code:   code,
			msg:    msg,
		})
	} else {
		errs = Err{Msg{
			err:    err,
			caller: getCaller(),
			code:   code,
			msg:    msg,
		}}
	}
	return errs
}
