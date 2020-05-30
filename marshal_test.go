package errors_test

import (
	"encoding/json"
	"testing"

	"github.com/bdlm/errors/v2"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = errors.New("test 1")
	err = errors.Wrap(err, "test 2")
	err = errors.Wrap(err, "test 3")
	byts, jsonerr := json.Marshal(err)

	assert.Nil(jsonerr, "jsonerr is not nil")
	assert.Equal(
		"[{\"caller\":\"#0 marshal_test.go:18 (github.com/bdlm/errors/v2_test.TestMarshalJSON)\",\"error\":\"test 3\"},{\"caller\":\"#1 marshal_test.go:17 (github.com/bdlm/errors/v2_test.TestMarshalJSON)\",\"error\":\"test 2\"},{\"caller\":\"#2 marshal_test.go:16 (github.com/bdlm/errors/v2_test.TestMarshalJSON)\",\"error\":\"test 1\"}]",
		string(byts),
		"JSON did not encode properly",
	)
}
