All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

- **Major**: backwards incompatible package updates
- **Minor**: feature additions, removal of deprecated features
- **Patch**: bug fixes, backward compatible model and function changes, etc.

# v2.1.5
#### Fixed
Eight defects, all found by a new contract test suite (`contract_test.go`) organised by obligation
rather than by function. Coverage 95.9%, race-clean.

* **Panics on nil.** `doc.go` promises every method works with nil values; four did not.
  `(*E).Error`, `(*E).Is` and `(*E).MarshalJSON` dereferenced a nil receiver, and `(*E).Is` also
  panicked whenever the frame carried no annotation of its own — `reflect.TypeOf(nil)` returns a
  NIL `reflect.Type`, and calling `Comparable` on that panics rather than answering false. The
  package-level `Is` had the same flaw. Both now share a `comparableErrors` helper.
* **`New` and `Wrap` corrupted any message containing a percent sign.** Both passed the message
  through `fmt.Errorf`, so `New("100% complete")` produced `"100%!c(MISSING)omplete"` and
  `"disk usage at 95%"` produced `"...95%!(NOVERB)"`. A message is now stored verbatim; only
  `Errorf`, and `Wrap` when given arguments, formats.
* **`MarshalJSON` panicked on a foreign wrapper in the chain.** `list()` includes any error
  implementing `Unwrap() error`, so a `fmt.Errorf("%w")` link reached an unguarded type assertion
  and a nil field read. Marshalling an error must never panic — it runs while something is already
  failing.
* **`Trace` severed the chain.** It held the wrapped error in the annotation slot and left `prev`
  nil, so `Unwrap` returned nothing and no sentinel below a trace was reachable. The error is now
  on the chain, and a frame with no message of its own renders transparently.
* **`Track` severed the chain for anything that was not already an `*E`.** `prev` came from a
  synthetic box whose own `prev` was nil, so a `fmt.Errorf("%w")` wrapper, a joined error, or any
  foreign wrapping type lost everything beneath it.
* **`Is` could find what `As` could not.** `WrapE`'s annotation is not on the `Unwrap` chain, and
  `Is` special-cased it while `As` did not — so which function you reached for changed the answer.
  `(*E)` now implements the standard `As(interface{}) bool` hook, and the package-level `As`
  searches the same slot.

#### Documentation
* `doc.go` documented `errors.Has`, which does not exist; `Is` searches the whole chain, which is
  what that section described. It also showed `errors.Wrap(err, err2)`, which does not compile —
  `Wrap` takes a message string, so wrapping with an error is `WrapE`.
* Stated the two guarantees now covered by tests: nothing panics for any combination of nil
  arguments and nil receivers, and `Is`/`As` agree with the standard library on every chain shape.

# v2.1.4 - 2026-08-20
#### Fixed
* **`Is` and `As` now traverse multi-errors.** An error may wrap several causes — `errors.Join`, or
  `fmt.Errorf` with more than one `%w` — and those implement `Unwrap() []error`, which the
  single-error `Unwrap` cannot express.
* `Unwrap` still returns `nil` for a multi-error, as `std/errors.Unwrap` does — there is no single
  previous error to return. This is now documented rather than implicit; callers walking a chain
  themselves must branch on `Unwrap() []error`.

# v2.1.3 - 2026-08-20
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
