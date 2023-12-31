package storage

import (
	"errors"
	"fmt"

	"github.com/Aleksao998/LightningUserVault/core/command/server/types"
	"github.com/Aleksao998/LightningUserVault/core/common"
	"github.com/Aleksao998/LightningUserVault/core/storage/keyvalue/pebble"
	"github.com/Aleksao998/LightningUserVault/core/storage/sql/postgresql"
	"go.uber.org/zap"
)

const (
	pebbleStorageRoute = "pebble-storage"
)

var errInvalidStorage = errors.New("invalid storage type")

type Storage interface {
	// Set stores a value and returns user ID and an error if any issue occurs during the operation
	Set(value string) (int64, error)

	// Get retrieves the value for a given user ID and returns an error if any issue occurs during the operation
	Get(key int64) (*common.User, error)

	// Close closes storage instance
	Close() error
}

type Config struct {
	StorageType types.StorageType
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
}

// GetStorage initializes and returns a storage instance based on the provided configuration
// The method supports multiple storage types including PEBBLE and POSTGRESQL
func GetStorage(logger *zap.Logger, config Config) (Storage, error) {
	switch config.StorageType {
	case types.PEBBLE:
		return pebble.NewStorage(pebbleStorageRoute, logger)
	case types.POSTGRESQL:
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBName, config.DBPass, config.DBName)

		return postgresql.NewStorage(logger, psqlInfo)
	default:
		return nil, errInvalidStorage
	}
}
