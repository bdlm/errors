package errors_test

/*
WARNING - changing the line numbers in this file will break the
examples.
*/

import (
	"fmt"

	errs "github.com/bdlm/errors"
)

func loadConfig() error {
	err := decodeConfig()
	return errs.Wrap(err, "service configuration could not be loaded", errs.ErrFatal)
}

func decodeConfig() error {
	err := readConfig()
	return errs.Wrap(err, "could not decode configuration data", errs.ErrInvalidJSON)
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return errs.Wrap(err, "could not read configuration file", errs.ErrEOF)
}

func someWork() error {
	return fmt.Errorf("failed to do work")
}
