package errors_test

import (
	"fmt"
	"testing"

	errs "github.com/mkenney/go-errors"
)

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return errs.Wrap(err, "could not read configuration file", errs.ErrEOF)
}

func decodeConfig() error {
	err := readConfig()
	return errs.Wrap(err, "could not decode configuration data", errs.ErrInvalidJSON)
}

func loadConfig() error {
	err := decodeConfig()
	return errs.Wrap(err, "service configuration could not be loaded", errs.ErrFatal)
}

func TestOutput(t *testing.T) {
	testErr := loadConfig()
	fmt.Printf("\n-------------- %%+v --------------\n\n%+v\n\n", testErr)
	fmt.Printf("\n-------------- %%#v --------------\n\n%#v\n\n", testErr)
	fmt.Printf("\n-------------- %%v ---------------\n\n%v\n\n", testErr)
	fmt.Printf("\n-------------- %%s ---------------\n\n%s\n\n", testErr)
}
