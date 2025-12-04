package router

import (
	"feedsystem_video_go/internal/handler"
	"feedsystem_video_go/internal/repository"
	"feedsystem_video_go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userHandler.CreateUser)
		userGroup.POST("/rename", userHandler.RenameByID)
		userGroup.POST("/changePassword", userHandler.ChangePassword)
		userGroup.POST("/findByID", userHandler.FindByID)
		userGroup.POST("/findByUsername", userHandler.FindByUsername)
	}

	return r
}
