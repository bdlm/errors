package errors

import (
	std "github.com/bdlm/std/error"
)

/*
ErrMsg defines the interface to error message data.
*/
type ErrMsg interface {
	Caller() std.Caller
	Code() std.Code
	Error() string
	Msg() string
	SetCode(std.Code) ErrMsg
	Trace() std.Trace
}

/*
Msg defines a single error message.
*/
type Msg struct {
	err    error
	caller std.Caller
	code   std.Code
	msg    string
	trace  std.Trace
}

/*
Caller implements ErrMsg.
*/
func (msg Msg) Caller() std.Caller {
	return msg.caller
}

/*
Code implements ErrMsg.
*/
func (msg Msg) Code() std.Code {
	return msg.code
}

/*
Error implements error.
*/
func (msg Msg) Error() string {
	return msg.String()
}

/*
Msg implements ErrMsg.
*/
func (msg Msg) Msg() string {
	return msg.msg
}

/*
SetCode implements ErrMsg.
*/
func (msg Msg) SetCode(code std.Code) ErrMsg {
	msg.code = code
	return msg
}

/*
String implements Stringer.
*/
func (msg Msg) String() string {
	if nil == msg.err {
		return msg.msg
	}
	return msg.err.Error()
}

/*
Trace implements ErrMsg.
*/
func (msg Msg) Trace() std.Trace {
	return msg.trace
}
