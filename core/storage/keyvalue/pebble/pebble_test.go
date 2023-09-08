package pebble

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/stretchr/testify/assert"
)

func createPebbleStorage() (string, *Storage, error) {
	// Create a temporary directory for Pebble storage
	tempDir, err := os.MkdirTemp("", "pebble-test")
	if err != nil {
		return "", nil, fmt.Errorf("error creating temporary directory: %w", err)
	}

	// Initialize a new Pebble storage instance
	store, err := NewStorage(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)

		return "", nil, fmt.Errorf("error creating new Pebble storage: %w", err)
	}

	return tempDir, store, nil
}

// TestPebbleStorageGet_NonExistentKey tests getting value with non-existent key.
func TestPebbleStorageGet_NonExistentKey(t *testing.T) {
	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		t.Fatalf("error creating pabble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	// Test Get on non-existent key
	nonExistentKey := int64(1)

	_, err = store.Get(nonExistentKey)
	if !assert.ErrorIs(t, err, pebble.ErrNotFound) {
		t.Errorf("Expected error not found when getting non-existent key")
	}
}

// TestPabbleStorage_WriteParallel tests writing in storage parallel.
func TestPabbleStorage_WriteParallel(t *testing.T) {
	var wg sync.WaitGroup

	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		t.Fatalf("error creating pabble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	wg.Add(50)

	for i := 0; i < 50; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()

			_, err := store.Set(fmt.Sprintf("user-%d", i))
			assert.Nil(t, err)
		}(i)
	}

	wg.Wait()

	// Get and validate all users
	for i := 1; i <= 50; i++ {
		retrievedValue, err := store.Get(int64(i))
		assert.Nil(t, err)

		assert.Equal(t, int64(i), retrievedValue.ID)
	}
}

// TestPebbleStorage_WriteRead tests writing and reading value from storage.
func TestPebbleStorage_WriteRead(t *testing.T) {
	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		t.Fatalf("error creating pabble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	// Test Set and Get methods for multiple key-value pairs
	for i := 1; i <= 10; i++ {
		var err error

		key := int64(i)
		value := fmt.Sprintf("user_%d", i)

		id, err := store.Set(value)
		if err != nil {
			t.Fatalf("Error setting value for key '%d': %v", key, err)
		}

		retrievedValue, err := store.Get(id)
		if err != nil {
			t.Fatalf("Error getting value for key '%d': %v", key, err)
		}

		assert.Equal(t, value, retrievedValue.Name)
		assert.Equal(t, key, retrievedValue.ID)
	}
}
