package postgresql

import (
	"errors"
	"testing"

	"github.com/Aleksao998/LightingUserVault/core/storage/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var errInternal = errors.New("internal db error")

// TestPostgres_GetValid tests the scenario where a user is successfully retrieved from the database
func TestPostgres_GetValid(t *testing.T) {
	t.Parallel()

	mockDB := &mocks.MockSQLdb{
		FirstFn: func(out interface{}, where ...interface{}) *gorm.DB {
			u, ok := out.(*User)
			if !ok {
				t.Fatalf("value is not of type *User")
			}
			u.ID = 1
			u.Name = "Mocked User"

			return &gorm.DB{}
		},
	}

	storage := &Storage{
		db: mockDB,
	}

	user, err := storage.Get(1)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "Mocked User", user.Name)
}

// TestPostgres_GetNotFound tests the scenario where the user is not found in the database
func TestPostgres_GetNotFound(t *testing.T) {
	t.Parallel()

	mockDB := &mocks.MockSQLdb{
		FirstFn: func(out interface{}, where ...interface{}) *gorm.DB {
			return &gorm.DB{Error: gorm.ErrRecordNotFound} // Return a mock *gorm.DB with gorm.ErrRecordNotFound error
		},
	}

	storage := &Storage{
		db: mockDB,
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
			return &gorm.DB{Error: errInternal} // Return a mock *gorm.DB with gorm.ErrRecordNotFound error
		},
	}

	storage := &Storage{
		db: mockDB,
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
			u.Name = "Mocked User"

			return &gorm.DB{}
		},
	}

	storage := &Storage{
		db: mockDB,
	}

	id, err := storage.Set("Mocked User")
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
		db: mockDB,
	}

	id, err := storage.Set("Mocked User")
	assert.Error(t, err)
	assert.Equal(t, errInternal, err)
	assert.Equal(t, int64(0), id)
}
