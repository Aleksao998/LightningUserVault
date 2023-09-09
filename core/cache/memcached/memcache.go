package memcached

import (
	"encoding/json"
	"strconv"

	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheCache struct {
	client MemcacheClient
}

// NewMemcacheCache initializes a new Memcache cache instance
func NewMemcacheCache(server string) (*MemcacheCache, error) {
	mc := memcache.New(server)

	err := mc.Ping()
	if err != nil {
		return nil, err
	}

	return &MemcacheCache{
		client: mc,
	}, nil
}

// Set stores a user in the Memcache cache
func (m *MemcacheCache) Set(key int64, value *common.User) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:   strconv.FormatInt(key, 10),
		Value: data,
	}

	return m.client.Set(item)
}

// Get retrieves a user from the Memcache cache
func (m *MemcacheCache) Get(key int64) (*common.User, error) {
	item, err := m.client.Get(strconv.FormatInt(key, 10))
	if err != nil {
		return nil, err
	}

	var user common.User
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
