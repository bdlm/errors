package errors

import (
	"errors"
	"fmt"
	"path"
	"runtime"
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
	msg := ""

	if _, ok := Codes[errs.Code()]; ok {
		msg = Codes[errs.Code()].External
	} else if len(errs) > 0 {
		msg = errs[len(errs)-1].Error()
	}

	return msg
}

/*
Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing

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
				msg = Codes[ErrUnknown]
			}
			switch {
			case s.Flag('+'):
				// Extended stack trace
				fmt.Fprintf(s, "#%d: `%s`\n", stackPos, runtime.FuncForPC(err.caller.Pc).Name())
				fmt.Fprintf(s, "\terror:   %s\n", err.msg)
				fmt.Fprintf(s, "\tline:    %s:%d\n", path.Base(err.caller.File), err.caller.Line)
				fmt.Fprintf(s, "\tcode:    %d - %s\n", err.code, msg.Internal)
				fmt.Fprintf(s, "\tentry:   %v\n", runtime.FuncForPC(err.caller.Pc).Entry())
				fmt.Fprintf(s, "\tmessage: %s\n\n", msg.External)

			case s.Flag('#'):
				// Condensed stack trace
				fmt.Fprintf(s, "#%d - \"%s\" %s:%d `%s` {%04d: %s}\n",
					stackPos,
					err.msg,
					path.Base(err.caller.File),
					err.caller.Line,
					runtime.FuncForPC(err.caller.Pc).Name(),
					err.code,
					msg.Internal,
				)

			default:
				// Loggable stack trace
				fmt.Fprintf(s, "#%d - \"%s\" %s:%d `%s` {%04d: %s} ",
					stackPos,
					err.msg,
					path.Base(err.caller.File),
					err.caller.Line,
					runtime.FuncForPC(err.caller.Pc).Name(),
					err.code,
					msg.Internal,
				)
			}
			stackPos++
		}
	default:
		// Externally-save error message
		fmt.Fprintf(s, "%04d: %s", errs.Code(), errs.Error())
	}
}

/*
With adds a new error to the stack without changing the leading cause.
*/
func (errs Err) With(err error) Err {
	// Can't include a nil...
	if nil == err {
		return errs
	}

	if 0 == len(errs) {
		errs = append(errs, Msg{
			err:    err,
			caller: getCaller(),
			code:   0,
			msg:    err.Error(),
		})
	} else {
		top := errs[len(errs)-1]
		errs = errs[:len(errs)-1]
		if msg, ok := err.(Err); ok {
			errs = append(errs, msg...)
		} else if msg, ok := err.(Msg); ok {
			errs = append(errs, msg)
		} else {
			errs = append(errs, Msg{
				err:    err,
				caller: getCaller(),
				code:   0,
				msg:    err.Error(),
			})
		}
		errs = append(errs, top)
	}

	return errs
}

/*
Wrap wraps an error into the stack.
*/
func Wrap(err error, msg string, code ...Code) Err {
	// Can't wrap a nil...
	if nil == err {
		return New(msg, code...)
	}

	var errs Err
	var errCode Code
	var ok bool

	if len(code) > 0 {
		errCode = code[0]
	}

	if errs, ok = err.(Err); ok {
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

	} else {
		errs = Err{Msg{
			err:    err,
			caller: getCaller(),
			code:   errCode,
			msg:    msg,
		}}
	}
	return errs
}
