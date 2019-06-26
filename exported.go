package errors

import (
	"errors"
	"fmt"
	"sync"
)

// Add creates a new stack (or updates a passed stack) the new msg added
// below the leading error. Usefule for adding debugging data to a system
// error.
func Add(err error, msg string, data ...interface{}) error {
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
func New(msg string, data ...interface{}) error {
	return newStackFromErr(fmt.Errorf(msg, data...))
}

// Track adds caller metadata to an error as it's passed back up the stack.
func Track(err error) error {
	if nil == err {
		return nil
	}

	stack := Stack{
		stack: []Error{},
		mux:   &sync.Mutex{},
	}

	switch e := err.(type) {
	case Stack:
		for _, v := range e.stack {
			stack.stack = append(stack.stack, v)
		}
		stack.stack[len(stack.stack)-1].err = nil
	}
	stack.stack = append(stack.stack, Error{
		err:    err,
		caller: getCaller(),
	})

	return stack
}

// Wrap returns a new error stack with a the leading error defined by msg.
func Wrap(err error, msg string, data ...interface{}) error {
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
