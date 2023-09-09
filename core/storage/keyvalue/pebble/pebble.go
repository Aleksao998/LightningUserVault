package pebble

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/cockroachdb/pebble"
	"go.uber.org/zap"
)

// nextIDKey represents key which will be used to save latest nextID on server stop.
const nextIDKey = "__nextID__"

type Storage struct {
	db     *pebble.DB
	logger *zap.Logger
	nextID int64
}

// NewStorage initializes a new Storage instance with a database at the given path.
func NewStorage(path string, logger *zap.Logger) (*Storage, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		logger.Error("Failed to open pebble database", zap.String("path", path), zap.Error(err))

		return nil, err
	}

	nextID := loadNextIDFromDB(db)
	logger.Debug("Loaded nextID from database", zap.Int64("nextID", nextID))

	return &Storage{
		db:     db,
		logger: logger,
		nextID: nextID,
	}, nil
}

// Set stores a value for a given key and returns an error if any issue occurs during the operation.
func (p *Storage) Set(value string) (int64, error) {
	id := p.getNextID()

	_, _, err := p.db.Get(id)
	if !errors.Is(err, pebble.ErrNotFound) {
		msg := fmt.Sprintf("Id already exists %d:%v", common.BytesToInt64(id), err)
		p.logger.Error(msg, zap.Int64("id", common.BytesToInt64(id)))

		// TODO check panic
		panic(fmt.Sprintf("Id already exists %d:%v", common.BytesToInt64(id), err))
	}

	err = p.db.Set(id, []byte(value), pebble.Sync)
	if err != nil {
		p.logger.Error("Failed to set value in database", zap.String("value", value), zap.Error(err))
	}

	return common.BytesToInt64(id), err
}

// Get retrieves the value for a given key and returns an error if any issue occurs during the operation.
func (p *Storage) Get(key int64) (*common.User, error) {
	value, closer, err := p.db.Get(common.Int64ToBytes(key))
	if err != nil {
		p.logger.Error("Failed to get value from database", zap.Int64("key", key), zap.Error(err))

		return nil, err
	}
	defer closer.Close()

	user := &common.User{
		ID:   key,
		Name: string(value),
	}

	p.logger.Debug("Retrieved user from database", zap.Int64("ID", user.ID), zap.String("Name", user.Name))

	return user, nil
}

// getNextID retrieves next id for a new user.
func (p *Storage) getNextID() []byte {
	id := atomic.AddInt64(&p.nextID, 1)

	return common.Int64ToBytes(id)
}

// saveNextID saves the current nextID value to the database.
func (p *Storage) saveNextID() error {
	data := common.Int64ToBytes(p.nextID)

	err := p.db.Set([]byte(nextIDKey), data, pebble.Sync)
	if err != nil {
		p.logger.Error("Failed to save nextID to database", zap.Int64("nextID", p.nextID), zap.Error(err))
	}

	p.logger.Debug("Successfully saved nextID to database", zap.Int64("nextID", p.nextID))

	return err
}

// Close closes the database connection and returns an error if any issue occurs during the operation.
func (p *Storage) Close() error {
	if err := p.saveNextID(); err != nil {
		p.logger.Warn("Could not save next id", zap.Error(err))
	}

	p.logger.Info("Closing database connection")

	return p.db.Close()
}

// loadNextIDFromDB loads the nextID from the database during initialization.
func loadNextIDFromDB(db *pebble.DB) int64 {
	data, closer, err := db.Get([]byte(nextIDKey))
	if err != nil {
		return 0
	}
	defer closer.Close()

	return common.BytesToInt64(data)
}
