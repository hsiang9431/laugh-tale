package internal

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	Version      = "unknown"
	SupportedAPI = "v1"
)

func Router(l *zap.Logger, cluster ClusterProvider) (*mux.Router, error) {
	rkc := retrieveKeyController{
		cluster: cluster,
		logger:  l,
	}
	router := mux.NewRouter()
	router.HandleFunc("/v1/key", rkc.RetrieveKey).Methods("GET")
	router.HandleFunc("/v1/version", versionHandler).Methods("GET")
	return router, nil
}
