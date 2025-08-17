package router

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-prd-manager/server/handler"
	"github.com/trancecho/mundo-prd-manager/server/middleware"
)

func GenerateRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := r.Group("/api", middleware.JWTAuthMiddleware(), middleware.ResponseMiddleware())
	{
		v1.POST("/product/generate", handler.ProductGenerate)
		v1.GET("/getChatID", handler.GetPersonalChatID)
		v1.DELETE("/deleteChatID", handler.DeleteChatID)
	}
}
