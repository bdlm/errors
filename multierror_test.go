package errors_test

import (
	std_errors "errors"
	"fmt"
	"testing"

	"github.com/bdlm/errors/v2"
)

var (
	errWanted = std_errors.New("the sentinel we are looking for")
	errOther  = std_errors.New("an unrelated cause")
)

// customErr uses a VALUE receiver so that *customErr is a pointer whose element type implements
// error, which is the shape bdlm's As(err, test error) requires of its target.
type customErr struct{ msg string }

func (c customErr) Error() string { return c.msg }

// TestIsTraversesMultiErrors covers errors that wrap SEVERAL causes -- errors.Join and fmt.Errorf
// with more than one %w. They implement Unwrap() []error, which the single-error Unwrap cannot
// express, so before this was handled the walk terminated at the join and Is answered false for a
// sentinel sitting directly inside it.
//
// This is not a hypothetical shape: golang-jwt joins its sentinels, so a parser rejection is a
// *fmt.wrapErrors and NO jwt error could be matched -- callers saw a credential rejection fall
// through to a 5xx.
func TestIsTraversesMultiErrors(t *testing.T) {
	for name, err := range map[string]error{
		"join, sentinel first":     std_errors.Join(errWanted, errOther),
		"join, sentinel last":      std_errors.Join(errOther, errWanted),
		"fmt.Errorf two %w":        fmt.Errorf("a: %w, b: %w", errOther, errWanted),
		"wrapped join":             errors.Wrap(std_errors.Join(errOther, errWanted), "context"),
		"twice-wrapped join":       errors.Wrap(errors.Wrap(std_errors.Join(errOther, errWanted), "inner"), "outer"),
		"join of a wrapped2 chain": std_errors.Join(errOther, errors.Wrap(errWanted, "context")),
		"nested joins":             std_errors.Join(errOther, std_errors.Join(errOther, errWanted)),
	} {
		if !errors.Is(err, errWanted) {
			t.Errorf("%s: Is = false, want true", name)
		}
		// The standard library is the reference implementation; agree with it.
		if std_errors.Is(err, errWanted) != errors.Is(err, errWanted) {
			t.Errorf("%s: disagrees with std_errors.Is", name)
		}
	}
}

// TestIsSearchesEveryBranch: a match in a later branch must still be found, and a sentinel that is
// genuinely absent must still report false -- the search must not become indiscriminate.
func TestIsSearchesEveryBranch(t *testing.T) {
	absent := std_errors.New("never wrapped")
	for name, err := range map[string]error{
		"join without it":    std_errors.Join(errOther, errWanted),
		"wrapped join":       errors.Wrap(std_errors.Join(errOther, errWanted), "context"),
		"deep nested branch": std_errors.Join(errOther, std_errors.Join(errOther, errors.Wrap(errWanted, "deep"))),
	} {
		if errors.Is(err, absent) {
			t.Errorf("%s: Is(absent) = true, want false", name)
		}
	}
	if errors.Is(nil, errWanted) || errors.Is(errWanted, nil) {
		t.Error("a nil operand must never match")
	}
}

// TestAsTraversesMultiErrors: As had the same defect from the same cause.
func TestAsTraversesMultiErrors(t *testing.T) {
	cause := customErr{msg: "concrete cause"}
	for name, err := range map[string]error{
		"join":         std_errors.Join(errOther, cause),
		"wrapped join": errors.Wrap(std_errors.Join(errOther, cause), "context"),
		"nested joins": std_errors.Join(errOther, std_errors.Join(errOther, cause)),
	} {
		var reference customErr
		if !std_errors.As(err, &reference) {
			t.Fatalf("%s: premise broken, std_errors.As cannot find it either", name)
		}

		found := &customErr{}
		if got := errors.As(err, found); nil == got || cause.msg != found.msg {
			t.Errorf("%s: As did not reach the concrete cause (got %v, found %q)", name, got, found.msg)
		}
	}
}
