package errors

import (
	"errors"
	"fmt"
)

// Add creates a new stack (or updates a passed stack) the new msg added
// below the leading error. Usefule for adding debugging data to a system
// error.
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

	errs.mux.Lock()
	last := errs.stack[len(errs.stack)-1]
	errs.stack[len(errs.stack)-1] = newError(fmt.Errorf(msg, data...))
	errs.stack = append(errs.stack, last)
	errs.mux.Unlock()

	return errs
}

// New creates a new error stack defined by msg.
func New(msg string, data ...interface{}) Stack {
	return newStackFromErr(fmt.Errorf(msg, data...))
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
	errs.stack = append(errs.stack, Error{
		err:    fmt.Errorf(msg, data...),
		caller: getCaller(),
	})

	return errs
}
