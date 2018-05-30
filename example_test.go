package errors_test

import (
	"fmt"
	"testing"

	errs "github.com/mkenney/go-errors"
)

func TestOutput(t *testing.T) {
	testErr := loadConfig()
	fmt.Printf("\n-------------- %%+v --------------\n\n%+v\n\n", testErr)
	fmt.Printf("\n-------------- %%#v --------------\n\n%#v\n\n", testErr)
	fmt.Printf("\n-------------- %%v ---------------\n\n%v\n\n", testErr)
	fmt.Printf("\n-------------- %%s ---------------\n\n%s\n\n", testErr)
}

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

func TestDoSteps(t *testing.T) {
	var errStack errs.Err

	err := doStep1()
	if nil != err {
		errStack = errStack.With(err)
	}

	err = doStep2()
	if nil != err {
		errStack = errStack.With(err)
	}

	err = doStep3()
	if nil != err {
		errStack = errStack.With(err)
	}

	err = doStep4()
	if nil != err {
		errStack = errStack.With(err)
	}

	fmt.Printf("\n-------------- %%+v --------------\n\n%+v\n\n", errStack)
	fmt.Printf("\n-------------- %%#v --------------\n\n%#v\n\n", errStack)
	fmt.Printf("\n-------------- %%v ---------------\n\n%v\n\n", errStack)
	fmt.Printf("\n-------------- %%s ---------------\n\n%s\n\n", errStack)
}

func doStep1() error {
	return fmt.Errorf("thing 1 failed")
}
func doStep2() error {
	return fmt.Errorf("thing 2 failed")
}
func doStep3() error {
	return fmt.Errorf("thing 3 failed")
}
func doStep4() error {
	return fmt.Errorf("thing 4 failed")
}
