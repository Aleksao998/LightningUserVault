package sql

import (
	"database/sql"

	"gorm.io/gorm"
)

// DBHandler provides an abstraction over GORM's database operations. It allows for
// consistent interactions with the database using GORM's methods and
// facilitates easier testing by enabling the mocking of these operations
type DBHandler interface {
	First(out interface{}, where ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	DB() (*sql.DB, error)
}
