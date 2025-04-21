package internal

import (
	"go_chat/internal/chat"
	"go_chat/internal/config"
	"go_chat/internal/user"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg := config.GetConfig()

	chatRepo := chat.NewChatRepository()
	chatService := chat.NewChatService(chatRepo)
	chatHandler := chat.NewChatHandler(chatService)

	userRepo := user.NewUserRepository()
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/users", userHandler.CreateUserHandler)
	router.GET("/users/:user_id", userHandler.GetUserHandler)
	router.PUT("/users/:user_id", userHandler.UpdateUserHandler)
	router.DELETE("/users/:user_id", userHandler.DeleteUserHandler)

	router.POST("/chats", chatHandler.CreateChatHandler)
	router.GET("/chats/:chat_id/messages", chatHandler.GetMessagesHandler)
	router.POST("/chats/:chat_id/messages", chatHandler.SendMessageHandler)

	router.Run(":" + cfg.Port)
}
