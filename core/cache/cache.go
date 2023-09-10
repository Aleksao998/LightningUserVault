package cache

import (
	"errors"
	"net"

	"github.com/Aleksao998/LightingUserVault/core/cache/memcache"
	"github.com/Aleksao998/LightingUserVault/core/command/server/types"
	"github.com/Aleksao998/LightingUserVault/core/common"
	"go.uber.org/zap"
)

var errInvalidCache = errors.New("invalid cache type")

// Cache represents a caching interface
type Cache interface {
	// Set stores a value in the cache with a given key.
	Set(key int64, value *common.User) error

	// Get retrieves a value from the cache using a given key.
	Get(key int64) (*common.User, error)
}

type Config struct {
	CacheType       types.CacheType
	MemcacheAddress *net.TCPAddr
	Enabled         bool
}

func GetCache(logger *zap.Logger, config Config) (Cache, error) {
	if !config.Enabled {
		logger.Debug("Cache disabled")

		return nil, nil
	}

	logger.Debug("Cache enabled")

	switch config.CacheType {
	case types.MEMCACHE:
		return memcache.NewMemcacheCache(logger, config.MemcacheAddress.String())
	default:
		return nil, errInvalidCache
	}
}
