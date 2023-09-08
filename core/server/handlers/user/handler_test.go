package userhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	errUserNotFound = errors.New("user not found")
	errInternal     = errors.New("internal error")
)

func TestUserHandler_GetValidUser(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{
		GetFn: func(key int64) (*common.User, error) {
			user := common.User{Name: "User-1", ID: 1}

			return &user, nil
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage)

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

func TestUserHandler_GetWithInvalidParams(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{}

	// Create test handler
	handler := NewUserHandler(mockStorage)

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

func TestUserHandler_GetMissingUser(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{
		GetFn: func(key int64) (*common.User, error) {
			return nil, errUserNotFound
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage)

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

func TestUserHandler_SetValidUser(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{
		SetFn: func(value string) (int64, error) {
			return 1, nil
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage)

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

func TestUserHandler_SetInternalError(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{
		SetFn: func(value string) (int64, error) {
			return 0, errInternal
		},
	}

	// Create test handler
	handler := NewUserHandler(mockStorage)

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

func TestUserHandler_SetMissingParams(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{}

	// Create test handler
	handler := NewUserHandler(mockStorage)

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

func TestUserHandler_SetInvalidJsonParams(t *testing.T) {
	t.Parallel()

	mockStorage := &mocks.MockStorage{}

	// Create test handler
	handler := NewUserHandler(mockStorage)

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
