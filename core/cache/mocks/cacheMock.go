package mocks

import (
	"github.com/Aleksao998/LightningUserVault/core/common"
)

type (
	SetDelegate func(key int64, value *common.User) error
	GetDelegate func(key int64) (*common.User, error)
)

type MockCache struct {
	SetFn SetDelegate
	GetFn GetDelegate
}

func (m *MockCache) Set(key int64, value *common.User) error {
	if m.SetFn != nil {
		return m.SetFn(key, value)
	}

	return nil
}

func (m *MockCache) Get(key int64) (*common.User, error) {
	if m.GetFn != nil {
		return m.GetFn(key)
	}

	return nil, nil
}
