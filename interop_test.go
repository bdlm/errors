package errors_test

import (
	std_errors "errors"
	"fmt"
	"io/fs"
	"testing"

	"github.com/bdlm/errors/v2"
)

// These tests cover interoperation with the STANDARD library's errors package, which is what the
// v2.1.2 behaviour broke. Each one failed before the Unwrap fix.

// Unwrap must yield the wrapped error itself, whatever its type. Re-boxing a foreign error as an *E
// meant no other concrete type was ever exposed, and the chain ended at the box because the box
// carried no prev of its own.
func TestUnwrapYieldsTheWrappedErrorItself(t *testing.T) {
	foreign := fs.ErrNotExist
	wrapped := errors.Wrap(foreign, "reading the config")

	got := std_errors.Unwrap(wrapped)
	if foreign != got {
		t.Fatalf("Unwrap returned %#v (%T), want the wrapped error itself", got, got)
	}
	if nil != std_errors.Unwrap(got) {
		t.Error("the chain should end at the wrapped error, not at a box around it")
	}
}

// errors.Is must reach a foreign sentinel through any depth of bdlm wrapping. Before the fix only
// sentinels created by THIS package were reachable, because only those survived Unwrap.
func TestStdIsReachesAForeignSentinel(t *testing.T) {
	for depth, err := range map[string]error{
		"one wrap":   errors.Wrap(fs.ErrNotExist, "a"),
		"two wraps":  errors.Wrap(errors.Wrap(fs.ErrNotExist, "a"), "b"),
		"four wraps": errors.Wrap(errors.Wrap(errors.Wrap(errors.Wrap(fs.ErrNotExist, "a"), "b"), "c"), "d"),
	} {
		if !std_errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%s: errors.Is could not reach the sentinel", depth)
		}
	}
}

// A sentinel reached through a foreign wrapper in the middle of the chain. This is the shape that
// matters in practice: a library wraps with fmt.Errorf("%w"), and this package wraps that.
func TestStdIsReachesThroughAForeignWrapper(t *testing.T) {
	mixed := errors.Wrap(fmt.Errorf("opening: %w", fs.ErrNotExist), "loading config")
	if !std_errors.Is(mixed, fs.ErrNotExist) {
		t.Error("errors.Is could not reach a sentinel behind a fmt.Errorf wrapper")
	}
}

// errors.As must find a wrapped concrete type. This is the one that broke grpc: status.FromError is
// errors.As for an interface{ GRPCStatus() *Status }, so before the fix a wrapped status error could
// never be found and every wrapped failure surfaced as codes.Unknown -- a 500 for what was often a
// client error.
type notFoundError struct{ path string }

func (e *notFoundError) Error() string { return "not found: " + e.path }

func TestStdAsFindsAWrappedConcreteType(t *testing.T) {
	wrapped := errors.Wrap(errors.Wrap(&notFoundError{path: "/etc/app.conf"}, "reading"), "starting up")

	var target *notFoundError
	if !std_errors.As(wrapped, &target) {
		t.Fatal("errors.As could not find the wrapped concrete type")
	}
	if "/etc/app.conf" != target.path {
		t.Errorf("target.path = %q, want the wrapped value's field", target.path)
	}
}

// THE grpc CASE, which is what makes this fix load-bearing rather than tidy.
//
// A modern grpc status.FromError is, in its own words, errors.As(err, &interface{ GRPCStatus()
// *Status }). So a status error wrapped by this package could never be found: every link the walk
// yielded was an *E, and *E does not implement that interface. Every wrapped failure therefore
// surfaced as codes.Unknown -- an HTTP 500 -- including failures that were plainly the caller's
// fault, such as an unacceptable auth token.
//
// This asserts the MECHANISM rather than calling status.FromError, deliberately: this module still
// vendors grpc v1.29.1, whose FromError is a bare type assertion with no errors.As at all, so it
// cannot see through ANY wrapper and would fail this test for a reason that has nothing to do with
// the fix. The interface below has the same shape grpc looks for, and the assertion below is the same
// one grpc makes.
type fakeStatus struct{ code int }

func (s *fakeStatus) Error() string           { return fmt.Sprintf("rpc error: code = %d", s.code) }
func (s *fakeStatus) GRPCStatus() *fakeStatus { return s }

func TestAnInterfaceTargetIsFoundThroughTheChain(t *testing.T) {
	fromKeyfunc := &fakeStatus{code: 16} // 16 == Unauthenticated
	asInterceptorSeesIt := errors.Wrap(fromKeyfunc, "unable to resolve the JWT signing key")

	var target interface{ GRPCStatus() *fakeStatus }
	if !std_errors.As(asInterceptorSeesIt, &target) {
		t.Fatal("errors.As could not find the interface implementation through the wrap")
	}
	if 16 != target.GRPCStatus().code {
		t.Errorf("found the wrong value: code = %d", target.GRPCStatus().code)
	}

	// And through several layers, since a real chain is rarely one deep.
	deeper := errors.Wrap(errors.Wrap(asInterceptorSeesIt, "validating"), "authenticating")
	var again interface{ GRPCStatus() *fakeStatus }
	if !std_errors.As(deeper, &again) {
		t.Error("errors.As could not find it through four wraps")
	}
}

// Error() carries the cause, so a single rendered string says what actually went wrong. Previously it
// returned only this frame's message, so "unable to resolve the JWT signing key" reached the caller
// with the reason -- already recorded one link down -- discarded.
func TestErrorIncludesTheCause(t *testing.T) {
	err := errors.Wrap(errors.Wrap(fs.ErrNotExist, "reading the config"), "starting up")
	want := "starting up: reading the config: file does not exist"
	if want != err.Error() {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
	// An unwrapped error is unchanged -- no trailing separator.
	if "plain" != errors.New("plain").Error() {
		t.Errorf("an unwrapped error's message changed: %q", errors.New("plain").Error())
	}
}

// The package's own Is must not stop at the first link that has an Is method. It used to return that
// method's answer directly, and every *E has one, so a sentinel further down was unreachable.
func TestPackageIsDoesNotTerminateTheWalk(t *testing.T) {
	if !errors.Is(errors.Wrap(errors.Wrap(fs.ErrNotExist, "a"), "b"), fs.ErrNotExist) {
		t.Error("package Is could not reach the sentinel")
	}
	if errors.Is(errors.Wrap(fs.ErrNotExist, "a"), fs.ErrClosed) {
		t.Error("package Is matched a sentinel that is not in the chain")
	}
}

// Guard rails: nil handling must not change.
func TestNilHandling(t *testing.T) {
	if nil != std_errors.Unwrap(errors.New("no cause")) {
		t.Error("an error with no cause should unwrap to nil")
	}
	if errors.Is(nil, fs.ErrNotExist) || errors.Is(errors.New("x"), nil) {
		t.Error("nil comparisons should be false")
	}
}
