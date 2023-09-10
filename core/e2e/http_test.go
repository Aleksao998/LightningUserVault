package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/e2e/freamwork"
	"github.com/stretchr/testify/assert"
)

func TestE2E_SetAndGet(t *testing.T) {
	// Initialize and start the test server using the framework
	testServer := freamwork.NewTestServerAndStart(t)

	// Create a request to set a new user
	user := common.User{Name: "John Doe"}
	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)

	resp, err := http.Post("http://"+testServer.Config.ServerAddress.String()+"/user", "application/json", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response to validate
	var newUser common.User
	err = json.NewDecoder(resp.Body).Decode(&newUser)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", newUser.Name)
	assert.Equal(t, int64(1), newUser.ID)

	// Make a request to get user with ID 1
	resp, err = http.Get("http://" + testServer.Config.ServerAddress.String() + "/user/1")
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)

	// Teardown logic after all tests
	freamwork.CleanupStorage()
}

func TestE2E_SetTwoItemsWithRestart(t *testing.T) {
	// Initialize and start the test server using the framework
	testServer := freamwork.NewTestServerAndStart(t)

	// Create a request to set the first user
	user1 := common.User{Name: "John Doe"}
	user1JSON, err := json.Marshal(user1)
	assert.NoError(t, err)

	resp, err := http.Post("http://"+testServer.Config.ServerAddress.String()+"/user", "application/json", bytes.NewBuffer(user1JSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response to validate
	var newUser common.User
	err = json.NewDecoder(resp.Body).Decode(&newUser)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", newUser.Name)
	assert.Equal(t, int64(1), newUser.ID)

	// Stop the server
	testServer.Stop()

	// Start the server again
	testServer = freamwork.NewTestServerAndStart(t)

	// Create a request to set the second user
	user2 := common.User{Name: "Jane Smith"}
	user2JSON, err := json.Marshal(user2)
	assert.NoError(t, err)

	resp, err = http.Post("http://"+testServer.Config.ServerAddress.String()+"/user", "application/json", bytes.NewBuffer(user2JSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Make a request to get the second user
	resp, err = http.Get("http://" + testServer.Config.ServerAddress.String() + "/user/2")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response
	var retrievedUser common.User
	err = json.NewDecoder(resp.Body).Decode(&retrievedUser)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), retrievedUser.ID)
	assert.Equal(t, "Jane Smith", retrievedUser.Name)

	// Teardown logic after all tests
	freamwork.CleanupStorage()
}
