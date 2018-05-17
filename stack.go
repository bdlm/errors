package errors

import (
	"errors"
	"fmt"
	"path"
	"runtime"
)

/*
ErrStack defines an error heap
*/
type ErrStack []Msg

/*
New returns an error with caller information for debugging.
*/
func New(msg string, code Code) ErrStack {
	return ErrStack{Msg{
		err:    errors.New(msg),
		caller: getCaller(),
		code:   code,
		msg:    msg,
	}}
}

/*
Callers returns an array of callers
*/
func (stack ErrStack) Callers() []Caller {
	callers := []Caller{}
	for _, msg := range stack {
		if msg.caller.Ok {
			callers = append(callers, msg.caller)
		}
	}
	return callers
}

/*
Cause returns cause of an error stack.
*/
func (stack ErrStack) Cause() Msg {
	if len(stack) > 0 {
		return stack[0]
	}
	return Msg{}
}

/*
Code returns the most recent error code
*/
func (stack ErrStack) Code() Code {
	code := ErrUnspecified
	if len(stack) > 0 {
		code = stack[len(stack)-1].code
	}
	return code
}

/*
Error implements the error interface
*/
func (stack ErrStack) Error() string {
	meta, ok := Codes[stack.Code()]
	if !ok {
		meta = Codes[ErrUnspecified]
	}
	return meta.Ext
}

/*
Format implements fmt.Formatter.

Format formats the stack trace output.
*/
func (stack ErrStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmtStr := ""
		for k, err := range stack {
			msg, ok := Codes[err.code]
			if !ok {
				msg = Codes[ErrUnspecified]
			}
			switch {
			case s.Flag('+'):
				// Detailed stack trace
				fmtStr = `(%d) %s:%d %s
	Code: %d
	Mesg: %s
	Text: %s
	Http: %d
`

			default:
				// Condensed stack trace
				fmtStr = "(%d) %s:%d %s - %d:%s  '%s' Status %d\n"
			}

			fmt.Fprintf(s, "%s", fmt.Sprintf(
				fmtStr,
				k,
				path.Base(err.caller.File),
				err.caller.Line,
				runtime.FuncForPC(err.caller.Pc).Name(),
				err.code,
				err.msg,
				msg.Int,
				msg.HTTPStatus,
			))
		}
	case 's':
		// Simple error messages
		fmt.Fprintf(s, "%s", stack.Error())
	}
}

/*
With adds a new error to the stack
*/
func (stack ErrStack) With(err error) ErrStack {
	if msg, ok := err.(ErrStack); ok {
		stack = append(stack, msg...)
	} else if msg, ok := err.(Msg); ok {
		stack = append(stack, msg)
	} else {
		stack = append(stack, Msg{
			err:    err,
			caller: getCaller(),
			code:   0,
			msg:    "",
		})
	}
	return stack
}

/*
Wrap wraps an error in an ErrStack.
*/
func Wrap(err error, msg string, code Code) ErrStack {
	// Can't wrap a nil...
	if nil == err {
		return nil
	}
	var stack ErrStack
	var ok bool
	if stack, ok = err.(ErrStack); ok {
		stack = stack.With(Msg{
			err:    err,
			caller: getCaller(),
			code:   code,
			msg:    msg,
		})
	} else {
		stack = ErrStack{Msg{
			err:    err,
			caller: getCaller(),
			code:   code,
			msg:    msg,
		}}
	}
	return stack
}
