package errors

import (
	"fmt"
	"sync"
)

// Add creates a new stack (or updates a passed stack) with the new msg
// added behind the lead error. Useful for associating a set of errors in a
// distributed system.
func Add(e error, msg string, data ...interface{}) Error {
	if nil == e {
		return nil
	}

	var errs stack
	switch err := e.(type) {
	case stack:
		errs = err
	default:
		errs = newStackFromErr(e)
	}

	errs.mux.Lock()
	last := errs.stack[len(errs.stack)-1]
	errs.stack[len(errs.stack)-1] = newErr(fmt.Errorf(msg, data...))
	errs.stack = append(errs.stack, last)
	errs.mux.Unlock()

	return errs
}

// New creates a new error stack defined by msg.
func New(msg string, data ...interface{}) Error {
	return newStack(msg, data...)
}

// Track adds caller metadata to an error as it's passed back up the stack.
func Track(e error) error {
	if nil == e {
		return nil
	}

	err := stack{stack: []err{}, mux: &sync.Mutex{}}
	switch typ := e.(type) {
	case stack:
		for _, v := range typ.stack {
			err.append(v)
		}
	}
	err.append(newErr(e))

	return err
}

// Wrap returns a new error stack with a the leading error defined by msg.
func Wrap(e error, msg string, data ...interface{}) error {
	if nil == e {
		return nil
	}

	// Merge any passed stacks, Error instances, or other errors.
	err := stack{stack: []err{}, mux: &sync.Mutex{}}
	switch typ := e.(type) {
	case stack:
		for _, v := range typ.stack {
			err.append(v)
		}
	default:
		err = newStackFromErr(typ)
	}

	// Add the new message to the stack.
	err.append(newErr(fmt.Errorf(msg, data...)))

	return err
}
