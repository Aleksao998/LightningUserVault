package mocks

import (
	"database/sql"

	"gorm.io/gorm"
)

type (
	FirstDelegate  func(out interface{}, where ...interface{}) *gorm.DB
	CreateDelegate func(value interface{}) *gorm.DB
	DBDelegate     func() (*sql.DB, error)
)

type MockSQLdb struct {
	FirstFn  FirstDelegate
	CreateFn CreateDelegate
	DBFn     DBDelegate
}

func (m *MockSQLdb) First(out interface{}, where ...interface{}) *gorm.DB {
	if m.FirstFn != nil {
		return m.FirstFn(out, where)
	}

	return nil
}

func (m *MockSQLdb) Create(value interface{}) *gorm.DB {
	if m.CreateFn != nil {
		return m.CreateFn(value)
	}

	return nil
}

func (m *MockSQLdb) DB() (*sql.DB, error) {
	if m.DBFn != nil {
		return m.DBFn()
	}

	return nil, nil
}
