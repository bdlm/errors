package errors

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestMsg(t *testing.T) {
	err := New(ErrUnknown, "new")
	if err.Error() != err.Msg() {
		t.Errorf("Expected '%s', received '%s'", err.Error(), err.Msg())
	}

	err = Err{mux: &sync.Mutex{}}
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
	err = Err{mux: &sync.Mutex{}}
	err = err.With(errors.New("error 1"), "msg 1")
	if 1 != len(err.errs) {
		t.Errorf("Expected 1, received %d", len(err.errs))
	}

	// Wrapping with a non-empty stack inserts the error after leading
	// causer
	err = err.With(errors.New("error 2"), "msg 2")
	err = err.With(errors.New("error 3"), "msg 3")
	if 3 != len(err.errs) {
		t.Errorf("Expected 3, received %d", len(err.errs))
	}

	if "error 1" != err.errs[2].Error() {
		t.Errorf("Expected 'error 1', received %s", err.errs[2].Error())
	}
	if "msg 1" != err.errs[2].Msg() {
		t.Errorf("Expected 'msg 1', received %s", err.errs[2].Msg())
	}

	// inserted...
	if "error 3" != err.errs[1].Error() {
		t.Errorf("Expected 'error 3', received %s", err.errs[1].Error())
	}
	if "msg 3" != err.errs[1].Msg() {
		t.Errorf("Expected 'msg 3', received %s", err.errs[1].Msg())
	}

	if "error 2" != err.errs[0].Error() {
		t.Errorf("Expected 'error 2', received %s", err.errs[0].Error())
	}
	if "msg 2" != err.errs[0].Msg() {
		t.Errorf("Expected 'msg 2', received %s", err.errs[0].Msg())
	}

	// Wrapping an individual Msg creates two stack entries
	err = Err{mux: &sync.Mutex{}}
	msg := Msg{
		err:    nil,
		caller: getCaller(),
		code:   1,
		msg:    "msg 2",
	}
	err = err.With(errors.New("error 1"), "msg 1")
	err = err.With(msg, "msg 1")
	if 3 != len(err.errs) {
		t.Errorf("Expected 3, received %d", len(err.errs))
	}

	// Standard behavior
	err = New(0, "err 1")
	err = err.With(errors.New("error 2"), "msg 2")
	if 2 != len(err.errs) {
		t.Errorf("Expected 2, received %d", len(err.errs))
	}
	if "err 1" != err.Error() {
		t.Errorf("Expected 'err 1', received %s", err.Error())
	}
}
