package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type retrieveKeyController struct {
	cluster ClusterProvider
	logger  *zap.Logger
}

func (rkc retrieveKeyController) RetrieveKey(w http.ResponseWriter, r *http.Request) {

}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	versionMessage := fmt.Sprintln("Ohara key retrieving service")
	versionMessage += fmt.Sprintln("Consolidates the truth with no fear")
	versionMessage += fmt.Sprintln("Server version: " + Version)
	versionMessage += fmt.Sprintln("API version: " + SupportedAPI)
	writeResponse(w, http.StatusOK, versionMessage)
}

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
