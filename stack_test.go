package errors_test

import (
	"testing"

	"github.com/bdlm/errors"
	"github.com/bdlm/log"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{
		//DisableTTY: true,
		ForceTTY: true,
		FieldMap: log.FieldMap{
			"data": "_",
		},
		//EnableTrace: true,
	})
}

func TestErrors(t *testing.T) {
	assert := assert.New(t)

	e := errors.New("error 1")
	e = errors.Wrap(e, "error 2")
	e = errors.Wrap(e, "error 3")

	log.WithError(e).Info("log test")

	//byts, _ := json.Marshal(e)
	//assert.Equal(2, 1, string(byts))

	assert.Equal(2, 1, "%#v", e)
}
