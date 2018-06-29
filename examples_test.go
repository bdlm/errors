package errors_test

import (
	"errors"
	"fmt"

	errs "github.com/bdlm/errors"
)

var errEOF = fmt.Errorf("read: end of input")
var otherErr = fmt.Errorf("some other process failed")

func ExampleNew() {
	var err error

	// If an error code isn't used or doesn't have a corresponding
	// ErrCode defined, the error message is returned.
	err = errs.New(0, "this is an error message")
	fmt.Println(err)

	// If an error with a corresponding ErrCode is specified, the
	// user-safe error string mapped to the error code is returned,
	// along with the code.
	err = errs.New(errs.ErrFatal, "this is an error message")
	fmt.Println(err)

	// Output: 0000: An unknown error occurred
	//0001: A fatal error occurred

}

func ExampleWrap_backtrace() {
	// To build up an error stack, add context to each error before
	// returning it up the call stack.
	err := loadConfig()
	if nil != err {
		err = errs.Wrap(err, 0, "failed to load configuration")
	}

	// The %v formatting verb can be used to print out the stack trace
	// in various ways. The %v verb is the default and prints out the
	// standard error message.
	fmt.Println(err)

	// The %-v verb is useful for logging and prints the trace on a
	// single line.
	fmt.Printf("%-v\n\n", err)

	// The %#v verb prints each cause in the stack on a separate line.
	fmt.Printf("%#v\n\n", err)

	// The %+v verb prints a verbose detailed backtrace intended for
	// human consumption.
	fmt.Printf("%+v\n\n", err)

	// Output: 0000: An unknown error occurred
	//#4 - "failed to load configuration" examples_test.go:37 `github.com/bdlm/errors_test.ExampleWrap_backtrace` {0000: failed to load configuration} #3 - "service configuration could not be loaded" mocks_test.go:30 `github.com/bdlm/errors_test.loadConfig` {1000: the configuration is invalid} #2 - "could not decode configuration data" mocks_test.go:35 `github.com/bdlm/errors_test.decodeConfig` {0208: could not decode configuration data} #1 - "could not read configuration file" mocks_test.go:40 `github.com/bdlm/errors_test.readConfig` {0100: unexpected EOF} #0 - "read: end of input" mocks_test.go:40 `github.com/bdlm/errors_test.readConfig` {0000: read: end of input}
	//
	//#4 - "failed to load configuration" examples_test.go:37 `github.com/bdlm/errors_test.ExampleWrap_backtrace` {0000: failed to load configuration}
	//#3 - "service configuration could not be loaded" mocks_test.go:30 `github.com/bdlm/errors_test.loadConfig` {1000: the configuration is invalid}
	//#2 - "could not decode configuration data" mocks_test.go:35 `github.com/bdlm/errors_test.decodeConfig` {0208: could not decode configuration data}
	//#1 - "could not read configuration file" mocks_test.go:40 `github.com/bdlm/errors_test.readConfig` {0100: unexpected EOF}
	//#0 - "read: end of input" mocks_test.go:40 `github.com/bdlm/errors_test.readConfig` {0000: read: end of input}
	//
	//#4: `github.com/bdlm/errors_test.ExampleWrap_backtrace`
	//	error:   failed to load configuration
	//	line:    examples_test.go:37
	//	code:    0000: failed to load configuration
	//	message: 0000: An unknown error occurred
	//#3: `github.com/bdlm/errors_test.loadConfig`
	//	error:   service configuration could not be loaded
	//	line:    mocks_test.go:30
	//	code:    1000: the configuration is invalid
	//	message: 1000: Configuration not valid
	//#2: `github.com/bdlm/errors_test.decodeConfig`
	//	error:   could not decode configuration data
	//	line:    mocks_test.go:35
	//	code:    0208: could not decode configuration data
	//	message: 0208: could not decode configuration data
	//#1: `github.com/bdlm/errors_test.readConfig`
	//	error:   could not read configuration file
	//	line:    mocks_test.go:40
	//	code:    0100: unexpected EOF
	//	message: 0100: End of input
	//#0: `github.com/bdlm/errors_test.readConfig`
	//	error:   read: end of input
	//	line:    mocks_test.go:40
	//	code:    0000: read: end of input
	//	message: 0000: An unknown error occurred
}

func ExampleFrom() {
	// Converting an error from another package into an error stack is
	// straightforward.
	err := errors.New("my error")
	if _, ok := err.(errs.Err); !ok {
		err = errs.From(0, err)
	}

	fmt.Println(err)
	fmt.Println(err.(errs.Err).Detail())

	// Output: 0000: An unknown error occurred
	//my error
}

func ExampleErr_With() {
	// To add to an error stack without modifying the leading cause, add
	// additional errors to the stack with the With() method.
	err := loadConfig()
	if nil != err {
		if e, ok := err.(errs.Err); nil != err && ok {
			err = e.With(errs.New(0, "failed to load configuration"), "loadConfig returned an error")
		} else {
			err = errs.From(0, err)
		}
	}

	fmt.Println(err)
	// Output: 1000: Configuration not valid
}

func ExampleDetail() {
	err := loadConfig()
	if nil != err {
		err = errs.Wrap(err, 0, "failed to load configuration")
	}

	// The single-line condensed stack trace is also availabe via the
	// Detail() method
	fmt.Println(err.(errs.Err).Detail())

	// Output: failed to load configuration
}

func ExampleHTTPStatus() {
	err := loadConfig()
	if nil != err {
		err = errs.Wrap(err, ConfigurationNotValid, "failed to load configuration")
	}

	// HTTPStatus() returns the HTTP status code associated with an error
	// code, if any.
	fmt.Println(err.(errs.Err).HTTPStatus())

	// Output: 500
}
