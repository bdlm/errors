package errors

import (
	"fmt"
	"runtime"
)

type Caller interface {
	File() string
	Line() int
	Ok() bool
	Pc() uintptr
	String() string
}

/*
Call holds runtime.Caller data
*/
type Call struct {
	loaded bool
	file   string
	line   int
	ok     bool
	pc     uintptr
}

func (call Call) File() string {
	return call.file
}
func (call Call) Line() int {
	return call.line
}
func (call Call) Ok() bool {
	return call.ok
}
func (call Call) Pc() uintptr {
	return call.pc
}

/*
String implements the Stringer interface
*/
func (call Call) String() string {
	return fmt.Sprintf(
		"%s:%d %s",
		call.file,
		call.line,
		runtime.FuncForPC(call.pc).Name(),
	)
}

func getCaller() Caller {
	call := Call{}
	call.pc, call.file, call.line, call.ok = runtime.Caller(2)
	//a := 0
	//for {
	//	if call.pc, call.file, call.line, call.ok = runtime.Caller(a + 2); call.ok {
	//		break
	//	} else {
	//		break
	//	}
	//	a++
	//}
	return call
}
