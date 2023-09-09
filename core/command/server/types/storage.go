package types

import (
	"fmt"
	"strings"
)

// Define the StorageType type and its possible values
type StorageType string

const (
	PEBBLE     StorageType = "PEBBLE"
	POSTGRESQL StorageType = "POSTGRESQL"
)

// StorageType converts a string to its corresponding StorageType
func ConvertStringToStorageType(s string) (StorageType, error) {
	switch strings.ToUpper(s) {
	case string(PEBBLE):
		return PEBBLE, nil
	case string(POSTGRESQL):
		return POSTGRESQL, nil
	default:
		return "", fmt.Errorf("invalid storage type: %s", s)
	}
}
