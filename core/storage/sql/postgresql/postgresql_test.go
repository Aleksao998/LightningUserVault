package postgresql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// createMockPostgreSQLStorage initializes a mock database connection and returns the mock DB and its associated mock object.
func createMockPostgreSQLStorage() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open mock sql db: %w", err)
	}

	return db, mock, nil
}

// TestPostgreSQLStorage_Get tests the functionality of retrieving a user from the storage.
func TestPostgreSQLStorage_Get(t *testing.T) {
	t.Parallel()

	// Initialize mock PostgreSQL storage
	db, mock, err := createMockPostgreSQLStorage()
	if err != nil {
		t.Fatalf("error creating mock PostgreSQL storage: %v", err)
	}
	defer db.Close()

	// Define expected columns and mock database responses
	columns := []string{"id", "name"}
	mock.ExpectPrepare(`^SELECT id, name FROM users WHERE id = \$1$`).ExpectQuery().
		WithArgs(1). // user id
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "User1"))

	mock.ExpectPrepare(`^SELECT id, name FROM users WHERE id = \$1$`).ExpectQuery().
		WithArgs(2). // user id
		WillReturnError(sql.ErrNoRows)

	// Prepare the getUser query statement
	getUserStmt, err := db.Prepare(getUserQuery)
	if err != nil {
		t.Fatalf("error preparing getUser query: %v", err)
	}

	// Initialize the storage with the mock DB and prepared statement
	storage := &Storage{
		db:          db,
		getUserStmt: getUserStmt,
	}

	// Test retrieving a user with ID 1
	user, err := storage.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, "User1", user.Name)

	// Test retrieving a user with ID 2 (expected to not exist)
	user, err = storage.Get(2)
	assert.Nil(t, user)
	assert.Error(t, errUserNotFound, err)
}

// TestPostgreSQLStorage_Set tests the functionality of storing a user in the storage.
func TestPostgreSQLStorage_Set(t *testing.T) {
	t.Parallel()

	// Initialize mock PostgreSQL storage
	db, mock, err := createMockPostgreSQLStorage()
	if err != nil {
		t.Fatalf("error creating mock PostgreSQL storage: %v", err)
	}
	defer db.Close()

	// Mock the insertion of a user and returning its ID
	userName := "User1"
	expectedID := int64(1)
	mock.ExpectPrepare(`^INSERT INTO users \(name\) VALUES \(\$1\) RETURNING id$`).ExpectQuery().
		WithArgs(userName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	// Prepare the addUser query statement
	addUserStmt, err := db.Prepare(addUserQuery)
	if err != nil {
		t.Fatalf("error preparing addUser query: %v", err)
	}

	// Initialize the storage with the mock DB and prepared statement
	storage := &Storage{
		db:          db,
		addUserStmt: addUserStmt,
	}

	// Test storing a user and retrieving its ID
	id, err := storage.Set(userName)
	assert.Nil(t, err)
	assert.Equal(t, expectedID, id)
}
