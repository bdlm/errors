package errors_test

import (
	"encoding/json"
	go_errors "errors"
	"fmt"
	"testing"

	"github.com/bdlm/errors/v2"
	std_errors "github.com/bdlm/std/v2/errors"

	"github.com/stretchr/testify/assert"
)

func TestMarshaller(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = errors.New("test 1")
	byts, err := json.Marshal(err)
	assert.Equal(nil, err, "err is not nil")
	assert.Equal(
		"[{\"caller\":\"#0 error_test.go:19 (github.com/bdlm/errors/v2_test.TestMarshaller)\",\"error\":\"test 1\"}]",
		string(byts),
		"JSON did not encode properly",
	)
}

func TestTrace(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Trace(err)

	assert.Equal(33, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors/v2_test.TestTrace", err.Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Trace(nil)
	assert.True(nil == err, "err is not nil")
}

func TestTrack(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Track(err)

	assert.Equal(45, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors/v2_test.TestTrack", err.Caller().Func(), "caller did not reflect the correct function name")

	assert.Equal(46, errors.Unwrap(err).(std_errors.Caller).Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors/v2_test.TestTrack", errors.Unwrap(err).(std_errors.Caller).Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Track(nil)
	assert.True(nil == err, "err is not nil")
}

func TestCallerString(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	caller := errors.Caller(err)

	assert.Equal("github.com/bdlm/errors/v2_test.TestCallerString:61", caller.(fmt.Stringer).String(), "caller.String() did not return the correct output")
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
	assert.True(errors.Is(err, testErr), "err is testErr")

	testErr = go_errors.New("test 1")
	err = errors.Wrap(testErr, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Is(err, testErr), "err is testErr")

	typedTestErr := go_errors.New("test 1")
	typedErr := typedTestErr
	assert.True(errors.Is(typedErr, typedTestErr), "typedErr is not typedTestErr")

	testErr = errors.New("test 1")
	err = errors.New("test 1")
	assert.False(errors.Is(err, testErr), "err is not testErr")

	testErr = errors.New("test 1")
	err = errors.New("test 1")
	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.False(errors.Is(err, testErr), "err is not testErr")

	testErr = errors.New("test 1")
	err = errors.New("test 1")
	err = errors.WrapE(err, testErr)
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Is(err, testErr), "err is not testErr")

	testErr = go_errors.New("test 1")
	err = go_errors.New("test 1")
	assert.False(errors.Is(err, testErr), "err is not testErr")

	testErr = go_errors.New("test 1")
	err = go_errors.New("test 1")
	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.False(errors.Is(err, testErr), "err is not testErr")

	testErr = go_errors.New("test 1")
	err = go_errors.New("test 1")
	err = errors.WrapE(err, testErr)
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Is(err, testErr), "err is not testErr")

	assert.False(errors.Is(nil, testErr), "nil did not evaluate to false")
}

func TestEIs(t *testing.T) {
	assert := assert.New(t)

	testErr := errors.New("test 1")
	err := testErr
	assert.True(err.Is(testErr), "err is not testErr")

	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(err.Is(testErr), "err is testErr")

	testGoErr := go_errors.New("test 1")
	err = errors.Wrap(testGoErr, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.True(err.Is(testGoErr), "err is testGoErr")

	testErr = errors.New("test 1")
	err = errors.New("test 1")
	assert.False(err.Is(testErr), "err is not testErr")

	testErr = errors.New("test 1")
	err = errors.New("test 1")
	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	assert.False(err.Is(testErr), "err is testErr")

	testGoErr = go_errors.New("test 1")
	err = errors.New("test 1")
	err = errors.WrapE(err, testGoErr)
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Is(err, testGoErr), "err is testGoErr")

	testGoErr = go_errors.New("test 1")
	testGoErr2 := go_errors.New("test 1")
	err = errors.WrapE(testGoErr2, testGoErr)
	err = errors.Wrap(err, "test 3")
	assert.True(errors.Is(err, testGoErr), "err is testGoErr")

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

	err = errors.New("test 1")
	assert.Nil(errors.Unwrap(err))

	testErr = errors.New("test 1")
	err = testErr
	assert.True(errors.Is(err, testErr), "err is not testErr")

	err = errors.Wrap(err, "test 2")
	assert.True(errors.Is(err, testErr), "err is testErr")

	err = errors.Unwrap(err)
	assert.True(errors.Is(err, testErr), "err is not testErr")
}

func TestEUnwrap(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	assert.Nil(err.Unwrap())

	testErr := errors.New("test 1")
	err = testErr
	assert.True(err.Is(testErr), "err is not testErr")

	err = errors.Wrap(err, "test 2")
	assert.True(err.Is(testErr), "err is testErr")

	err = err.Unwrap().(*errors.E)
	assert.True(err.Is(testErr), "err is not testErr")
}
