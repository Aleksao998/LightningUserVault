package postgresql

import (
	"errors"

	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var errUserNotFound = errors.New("user not found")

type User struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255)"`
}

type Storage struct {
	db sql.DBHandler
}

// NewStorage initializes a new Storage instance with a database at the given path
func NewStorage(connStr string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate will ONLY create tables, missing columns and missing indexes
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

// Get retrieves the user for a given ID
func (p *Storage) Get(id int64) (*common.User, error) {
	var user common.User

	result := p.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errUserNotFound
		}

		return nil, result.Error
	}

	return &user, nil
}

// Set stores a user with the given name and returns the ID
func (p *Storage) Set(name string) (int64, error) {
	user := User{Name: name}

	result := p.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

// Close closes the database connection
func (p *Storage) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
