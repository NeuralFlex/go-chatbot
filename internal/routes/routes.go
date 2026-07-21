package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-chatbot/internal/handlers"
	"go-chatbot/internal/middleware"
)

func Register(router *gin.Engine, conv *handlers.ConversationsHandler) {
	router.GET("/ping", handlers.PingHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	chat := router.Group("")
	chat.Use(middleware.RequireUser)
	chat.GET("/conversations", conv.List)
	chat.GET("/conversations/:id", conv.Get)
	chat.POST("/chat", conv.Send)
}
