package resource

import (
	"laugh-tale/pkg/common/crypto"
	"laugh-tale/pkg/kozuki/service/internal/datastore"
	"laugh-tale/pkg/kozuki/types"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// for serving apis used by roger and ohara
type KeyController struct {
	KeyStore datastore.KeyStore
	Logger   *zap.Logger
}

var (
	// The length of random base 64 string
	ImplantKeyLength = 45
	DecryptKeyLength = 45
)

var validSHA256Regex = regexp.MustCompile("^sha256\\:[a-f0-9]{64}$")
var validUUIDRegex = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$")

func (kc KeyController) CreateEncryptionKey(w http.ResponseWriter, r *http.Request) {
	newKey := types.Key{
		ImplantKey: crypto.RandB64String(ImplantKeyLength),
		DecryptKey: crypto.RandB64String(DecryptKeyLength),
	}
	createdKey, err := kc.KeyStore.Create(newKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to create encryption key")
		kc.Logger.Error("failed to create encryption key",
			zap.Any("key", newKey),
			zap.Error(err))
		return
	}
	retKey := types.Key{ID: createdKey.ID}
	err = writeJSONResponse(w, retKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to create encryption key")
		kc.Logger.Error("failed to marshal encryption key",
			zap.Any("key", retKey),
			zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	kc.Logger.Info("encryption key created",
		zap.String("ip", r.Host),
		zap.String("key_id", retKey.ID.String()))
}

func (kc KeyController) BindKeyToImageID(w http.ResponseWriter, r *http.Request) {
	imageID := r.PostForm.Get("image_id")
	if imageID == "" || !validSHA256Regex.MatchString(imageID) {
		writeResponse(w, http.StatusBadRequest, "invalid image id")
		kc.Logger.Warn("bad binding request, invalid image id",
			zap.String("ip", r.Host),
			zap.String("image_id", imageID))
		return
	}
	keyID := r.PostForm.Get("key_id")
	if keyID == "" || !validUUIDRegex.MatchString(keyID) {
		writeResponse(w, http.StatusBadRequest, "invalid key id")
		kc.Logger.Warn("bad binding request, invalid key id",
			zap.String("ip", r.Host),
			zap.String("key_id", keyID))
		return
	}
	// find key with id
	key, err := kc.KeyStore.Retrieve(types.Key{ID: uuid.MustParse(keyID)})
	if err != nil {
		writeResponse(w, http.StatusNotFound, "key not found")
		kc.Logger.Warn("failed binding request, key look up failed",
			zap.String("ip", r.Host),
			zap.String("key_id", keyID),
			zap.Error(err))
		return
	}
	key.ImageID = imageID
	_, err = kc.KeyStore.Update(key)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "key not updated")
		kc.Logger.Error("failed to update key",
			zap.Any("key", key),
			zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	kc.Logger.Info("key bound to image and tag",
		zap.String("ip", r.Host),
		zap.String("key_id", key.ID.String()),
		zap.String("image_id", imageID))
}

func (kc KeyController) RetrieveKeyByImageID(w http.ResponseWriter, r *http.Request) {
	imageID := mux.Vars(r)["image_id"]
	if imageID == "" || !validSHA256Regex.MatchString(imageID) {
		writeResponse(w, http.StatusBadRequest, "invalid image id")
		kc.Logger.Warn("bad retrieve key request, invalid image id",
			zap.String("ip", r.Host),
			zap.String("image_id", imageID))
		return
	}
	// find key with image id
	retKey, err := kc.KeyStore.Retrieve(types.Key{ImageID: imageID})
	if err != nil {
		writeResponse(w, http.StatusNotFound, "no key bound to provided ID")
		kc.Logger.Warn("failed retrieve key request, key look up failed",
			zap.String("ip", r.Host),
			zap.String("image_id", imageID), zap.Error(err),
			zap.Error(err))
		return
	}
	err = writeJSONResponse(w, retKey)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"failed to marshal encryption key")
		kc.Logger.Error("failed to marshal encryption key",
			zap.Any("key", retKey),
			zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	kc.Logger.Info("key released to client",
		zap.String("ip", r.Host),
		zap.String("key_id", retKey.ID.String()))
}
