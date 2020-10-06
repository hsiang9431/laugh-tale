package resource

import (
	"fmt"
	"net/http"
)

var (
	Version      = "unknown"
	SupportedAPI = "v1"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	versionMessage := fmt.Sprintln("Kozuki key manage service")
	versionMessage += fmt.Sprintln("Serving with the lunar shine")
	versionMessage += fmt.Sprintln("Server version: " + Version)
	versionMessage += fmt.Sprintln("API version: " + SupportedAPI)
	writeResponse(w, http.StatusOK, versionMessage)
}
