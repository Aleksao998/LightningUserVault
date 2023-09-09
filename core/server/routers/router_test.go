package routers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cacheMock "github.com/Aleksao998/LightingUserVault/core/cache/mocks"
	"github.com/Aleksao998/LightingUserVault/core/common"
	storageMock "github.com/Aleksao998/LightingUserVault/core/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestRouter_GetUser tests the successful retrieval of a user via the router's endpoint
func TestRouter_GetUser(t *testing.T) {
	mockStorage := &storageMock.MockStorage{}
	mockCache := &cacheMock.MockCache{
		GetFn: func(key int64) (*common.User, error) {
			user := common.User{Name: "User-1", ID: 1}

			return &user, nil
		},
	}

	// Create test handler
	router := InitRouter(zap.NewNop(), mockStorage, mockCache)

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

// TestRouter_SetUser tests the successful setting of a user via the router's endpoint
func TestRouter_SetUser(t *testing.T) {
	mockStorage := &storageMock.MockStorage{
		SetFn: func(value string) (int64, error) {
			return 1, nil
		},
	}
	mockCache := &cacheMock.MockCache{}

	// Create test handler
	router := InitRouter(zap.NewNop(), mockStorage, mockCache)

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
