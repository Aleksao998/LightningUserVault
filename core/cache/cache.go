package cache

import (
	"github.com/Aleksao998/LightingUserVault/core/cache/memcached"
	"github.com/Aleksao998/LightingUserVault/core/common"
	"go.uber.org/zap"
)

// Cache represents a caching interface
type Cache interface {
	// Set stores a value in the cache with a given key.
	Set(key int64, value *common.User) error

	// Get retrieves a value from the cache using a given key.
	Get(key int64) (*common.User, error)
}

func GetCache(logger *zap.Logger, path string) (Cache, error) {
	return memcached.NewMemcacheCache(logger, path)
}
