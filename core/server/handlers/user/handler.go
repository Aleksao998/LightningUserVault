package userhandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Aleksao998/LightingUserVault/core/cache"
	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"github.com/gin-gonic/gin"
)

var (
	errInvalidUserID       = errors.New("invalid user ID")
	errInvalidUserName     = errors.New("invalid user name")
	errInvalidReqJSONParam = errors.New("request is invalid json")
)

type UserHandler struct {
	vault storage.Storage
	cache cache.Cache
}

// NewUserHandler creates a new UserHandler with the given storage
func NewUserHandler(storage storage.Storage, cache cache.Cache) *UserHandler {
	return &UserHandler{
		vault: storage,
		cache: cache,
	}
}

// @Summary Get user by ID
// @Description Retrieve user details by user ID
// @ID get-user-by-id
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} common.User
// @Failure 400 {object} common.ErrorResponse
// @Router /user/{id} [get]
func (h *UserHandler) GetHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: errInvalidUserID.Error()})

		return
	}

	// Try to get the user from cache first
	user, err := h.cache.Get(id)
	if err == nil {
		c.JSON(http.StatusOK, user)

		return
	}

	// If not in cache, get from vault
	user, err = h.vault.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{Error: err.Error()})

		return
	}

	// Store the fetched user in cache
	err = h.cache.Set(id, user)
	if err != nil {
		//  TODO log
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Set a new user
// @Description Add a new user and return their ID
// @ID set-user
// @Accept  json
// @Produce json
// @Param user body common.User true "User object"
// @Success 200 {object} common.User
// @Failure 400 {object} common.ErrorResponse
// @Router /user [post]
func (h *UserHandler) SetHandler(c *gin.Context) {
	var user common.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: errInvalidReqJSONParam.Error()})

		return
	}

	if user.Name == "" {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: errInvalidUserName.Error()})

		return
	}

	id, err := h.vault.Set(user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})

		return
	}

	user.ID = id

	c.JSON(http.StatusOK, user)
}
