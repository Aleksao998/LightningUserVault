package userhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cacheMock "github.com/Aleksao998/LightingUserVault/core/cache/mocks"
	"github.com/Aleksao998/LightingUserVault/core/common"
	storageMock "github.com/Aleksao998/LightingUserVault/core/storage/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	errUserNotFound   = errors.New("user not found")
	errInternal       = errors.New("internal error")
	errUserNotInCache = errors.New("user not in cache")
)

// TestUserHandler_GetWithInvalidParams tests the behavior of the GetHandler when provided with invalid parameters
func TestUserHandler_GetWithInvalidParams(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{}
	mockCache := &cacheMock.MockCache{}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)

	// Set the "id" parameter
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "asd"})

	// Call the GetHandler function
	handler.GetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var jsonError common.ErrorResponse

	err := json.Unmarshal(w.Body.Bytes(), &jsonError)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, errInvalidUserID.Error(), jsonError.Error)
}

// TestUserHandler_GetValidUserCache tests the successful retrieval of a user from the cache
func TestUserHandler_GetValidUserCache(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{}
	mockCache := &cacheMock.MockCache{
		GetFn: func(key int64) (*common.User, error) {
			user := common.User{Name: "User-1", ID: 1}

			return &user, nil
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)

	// Set the "id" parameter
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	// Call the GetHandler function
	handler.GetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var user common.User

	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "User-1", user.Name)
}

// TestUserHandler_GetValidUserDB tests the successful retrieval of a user from the database when it's not found in the cache
func TestUserHandler_GetValidUserDB(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{
		GetFn: func(key int64) (*common.User, error) {
			user := common.User{Name: "User-1", ID: 1}

			return &user, nil
		},
	}
	mockCache := &cacheMock.MockCache{
		GetFn: func(key int64) (*common.User, error) {
			return nil, errUserNotInCache
		},
		SetFn: func(key int64, value *common.User) error {
			return nil
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)

	// Set the "id" parameter
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	// Call the GetHandler function
	handler.GetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var user common.User

	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "User-1", user.Name)
}

// TestUserHandler_GetMissingUser tests the behavior of the GetHandler when the user is missing both in the cache and the database
func TestUserHandler_GetMissingUser(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{
		GetFn: func(key int64) (*common.User, error) {
			return nil, errUserNotFound
		},
	}
	mockCache := &cacheMock.MockCache{
		GetFn: func(key int64) (*common.User, error) {
			return nil, errUserNotInCache
		},
		SetFn: func(key int64, value *common.User) error {
			return nil
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)

	// Set the "id" parameter
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	// Call the GetHandler function
	handler.GetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var jsonError common.ErrorResponse

	err := json.Unmarshal(w.Body.Bytes(), &jsonError)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, errUserNotFound.Error(), jsonError.Error)
}

// TestUserHandler_GetErrSaveCache tests the behavior of the GetHandler when there's an error saving the user to the cache
func TestUserHandler_GetErrSaveCache(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{
		GetFn: func(key int64) (*common.User, error) {
			return nil, errUserNotFound
		},
	}
	mockCache := &cacheMock.MockCache{
		GetFn: func(key int64) (*common.User, error) {
			return nil, errUserNotInCache
		},
		SetFn: func(key int64, value *common.User) error {
			return errInternal
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)

	// Set the "id" parameter
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	// Call the GetHandler function
	handler.GetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var jsonError common.ErrorResponse

	err := json.Unmarshal(w.Body.Bytes(), &jsonError)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, errUserNotFound.Error(), jsonError.Error)
}

// TestUserHandler_SetValidUser tests the successful setting of a valid user
func TestUserHandler_SetValidUser(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{
		SetFn: func(value string) (int64, error) {
			return 1, nil
		},
	}
	mockCache := &cacheMock.MockCache{}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a new HTTP request with a valid user JSON body
	userJSON := `{"Name": "User-1"}`

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(userJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the GetHandler function
	handler.SetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var user common.User

	err = json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "User-1", user.Name)
}

// TestUserHandler_SetInternalError tests the behavior of the SetHandler when there's an internal error
func TestUserHandler_SetInternalError(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{
		SetFn: func(value string) (int64, error) {
			return 0, errInternal
		},
	}
	mockCache := &cacheMock.MockCache{}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a new HTTP request with a user which does not exists
	userJSON := `{"Name": "User-1"}`

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(userJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the GetHandler function
	handler.SetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var jsonError common.ErrorResponse

	err = json.Unmarshal(w.Body.Bytes(), &jsonError)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, errInternal.Error(), jsonError.Error)
}

// TestUserHandler_SetMissingParams tests the behavior of the SetHandler when provided with missing parameters
func TestUserHandler_SetMissingParams(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{}
	mockCache := &cacheMock.MockCache{}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a new HTTP request with missing user name
	userJSON := `{}`

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(userJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the GetHandler function
	handler.SetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var jsonError common.ErrorResponse

	err = json.Unmarshal(w.Body.Bytes(), &jsonError)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, errInvalidUserName.Error(), jsonError.Error)
}

// TestUserHandler_SetInvalidJsonParams tests the behavior of the SetHandler when provided with an invalid JSON body
func TestUserHandler_SetInvalidJsonParams(t *testing.T) {
	t.Parallel()

	mockStorage := &storageMock.MockStorage{}
	mockCache := &cacheMock.MockCache{}

	// Create test handler
	handler := NewUserHandler(mockStorage, mockCache)

	// Create a new HTTP request with invalid JSON body
	userJSON := `{ "Name": "User-1`

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(userJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the GetHandler function
	handler.SetHandler(c)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var jsonError common.ErrorResponse

	err = json.Unmarshal(w.Body.Bytes(), &jsonError)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, errInvalidReqJSONParam.Error(), jsonError.Error)
}
