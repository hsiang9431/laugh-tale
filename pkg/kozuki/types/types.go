package types

import (
	"github.com/google/uuid"
)

// the key add to key manager
type Key struct {
	ID         uuid.UUID `json:"id"`
	ImageID    string    `json:"image_id"`
	ImplantKey string    `json:"implant_key"`
	DecryptKey string    `json:"decrypt_key"`
}
