package errors

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

/*
Err defines an error heap.
*/
type Err []ErrMsg

/*
New returns an error with caller information for debugging. `codes` is
optional. Although you can pass multiple codes, only the first is
accepted.
*/
func New(msg string, codes ...Code) Err {
	var code Code
	if len(codes) > 0 {
		code = codes[0]
	}
	return Err{Msg{
		err:    errors.New(msg),
		caller: getCaller(),
		code:   code,
		msg:    msg,
	}}
}

/*
Callers returns the call stack.
*/
func (errs Err) Callers() []Caller {
	callers := []Caller{}
	for _, msg := range errs {
		callers = append(callers, msg.Caller())
	}
	return callers
}

/*
Cause returns the root cause of an error stack.
*/
func (errs Err) Cause() error {
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

/*
Code returns the most recent error code.
*/
func (errs Err) Code() Code {
	code := ErrUnknown
	if len(errs) > 0 {
		code = errs[len(errs)-1].Code()
	}
	return code
}

/*
Error implements the error interface.
*/
func (errs Err) Error() string {
	str := ""
	if len(errs) > 0 {
		str = errs[len(errs)-1].Error()
	}
	return str
}

/*
String implements the stringer interface.
*/
func (errs Err) String() string {
	return fmt.Sprintf("%s", errs)
}

/*
Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing

Format formats the stack trace output. Several verbs are supported:
	%s  - Returns the user-safe error string mapped to the error code or
	    the error message if none is specified.

	%v  - Alias for %s

	%#v - Returns the full stack trace in a single line, useful for
		logging. Same as %#v with the newlines escaped.

	%-v - Returns a multi-line stack trace, one column-delimited line
	    per error.

	%+v - Returns a multi-line detailed stack trace with multiple lines
	      per error. Only useful for human consumption.
*/
func (errs Err) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})
		for k := len(errs) - 1; k >= 0; k-- {
			err := errs[k]
			msg, ok := Codes[err.Code()]
			if !ok {
				msg = ErrCode{
					Int: err.Error(),
					Ext: err.Error(),
				}
			}

			errMsgInt := fmt.Sprintf("%04d", err.Code())
			if "" != msg.Detail() {
				errMsgInt = fmt.Sprintf("%s: %s", errMsgInt, msg.Detail())
			} else {
				errMsgInt = fmt.Sprintf("%s: %s", errMsgInt, err.Error())
			}

			errMsgExt := fmt.Sprintf("%04d", err.Code())
			if "" != msg.String() {
				errMsgExt = fmt.Sprintf("%s: %s", errMsgExt, msg.String())
			} else {
				errMsgExt = fmt.Sprintf("%s: %s", errMsgExt, err.Error())
			}

			switch {
			case state.Flag('+'):
				// Extended stack trace
				fmt.Fprintf(str, "#%d: `%s`\n", k, runtime.FuncForPC(err.Caller().Pc()).Name())
				fmt.Fprintf(str, "\terror:   %s\n", err.Msg())
				fmt.Fprintf(str, "\tline:    %s:%d\n", path.Base(err.Caller().File()), err.Caller().Line())
				fmt.Fprintf(str, "\tcode:    %s\n", errMsgInt)
				fmt.Fprintf(str, "\tmessage: %s\n\n", errMsgExt)

			case state.Flag('#'):
				// Condensed stack trace
				fmt.Fprintf(str, "#%d - \"%s\" %s:%d `%s` {%s}\n",
					k,
					err.Msg(),
					path.Base(err.Caller().File()),
					err.Caller().Line(),
					runtime.FuncForPC(err.Caller().Pc()).Name(),
					errMsgInt,
				)

			case state.Flag('-'):
				// Inline stack trace
				fmt.Fprintf(str, "#%d - \"%s\" %s:%d `%s` {%s} ",
					k,
					err.Msg(),
					path.Base(err.Caller().File()),
					err.Caller().Line(),
					runtime.FuncForPC(err.Caller().Pc()).Name(),
					errMsgInt,
				)

			default:
				// Externally-safe error message
				fmt.Fprintf(state, errMsgExt)
				return
			}
		}
		fmt.Fprintf(state, "%s", strings.Trim(str.String(), " \n\t"))
	default:
		// Externally-safe error message
		fmt.Fprintf(state, errs.Error())
	}
}

/*
From creates a new error stack based on a provided error and returns it.
*/
func From(err error, codes ...Code) Err {
	var code Code

	if len(codes) > 0 {
		code = codes[0]
	}

	if e, ok := err.(Err); ok {
		if len(codes) > 0 && len(e) > 0 {
			e[len(e)-1].SetCode(code)
		}
		err = e
	} else {
		err = Err{Msg{
			err:    err,
			caller: getCaller(),
			code:   code,
			msg:    err.Error(),
		}}
	}
	return err.(Err)
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
Wrap wraps an error into a new stack led by msg.
*/
func Wrap(err error, msg string, code ...Code) Err {
	var errs Err
	var errCode Code

	// Can't wrap a nil...
	if nil == err {
		return New(msg, code...)
	}

	if len(code) > 0 {
		errCode = code[0]
	}

	if e, ok := err.(Err); ok {
		errs = append(errs, e...)
	} else if e, ok := err.(Msg); ok {
		errs = append(errs, e)
	} else {
		errs = Err{Msg{
			err:    err,
			caller: getCaller(),
			code:   0,
			msg:    err.Error(),
		}}
	}

	errs = append(errs, Msg{
		err:    errors.New(msg),
		caller: getCaller(),
		code:   errCode,
		msg:    msg,
	})

	return errs
}
