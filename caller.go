package errors

import (
	"fmt"
	"runtime"
	"strings"
)

/*
Caller holds runtime.Caller data
*/
type Caller struct {
	Pc   uintptr
	File string
	Line int
	Ok   bool
}

/*
String implements the Stringer interface
*/
func (caller Caller) String() string {
	return fmt.Sprintf(
		"%s:%d %s",
		caller.File,
		caller.Line,
		runtime.FuncForPC(caller.Pc).Name(),
	)
}

func getCaller() Caller {
	caller := Caller{}
	a := 0
	for {
		if caller.Pc, caller.File, caller.Line, caller.Ok = runtime.Caller(a + 2); caller.Ok {
			if !strings.Contains(caller.File, "github.com/ReturnPath/rp-auth/error") {
				break
			}
		} else {
			break
		}
		a++
	}
	return caller
}
