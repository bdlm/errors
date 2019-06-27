package errors

import (
	"fmt"
)

// Add creates a new stack (or updates a passed stack) with the new msg
// added behind the lead error. Useful for associating a set of errors in a
// distributed system.
func Add(e error, msg string, data ...interface{}) error {
	if nil == e {
		return nil
	}

	var clone stack
	switch typ := e.(type) {
	case stack:
		clone = typ.clone()
	default:
		clone = newStackFromErr(e)
	}
	ret := newEmptyStack()
	ret.stack = make([]err, len(clone.stack)+1)

	ret.stack[0] = clone.stack[0]
	ret.stack[1] = newErr(fmt.Errorf(msg, data...))
	for a := 1; a < len(clone.stack); a++ {
		ret.stack[a+1] = clone.stack[a]
	}

	return ret
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

	ret := newEmptyStack()
	switch typ := e.(type) {
	case stack:
		ret = typ.clone()
	}
	ret = ret.prepend(newErr(e))

	return ret
}

// Wrap returns a new error stack with a the leading error defined by msg.
func Wrap(e error, msg string, data ...interface{}) error {
	if nil == e {
		return nil
	}

	// Merge any passed stacks, Error instances, or other errors.
	ret := newEmptyStack()
	switch typ := e.(type) {
	case stack:
		ret = typ.clone()
	default:
		ret = newStackFromErr(typ)
	}

	// Add the new message to the stack.
	ret = ret.prepend(newErr(fmt.Errorf(msg, data...)))

	return ret
}
