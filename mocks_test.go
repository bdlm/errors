package errors_test

/*
WARNING - changing the line numbers in this file will break the
examples.
*/

import (
	"fmt"

	errs "github.com/bdlm/errors"
	std "github.com/bdlm/std/error"
)

const (
	// Error codes below 1000 are reserved future use by the
	// "github.com/bdlm/errors" package.
	ConfigurationNotValid std.Code = iota + 1000
)

func init() {
	errs.Codes[ConfigurationNotValid] = errs.ErrCode{
		Ext:  "Configuration not valid",
		Int:  "the configuration is invalid",
		HTTP: 500,
	}
}

func loadConfig() error {
	err := decodeConfig()
	return errs.Wrap(err, ConfigurationNotValid, "service configuration could not be loaded")
}

func decodeConfig() error {
	err := readConfig()
	return errs.Wrap(err, errs.ErrInvalidJSON, "could not decode configuration data")
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return errs.Wrap(err, errs.ErrEOF, "could not read configuration file")
}

func someWork() error {
	return fmt.Errorf("failed to do work")
}
