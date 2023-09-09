package memcache

import (
	"encoding/json"
	"strconv"

	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

type MemcacheCache struct {
	client MemcacheClient
	logger *zap.Logger
}

// NewMemcacheCache initializes a new Memcache cache instance
func NewMemcacheCache(logger *zap.Logger, server string) (*MemcacheCache, error) {
	mc := memcache.New(server)

	err := mc.Ping()
	if err != nil {
		logger.Error("Failed to ping Memcache server", zap.String("server", server), zap.Error(err))

		return nil, err
	}

	logger.Debug("Successfully connected to Memcache server", zap.String("server", server))

	return &MemcacheCache{
		client: mc,
		logger: logger,
	}, nil
}

// Set stores a user in the Memcache cache
func (m *MemcacheCache) Set(key int64, value *common.User) error {
	data, err := json.Marshal(value)
	if err != nil {
		m.logger.Error("Failed to marshal user data", zap.Int64("key", key), zap.Error(err))

		return err
	}

	item := &memcache.Item{
		Key:   strconv.FormatInt(key, 10),
		Value: data,
	}

	err = m.client.Set(item)
	if err != nil {
		m.logger.Error("Failed to set user data in Memcache", zap.Int64("key", key), zap.Error(err))

		return err
	}

	m.logger.Debug("Successfully stored user data in Memcache", zap.Int64("key", key))

	return nil
}

// Get retrieves a user from the Memcache cache
func (m *MemcacheCache) Get(key int64) (*common.User, error) {
	item, err := m.client.Get(strconv.FormatInt(key, 10))
	if err != nil {
		m.logger.Error("Failed to get user data from Memcache", zap.Int64("key", key), zap.Error(err))

		return nil, err
	}

	var user common.User
	if err := json.Unmarshal(item.Value, &user); err != nil {
		m.logger.Error("Failed to unmarshal user data", zap.Int64("key", key), zap.Error(err))

		return nil, err
	}

	m.logger.Debug("Successfully retrieved user data from Memcache", zap.Int64("key", key))

	return &user, nil
}
