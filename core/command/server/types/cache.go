package types

import (
	"fmt"
	"strings"
)

// Define the CacheType type and its possible values
type CacheType string

const (
	MEMCACHE CacheType = "MEMCACHE"
)

// ConvertStringToCacheType converts a string to its corresponding CacheType
func ConvertStringToCacheType(s string) (CacheType, error) {
	switch strings.ToUpper(s) {
	case string(MEMCACHE):
		return MEMCACHE, nil
	default:
		return "", fmt.Errorf("invalid cache type: %s", s)
	}
}
