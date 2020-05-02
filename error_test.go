package errors_test

import (
	"encoding/json"
	go_errors "errors"
	"fmt"
	"testing"

	"github.com/bdlm/errors/v2"

	"github.com/stretchr/testify/assert"
)

func TestMarshaller(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = errors.New("test 1")
	byts, err := json.Marshal(err)
	assert.Equal(nil, err, "err is not nil")
	assert.Equal(
		"[{\"caller\":\"#0 error_test.go:18 (github.com/bdlm/errors/v2_test.TestMarshaller)\",\"error\":\"test 1\"}]",
		string(byts),
		"JSON did not encode properly",
	)
}

func TestTrace(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Trace(err)

	assert.Equal(32, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors/v2_test.TestTrace", err.Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Trace(nil)
	assert.Equal(nil, err, "err is not nil")
}

func TestTrack(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Track(err)

	assert.Equal(44, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors/v2_test.TestTrack", err.Caller().Func(), "caller did not reflect the correct function name")

	assert.Equal(45, errors.Unwrap(err).Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors/v2_test.TestTrack", errors.Unwrap(err).Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Track(nil)
	assert.Equal(nil, err, "err is not nil")
}

func TestCallerString(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	caller := errors.Caller(err)

	assert.Equal("github.com/bdlm/errors/v2_test.TestCallerString:60", caller.(fmt.Stringer).String(), "caller.String() did not return the correct output")
}

func TestHas(t *testing.T) {
	var err error
	var testErr error
	assert := assert.New(t)

	testErr = errors.New("test 1")
	err = testErr
	assert.True(errors.Has(err, testErr), "err does not contain testErr")

	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Has(err, testErr), "err does not contain testErr")

	testErr = go_errors.New("test 1")
	err = errors.Wrap(testErr, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Has(err, testErr), "err does not contain testErr")

	testErr = go_errors.New("test 1")
	typedErr := errors.Wrap(testErr, "test 2")
	typedErr = errors.Wrap(typedErr, "test 3")
	assert.True(errors.Has(typedErr, testErr), "typedErr does not contain testErr")

	testErr = fmt.Errorf("test 1")
	err = errors.Wrap(testErr, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Has(err, testErr), "err does not contain testErr")

	assert.False(errors.Has(nil, testErr), "nil did not evaluate to false")
}

func TestEHas(t *testing.T) {
	assert := assert.New(t)

	testErr := errors.New("test 1")
	err := testErr
	assert.True(err.Has(testErr), "err does not contain testErr")

	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(err.Has(testErr), "err does not contain testErr")

	testGoErr := go_errors.New("test 1")
	err = errors.Wrap(testGoErr, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(err.Has(testGoErr), "err does not contain testGoErr")

	testGoErr = go_errors.New("test 1")
	typedErr := errors.Wrap(testGoErr, "test 2")
	typedErr = errors.Wrap(typedErr, "test 3")
	assert.True(typedErr.Has(testGoErr), "typedErr does not contain testGoErr")

	assert.False(err.Has(nil), "nil did not evaluate to false")
}

func TestIs(t *testing.T) {
	var err error
	var testErr error
	assert := assert.New(t)

	testErr = errors.New("test 1")
	err = testErr
	assert.True(errors.Is(err, testErr), "err is not testErr")

	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.False(errors.Is(err, testErr), "err is testErr")

	testErr = go_errors.New("test 1")
	err = errors.Wrap(testErr, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.False(errors.Is(err, testErr), "err is testErr")

	typedTestErr := go_errors.New("test 1")
	typedErr := typedTestErr
	assert.True(errors.Is(typedErr, typedTestErr), "typedErr is not typedTestErr")

	assert.False(errors.Is(nil, testErr), "nil did not evaluate to false")
}

func TestEIs(t *testing.T) {
	assert := assert.New(t)

	testErr := errors.New("test 1")
	err := testErr
	assert.True(err.Is(testErr), "err is not testErr")

	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.False(err.Is(testErr), "err is testErr")

	testGoErr := go_errors.New("test 1")
	err = errors.Wrap(testGoErr, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.False(errors.Is(err, testGoErr), "err is testGoErr")

	assert.False(err.Is(nil), "nil did not evaluate to false")
}

func TestErrorf(t *testing.T) {
	assert := assert.New(t)

	err := errors.Errorf("test %d", 1)
	assert.Equal("test 1", err.Error(), "err is not 'test 1'")
}

func TestCaller(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = errors.New("test 1")
	caller := errors.Caller(err)
	assert.NotEqual(nil, caller, "caller is nil")

	err = nil
	caller = errors.Caller(err)
	assert.Equal(nil, caller, "caller is not nil")
}

func TestUnwrap(t *testing.T) {
	var err error
	var testErr error
	assert := assert.New(t)

	testErr = errors.New("test 1")
	err = testErr
	assert.True(errors.Is(err, testErr), "err is not testErr")

	err = errors.Wrap(err, "test 2")
	assert.False(errors.Is(err, testErr), "err is testErr")

	err = errors.Unwrap(err)
	assert.True(errors.Is(err, testErr), "err is not testErr")
}

func TestEUnwrap(t *testing.T) {
	assert := assert.New(t)

	testErr := errors.New("test 1")
	err := testErr
	assert.True(err.Is(testErr), "err is not testErr")

	err = errors.Wrap(err, "test 2")
	assert.False(err.Is(testErr), "err is testErr")

	err = err.Unwrap()
	assert.True(err.Is(testErr), "err is not testErr")
}
