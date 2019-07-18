package errors_test

import (
	go_errors "errors"
	"testing"

	"github.com/bdlm/errors"
	"github.com/bdlm/log"

	"github.com/stretchr/testify/assert"
)

var testError = errors.Errorf("Test error")

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

func TestEHas(t *testing.T) {
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

func TestEIs(t *testing.T) {
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

func TestETrace(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Trace(err)

	assert.Equal(83, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors_test.TestETrace", err.Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Trace(nil)
	assert.Equal(nil, err, "err is not nil")
}

func TestETrack(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("test 1")
	err = errors.Track(err)

	assert.Equal(95, err.Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors_test.TestETrack", err.Caller().Func(), "caller did not reflect the correct function name")

	assert.Equal(96, errors.Unwrap(err).Caller().Line(), "caller did not reflect the correct line number")
	assert.Equal("github.com/bdlm/errors_test.TestETrack", errors.Unwrap(err).Caller().Func(), "caller did not reflect the correct function name")

	err = errors.Track(nil)
	assert.Equal(nil, err, "err is not nil")
}
