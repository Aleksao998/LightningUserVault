package routers

import (
	"bytes"
	"encoding/json"
	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_GetUser(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{
		GetFn: func(key int64) (*common.User, error) {
			user := common.User{Name: "User-1", ID: 1}
			return &user, nil
		},
	}

	// Create test handler
	router := InitRouter(mockStorage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var user common.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "User-1", user.Name)
}

func TestRouter_SetUser(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{
		SetFn: func(value string) (int64, error) {
			return 1, nil
		},
	}

	// Create test handler
	router := InitRouter(mockStorage)

	// Create a mock user data for the POST request
	userData := map[string]interface{}{
		"id":   1,
		"name": "User-1",
		// Add other fields as necessary
	}
	userDataBytes, _ := json.Marshal(userData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/", bytes.NewBuffer(userDataBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var user common.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "User-1", user.Name)
}
