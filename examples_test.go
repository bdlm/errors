package errors_test

import (
	"encoding/json"
	"fmt"

	errors "github.com/bdlm/errors"
)

var errEOF = fmt.Errorf("read: end of input")
var otherErr = fmt.Errorf("some other process failed")

func ExampleNew() {
	err := errors.New("this is an error message")

	fmt.Println(err)
	// Output: this is an error message
}

func ExampleE_Format_string() {
	err := loadConfig()
	fmt.Println(err)
	// Output: service configuration could not be loaded
}

func ExampleE_Format_stringPreformat() {
	err := loadConfig()
	fmt.Printf("% v", err)
	// Output: service configuration could not be loaded
}

func ExampleE_Format_stringDetail() {
	err := loadConfig()
	fmt.Printf("%-v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig);
}

func ExampleE_Format_stringTrace() {
	err := loadConfig()
	fmt.Printf("%+v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig); could not decode configuration data - #1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig); could not read configuration file - #2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig);
}

func ExampleE_Format_stringDetailPreformat() {
	err := loadConfig()
	fmt.Printf("% -v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig);
}
func ExampleE_Format_stringTracePreformat() {
	err := loadConfig()
	fmt.Printf("% +v", err)
	// Output: service configuration could not be loaded - #0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig);
	// could not decode configuration data - #1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig);
	// could not read configuration file - #2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig);
}

func ExampleE_Format_json() {
	err := loadConfig()
	fmt.Printf("%#v", err)
	// Output: [{"error":"service configuration could not be loaded"}]
}

func ExampleE_Format_jsonPreformat() {
	err := loadConfig()
	fmt.Printf("% #v", err)
	// Output: [
	//     {
	//         "error": "service configuration could not be loaded"
	//     }
	// ]
}

func ExampleE_Format_jsonDetail() {
	err := loadConfig()
	fmt.Printf("%#-v", err)
	// Output: [{"caller":"#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)","error":"service configuration could not be loaded"}]
}

func ExampleE_Format_jsonDetailPreformat() {
	err := loadConfig()
	fmt.Printf("% #-v", err)
	// Output: [
	//     {
	//         "caller": "#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)",
	//         "error": "service configuration could not be loaded"
	//     }
	// ]
}

func ExampleE_Format_jsonTrace() {
	err := loadConfig()
	fmt.Printf("%#+v", err)
	// Output: [{"caller":"#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)","error":"service configuration could not be loaded"},{"caller":"#1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig)","error":"could not decode configuration data"},{"caller":"#2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig)","error":"could not read configuration file"}]
}

func ExampleE_Format_jsonTracePreformat() {
	err := loadConfig()
	fmt.Printf("% #+v", err)
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

func ExampleE_MarshalJSON_marshal() {
	err := loadConfig()
	jsn, _ := json.Marshal(err)

	fmt.Println(string(jsn))
	// Output: [{"caller":"#0 mocks_test.go:16 (github.com/bdlm/errors_test.loadConfig)","error":"service configuration could not be loaded"},{"caller":"#1 mocks_test.go:21 (github.com/bdlm/errors_test.decodeConfig)","error":"could not decode configuration data"},{"caller":"#2 mocks_test.go:26 (github.com/bdlm/errors_test.readConfig)","error":"could not read configuration file"}]
}

func ExampleE_MarshalJSON_marshalIndent() {
	err := loadConfig()
	jsn, _ := json.MarshalIndent(err, "", "    ")

	fmt.Println(string(jsn))
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

func ExampleUnwrap() {
	err1 := errors.New("error 1")
	err2 := errors.Wrap(err1, "error 2")
	err := errors.Unwrap(err2)

	fmt.Println(err)
	// Output: error 1
}

func ExampleWrap() {
	// Wrap an error with additional metadata.
	err := loadConfig()
	err = errors.Wrap(err, "loadConfig returned an error")

	fmt.Println(err)
	// Output: loadConfig returned an error
}

func ExampleWrapE() {
	// Wrap an error with another error.
	err := loadConfig()
	if nil != err {
		retryErr := tryAgain()
		if nil != retryErr {
			err = errors.WrapE(err, retryErr)
		}
	}

	fmt.Println(err)
	// Output: retry failed
}
