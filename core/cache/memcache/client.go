package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheClient interface {
	Set(item *memcache.Item) error
	Get(key string) (item *memcache.Item, err error)
}
