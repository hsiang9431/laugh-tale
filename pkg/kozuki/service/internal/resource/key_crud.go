package resource

import (
	"encoding/json"
	"io/ioutil"
	"laugh-tale/pkg/kozuki/service/internal/datastore"
	"laugh-tale/pkg/kozuki/types"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// for crud operations on keys
type KeyControllerCRUD struct {
	KeyStore datastore.KeyStore
	Logger   *zap.Logger
}

var (
	zapKeyCRUDField = zap.String("management-crud", "key")

	zapCRUDCreate   = zap.String("crud", "create")
	zapCRUDRetrieve = zap.String("crud", "retrieve")
	zapCRUDUpdate   = zap.String("crud", "update")
	zapCRUDDelete   = zap.String("crud", "delete")
)

// http POST
func (kcCRUD KeyControllerCRUD) Create(w http.ResponseWriter, r *http.Request) {
	keyJsonB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "invalid request body")
		kcCRUD.Logger.Warn("bad key crud request, invalid request body",
			zap.String("ip", r.Host),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDCreate)
		return
	}
	newKey := types.Key{}
	err = json.Unmarshal(keyJsonB, &newKey)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "invalid request body")
		kcCRUD.Logger.Warn("bad key crud request, failed to parse body json blob",
			zap.String("ip", r.Host),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDCreate)
		return
	}
	retKey, err := kcCRUD.KeyStore.Create(newKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to create encryption key")
		kcCRUD.Logger.Error("failed to create encryption key",
			zap.Any("key", newKey),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDCreate)
		return
	}
	err = writeJSONResponse(w, retKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to return encryption key")
		kcCRUD.Logger.Error("failed to marshal encryption key",
			zap.Any("key", retKey),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDCreate)
		return
	}
	w.WriteHeader(http.StatusCreated)
	kcCRUD.Logger.Info("encryption key created",
		zap.String("ip", r.Host),
		zap.String("key_id", retKey.ID.String()),
		zapKeyCRUDField,
		zapCRUDCreate)
}

// http GET
func (kcCRUD KeyControllerCRUD) Retrieve(w http.ResponseWriter, r *http.Request) {
	keyID := r.URL.Query().Get("key_id")
	imageID := r.URL.Query().Get("image_id")
	if imageID == "" || !validSHA256Regex.MatchString(imageID) {
		if keyID == "" || !validUUIDRegex.MatchString(keyID) {
			writeResponse(w, http.StatusBadRequest, "invalid key or image id")
			kcCRUD.Logger.Warn("bad key crud request, invalid key or image id",
				zap.String("ip", r.Host),
				zap.String("image_id", imageID),
				zap.String("key_id", keyID),
				zapKeyCRUDField,
				zapCRUDRetrieve)
			return
		}
	}
	retKey, err := kcCRUD.KeyStore.Retrieve(types.Key{ID: uuid.MustParse(keyID), ImageID: imageID})
	if err != nil {
		writeResponse(w, http.StatusNotFound, "key not found")
		kcCRUD.Logger.Warn("failed binding request, key look up failed",
			zap.String("ip", r.Host),
			zap.String("image_id", imageID),
			zap.String("key_id", keyID),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDRetrieve)
		return
	}
	err = writeJSONResponse(w, retKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "failed to marshal key")
		kcCRUD.Logger.Error("failed to marshal key",
			zap.Any("key", retKey),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDRetrieve)
		return
	}
	w.WriteHeader(http.StatusOK)
	kcCRUD.Logger.Info("encryption key retrieved",
		zap.String("ip", r.Host),
		zap.String("key_id", retKey.ID.String()),
		zapKeyCRUDField,
		zapCRUDRetrieve)
}

// http PATCH
func (kcCRUD KeyControllerCRUD) Update(w http.ResponseWriter, r *http.Request) {
	keyJsonB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "invalid request body")
		kcCRUD.Logger.Warn("bad key crud request, invalid request body",
			zap.String("ip", r.Host),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDUpdate)
		return
	}
	newKey := types.Key{}
	err = json.Unmarshal(keyJsonB, &newKey)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "invalid request body")
		kcCRUD.Logger.Warn("bad key crud request, failed to parse body json blob",
			zap.String("ip", r.Host),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDUpdate)
		return
	}
	retKey, err := kcCRUD.KeyStore.Update(newKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to updated encryption key")
		kcCRUD.Logger.Error("failed to updated encryption key",
			zap.Any("key", newKey),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDUpdate)
		return
	}
	err = writeJSONResponse(w, retKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to return encryption key")
		kcCRUD.Logger.Error("failed to marshal encryption key",
			zap.Any("key", retKey),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDUpdate)
		return
	}
	w.WriteHeader(http.StatusOK)
	kcCRUD.Logger.Info("encryption key updated",
		zap.String("ip", r.Host),
		zap.String("key_id", retKey.ID.String()),
		zapKeyCRUDField,
		zapCRUDUpdate)
}

// http DELETE
func (kcCRUD KeyControllerCRUD) Delete(w http.ResponseWriter, r *http.Request) {
	keyID := r.URL.Query().Get("key_id")
	if keyID == "" || !validUUIDRegex.MatchString(keyID) {
		writeResponse(w, http.StatusBadRequest, "invalid key id")
		kcCRUD.Logger.Warn("bad delete request, invalid key id",
			zap.String("ip", r.Host),
			zap.String("key_id", keyID),
			zapKeyCRUDField,
			zapCRUDDelete)
		return
	}
	err := kcCRUD.KeyStore.Delete(types.Key{ID: uuid.MustParse(keyID)})
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to delete encryption key")
		kcCRUD.Logger.Error("failed to delete encryption key",
			zap.String("ip", r.Host),
			zap.Any("key_id", keyID),
			zap.Error(err),
			zapKeyCRUDField,
			zapCRUDDelete)
		return
	}
	w.WriteHeader(http.StatusOK)
	kcCRUD.Logger.Info("key deleted by client",
		zap.String("ip", r.Host),
		zap.String("key_id", keyID),
		zapKeyCRUDField,
		zapCRUDDelete)
}
