package mocks

import "github.com/Aleksao998/LightingUserVault/core/common"

type (
	getDelegate   func(key int64) (*common.User, error)
	setDelegate   func(value string) (int64, error)
	closeDelegate func() error
)

type MockStorage struct {
	GetFn   getDelegate
	SetFn   setDelegate
	CloseFn closeDelegate
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

func (m *MockStorage) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}

	return nil
}
