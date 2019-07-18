package errors_test

import (
	"encoding/json"
	go_errors "errors"
	"testing"

	"github.com/bdlm/errors"
	"github.com/bdlm/log"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		//DisableTTY: true,
		//ForceTTY: true,
		FieldMap: log.FieldMap{
			"data": "_",
		},
		//EnableTrace: true,
	})
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

	assert.False(errors.Has(nil, testErr), "nil did not evaluate to false")
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

func TestTrace(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Trace(err)

	assert.Equal(81, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors_test.TestTrace", err.Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Trace(nil)
	assert.Equal(nil, err, "err is not nil")
}

func TestTrack(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Track(err)

	assert.Equal(93, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors_test.TestTrack", err.Caller().Func(), "caller did not reflect the correct function name")

	assert.Equal(94, errors.Unwrap(err).Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors_test.TestTrack", errors.Unwrap(err).Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Track(nil)
	assert.Equal(nil, err, "err is not nil")
}

func TestErrorf(t *testing.T) {
	assert := assert.New(t)

	err := errors.Errorf("test %d", 1)
	assert.Equal("test 1", err.Error(), "err is not 'test 1'")
}

func TestGetCaller(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = errors.New("test 1")
	caller := errors.GetCaller(err)
	assert.NotEqual(nil, caller, "caller is nil")

	err = nil
	caller = errors.GetCaller(err)
	assert.Equal(nil, caller, "caller is not nil")
}

func TestMarshaller(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = errors.New("test 1")
	byts, err := json.Marshal(err)
	assert.Equal(nil, err, "err is not nil")
	assert.Equal(
		"[{\"caller\":\"#0 error_test.go:130 (github.com/bdlm/errors_test.TestMarshaller)\",\"error\":\"test 1\"}]",
		string(byts),
		"JSON did not encode properly",
	)
}
