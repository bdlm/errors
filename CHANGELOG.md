All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

- **Major**: backwards incompatible package updates
- **Minor**: feature additions
- **Patch**: bug fixes, backward compatible protobuf model changes, etc.

# v2.0.0 - 2020-05-01
`v2.0.0` is the production release of the `v0.2.0` development branch.

#### Added
* `go.mod`
* `github.com/bdlm/std` error interface

#### Changed
* replace interfaces with `github.com/bdlm/std/v2/errors` implementations

#### Removed
*

# v0.2.1 - 2019-07-18
#### Added
* Documentation and examples

#### Changed
* Minor formatting updates and cleanup


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
#### Changed
* Fixes a message formatting error


# v0.1.2 - 2018-09-09
#### Changed
* Fixes issues with concurrent writes


# v0.1.1 - 2018-08-22
#### Added
* Implemented a `Trace` method


# v0.1.0 - 2018-06-23
#### Added
* Initial release
