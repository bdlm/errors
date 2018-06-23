package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestMsg(t *testing.T) {
	err := New(ErrUnknown, "new")
	if err.Error() != err.Msg() {
		t.Errorf("Expected '%s', received '%s'", err.Error(), err.Msg())
	}

	err = Err{}
	err = err.With(errors.New("error 1"), "msg 1")
	if "error 1" != err.Error() {
		t.Errorf("Expected 'error 1', received %s", err.Error())
	}
	if "msg 1" != err.Msg() {
		t.Errorf("Expected 'msg 1', received %s", err.Msg())
	}
}

func TestWith(t *testing.T) {
	// Can't wrap a nil
	err := New(0, "new")
	err2 := err.With(nil, "with")
	if fmt.Sprintf("%+v", err) != fmt.Sprintf("%+v", err2) {
		t.Errorf("Expected %v, received %v", fmt.Sprintf("%+v", err), fmt.Sprintf("%+v", err2))
	}

	// Wrapping with an empty stack makes the error the leading causer
	err = Err{}
	err = err.With(errors.New("error 1"), "msg 1")
	if 1 != len(err) {
		t.Errorf("Expected 1, received %d", len(err))
	}

	// Wrapping with a non-empty stack inserts the error after leading
	// causer
	err = err.With(errors.New("error 2"), "msg 2")
	err = err.With(errors.New("error 3"), "msg 3")
	if 3 != len(err) {
		t.Errorf("Expected 3, received %d", len(err))
	}

	if "error 1" != err[2].Error() {
		t.Errorf("Expected 'error 1', received %s", err[2].Error())
	}
	if "msg 1" != err[2].Msg() {
		t.Errorf("Expected 'msg 1', received %s", err[2].Msg())
	}

	// inserted...
	if "error 3" != err[1].Error() {
		t.Errorf("Expected 'error 3', received %s", err[1].Error())
	}
	if "msg 3" != err[1].Msg() {
		t.Errorf("Expected 'msg 3', received %s", err[1].Msg())
	}

	if "error 2" != err[0].Error() {
		t.Errorf("Expected 'error 2', received %s", err[0].Error())
	}
	if "msg 2" != err[0].Msg() {
		t.Errorf("Expected 'msg 2', received %s", err[0].Msg())
	}

	// Wrapping an individual Msg creates two stack entries
	err = Err{}
	msg := Msg{
		err:    nil,
		caller: getCaller(),
		code:   1,
		msg:    "msg 2",
	}
	err = err.With(errors.New("error 1"), "msg 1")
	err = err.With(msg, "msg 1")
	if 3 != len(err) {
		t.Errorf("Expected 3, received %d", len(err))
	}

	// Standard behavior
	err = New(0, "err 1")
	err = err.With(errors.New("error 2"), "msg 2")
	if 2 != len(err) {
		t.Errorf("Expected 2, received %d", len(err))
	}
	if "err 1" != err.Error() {
		t.Errorf("Expected 'err 1', received %s", err.Error())
	}
}
