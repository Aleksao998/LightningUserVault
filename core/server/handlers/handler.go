package handlers

import (
	"github.com/Aleksao998/LightingUserVault/core/common"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	vault storage.Storage
}

func NewUserHandler(storage storage.Storage) *UserHandler {
	return &UserHandler{
		vault: storage,
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
func (uh *UserHandler) GetUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := uh.vault.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
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
func (uh *UserHandler) SetUserHandler(c *gin.Context) {
	var user common.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid request body"})
		return
	}

	id, err := uh.vault.Set(user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	user.ID = id

	c.JSON(http.StatusOK, user)
}
