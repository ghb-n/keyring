package routes

import (
	"ghbn/gin-study/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	jwtSecret string,
) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
}
