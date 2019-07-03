package errors_test

import (
	"fmt"

	errs "github.com/bdlm/errors"
)

var errEOF = fmt.Errorf("read: end of input")
var otherErr = fmt.Errorf("some other process failed")

func ExampleNew() {
	err := errs.New("this is an error message")

	fmt.Println(err)
	// Output: this is an error message
}

func ExampleFormat_string() {
	err := loadConfig()
	fmt.Println(err)
	// Output: service configuration could not be loaded
}

func ExampleFormat_string_formatted() {
	err := loadConfig()
	fmt.Printf("% v", err)
	// Output: service configuration could not be loaded
}

func ExampleFormat_string_detail() {
	err := loadConfig()
	fmt.Printf("%-v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig);
}

func ExampleFormat_string_detail_formatted() {
	err := loadConfig()
	fmt.Printf("% -v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig);
}

func ExampleFormat_string_trace() {
	err := loadConfig()
	fmt.Printf("%+v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig); could not decode configuration data - #1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig); could not read configuration file - #2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig);
}

func ExampleFormat_string_trace_formatted() {
	err := loadConfig()
	fmt.Printf("% +v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig);
	//could not decode configuration data - #1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig);
	//could not read configuration file - #2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig);
}

func ExampleFormat_json() {
	err := loadConfig()
	fmt.Printf("%#v", err)
	// Output: [{"error":"service configuration could not be loaded"},{"error":"could not decode configuration data"},{"error":"could not read configuration file"}]
}

func ExampleFormat_json_formatted() {
	err := loadConfig()
	fmt.Printf("% #v", err)
	// Output: [
	//     {
	//         "error": "service configuration could not be loaded"
	//     },
	//     {
	//         "error": "could not decode configuration data"
	//     },
	//     {
	//         "error": "could not read configuration file"
	//     }
	// ]
}

func ExampleFormat_json_detail() {
	err := loadConfig()
	fmt.Printf("%#-v", err)
	// Output: [{"caller":"#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)","error":"service configuration could not be loaded"},{"caller":"#1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig)","error":"could not decode configuration data"},{"caller":"#2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig)","error":"could not read configuration file"}]
}

func ExampleFormat_json_detail_formatted() {
	err := loadConfig()
	fmt.Printf("% #-v", err)
	// Output: [
	//     {
	//         "caller": "#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)",
	//         "error": "service configuration could not be loaded"
	//     },
	//     {
	//         "caller": "#1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig)",
	//         "error": "could not decode configuration data"
	//     },
	//     {
	//         "caller": "#2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig)",
	//         "error": "could not read configuration file"
	//     }
	// ]
}

func ExampleFormat_json_trace() {
	err := loadConfig()
	fmt.Printf("%#+v", err)
	// Output: [{"error":"service configuration could not be loaded","trace":["#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)","#1 examples_test.go:105 (github.com/bdlm/errors_test.ExampleFormat_json_trace)","#2 example.go:121 (testing.runExample)","#3 example.go:45 (testing.runExamples)","#4 testing.go:1073 (testing.(*M).Run)","#5 _testmain.go:74 (main.main)","#6 proc.go:200 (runtime.main)","#7 asm_amd64.s:1337 (runtime.goexit)"]},{"error":"could not decode configuration data","trace":["#0 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig)","#1 mocks_test.go:15 (github.com/bdlm/errors_test.loadConfig)","#2 examples_test.go:105 (github.com/bdlm/errors_test.ExampleFormat_json_trace)","#3 example.go:121 (testing.runExample)","#4 example.go:45 (testing.runExamples)","#5 testing.go:1073 (testing.(*M).Run)","#6 _testmain.go:74 (main.main)","#7 proc.go:200 (runtime.main)","#8 asm_amd64.s:1337 (runtime.goexit)"]},{"error":"could not read configuration file","trace":["#0 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig)","#1 mocks_test.go:20 (github.com/bdlm/errors_test.decodeConfig)","#2 mocks_test.go:15 (github.com/bdlm/errors_test.loadConfig)","#3 examples_test.go:105 (github.com/bdlm/errors_test.ExampleFormat_json_trace)","#4 example.go:121 (testing.runExample)","#5 example.go:45 (testing.runExamples)","#6 testing.go:1073 (testing.(*M).Run)","#7 _testmain.go:74 (main.main)","#8 proc.go:200 (runtime.main)","#9 asm_amd64.s:1337 (runtime.goexit)"]}]
}

func ExampleFormat_json_trace_formatted() {
	err := loadConfig()
	fmt.Printf("% #+v", err)
	// Output: [
	//     {
	//         "error": "service configuration could not be loaded",
	//         "trace": [
	//             "#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)",
	//             "#1 examples_test.go:111 (github.com/bdlm/errors_test.ExampleFormat_json_trace_debug)",
	//             "#2 example.go:121 (testing.runExample)",
	//             "#3 example.go:45 (testing.runExamples)",
	//             "#4 testing.go:1073 (testing.(*M).Run)",
	//             "#5 _testmain.go:74 (main.main)",
	//             "#6 proc.go:200 (runtime.main)",
	//             "#7 asm_amd64.s:1337 (runtime.goexit)"
	//         ]
	//     },
	//     {
	//         "error": "could not decode configuration data",
	//         "trace": [
	//             "#0 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig)",
	//             "#1 mocks_test.go:15 (github.com/bdlm/errors_test.loadConfig)",
	//             "#2 examples_test.go:111 (github.com/bdlm/errors_test.ExampleFormat_json_trace_debug)",
	//             "#3 example.go:121 (testing.runExample)",
	//             "#4 example.go:45 (testing.runExamples)",
	//             "#5 testing.go:1073 (testing.(*M).Run)",
	//             "#6 _testmain.go:74 (main.main)",
	//             "#7 proc.go:200 (runtime.main)",
	//             "#8 asm_amd64.s:1337 (runtime.goexit)"
	//         ]
	//     },
	//     {
	//         "error": "could not read configuration file",
	//         "trace": [
	//             "#0 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig)",
	//             "#1 mocks_test.go:20 (github.com/bdlm/errors_test.decodeConfig)",
	//             "#2 mocks_test.go:15 (github.com/bdlm/errors_test.loadConfig)",
	//             "#3 examples_test.go:111 (github.com/bdlm/errors_test.ExampleFormat_json_trace_debug)",
	//             "#4 example.go:121 (testing.runExample)",
	//             "#5 example.go:45 (testing.runExamples)",
	//             "#6 testing.go:1073 (testing.(*M).Run)",
	//             "#7 _testmain.go:74 (main.main)",
	//             "#8 proc.go:200 (runtime.main)",
	//             "#9 asm_amd64.s:1337 (runtime.goexit)"
	//         ]
	//     }
	// ]
}

func ExampleErr_WrapE() {
	// To add to an error with another error.
	err := loadConfig()
	if nil != err {
		retryErr := tryAgain()
		if nil != retryErr {
			err = errs.WrapE(err, retryErr)
		}
	}

	fmt.Println(err)
	// Output: retry failed
}

func ExampleErr_Wrap() {
	// To add to an error with another error.
	err := loadConfig()
	err = errs.WrapE(err, fmt.Errorf("loadConfig returned an error"))

	fmt.Println(err)
	// Output: loadConfig returned an error
}
