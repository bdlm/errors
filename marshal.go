package errors

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

// MarshalJSON implements the json.Marshaller interface.
func (e *E) MarshalJSON() ([]byte, error) {
	if nil == e {
		return []byte("null"), nil
	}
	var lastE, nextE error
	var key int
	jsonData := []map[string]interface{}{}

	for key, nextE = range list(e) {
		data := map[string]interface{}{}
		err, ok := nextE.(*E)
		if ok && nil != err.Caller() {
			data["caller"] = fmt.Sprintf("#%d %s:%d (%s)",
				key,
				path.Base(err.Caller().File()),
				err.Caller().Line(),
				runtime.FuncForPC(err.Caller().Pc()).Name(),
			)
		} else {
			data["caller"] = fmt.Sprintf("#%d n/a",
				key,
			)
		}
		// Guarded on ok: list() includes ANY error implementing Unwrap() error, so a foreign
		// wrapper in the chain -- fmt.Errorf("%w") being the common one -- lands here with err
		// nil, and reading err.prev panicked. Marshalling an error must never panic; that is the
		// path a logger takes while already handling a failure.
		if ok {
			lastE = err.prev
		} else {
			lastE = nil
		}
		// frameMessage, not Error(): each entry is one frame, and Error() now carries the wrapped
		// chain -- so using it here would repeat the whole tail in every entry of the array.
		if "" != frameMessage(nextE) {
			data["error"] = frameMessage(nextE)
		}
		jsonData = append(jsonData, data)
	}

	if nil != lastE {
		data := map[string]interface{}{}
		err, ok := lastE.(*E)
		if ok && nil != err.Caller() {
			data["caller"] = fmt.Sprintf("#%d %s:%d (%s)",
				key+1,
				path.Base(err.Caller().File()),
				err.Caller().Line(),
				runtime.FuncForPC(err.Caller().Pc()).Name(),
			)
		} else {
			data["caller"] = fmt.Sprintf("#%d n/a",
				key+1,
			)
		}
		if "" != frameMessage(lastE) {
			data["error"] = frameMessage(lastE)
		}
		jsonData = append(jsonData, data)
	}

	return json.Marshal(jsonData)
}
