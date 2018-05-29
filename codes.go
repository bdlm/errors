package errors

func init() {
	// Internal errors
	//Codes[ErrUnspecified] = Metadata{"Internal Server Error", "error code unspecified", 0}
	Codes[ErrUnknown] = Metadata{"Internal Server Error", "an unknown error occurred", 500}
	Codes[ErrFatal] = Metadata{"Internal Server Error", "a fatal error occurred", 500}

	// I/O errors
	Codes[ErrEOF] = Metadata{"End of input", "unexpected EOF", 400}

	// Encoding errors
	Codes[ErrInvalidJSON] = Metadata{"Invalid JSON Data", "invalid JSON data could not be decoded", 400}

	// Server errors
	Codes[ErrInvalidHTTPMethod] = Metadata{"Invalid HTTP Method", "an invalid HTTP method was requested", 400}
}

/*
Code defines an error code type.
*/
type Code int

/*
Metadata contains metadata that can be associated with an error code.
*/
type Metadata struct {
	// External (user) facing error text.
	External string
	// Internal only (logs) error text.
	Internal string
	// HTTP status that should be used for the associated error code.
	HTTPStatus int
}

/*
Codes contains a map of error codes to metadata
*/
var Codes = map[Code]Metadata{}

/*
Internal errors
*/
const (
	// ErrUnspecified - 0: The error code was unspecified.
	ErrUnspecified Code = iota
	// ErrUnknown - 1: An unknown error occurred.
	ErrUnknown
	// ErrFatal - 2: An fatal error occurred.
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
