package errors_test

import (
	"encoding/json"
	std_errors "errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bdlm/errors/v2"
)

// This file is the package's contract, expressed as tests. It is deliberately organised by
// OBLIGATION rather than by function, because the obligations are what callers depend on:
//
//   1. nothing panics, for any input, including nil receivers and nil arguments
//   2. a message survives construction verbatim
//   3. the chain is walkable, and agrees with the standard library at every depth
//   4. Is and As find what is there and only what is there
//   5. the rendering paths (Error, Format, MarshalJSON) never lose or invent information
//
// doc.go states the package "aims to implement the Error Inspection and Error Values Go2 draft
// designs" and that "all package methods work with any error type as well as nil values". Those
// two sentences are the specification these tests hold it to.

var (
	sentinel = std_errors.New("sentinel")
	other    = std_errors.New("other")
)

// custom is a foreign error type: not created by this package, with its own Is hook. Real chains
// contain these, and they are where an error package's chain walking usually breaks.
type custom struct{ msg string }

func (c *custom) Error() string { return c.msg }

// valueErr uses a VALUE receiver, so *valueErr is a pointer whose element implements error --
// the shape bdlm's As(err, test error) requires of its target, and one the standard library's
// As(err, any) accepts too. It exists so the two can be compared on identical input.
//
// Worth noting on its own: bdlm's As signature differs from the standard library's
// (As(err, test error) error vs As(err error, target any) bool), which is why a target type has to
// be chosen with both in mind.
type valueErr struct{ msg string }

func (v valueErr) Error() string { return v.msg }

// customMatching has an Is hook that matches sentinel, the pattern the standard library documents.
type customMatching struct{}

func (customMatching) Error() string   { return "matches sentinel" }
func (customMatching) Is(t error) bool { return t == sentinel }

// ---------------------------------------------------------------------------------------------
// 1. Nothing panics. Every exported entry point, with the inputs a caller can actually supply.
// ---------------------------------------------------------------------------------------------

func TestNilSafety(t *testing.T) {
	var nilE *errors.E

	cases := map[string]func(){
		"New empty":            func() { _ = errors.New("") },
		"Errorf empty":         func() { _ = errors.Errorf("") },
		"Wrap nil":             func() { _ = errors.Wrap(nil, "msg") },
		"Wrap nil empty msg":   func() { _ = errors.Wrap(nil, "") },
		"WrapE nil nil":        func() { _ = errors.WrapE(nil, nil) },
		"WrapE err nil":        func() { _ = errors.WrapE(sentinel, nil) },
		"WrapE nil err":        func() { _ = errors.WrapE(nil, sentinel) },
		"Unwrap nil":           func() { _ = errors.Unwrap(nil) },
		"Is nil nil":           func() { _ = errors.Is(nil, nil) },
		"Is nil target":        func() { _ = errors.Is(sentinel, nil) },
		"Is nil err":           func() { _ = errors.Is(nil, sentinel) },
		"As nil nil":           func() { _ = errors.As(nil, nil) },
		"As nil err":           func() { _ = errors.As(nil, sentinel) },
		"Caller nil":           func() { _ = errors.Caller(nil) },
		"Trace nil":            func() { _ = errors.Trace(nil) },
		"Track nil":            func() { _ = errors.Track(nil) },
		"nil E Error":          func() { _ = nilE.Error() },
		"nil E Unwrap":         func() { _ = nilE.Unwrap() },
		"nil E Caller":         func() { _ = nilE.Caller() },
		"nil E Is":             func() { _ = nilE.Is(sentinel) },
		"nil E MarshalJSON":    func() { _, _ = nilE.MarshalJSON() },
		"nil E format %v":      func() { _ = fmt.Sprintf("%v", nilE) },
		"nil E format %+v":     func() { _ = fmt.Sprintf("%+v", nilE) },
		"nil E format %#v":     func() { _ = fmt.Sprintf("%#v", nilE) },
		"WrapE nil err Error":  func() { _ = errors.WrapE(sentinel, nil).Error() },
		"WrapE nil err Is":     func() { _ = errors.WrapE(sentinel, nil).Is(sentinel) },
		"WrapE nil err format": func() { _ = fmt.Sprintf("%+v", errors.WrapE(sentinel, nil)) },
	}
	for name, call := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); nil != r {
					t.Fatalf("panicked: %v", r)
				}
			}()
			call()
		})
	}
}

