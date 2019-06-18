package errors

import (
	"errors"
	"fmt"
)

// Add creates a new stack with the new error behind the specified error.
func Add(err error, msg string, data ...interface{}) Stack {
	var errs Stack

	if nil == err {
		err = errors.New("")
	}

	// Merge any passed stacks, Error instances, or other errors.
	switch e := err.(type) {
	case Stack:
		errs = e
	default:
		errs = newStackFromErr(err)
	}

	pop := errs.last()
	errs.mux.Lock()
	errs.stack[len(errs.stack)-1] = Error{
		err:    fmt.Errorf(msg, data...),
		caller: getCaller(),
	}
	errs.mux.Unlock()
	errs.append(pop)

	return errs
}

// Wrap creates a new stack with a the leading error defined by msg.
func Wrap(err error, msg string, data ...interface{}) Stack {
	var errs Stack

	if nil == err {
		err = errors.New("")
	}

	// Merge any passed stacks, Error instances, or other errors.
	switch e := err.(type) {
	case Stack:
		errs = e
	default:
		errs = newStackFromErr(err)
	}

	// Add the new message to the stack.
	errs.append(Error{
		err:    fmt.Errorf(msg, data...),
		caller: getCaller(),
	})

	return errs
}
