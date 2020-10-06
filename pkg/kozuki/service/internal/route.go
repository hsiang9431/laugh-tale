package internal

import (
	"laugh-tale/pkg/kozuki/service/internal/datastore"
	"laugh-tale/pkg/kozuki/service/internal/resource"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	Version      = "unknown"
	SupportedAPI = "v1"
)

func ClientRouter(l *zap.Logger, ks datastore.KeyStore) (*mux.Router, error) {
	resource.Version = Version
	resource.SupportedAPI = SupportedAPI
	kc := resource.KeyController{
		KeyStore: ks,
		Logger:   l,
	}
	router := mux.NewRouter()
	router.HandleFunc("/v1/key/create", kc.CreateEncryptionKey).Methods("POST")
	router.HandleFunc("/v1/key/bind", kc.BindKeyToImageID).Methods("POST")
	router.HandleFunc("/v1/key/{image_id}", kc.RetrieveKeyByImageID).Methods("GET")
	return router, nil
}

func CRUDRouter(l *zap.Logger, ks datastore.KeyStore) (*mux.Router, error) {
	resource.Version = Version
	resource.SupportedAPI = SupportedAPI
	kcCRUD := resource.KeyControllerCRUD{
		KeyStore: ks,
		Logger:   l,
	}
	router := mux.NewRouter()
	router.HandleFunc("/key", kcCRUD.Create).Methods("POST")
	router.HandleFunc("/key", kcCRUD.Retrieve).Methods("GET")
	router.HandleFunc("/key", kcCRUD.Update).Methods("PATCH")
	router.HandleFunc("/key", kcCRUD.Delete).Methods("DELETE")
	return router, nil
}
