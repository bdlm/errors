# go-errors

Go errors is inspired by `pkg/errors` and uses a similar API:

```go
import errs "github.com/mkenney/go-errors"
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

`%s`:
```
An unknown error occurred
```

`%v`:
```
(2) err_test.go:14 github.com/mkenney/go-errors_test.TestConstructor - 1:test message 3 'An unknown error occurred' Status 500
(1) err_test.go:13 github.com/mkenney/go-errors_test.TestConstructor - 1:test message 2 'An unknown error occurred' Status 500
(0) err_test.go:11 github.com/mkenney/go-errors_test.TestConstructor - 0:test message 1 'Error code unspecified' Status 500
```

`%+v`:
```
(2) err_test.go:14 github.com/mkenney/go-errors_test.TestConstructor
        Code: 1
        Mesg: test message 3
        Text: An unknown error occurred
        Http: 500
(1) err_test.go:13 github.com/mkenney/go-errors_test.TestConstructor
        Code: 1
        Mesg: test message 2
        Text: An unknown error occurred
        Http: 500
(0) err_test.go:11 github.com/mkenney/go-errors_test.TestConstructor
        Code: 0
        Mesg: test message 1
        Text: Error code unspecified
        Http: 500
```