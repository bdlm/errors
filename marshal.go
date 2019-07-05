package errors

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

// MarshalJSON implements the json.Marshaller interface.
func (e E) MarshalJSON() ([]byte, error) {
	jsonData := []map[string]interface{}{}

	for a, b := range list(e) {
		err, ok := b.(E)
		if !ok {
			break
		}
		data := map[string]interface{}{}
		data["caller"] = fmt.Sprintf("#%d %s:%d (%s)",
			a,
			path.Base(err.Caller().File()),
			err.Caller().Line(),
			runtime.FuncForPC(err.Caller().Pc()).Name(),
		)
		if "" != err.Error() {
			data["error"] = err.Error()
		}
		jsonData = append(jsonData, data)
	}

	return json.Marshal(jsonData)
}
