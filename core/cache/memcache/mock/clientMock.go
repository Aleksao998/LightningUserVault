package mock

import "github.com/bradfitz/gomemcache/memcache"

type (
	SetDelegate func(item *memcache.Item) error
	GetDelegate func(key string) (item *memcache.Item, err error)
)

type MockClient struct {
	SetFn SetDelegate
	GetFn GetDelegate
}

func (m *MockClient) Set(item *memcache.Item) error {
	if m.SetFn != nil {
		return m.SetFn(item)
	}

	return nil
}

func (m *MockClient) Get(key string) (item *memcache.Item, err error) {
	if m.GetFn != nil {
		return m.GetFn(key)
	}

	return nil, nil
}
