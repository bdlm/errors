package errors

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strings"
	"sync"

	std "github.com/bdlm/std/error"
)

/*
Err defines an error heap.
*/
type Err struct {
	errs []ErrMsg
	mux  *sync.Mutex
}

/*
New returns an error with caller information for debugging.
*/
func New(code std.Code, msg string, data ...interface{}) Err {
	return Err{
		errs: []ErrMsg{Msg{
			err:    fmt.Errorf(msg, data...),
			caller: getCaller(),
			code:   code,
			msg:    msg,
			trace:  getTrace(),
		}},
		mux: &sync.Mutex{},
	}
}

/*
Caller returns the most recent error caller.
*/
func (err Err) Caller() std.Caller {
	var caller std.Caller
	if len(err.errs) > 0 {
		caller = err.errs[len(err.errs)-1].Caller()
	}
	return caller
}

/*
Cause returns the root cause of an error stack.
*/
func (err Err) Cause() error {
	if len(err.errs) > 0 {
		return err.errs[0]
	}
	return nil
}

/*
Code returns the most recent error code.
*/
func (err Err) Code() std.Code {
	code := ErrUnknown
	if len(err.errs) > 0 {
		code = err.errs[len(err.errs)-1].Code()
	}
	return code
}

/*
Detail implements the Coder interface. Detail returns the single-line
stack trace.
*/
func (err Err) Detail() string {
	if len(err.errs) > 0 {
		if code, ok := Codes[err.Code()]; ok {
			if "" != code.Detail() {
				return code.Detail()
			}
			return err.Error()
		}
	}
	return ""
}

/*
Error implements the error interface.
*/
func (err Err) Error() string {
	str := ""
	if len(err.errs) > 0 {
		str = err.errs[len(err.errs)-1].Error()
	}
	return str
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
func (err Err) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})
		for k := len(err.errs) - 1; k >= 0; k-- {
			err := err.errs[k]
			code, ok := Codes[err.Code()]
			if !ok {
				code = ErrCode{
					Int: err.Error(),
					Ext: err.Error(),
				}
			}

			errMsgInt := fmt.Sprintf("%04d", err.Code())
			if "" != code.Detail() {
				errMsgInt = fmt.Sprintf("%s: %s", errMsgInt, code.Detail())
			} else {
				errMsgInt = fmt.Sprintf("%s: %s", errMsgInt, err.Error())
			}

			errMsgExt := fmt.Sprintf("%04d", err.Code())
			if "" != code.String() {
				errMsgExt = fmt.Sprintf("%s: %s", errMsgExt, code.String())
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
				fmt.Fprintf(str, "\tmessage: %s\n", errMsgExt)

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
		fmt.Fprintf(state, err.Error())
	}
}

/*
From creates a new error stack based on a provided error and returns it.
*/
func From(code std.Code, err error) Err {
	if e, ok := err.(Err); ok {
		e.errs[len(e.errs)-1].SetCode(code)
		err = e
	} else {
		err = Err{
			errs: []ErrMsg{Msg{
				err:    err,
				caller: getCaller(),
				code:   code,
				msg:    err.Error(),
			}},
			mux: &sync.Mutex{},
		}
	}
	return err.(Err)
}

/*
HTTPStatus returns the associated HTTP status code, if any. Otherwise,
returns 200.
*/
func (err Err) HTTPStatus() int {
	status := http.StatusOK
	if len(err.errs) > 0 {
		if code, ok := Codes[err.errs[len(err.errs)-1].Code()]; ok {
			status = code.HTTPStatus()
		}
	}
	return status
}

/*
Len returns the size of the error stack.
*/
func (err Err) Len() int {
	return len(err.errs)
}

/*
Msg returns the error message.
*/
func (err Err) Msg() string {
	str := ""
	if len(err.errs) > 0 {
		str = err.errs[len(err.errs)-1].Msg()
	}
	return str
}

/*
String implements the stringer and Coder interfaces.
*/
func (err Err) String() string {
	return fmt.Sprintf("%s", err)
}

/*
Trace returns the call stack.
*/
func (err Err) Trace() std.Trace {
	var callers std.Trace
	for _, msg := range err.errs {
		callers = append(callers, msg.Caller())
	}
	return callers
}

/*
With adds a new error to the stack without changing the leading cause.
*/
func (err Err) With(e error, msg string, data ...interface{}) Err {
	// Can't include a nil...
	if nil == e {
		return err
	}

	if 0 == len(err.errs) {
		err = err.Push(Msg{
			err:    e,
			caller: getCaller(),
			code:   0,
			msg:    fmt.Sprintf(msg, data...),
		})
	} else {
		top := err.errs[len(err.errs)-1]
		err.errs = err.errs[:len(err.errs)-1]
		if msgs, ok := e.(Err); ok {
			err = err.Push(Msg{
				err:    fmt.Errorf(msg, data...),
				caller: getCaller(),
				code:   0,
				msg:    fmt.Sprintf(msg, data...),
			})
			err = err.Push(msgs.errs...)
		} else if msgs, ok := e.(Msg); ok {
			err = err.Push(Msg{
				err:    fmt.Errorf(msg, data...),
				caller: getCaller(),
				code:   0,
				msg:    err.Error(),
			}, msgs)
		} else {
			err = err.Push(Msg{
				err:    e,
				caller: getCaller(),
				code:   0,
				msg:    fmt.Sprintf(msg, data...),
			})
		}
		err = err.Push(top)
	}

	return err
}

/*
Push append an ErrMsg to the lst.
*/
func (err Err) Push(e ...ErrMsg) Err {
	err.mux.Lock()
	err.errs = append(err.errs, e...)
	err.mux.Unlock()
	return err
}

/*
Wrap wraps an error into a new stack led by msg.
*/
func Wrap(err error, code std.Code, msg string, data ...interface{}) Err {
	var errs = Err{
		errs: []ErrMsg{},
		mux:  &sync.Mutex{},
	}

	// Can't wrap a nil...
	if nil == err {
		return New(code, msg)
	}

	if e, ok := err.(Err); ok {
		errs = errs.Push(e.errs...)
	} else if e, ok := err.(Msg); ok {
		errs = errs.Push(e)
	} else {
		errs = Err{
			errs: []ErrMsg{Msg{
				err:    err,
				caller: getCaller(),
				code:   0,
				msg:    err.Error(),
			}},
			mux: &sync.Mutex{},
		}
	}

	errs = errs.Push(Msg{
		err:    fmt.Errorf(msg, data...),
		caller: getCaller(),
		code:   code,
		msg:    msg,
	})

	return errs
}
