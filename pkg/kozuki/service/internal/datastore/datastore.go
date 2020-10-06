package datastore

import (
	"errors"
	"laugh-tale/pkg/kozuki/types"
)

var ErrNotExists = errors.New("no record found")

type PostgresDBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string

	TLSMode string
	TLSCert string
}

type KeyStore interface {
	Create(types.Key) (types.Key, error)
	Retrieve(types.Key) (types.Key, error)
	Update(types.Key) (types.Key, error)
	Delete(types.Key) error
}
