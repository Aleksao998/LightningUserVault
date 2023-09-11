package memcache

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Aleksao998/LightningUserVault/core/cache/memcache/mock"
	"github.com/Aleksao998/LightningUserVault/core/common"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	errInternal = errors.New("internal error")
	errClient   = errors.New("client error")
)

// TestMemcache_SetValid tests the scenario where a user is successfully set in the Memcache cache
func TestMemcache_SetValid(t *testing.T) {
	cache := &MemcacheCache{
		client: &mock.MockClient{
			SetFn: func(item *memcache.Item) error {
				return nil
			},
		},
		logger: zap.NewNop(),
	}

	user := &common.User{Name: "Valid User"}
	err := cache.Set(1, user)
	assert.Nil(t, err)
}

// TestMemcache_SetInternalError tests the scenario where an internal error occurs while setting a user in the Memcache cache
func TestMemcache_SetInternalError(t *testing.T) {
	cache := &MemcacheCache{
		client: &mock.MockClient{
			SetFn: func(item *memcache.Item) error {
				return errInternal
			},
		},
		logger: zap.NewNop(),
	}

	user := &common.User{Name: "Another User"}
	err := cache.Set(1, user)
	assert.Error(t, err)
	assert.Equal(t, errInternal, err)
}

// TestMemcache_GetClientError tests the scenario where the Memcache client returns an error when trying to get a user from the cache
func TestMemcache_GetClientError(t *testing.T) {
	cache := &MemcacheCache{
		client: &mock.MockClient{
			GetFn: func(key string) (*memcache.Item, error) {
				return nil, errClient
			},
		},
		logger: zap.NewNop(),
	}

	user, err := cache.Get(1)
	assert.Error(t, err)
	assert.Equal(t, errClient, err)
	assert.Nil(t, user)
}

// TestMemcache_GetValid tests the successful retrieval of a user from the Memcache cache
func TestMemcache_GetValid(t *testing.T) {
	validUser := &common.User{Name: "Valid User"}

	cache := &MemcacheCache{
		client: &mock.MockClient{
			GetFn: func(key string) (*memcache.Item, error) {
				data, _ := json.Marshal(validUser)

				return &memcache.Item{Value: data}, nil
			},
		},
		logger: zap.NewNop(),
	}

	user, err := cache.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, validUser.Name, user.Name)
}
