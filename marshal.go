package errors

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

// MarshalJSON implements the json.Marshaller interface.
func (e *E) MarshalJSON() ([]byte, error) {
	var lastE, nextE error
	var key int
	jsonData := []map[string]interface{}{}

	for key, nextE = range list(e) {
		data := map[string]interface{}{}
		err, ok := nextE.(*E)
		if ok {
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
		lastE = err.prev
		if "" != nextE.Error() {
			data["error"] = err.Error()
		}
		jsonData = append(jsonData, data)
	}

	if nil != lastE {
		data := map[string]interface{}{}
		err, ok := lastE.(*E)
		if ok {
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
		if "" != lastE.Error() {
			data["error"] = lastE.Error()
		}
		jsonData = append(jsonData, data)
	}

	return json.Marshal(jsonData)
}
