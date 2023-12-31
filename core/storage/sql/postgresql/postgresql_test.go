package postgresql

import (
	"errors"
	"testing"

	"github.com/Aleksao998/LightningUserVault/core/common"
	"github.com/Aleksao998/LightningUserVault/core/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var errInternal = errors.New("internal db error")

const mockUserName = "Mocked User"

// TestPostgres_GetValid tests the scenario where a user is successfully retrieved from the database
func TestPostgres_GetValid(t *testing.T) {
	t.Parallel()

	mockDB := &mocks.MockSQLdb{
		FirstFn: func(out interface{}, where ...interface{}) *gorm.DB {
			u, ok := out.(*common.User)
			if !ok {
				t.Fatalf("value is not of type *User")
			}
			u.ID = 1
			u.Name = mockUserName

			return &gorm.DB{}
		},
	}

	storage := &Storage{
		db:     mockDB,
		logger: zap.NewNop(),
	}

	user, err := storage.Get(1)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, mockUserName, user.Name)
}

// TestPostgres_GetNotFound tests the scenario where the user is not found in the database
func TestPostgres_GetNotFound(t *testing.T) {
	t.Parallel()

	mockDB := &mocks.MockSQLdb{
		FirstFn: func(out interface{}, where ...interface{}) *gorm.DB {
			return &gorm.DB{Error: gorm.ErrRecordNotFound}
		},
	}

	storage := &Storage{
		db:     mockDB,
		logger: zap.NewNop(),
	}

	user, err := storage.Get(1)
	assert.Error(t, err)
	assert.Equal(t, errUserNotFound, err)
	assert.Nil(t, user)
}

// TestPostgres_GetInternalError tests the scenario where an internal error occurs while retrieving the user
func TestPostgres_GetInternalError(t *testing.T) {
	t.Parallel()

	mockDB := &mocks.MockSQLdb{
		FirstFn: func(out interface{}, where ...interface{}) *gorm.DB {
			return &gorm.DB{Error: errInternal}
		},
	}

	storage := &Storage{
		db:     mockDB,
		logger: zap.NewNop(),
	}

	user, err := storage.Get(1)
	assert.Error(t, err)
	assert.Equal(t, errInternal, err)
	assert.Nil(t, user)
}

// TestPostgres_SetSuccessfully tests the scenario where a user is successfully added to the database
func TestPostgres_SetSuccessfully(t *testing.T) {
	t.Parallel()

	mockDB := &mocks.MockSQLdb{
		CreateFn: func(value interface{}) *gorm.DB {
			u, ok := value.(*User)
			if !ok {
				t.Fatalf("value is not of type *User")
			}
			u.ID = 1
			u.Name = mockUserName

			return &gorm.DB{}
		},
	}

	storage := &Storage{
		db:     mockDB,
		logger: zap.NewNop(),
	}

	id, err := storage.Set(mockUserName)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id)
}

// TestPostgres_SetError tests the scenario where an error occurs while adding a user to the database
func TestPostgres_SetError(t *testing.T) {
	t.Parallel()

	mockDB := &mocks.MockSQLdb{
		CreateFn: func(value interface{}) *gorm.DB {
			return &gorm.DB{Error: errInternal}
		},
	}

	storage := &Storage{
		db:     mockDB,
		logger: zap.NewNop(),
	}

	id, err := storage.Set(mockUserName)
	assert.Error(t, err)
	assert.Equal(t, errInternal, err)
	assert.Equal(t, int64(0), id)
}
