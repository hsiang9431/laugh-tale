package resource

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func writeResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func writeJSONResponse(w http.ResponseWriter, obj interface{}) error {
	if w == nil || obj == nil {
		return nil
	}
	jsonB, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "object marshal failed")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonB)
	return nil
}
