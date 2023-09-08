package pebble

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzPebbleStorageWriteRead(f *testing.F) {
	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		f.Fatalf("error creating pabble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	f.Fuzz(func(t *testing.T, value string) {
		t.Parallel()

		// Test Set method
		id, err := store.Set(value)
		if err != nil {
			t.Fatalf("Error setting value: %s:%v", value, err)
		}

		// Test Get method
		retrievedUser, err := store.Get(id)
		if err != nil {
			t.Fatalf("Error getting value for key '%d': %v", id, err)
		}

		assert.Equal(t, value, retrievedUser.Name)
	})
}
