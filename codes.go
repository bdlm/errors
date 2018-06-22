package errors

func init() {
	// Internal errors
	//Codes[ErrUnspecified] = ErrCode{"Error Unspecified", "error code unspecified", 500}
	Codes[ErrUnknown] = ErrCode{Int: "unknown error", HTTP: 500}
	Codes[ErrFatal] = ErrCode{"Internal Server Error", "fatal error", 500}

	// I/O errors
	Codes[ErrEOF] = ErrCode{"End of input", "unexpected EOF", 400}

	// Encoding errors
	Codes[ErrInvalidJSON] = ErrCode{"Invalid JSON Data", "invalid JSON data could not be decoded", 400}

	// Server errors
	Codes[ErrInvalidHTTPMethod] = ErrCode{"Invalid HTTP Method", "an invalid HTTP method was requested", 400}
}

/*
Coder defines an interface for an error code.
*/
type Coder interface {
	// Internal only (logs) error text.
	Detail() string
	// HTTP status that should be used for the associated error code.
	HTTPStatus() int
	// External (user) facing error text.
	String() string
}

/*
Code defines an error code type.
*/
type Code int

/*
Codes contains a map of error codes to metadata
*/
var Codes = make(map[Code]Coder)

/*
ErrCode implements coder
*/
type ErrCode struct {
	// External (user) facing error text.
	Ext string
	// Internal only (logs) error text.
	Int string
	// HTTP status that should be used for the associated error code.
	HTTP int
}

/*
Detail returns the internal error message, if any.
*/
func (code ErrCode) Detail() string {
	return code.Int
}

/*
String implements stringer. String returns the external error message,
if any.
*/
func (code ErrCode) String() string {
	return code.Ext
}

/*
HTTPStatus returns the associated HTTP status code, if any. Otherwise,
returns 200.
*/
func (code ErrCode) HTTPStatus() int {
	if 0 == code.HTTP {
		return 200
	}
	return code.HTTP
}

/*
Internal errors
*/
const (
	// ErrUnknown - 0: An unknown error occurred.
	ErrUnknown Code = iota
	// ErrFatal - 1: An fatal error occurred.
	ErrFatal
)

/*
I/O errors
*/
const (
	// ErrEOF - 100: An invalid HTTP method was requested.
	ErrEOF Code = iota + 100
)

/*
Encoding errors
*/
const (
	// ErrInvalidJSON - 200: Invalid JSON data could not be decoded.
	ErrInvalidJSON Code = iota + 200
)

/*
Server errors
*/
const (
	// ErrInvalidHTTPMethod - 300: An invalid HTTP method was requested.
	ErrInvalidHTTPMethod Code = iota + 300
)
