# go-errors

Go errors is inspired by `pkg/errors` and uses a similar API:

```go
import errs "github.com/mkenney/go-errors"
```

## Error stacks

An error stack is an array of errors.

### Create a new stack

```go
if !someValidationCall() {
    err := errs.New("validation failed")
}
```

### Base a new stack off any error

```go
err := someValidationCall()
err := errs.Wrap("validation failed")
```

## Define error codes

See [`codes.go`](https://github.com/mkenney/go-errors/blob/master/codes.go)

```go
const (
	// Error codes below 1000 are reserved for future use.
	UserError errs.Code = iota + 1000
	InternalError
)
func init() {
	errs.Codes[UserError] = errs.Metadata{
		"A user error occurred",
		"bad user input",
		400,
	}
	errs.Codes[InternalError] = errs.Metadata{
		"An internal server occurred",
		"A service error occurred",
		500,
	}
}
func SomeFunc() error {
	return errs.New("Some internal thing broke", InternalError)
}
```

## Define a new error with an error code

Creating a new error defines the root of a backtrace.
```go
_, err := ioutil.ReadAll(r)
if err != nil {
	return errs.New("read failed", errs.ErrUnknown)
}
```


## Adding context to an error

The errors.Wrap function returns a new error that adds context to the original error and adds to the error stack trace:
```go
_, err := ioutil.ReadAll(r)
if err != nil {
	return errs.Wrap(err, "read failed", errs.ErrUnknown)
}
```

In this case, if the original `err` is not an instance of `Stack`, that error becomes the root of the error stack.

## Root cause of an error stack

Retrieving the root cause of an error stack is straightforward:
```go
log.Println(err.(errs.Stack).Cause())
```


 It implements the Formatter interface to provide a stack trace with the `%v` verb. `%+v` formats a more detailed trace.

Standard error output `%s`:
```
Internal Server Error
```

Single-line stack trace `%v`:
```
0 - err_test.go:38    github.com/mkenney/go-errors_test.TestOutput    1:An unknown error occurred    could not read configuration    \n 1 - err_test.go:37    github.com/mkenney/go-errors_test.TestOutput    0:Error code unspecified    failed to read data stream    \n 2 - err_test.go:36    github.com/mkenney/go-errors_test.TestOutput    0:Error code unspecified    read: end of input    \n
```

Multi-line condensed stack trace `%#v`:
```
0 - err_test.go:38    github.com/mkenney/go-errors_test.TestOutput    1:An unknown error occurred    could not read configuration
1 - err_test.go:37    github.com/mkenney/go-errors_test.TestOutput    0:Error code unspecified    failed to read data stream
2 - err_test.go:36    github.com/mkenney/go-errors_test.TestOutput    0:Error code unspecified    read: end of input
```

Multi-line detailed stack trace `%+v`:
```
0: github.com/mkenney/go-errors_test.TestOutput
        line: err_test.go: 38
        code: 1: An unknown error occurred
        mesg: could not read configuration

1: github.com/mkenney/go-errors_test.TestOutput
        line: err_test.go: 37
        code: 0: Error code unspecified
        mesg: failed to read data stream

2: github.com/mkenney/go-errors_test.TestOutput
        line: err_test.go: 36
        code: 0: Error code unspecified
        mesg: read: end of input
```