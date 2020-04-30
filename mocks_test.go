package errors_test

/*
WARNING - changing the line numbers in this file will break the
examples.
*/

import (
	"fmt"

	errs "github.com/bdlm/errors/v2"
)

func loadConfig() error {
	err := decodeConfig()
	return errs.Wrap(err, "service configuration could not be loaded")
}

func decodeConfig() error {
	err := readConfig()
	return errs.Wrap(err, "could not decode configuration data")
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return errs.Wrap(err, "could not read configuration file")
}

func someWork() error {
	return fmt.Errorf("failed to do work")
}

func tryAgain() error {
	return errs.Wrap(loadConfig(), "retry failed")
}
