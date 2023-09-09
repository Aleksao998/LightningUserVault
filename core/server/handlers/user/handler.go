package userhandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Aleksao998/LightingUserVault/core/cache"
	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	errInvalidUserID       = errors.New("invalid user ID")
	errInvalidUserName     = errors.New("invalid user name")
	errInvalidReqJSONParam = errors.New("request is invalid json")
)

type Config struct {
	CacheEnabled bool
}

type UserHandler struct {
	vault  storage.Storage
	cache  cache.Cache
	logger *zap.Logger
	config Config
}

// NewUserHandler creates a new UserHandler with the given storage
func NewUserHandler(logger *zap.Logger, storage storage.Storage, cache cache.Cache, config Config) *UserHandler {
	return &UserHandler{
		vault:  storage,
		cache:  cache,
		logger: logger,
		config: config,
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
		h.logger.Warn("Invalid user ID received", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: errInvalidUserID.Error()})

		return
	}

	if h.config.CacheEnabled {
		// Try to get the user from cache first
		user, err := h.cache.Get(id)
		if err == nil {
			h.logger.Debug("User fetched from cache", zap.Int64("id", id))
			c.JSON(http.StatusOK, user)

			return
		}
	}

	// If not in cache, get from vault
	user, err := h.vault.Get(id)
	if err != nil {
		h.logger.Error("Failed to fetch user from vault", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, common.ErrorResponse{Error: err.Error()})

		return
	}

	if h.config.CacheEnabled {
		// Store the fetched user in cache
		err = h.cache.Set(id, user)
		if err != nil {
			h.logger.Error("Failed to set user in cache", zap.Int64("id", id), zap.Error(err))
		} else {
			h.logger.Debug("User stored in cache", zap.Int64("id", id))
		}
	}

	h.logger.Info("Returning user data", zap.Int64("id", id))
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
		h.logger.Warn("Invalid JSON received", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: errInvalidReqJSONParam.Error()})

		return
	}

	if user.Name == "" {
		h.logger.Warn("Invalid user name received", zap.String("name", user.Name))
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: errInvalidUserName.Error()})

		return
	}

	id, err := h.vault.Set(user.Name)
	if err != nil {
		h.logger.Error("Failed to set user in vault", zap.String("name", user.Name), zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})

		return
	}

	user.ID = id
	h.logger.Info("User successfully stored", zap.Int64("id", id), zap.String("name", user.Name))
	c.JSON(http.StatusOK, user)
}
