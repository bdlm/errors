# go-errors

Go errors is inspired by `pkg/errors` and uses a similar API. It implements the Formatter interface to provide a stack trace with the `%v` verb. `%+v` formats a more detailed trace.

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