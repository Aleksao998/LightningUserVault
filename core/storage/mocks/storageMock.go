package mocks

import "github.com/Aleksao998/LightingUserVault/core/common"

type (
	getFnDelegate func(key int64) (*common.User, error)
	setFnDelegate func(value string) (int64, error)
)

type MockStorage struct {
	GetFn getFnDelegate
	SetFn setFnDelegate
}

func (m *MockStorage) Get(key int64) (*common.User, error) {
	if m.GetFn != nil {
		return m.GetFn(key)
	}

	return nil, nil
}

func (m *MockStorage) Set(value string) (int64, error) {
	if m.SetFn != nil {
		return m.SetFn(value)
	}

	return 0, nil
}
