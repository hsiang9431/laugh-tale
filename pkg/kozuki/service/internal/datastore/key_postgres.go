package datastore

import (
	"fmt"
	"laugh-tale/pkg/kozuki/types"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type key struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid"`
	CreatedAt  time.Time `json:"created"`
	ImageID    string    `json:"image_id"`
	ImplantKey string    `json:"implant_key"`
	DecryptKey string    `json:"decrypt_key"`
}

type pgdbKeyStore struct {
	logger *zap.Logger
	db     *gorm.DB
}

var zapPGDBField = zap.String("database", "postgres")

func NewPostgresDB(l *zap.Logger, c *PostgresDBConfig) (KeyStore, error) {
	if l == nil || c == nil {
		return nil, errors.New("postgres: neither logger nor database configration can be nil")
	}
	if c.Host == "" ||
		c.Port == "" ||
		c.Username == "" ||
		c.Password == "" ||
		c.DBName == "" {
		return nil, errors.New("postgres: invalid config")
	}
	dbConn, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s cfg.SslMode=%s%s",
			c.Host, c.Port, c.Username, c.DBName, c.Password, c.TLSMode, c.TLSCert))
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to connect to database")
	}
	if dbConn == nil {
		return nil, errors.New("postgres: db connection is nil")
	}
	dbConn.AutoMigrate(key{})
	l.Info("connected to database",
		zap.String("db-host", c.Host+":"+c.Port),
		zapPGDBField)
	return &pgdbKeyStore{logger: l, db: dbConn}, nil
}

func (ks *pgdbKeyStore) Create(k types.Key) (types.Key, error) {
	k.ID = uuid.New()
	dbKey := key{
		ID:         k.ID,
		ImageID:    k.ImageID,
		ImplantKey: k.ImplantKey,
		DecryptKey: k.DecryptKey,
	}
	if err := ks.db.Create(&dbKey).Error; err != nil {
		return types.Key{}, errors.Wrap(err, "postgres: failed to create key in db")
	}
	ks.logger.Info("key created in database",
		zap.String("key_id", k.ID.String()),
		zapPGDBField)
	return k, nil
}

func (ks *pgdbKeyStore) Retrieve(k types.Key) (types.Key, error) {
	dbKey := key{}
	if err := ks.db.First(&dbKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return types.Key{}, ErrNotExists
		}
		return types.Key{}, errors.Wrap(err, "postgres: failed to create key in db")
	}
	ks.logger.Info("key updated in database",
		zap.String("key_id", k.ID.String()),
		zapPGDBField)
	return types.Key{
		ID:         dbKey.ID,
		ImageID:    dbKey.ImageID,
		ImplantKey: dbKey.ImplantKey,
		DecryptKey: dbKey.DecryptKey,
	}, nil
}

func (ks *pgdbKeyStore) Update(k types.Key) (types.Key, error) {
	if k.ID.String() == "" {
		return types.Key{}, errors.New("postgres: cannot update with empty uuid")
	}
	dbKey := key{
		ID:         k.ID,
		ImageID:    k.ImageID,
		ImplantKey: k.ImplantKey,
		DecryptKey: k.DecryptKey,
	}
	if err := ks.db.Model(&key{}).Updates(dbKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return types.Key{}, ErrNotExists
		}
		return types.Key{}, errors.Wrap(err, "postgres: failed to update key in db")
	}
	ks.logger.Info("key updated in database",
		zap.String("key_id", k.ID.String()),
		zapPGDBField)
	return types.Key{
		ID:         dbKey.ID,
		ImageID:    dbKey.ImageID,
		ImplantKey: dbKey.ImplantKey,
		DecryptKey: dbKey.DecryptKey,
	}, nil
}

func (ks *pgdbKeyStore) Delete(k types.Key) error {
	if err := ks.db.Delete(key{ID: k.ID}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotExists
		}
		return errors.Wrap(err, "postgres: failed to delete key in db")
	}
	ks.logger.Info("key deleted in database",
		zap.String("key_id", k.ID.String()),
		zapPGDBField)
	return nil
}
