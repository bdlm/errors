package errors

/*
Msg defines a single error message.
*/
type Msg struct {
	err    error
	caller Caller
	code   Code
	msg    string
}

/*
String implements the Stringer interface
*/
func (msg Msg) String() string {
	return msg.err.Error()
}

/*
Error implements the error interface
*/
func (msg Msg) Error() string {
	return msg.String()
}