// ---------------------------------------------------------------------------------------------
// 2. A message survives construction verbatim.
// ---------------------------------------------------------------------------------------------

// TestMessagesAreNotReinterpretedAsFormatStrings: New and Wrap take a MESSAGE. A caller passing a
// message that happens to contain a percent sign -- a URL-encoded value, a percentage, a Windows
// path -- must get that message back, not fmt's opinion of it.
func TestMessagesAreNotReinterpretedAsFormatStrings(t *testing.T) {
	for _, msg := range []string{
		"100% complete",
		"disk usage at 95%",
		"unexpected %s in input",
		"query returned %d rows",
		"literal %% sign",
	} {
		t.Run(msg, func(t *testing.T) {
			if got := errors.New(msg).Error(); msg != got {
				t.Errorf("New(%q).Error() = %q", msg, got)
			}
			if got := errors.Wrap(nil, msg).Error(); msg != got {
				t.Errorf("Wrap(nil, %q).Error() = %q", msg, got)
			}
		})
	}
}

// TestErrorfFormats: Errorf is the one that SHOULD interpret its argument as a format string.
func TestErrorfFormats(t *testing.T) {
	if got, want := errors.Errorf("%d of %d", 3, 7).Error(), "3 of 7"; want != got {
		t.Errorf("Errorf = %q, want %q", got, want)
	}
	if got, want := errors.Wrap(nil, "%d rows", 42).Error(), "42 rows"; want != got {
		t.Errorf("Wrap with args = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------------------------
// 3. The chain is walkable and agrees with the standard library.
// ---------------------------------------------------------------------------------------------

// TestUnwrapContract: Unwrap yields the wrapped error ITSELF, so a caller can type-assert it.
func TestUnwrapContract(t *testing.T) {
	target := &custom{msg: "foreign"}
	wrapped := errors.Wrap(target, "outer")

	got := std_errors.Unwrap(wrapped)
	if got != error(target) {
		t.Fatalf("Unwrap returned %#v, want the original *custom", got)
	}
	if _, ok := got.(*custom); !ok {
		t.Errorf("Unwrap yielded %T, not the concrete wrapped type", got)
	}
	if nil != errors.Unwrap(errors.New("leaf")) {
		t.Error("a leaf error must unwrap to nil")
	}
}

// TestChainDepthAgreesWithStdlib walks a deep mixed chain and requires this package and the
// standard library to give the same answer at every depth. Mixed on purpose: bdlm frames, a
// fmt.Errorf frame and a foreign type, which is what real chains look like.
func TestChainDepthAgreesWithStdlib(t *testing.T) {
	chain := errors.Wrap(
		fmt.Errorf("fmt layer: %w",
			errors.Wrap(
				errors.WrapE(sentinel, &custom{msg: "foreign"}),
				"inner bdlm")),
		"outer bdlm")

	if !std_errors.Is(chain, sentinel) {
		t.Error("std errors.Is cannot reach the sentinel through the mixed chain")
	}
	if !errors.Is(chain, sentinel) {
		t.Error("package Is cannot reach the sentinel through the mixed chain")
	}
	if std_errors.Is(chain, other) != errors.Is(chain, other) {
		t.Error("package Is and std errors.Is disagree on an absent sentinel")
	}
	var target *custom
	if !std_errors.As(chain, &target) {
		t.Error("std errors.As cannot reach the foreign type")
	}
}

// TestWrappedByFmtIsReachable: fmt.Errorf("%w") of a bdlm error must stay inspectable, since that
// is how most Go code wraps.
func TestWrappedByFmtIsReachable(t *testing.T) {
	inner := errors.Wrap(sentinel, "bdlm frame")
	outer := fmt.Errorf("fmt frame: %w", inner)

	if !std_errors.Is(outer, sentinel) {
		t.Error("sentinel unreachable through fmt.Errorf")
	}
	var e *errors.E
	if !std_errors.As(outer, &e) {
		t.Error("the *E frame is not reachable via std errors.As")
	}
}

// ---------------------------------------------------------------------------------------------
// 4. Is and As find what is there, and only what is there.
// ---------------------------------------------------------------------------------------------

func TestIsFindsSentinelAtEveryDepth(t *testing.T) {
	depths := map[string]error{
		"depth 0": sentinel,
		"depth 1": errors.Wrap(sentinel, "a"),
		"depth 2": errors.Wrap(errors.Wrap(sentinel, "a"), "b"),
		"depth 5": errors.Wrap(errors.Wrap(errors.Wrap(errors.Wrap(errors.Wrap(sentinel, "a"), "b"), "c"), "d"), "e"),
	}
	for name, err := range depths {
		t.Run(name, func(t *testing.T) {
			if !errors.Is(err, sentinel) {
				t.Error("package Is missed it")
			}
			if !std_errors.Is(err, sentinel) {
				t.Error("std errors.Is missed it")
			}
			if errors.Is(err, other) {
				t.Error("matched an absent sentinel")
			}
		})
	}
}

// TestIsRespectsForeignIsHooks: a foreign type's Is hook answering false must NOT end the walk,
// because the sentinel may still be deeper in the chain.
func TestIsRespectsForeignIsHooks(t *testing.T) {
	// customMatching.Is(sentinel) is true -- found via the hook, not by identity.
	if !errors.Is(errors.Wrap(customMatching{}, "wrapped"), sentinel) {
		t.Error("a foreign Is hook that matches was not consulted")
	}
	// A foreign type whose hook says false must not hide a deeper match. `withFalseIs` wraps the
	// sentinel and answers false for everything.
	hidden := errors.Wrap(&falseIs{inner: sentinel}, "wrapped")
	if !std_errors.Is(hidden, sentinel) {
		t.Fatal("premise: std errors.Is should reach it via Unwrap")
	}
	if !errors.Is(hidden, sentinel) {
		t.Error("a foreign Is hook answering false ended the walk and hid a deeper match")
	}
}

// falseIs answers false to every Is while still wrapping a real cause.
type falseIs struct{ inner error }

func (f *falseIs) Error() string { return "false-is: " + f.inner.Error() }
func (f *falseIs) Is(error) bool { return false }
func (f *falseIs) Unwrap() error { return f.inner }

func TestAsTargets(t *testing.T) {
	chain := errors.Wrap(errors.Wrap(&custom{msg: "found me"}, "a"), "b")

	var concrete *custom
	if !std_errors.As(chain, &concrete) || "found me" != concrete.msg {
		t.Error("std errors.As could not reach the concrete foreign type")
	}

	var self *errors.E
	if !std_errors.As(chain, &self) {
		t.Error("std errors.As could not match *E itself")
	}

	var absent *falseIs
	if std_errors.As(chain, &absent) {
		t.Error("As matched a type that is not in the chain")
	}
}

// TestAsRejectsInvalidTargets: the package's own As must not panic on a target the standard
// library would reject outright.
func TestAsRejectsInvalidTargets(t *testing.T) {
	defer func() {
		if r := recover(); nil != r {
			t.Fatalf("panicked on an invalid target: %v", r)
		}
	}()
	if nil != errors.As(errors.New("x"), nil) {
		t.Error("As with a nil target should find nothing")
	}
	if nil != errors.As(errors.New("x"), sentinel) {
		t.Error("As with a non-pointer target should find nothing")
	}
}

// ---------------------------------------------------------------------------------------------
// 5. The rendering paths never lose or invent information.
// ---------------------------------------------------------------------------------------------

// TestErrorIncludesTheWholeChain: Error() is what an HTTP or gRPC body carries, so a cause
// recorded one link down must not be silently dropped.
func TestErrorIncludesTheWholeChain(t *testing.T) {
	err := errors.Wrap(errors.Wrap(sentinel, "middle"), "outer")
	got := err.Error()
	for _, want := range []string{"outer", "middle", "sentinel"} {
		if !strings.Contains(got, want) {
			t.Errorf("Error() = %q, missing %q", got, want)
		}
	}
}

// TestFormatVerbs: every documented verb produces something, none panics, and the trace verbs
// include more than the plain ones.
func TestFormatVerbs(t *testing.T) {
	err := errors.Wrap(errors.Wrap(sentinel, "middle"), "outer")
	for _, verb := range []string{"%s", "%v", "%q", "%-v", "%+v", "%#v", "%#+v", "% +v"} {
		t.Run(verb, func(t *testing.T) {
			out := fmt.Sprintf(verb, err)
			if "" == out {
				t.Errorf("%s produced nothing", verb)
			}
			if strings.Contains(out, "%!") {
				t.Errorf("%s produced a formatting error: %q", verb, out)
			}
		})
	}
	if plain, trace := fmt.Sprintf("%v", err), fmt.Sprintf("%+v", err); len(trace) <= len(plain) {
		t.Errorf("%%+v (%d bytes) should carry more than %%v (%d bytes)", len(trace), len(plain))
	}
}

// TestMarshalJSONWithForeignWrapperInChain: a chain containing a non-*E error that itself wraps
// -- fmt.Errorf("%w") being the common case -- must marshal, not panic.
func TestMarshalJSONWithForeignWrapperInChain(t *testing.T) {
	defer func() {
		if r := recover(); nil != r {
			t.Fatalf("MarshalJSON panicked on a foreign wrapper in the chain: %v", r)
		}
	}()
	err := errors.Wrap(fmt.Errorf("fmt wrapper: %w", sentinel), "outer")
	raw, marshalErr := json.Marshal(err)
	if nil != marshalErr {
		t.Fatalf("marshal: %v", marshalErr)
	}
	var entries []map[string]interface{}
	if err := json.Unmarshal(raw, &entries); nil != err {
		t.Fatalf("result is not the documented array shape: %v", err)
	}
	if 0 == len(entries) {
		t.Error("marshalled to an empty array")
	}
}

func TestMarshalJSONShape(t *testing.T) {
	raw, err := json.Marshal(errors.Wrap(errors.New("inner"), "outer"))
	if nil != err {
		t.Fatalf("marshal: %v", err)
	}
	var entries []map[string]interface{}
	if err := json.Unmarshal(raw, &entries); nil != err {
		t.Fatalf("unmarshal: %v", err)
	}
	if 2 > len(entries) {
		t.Fatalf("entries = %d, want at least 2 (one per frame): %s", len(entries), raw)
	}
	for i, entry := range entries {
		if _, ok := entry["caller"]; !ok {
			t.Errorf("entry %d has no caller", i)
		}
	}
}

// ---------------------------------------------------------------------------------------------
// 6. Caller data is captured, and points at the caller rather than into this package.
// ---------------------------------------------------------------------------------------------

func TestCallerPointsAtTheCallSite(t *testing.T) {
	err := errors.New("here")
	clr := errors.Caller(err)
	if nil == clr {
		t.Fatal("no caller recorded")
	}
	if !strings.Contains(clr.File(), "contract_test.go") {
		t.Errorf("caller file = %q, want this test file", clr.File())
	}
	if 0 == clr.Line() {
		t.Error("caller line is zero")
	}
}

func TestTraceAndTrackPreserveTheChain(t *testing.T) {
	for name, build := range map[string]func(error) error{
		"Trace": func(e error) error { return errors.Trace(e) },
		"Track": func(e error) error { return errors.Track(e) },
	} {
		t.Run(name, func(t *testing.T) {
			out := build(errors.Wrap(sentinel, "inner"))
			if !std_errors.Is(out, sentinel) {
				t.Error("the sentinel is no longer reachable")
			}
			if "" == out.Error() {
				t.Error("message lost")
			}
		})
	}
}

// ---------------------------------------------------------------------------------------------
// 7. Is and As agree with the standard library as a PROPERTY, across chain shapes, rather than in
//    a handful of hand-picked cases. Disagreement is the failure mode that matters: a caller
//    cannot tell which of the two functions to trust.
// ---------------------------------------------------------------------------------------------

func TestIsAndAsAgreeWithStdlibAcrossShapes(t *testing.T) {
	shapes := map[string]error{
		"bare sentinel":            sentinel,
		"wrapped once":             errors.Wrap(sentinel, "a"),
		"wrapped twice":            errors.Wrap(errors.Wrap(sentinel, "a"), "b"),
		"annotated with foreign":   errors.WrapE(sentinel, &custom{msg: "ann"}),
		"fmt over bdlm":            fmt.Errorf("f: %w", errors.Wrap(sentinel, "a")),
		"bdlm over fmt":            errors.Wrap(fmt.Errorf("f: %w", sentinel), "a"),
		"joined":                   std_errors.Join(other, sentinel),
		"bdlm over joined":         errors.Wrap(std_errors.Join(other, sentinel), "a"),
		"joined over bdlm":         std_errors.Join(other, errors.Wrap(sentinel, "a")),
		"nested joins":             std_errors.Join(other, std_errors.Join(other, sentinel)),
		"traced":                   errors.Trace(errors.Wrap(sentinel, "a")),
		"tracked":                  errors.Track(errors.Wrap(sentinel, "a")),
		"foreign hiding the cause": errors.Wrap(&falseIs{inner: sentinel}, "a"),
		"absent, wrapped":          errors.Wrap(other, "a"),
		"absent, deep":             errors.Wrap(errors.Wrap(errors.Wrap(other, "a"), "b"), "c"),
	}
	for name, err := range shapes {
		t.Run(name, func(t *testing.T) {
			std, pkg := std_errors.Is(err, sentinel), errors.Is(err, sentinel)
			if std != pkg {
				t.Errorf("Is disagreement: std=%v package=%v", std, pkg)
			}
			// The same chain, asked for a concrete type instead of a sentinel.
			var viaStd valueErr
			stdAs := std_errors.As(err, &viaStd)
			pkgAs := nil != errors.As(err, &valueErr{})
			if stdAs != pkgAs {
				t.Errorf("As disagreement on valueErr: std=%v package=%v", stdAs, pkgAs)
			}
		})
	}
}

// TestSentinelIsFoundWhereItIsAndNotWhereItIsNot: the shapes above assert agreement; this asserts
// the answers are actually CORRECT, so both could not be wrong together.
func TestSentinelIsFoundWhereItIsAndNotWhereItIsNot(t *testing.T) {
	present := []error{
		errors.Wrap(sentinel, "a"),
		errors.WrapE(sentinel, &custom{msg: "ann"}),
		errors.Wrap(std_errors.Join(other, sentinel), "a"),
		errors.Trace(errors.Wrap(sentinel, "a")),
		errors.Track(errors.Wrap(sentinel, "a")),
		errors.Wrap(&falseIs{inner: sentinel}, "a"),
	}
	for i, err := range present {
		if !errors.Is(err, sentinel) {
			t.Errorf("present[%d]: sentinel not found in %v", i, err)
		}
	}
	absent := []error{
		errors.New("unrelated"),
		errors.Wrap(other, "a"),
		errors.Wrap(errors.Wrap(other, "a"), "b"),
		std_errors.Join(other, other),
	}
	for i, err := range absent {
		if errors.Is(err, sentinel) {
			t.Errorf("absent[%d]: sentinel found where it is not: %v", i, err)
		}
	}
}

// ---------------------------------------------------------------------------------------------
// 8. Robustness: depth, and the rendering paths over a multi-error.
// ---------------------------------------------------------------------------------------------

// TestDeepChainDoesNotExhaustTheStack: Error, Is and the formatters all recurse, and a wrapped
// error can accumulate a frame per call-stack level in a long-running service.
func TestDeepChainDoesNotExhaustTheStack(t *testing.T) {
	var err error = sentinel
	for i := 0; i < 2000; i++ {
		err = errors.Wrap(err, fmt.Sprintf("frame %d", i))
	}
	if !errors.Is(err, sentinel) {
		t.Error("sentinel unreachable at depth 2000")
	}
	if !std_errors.Is(err, sentinel) {
		t.Error("std errors.Is could not reach depth 2000")
	}
	if "" == err.Error() {
		t.Error("Error() produced nothing at depth 2000")
	}
	if out := fmt.Sprintf("%+v", err); "" == out {
		t.Error("the trace verb produced nothing at depth 2000")
	}
	if _, marshalErr := json.Marshal(err); nil != marshalErr {
		t.Errorf("marshal at depth 2000: %v", marshalErr)
	}
}

// TestRenderingAMultiError: a joined error has no single previous error, so the single-Unwrap
// rendering paths must still produce something sensible rather than dropping the branches or
// panicking.
func TestRenderingAMultiError(t *testing.T) {
	err := errors.Wrap(std_errors.Join(other, sentinel), "outer")

	if got := err.Error(); !strings.Contains(got, "outer") {
		t.Errorf("Error() lost its own message: %q", got)
	}
	for _, verb := range []string{"%v", "%+v", "%#v"} {
		out := fmt.Sprintf(verb, err)
		if "" == out || strings.Contains(out, "%!") {
			t.Errorf("verb %s over a multi-error produced %q", verb, out)
		}
	}
	if _, marshalErr := json.Marshal(err); nil != marshalErr {
		t.Errorf("marshal over a multi-error: %v", marshalErr)
	}
}

// ---------------------------------------------------------------------------------------------
// 9. Concurrency. An error value is shared freely once created -- returned up a stack, logged on
//    one goroutine while inspected on another -- so reads must not mutate it. Run with -race.
// ---------------------------------------------------------------------------------------------

func TestConcurrentReadsAreSafe(t *testing.T) {
	err := errors.Wrap(errors.Wrap(sentinel, "middle"), "outer")
	const workers = 32
	done := make(chan struct{}, workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for j := 0; j < 50; j++ {
				_ = err.Error()
				_ = errors.Is(err, sentinel)
				_ = std_errors.Is(err, sentinel)
				_ = fmt.Sprintf("%+v", err)
				_, _ = json.Marshal(err)
				_ = errors.Unwrap(err)
				_ = errors.Caller(err)
			}
		}()
	}
	for i := 0; i < workers; i++ {
		<-done
	}
}

// ---------------------------------------------------------------------------------------------
// 10. The branches a happy-path test never reaches: foreign inputs to the decorators, and frames
//     that carry nothing.
// ---------------------------------------------------------------------------------------------

// withCaller is a foreign error that already carries caller data, which Trace and Track are
// documented to adopt rather than discard.
type withCaller struct{ inner error }

func (w *withCaller) Error() string { return w.inner.Error() }
func (w *withCaller) Unwrap() error { return w.inner }

func TestDecoratorsAcceptForeignErrors(t *testing.T) {
	cases := map[string]error{
		"plain stdlib error": sentinel,
		"foreign type":       &custom{msg: "foreign"},
		"foreign that wraps": &withCaller{inner: sentinel},
		"fmt wrapped":        fmt.Errorf("f: %w", sentinel),
		"joined":             std_errors.Join(other, sentinel),
	}
	for name, input := range cases {
		t.Run("Trace/"+name, func(t *testing.T) {
			out := errors.Trace(input)
			if nil == out {
				t.Fatal("Trace returned nil for a non-nil error")
			}
			if "" == out.Error() {
				t.Error("Trace lost the message")
			}
			if !std_errors.Is(out, sentinel) && !std_errors.Is(input, sentinel) {
				return // the input genuinely has no sentinel
			}
			if !std_errors.Is(out, sentinel) {
				t.Error("Trace severed the chain")
			}
		})
		t.Run("Track/"+name, func(t *testing.T) {
			out := errors.Track(input)
			if nil == out {
				t.Fatal("Track returned nil for a non-nil error")
			}
			if "" == out.Error() {
				t.Error("Track lost the message")
			}
			if std_errors.Is(input, sentinel) && !std_errors.Is(out, sentinel) {
				t.Error("Track severed the chain")
			}
		})
	}
}

// TestFrameWithNoMessageIsTransparent: a frame can exist purely to add caller data. It must not
// contribute an empty segment, a stray separator, or swallow what it wraps.
func TestFrameWithNoMessageIsTransparent(t *testing.T) {
	if got := errors.WrapE(nil, nil).Error(); "" != got {
		t.Errorf("an empty frame rendered %q, want the empty string", got)
	}
	if got, want := errors.WrapE(sentinel, nil).Error(), sentinel.Error(); want != got {
		t.Errorf("a frame with no annotation rendered %q, want %q", got, want)
	}
	traced := errors.Trace(errors.Wrap(sentinel, "inner"))
	if got := traced.Error(); strings.HasPrefix(got, ":") || strings.Contains(got, ": :") {
		t.Errorf("Trace introduced an empty segment: %q", got)
	}
	if got, want := traced.Error(), "inner: sentinel"; want != got {
		t.Errorf("Trace changed the rendered message: %q, want %q", got, want)
	}
}

// TestErrorfAndNewAreDistinct pins the one difference between them, since it is the whole reason
// both exist: Errorf interprets, New does not.
func TestErrorfAndNewAreDistinct(t *testing.T) {
	const raw = "50% of %d"
	if got := errors.New(raw).Error(); raw != got {
		t.Errorf("New must not interpret: got %q", got)
	}
	if got, want := errors.Errorf("50%% of %d", 8).Error(), "50% of 8"; want != got {
		t.Errorf("Errorf must interpret: got %q, want %q", got, want)
	}
}
