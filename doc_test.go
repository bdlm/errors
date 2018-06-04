package errors_test

import (
	"errors"
	"fmt"

	errs "github.com/bdlm/errors"
)

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return errs.Wrap(err, "could not read configuration file", errs.ErrEOF)
}

func decodeConfig() error {
	err := readConfig()
	return errs.Wrap(err, "could not decode configuration data", errs.ErrInvalidJSON)
}

func loadConfig() error {
	err := decodeConfig()
	return errs.Wrap(err, "service configuration could not be loaded", errs.ErrFatal)
}

func someWork() error {
	return fmt.Errorf("failed to do work")
}

var errEOF = fmt.Errorf("read: end of input")
var otherErr = fmt.Errorf("some other process failed")

func ExampleNew() {
	var err error

	// If an error code isn't used or doesn't have a corresponding
	// ErrCode defined, the error message is returned.
	err = errs.New("this is an error message")
	fmt.Println(err)

	// If an error with a corresponding ErrCode is specified, the
	// user-safe error string mapped to the error code is returned,
	// along with the code.
	err = errs.New("this is an error message", errs.ErrFatal)
	fmt.Println(err)

	// Output: 0000: this is an error message
	//0001: Internal Server Error
}

func ExampleWrap_backtrace() {
	// To build up an error stack, add context to each error before
	// returning it up the call stack.
	err := loadConfig()
	if nil != err {
		err = errs.Wrap(err, "failed to load configuration")
	}

	// The %v formatting verb can be used to print out the stack trace
	// in various ways. The %v verb is the default and prints out the
	// standard error message.
	fmt.Println(err)

	// The %-v verb is useful for logging and prints the trace on a
	// single line.
	fmt.Printf("%-v\n", err)

	// The %#v verb prints each cause in the stack on a separate line.
	fmt.Printf("%#v\n", err)

	// The %+v verb prints a verbose detailed backtrace intended for
	// human consumption.
	fmt.Printf("%+v\n", err)

	// Output: 0000: failed to load configuration
	//#4 - "failed to load configuration" doc_test.go:55 `github.com/bdlm/errors_test.ExampleWrap_backtrace` {0000: unknown error} #3 - "service configuration could not be loaded" doc_test.go:22 `github.com/bdlm/errors_test.loadConfig` {0001: fatal error} #2 - "could not decode configuration data" doc_test.go:17 `github.com/bdlm/errors_test.decodeConfig` {0200: invalid JSON data could not be decoded} #1 - "could not read configuration file" doc_test.go:12 `github.com/bdlm/errors_test.readConfig` {0100: unexpected EOF} #0 - "read: end of input" doc_test.go:12 `github.com/bdlm/errors_test.readConfig` {0000: unknown error}
	//#4 - "failed to load configuration" doc_test.go:55 `github.com/bdlm/errors_test.ExampleWrap_backtrace` {0000: unknown error}
	//#3 - "service configuration could not be loaded" doc_test.go:22 `github.com/bdlm/errors_test.loadConfig` {0001: fatal error}
	//#2 - "could not decode configuration data" doc_test.go:17 `github.com/bdlm/errors_test.decodeConfig` {0200: invalid JSON data could not be decoded}
	//#1 - "could not read configuration file" doc_test.go:12 `github.com/bdlm/errors_test.readConfig` {0100: unexpected EOF}
	//#0 - "read: end of input" doc_test.go:12 `github.com/bdlm/errors_test.readConfig` {0000: unknown error}
	//#4: `github.com/bdlm/errors_test.ExampleWrap_backtrace`
	//	error:   failed to load configuration
	//	line:    doc_test.go:55
	//	code:    0000: unknown error
	//	message: 0000: failed to load configuration
	//
	//#3: `github.com/bdlm/errors_test.loadConfig`
	//	error:   service configuration could not be loaded
	//	line:    doc_test.go:22
	//	code:    0001: fatal error
	//	message: 0001: Internal Server Error
	//
	//#2: `github.com/bdlm/errors_test.decodeConfig`
	//	error:   could not decode configuration data
	//	line:    doc_test.go:17
	//	code:    0200: invalid JSON data could not be decoded
	//	message: 0200: Invalid JSON Data
	//
	//#1: `github.com/bdlm/errors_test.readConfig`
	//	error:   could not read configuration file
	//	line:    doc_test.go:12
	//	code:    0100: unexpected EOF
	//	message: 0100: End of input
	//
	//#0: `github.com/bdlm/errors_test.readConfig`
	//	error:   read: end of input
	//	line:    doc_test.go:12
	//	code:    0000: unknown error
	//	message: 0000: read: end of input
}

func ExampleFrom() {
	// Converting an error from another package into an error stack is
	// straightforward.
	err := errors.New("my error")
	if _, ok := err.(errs.Err); !ok {
		err = errs.From(err)
	}

	fmt.Println(err)

	// Output: 0000: my error
}

func ExampleErr_With() {
	// To add to an error stack without modifying the leading cause, add
	// additional errors to the stack with the With() method.
	err := loadConfig()
	if nil != err {
		if e, ok := err.(errs.Err); nil != err && ok {
			err = e.With(errors.New("failed to load configuration"))
		} else {
			err = errs.From(err)
		}
	}

	fmt.Println(err)
	// Output: 0001: Internal Server Error
}

func ExampleErr_Format() {
	err := loadConfig()
	if nil != err {
		err = errs.Wrap(err, "could not load config")
	}

	fmt.Printf("%s", err)
	fmt.Printf("%v", err)
	fmt.Printf("%#v", err)
	fmt.Printf("%+v", err)
}
