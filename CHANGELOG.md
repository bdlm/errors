All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

- **Major**: backwards incompatible package updates
- **Minor**: feature additions, removal of deprecated features
- **Patch**: bug fixes, backward compatible model and function changes, etc.

# v2.1.4
#### Fixed
* **`Is` and `As` now traverse multi-errors.** An error may wrap several causes — `errors.Join`, or
  `fmt.Errorf` with more than one `%w` — and those implement `Unwrap() []error`, which the
  single-error `Unwrap` cannot express: it returns `nil` for them. Both walks therefore terminated at
  the join and every cause beneath it was unreachable, so `Is` answered `false` for a sentinel sitting
  *directly inside* the error it was handed, unwrapped. `golang-jwt` joins its sentinels, so no `jwt`
  error was matchable at all and callers classifying on `jwt.ErrTokenMalformed` saw a caller's bad
  credential fall through to a 5xx. Both now branch over `Unwrap() []error` and search the whole tree,
  matching the standard library.
* `Unwrap` still returns `nil` for a multi-error, as `std/errors.Unwrap` does — there is no single
  previous error to return. This is now documented rather than implicit; callers walking a chain
  themselves must branch on `Unwrap() []error`.

# v2.1.3-rc1
#### Added
* Interoperation tests against the standard library's `errors` package (`interop_test.go`).

#### Changed
* **`func (e *E) Unwrap() error` now returns the wrapped error itself.** It previously re-boxed a
  wrapped error that did not implement `std/errors.Error` as `&E{err: prev}`. That box carried no
  `prev`, so the chain ended there and the original error value was never handed to the caller — its
  own `Unwrap`, `Is` and `As` were never consulted. The consequence was that `errors.As` could only
  ever match `*E`, since no other concrete type was ever yielded, so anything built on `errors.As`
  could not see through a wrap. gRPC is the sharpest example: `status.FromError` is `errors.As` for an
  `interface{ GRPCStatus() *Status }`, so a wrapped status error was reported as `codes.Unknown` — an
  HTTP 500 — including for failures that were plainly the caller's fault, such as an unacceptable auth
  token.
* **`func Is(err, test error) bool` no longer ends the walk early.** It returned the first link's `Is`
  answer directly; every `*E` has an `Is` method, so a `false` there ended the search and a sentinel
  further down was unreachable. A negative answer from one link now continues to the next.
  `func (e *E) Is(error) bool` is unchanged and still searches the chain — the standard library calls
  that hook as `ok && x.Is(target)`, so its answer never terminated `errors.Is`'s own walk.
* **`func (e *E) Error() string` includes the wrapped message**, `"outer: inner"`, as
  `fmt.Errorf("%w")` produces. Returning only the frame's own message discarded everything beneath it
  wherever an error is rendered through `Error()` rather than a `%+v` verb — which is most of the
  places that reach a user, including an HTTP or gRPC response body. The trace and JSON formats print
  one line per frame and now use each frame's own message, so their output is unchanged.

#### Removed
* n/a

# v2.1.2 - 2021-05-12
#### Added
* n/a

#### Changed
* Updated `func (e *E) Error() string` implementation to handle nil pointer references.

#### Removed
* n/a

# v2.1.1 - 2020-05-30
#### Added
* n/a

#### Changed
* Updated various docs

#### Removed
* n/a

# v2.1.0 - 2020-05-30
#### Added
* Additional documentation
* Pre-generics implementation of `As`, returns the error instance when found

#### Changed
* Upgrade to `bdlm/std` v2.1.0
* Cleaner type checking in `Is` methods

#### Removed
* deprecated `Has` methods, use `Is` instead

# v2.0.1 - 2020-05-01
#### Added
* n/a

#### Changed
* Refactoring to better implement the `std/errors.Error` interface

#### Removed
* n/a

# v2.0.0 - 2020-05-01
`v2.0.0` is the production release of the `v0.2.0` development branch.

#### Added
* `go.mod`
* `github.com/bdlm/std/v2/errors` interfaces

#### Changed
* licence changed from BSD to MIT
* replace interfaces with `github.com/bdlm/std/v2/errors` implementations
* simplified formatting and marshalling logic
* renamed `GetCaller(error) std_err.Caller` to `Caller(error) std_err.Caller`

#### Removed
* unused code

# v0.2.1 - 2019-07-18
#### Added
* Documentation and examples

#### Changed
* Minor formatting updates and cleanup

#### Removed
* n/a


# v0.2.0 - 2019-07-18
This release is a full rewrite of the `errors` package. See the [README](README.md) for further details.
#### Added
* `Caller` interface
* `Error` interface
* Exported methods
  - `Errorf(msg string, data ...interface{}) Error`
  - `GetCaller(err error) Caller`
  - `Has(err, test error) bool`
  - `Is(err, test error) bool`
  - `Trace(e error) Error`
  - `Track(e error) Error`

#### Changed
* Exported methods
  - `New(code std.Code, msg string, data ...interface{}) *Err`
    - `New(msg string) Error`
  - `Wrap(err error, code std.Code, msg string, data ...interface{}) *Err`
    - `Wrap(e error, msg string, data ...interface{}) Error`
#### Removed
* Exported methods
  - `From(code std.Code, err error) *Err`
* Support for error codes
* Support for sanitized vs raw error messages
* Support for HTTP status codes


# v0.1.3 - 2018-10-02
#### Added
* n/a

#### Changed
* Fixes a message formatting error

#### Removed
* n/a

# v0.1.2 - 2018-09-09
#### Added
* n/a

#### Changed
* Fixes issues with concurrent writes

#### Removed
* n/a

# v0.1.1 - 2018-08-22
#### Added
* Implemented a `Trace` method

#### Changed
* n/a

#### Removed
* n/a

# v0.1.0 - 2018-06-23
#### Added
* Initial release

#### Changed
* n/a

#### Removed
* n/a
