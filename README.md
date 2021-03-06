# errors

<a href="https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#mature"><img src="https://img.shields.io/badge/stability-mature-008000.svg" alt="Mature"></a> This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). This package is considered mature, you should expect package stability in <strong>Minor</strong> and <strong>Patch</strong> version releases

- **Major**: backwards incompatible package updates
- **Minor**: feature additions, removal of deprecated features
- **Patch**: bug fixes, backward compatible model and function changes, etc.

**[CHANGELOG](CHANGELOG.md)**<br>

<a href="https://github.com/bdlm/errors/blob/master/CHANGELOG.md"><img src="https://img.shields.io/github/v/release/bdlm/errors" alt="Release"></a>
<a href="https://pkg.go.dev/github.com/bdlm/errors/v2#pkg-examples"><img src="https://godoc.org/github.com/bdlm/errors?status.svg" alt="GoDoc"></a>
<a href="https://travis-ci.org/bdlm/errors"><img src="https://travis-ci.org/bdlm/errors.svg?branch=master" alt="Build status"></a>
<a href="https://codecov.io/gh/bdlm/errors"><img src="https://img.shields.io/codecov/c/github/bdlm/errors/master.svg" alt="Coverage status"></a>
<a href="https://goreportcard.com/report/github.com/bdlm/errors"><img src="https://goreportcard.com/badge/github.com/bdlm/errors" alt="Go Report Card"></a>
<a href="https://github.com/bdlm/errors/issues"><img src="https://img.shields.io/github/issues-raw/bdlm/errors.svg" alt="Github issues"></a>
<a href="https://github.com/bdlm/errors/pulls"><img src="https://img.shields.io/github/issues-pr/bdlm/errors.svg" alt="Github pull requests"></a>
<a href="https://github.com/bdlm/errors/blob/master/LICENSE"><img src="https://img.shields.io/github/license/bdlm/errors.svg" alt="MIT"></a>

`errors` provides simple, concise, useful error handling and annotation. This package aims to implement the [Error Inspection](https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md) and [Error Values](https://go.googlesource.com/proposal/+/master/design/go2draft-error-values-overview.md) Go2 [draft designs](https://go.googlesource.com/proposal/+/master/design/go2draft.md).

One of the biggest frustrations with Go error handling is the lack of forensic and meta information errors should provide. By default errors are just a string and possibly a type. They can't tell you where they occurred or the path through the call stack they followed. The error implementation in Go is robust enough to control program flow but it's not very efficient for troubleshooting or analysis.

Since the idom in Go is that we pass the error back up the stack anyway:
```go
if nil != err {
	return err
}
```
it's trivial to make errors much more informative with a simple error package. `bdlm/errors` makes this easy and supports tracing the call stack and the error callers with relative ease. Custom error types are also fully compatible with this package and can be used freely.

## Install

```
go get github.com/bdlm/errors/v2
```

## Quick start

See the [documentation](https://pkg.go.dev/github.com/bdlm/errors#pkg-examples) for more examples. All package methods work with any `error` type as well as `nil` values, and error instances implement the [`Unwrap`](https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md), [`Is`](https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md), [`Marshaler`](https://golang.org/pkg/encoding/json/#Marshaler), and [`Formatter`](https://golang.org/pkg/fmt/#Formatter) interfaces as well as the [`github.com/bdlm/std/errors`](https://github.com/bdlm/std/tree/master/errors) interfaces.

#### Create an error
```go
var MyError = errors.New("My error")
```

#### Create an error using formatting verbs
```go
var MyError = errors.Errorf("My error #%d", 1)
```

#### Wrap an error
```go
if nil != err {
	return errors.Wrap(err, "the operation failed")
}
```

#### Wrap an error with another error
```go
err := try1()
if nil != err {
	err2 := try2()
	if nil != err2 {
		return errors.WrapE(err, err2)
	}
	return err
}
```

#### Get the previous error, if any
```go
err := doWork()
if prevErr := errors.Unwrap(err); nil != prevErr {
	...
}
```

#### Test to see if a specific error type exists anywhere in an error stack
```go
var MyError = errors.New("My error")
func main() {
	err := doWork()
	if errors.Is(err, MyError) {
		...
	}
}
```

#### Iterate through an error stack
```go
err := doWork()
for nil != err {
	fmt.Println(err)
	err = errors.Unwrap(err)
}
```

#### Formatting verbs
`errors` implements the `%s` and `%v` [`fmt.Formatter`](https://golang.org/pkg/fmt/#hdr-Printing) formatting verbs and several modifier flags:

##### Verbs
* `%s` - Returns the error string of the last error added.
* `%v` - Alias for `%s`

##### Flags
* `#` - JSON formatted output, useful for logging
* `-` - Output caller details, useful for troubleshooting
* `+` - Output full error stack details, useful for debugging
* ` ` - (space) Add whitespace formatting for readability, useful for development

##### Examples
`fmt.Printf("%s", err)`
```
An error occurred
```
`fmt.Printf("%v", err)`
```
An error occurred
```
`fmt.Printf("%-v", err)`
```
#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors) - An error occurred
```
`fmt.Printf("%+v", err)`
```
#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors) - An error occurred #1 stack_test.go:39 (github.com/bdlm/error_test.TestErrors) - An error occurred
```
`fmt.Printf("%#v", err)`
```json
{"error":"An error occurred"}
```
`fmt.Printf("%#-v", err)`
```json
{"caller":"#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"}
```
`fmt.Printf("%#+v", err)`
```json
[{"caller":"#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"},{"caller":"#1 stack_test.go:39 (github.com/bdlm/error_test.TestErrors)","error":"An error occurred"}]
```
`fmt.Printf("% #-v", err)`
```json
{
    "caller": "#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)",
    "error": "An error occurred"
}
```
`fmt.Printf("% #+v", err)`
```json
[
    {
        "caller": "#0 stack_test.go:40 (github.com/bdlm/error_test.TestErrors)",
        "error": "An error occurred"
    },
    {
        "caller": "#1 stack_test.go:39 (github.com/bdlm/error_test.TestErrors)",
        "error": "An error occurred"
    }
]
```

#

See the [documentation](https://godoc.org/github.com/bdlm/errors#pkg-examples) for more examples.
