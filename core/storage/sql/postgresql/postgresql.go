package postgresql

import (
	"errors"

	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage/sql"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var errUserNotFound = errors.New("user not found")

type User struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255)"`
}

type Storage struct {
	db     sql.DBHandler
	logger *zap.Logger
}

// NewStorage initializes a new Storage instance with a database at the given path
func NewStorage(logger *zap.Logger, connStr string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to open PostgreSQL database", zap.String("connectionString", connStr), zap.Error(err))

		return nil, err
	}

	// AutoMigrate will ONLY create tables, missing columns and missing indexes
	err = db.AutoMigrate(&User{})
	if err != nil {
		logger.Error("Failed to auto-migrate User schema", zap.Error(err))

		return nil, err
	}

	logger.Info("Successfully initialized PostgreSQL storage")

	return &Storage{
		db:     db,
		logger: logger,
	}, nil
}

// Get retrieves the user for a given ID
func (p *Storage) Get(id int64) (*common.User, error) {
	var user common.User

	result := p.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			p.logger.Warn("User not found", zap.Int64("ID", id))

			return nil, errUserNotFound
		}

		p.logger.Error("Failed to retrieve user from database", zap.Int64("ID", id), zap.Error(result.Error))

		return nil, result.Error
	}

	p.logger.Debug("Successfully retrieved user from database", zap.Int64("ID", user.ID))

	return &user, nil
}

// Set stores a user with the given name and returns the ID
func (p *Storage) Set(name string) (int64, error) {
	user := User{Name: name}

	result := p.db.Create(&user)
	if result.Error != nil {
		p.logger.Error("Failed to store user in database", zap.String("Name", name), zap.Error(result.Error))

		return 0, result.Error
	}

	p.logger.Debug("Successfully stored user in database", zap.Int64("ID", user.ID))

	return user.ID, nil
}

// Close closes the database connection
func (p *Storage) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		p.logger.Error("Failed to get underlying SQL database instance", zap.Error(err))

		return err
	}

	err = sqlDB.Close()
	if err != nil {
		p.logger.Error("Failed to close PostgreSQL database connection", zap.Error(err))

		return err
	}

	p.logger.Info("Successfully closed PostgreSQL database connection")

	return nil
}
