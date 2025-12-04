package handler

import (
	"feedsystem_video_go/internal/api"
	"feedsystem_video_go/internal/model"
	"feedsystem_video_go/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req api.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.CreateUser(&model.User{
		Username: req.Username,
		Password: req.Password,
	}); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user created"})
}

func (h *UserHandler) RenameByID(c *gin.Context) {
	var req api.RenameByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.RenameByID(req.ID, req.NewUsername); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user renamed"})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req api.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.ChangePassword(req.ID, req.NewPassword); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "password changed"})
}

func (h *UserHandler) FindByID(c *gin.Context) {
	var req api.FindByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user, err := h.userService.FindByID(req.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, user)
	}
}

func (h *UserHandler) FindByUsername(c *gin.Context) {
	var req api.FindByUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user, err := h.userService.FindByUsername(req.Username); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, user)
	}
}
