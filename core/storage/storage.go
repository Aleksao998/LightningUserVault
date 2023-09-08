package storage

import (
	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage/keyvalue/pebble"
)

// Storage represents a database interface.
type Storage interface {
	// Set stores a value and returns user ID and an error if any issue occurs during the operation
	Set(value string) (int64, error)

	// Get retrieves the value for a given user ID and returns an error if any issue occurs during the operation
	Get(key int64) (*common.User, error)

	// Close closes database connection and returns an error if any issue occurs during the operation
	Close() error
}

func GetStorage() (Storage, error) {
	return pebble.NewStorage("test-storage")
}
