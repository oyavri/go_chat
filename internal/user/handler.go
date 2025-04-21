package user

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// POST /users
func (h *UserHandler) CreateUserHandler(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	user, err := h.service.CreateUser(req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// GET /users/:user_id
func (h *UserHandler) GetUserHandler(ctx *gin.Context) {
	var req GetUserRequest

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	user, err := h.service.GetUserById(req.UserId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DELETE /users/:user_id
func (h *UserHandler) DeleteUserHandler(ctx *gin.Context) {
	var req DeleteUserRequest

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	err := h.service.DeleteUserById(req.UserId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// PUT /users/:user_id
func (h *UserHandler) UpdateUserHandler(ctx *gin.Context) {
	var req UpdateUserRequest

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	user, err := h.service.UpdateUserById(req.UserId, req.NewUsername, req.NewEmail)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
