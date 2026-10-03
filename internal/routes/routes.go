package routes

import (
	"ghbn/gin-study/internal/handlers"
	"ghbn/gin-study/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	refreshHandler *handlers.RefreshHandler,
	jwtSecret string,
) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/request", refreshHandler.Refresh)
	r.POST("/logout", refreshHandler.Logout)
	api := r.Group("/api")
	api.Use(middlewares.Auth(jwtSecret))
	{
		api.GET("/me", userHandler.GetMe)
		api.POST("/user/update", userHandler.UpdateMe)
		api.POST("/user/delete", userHandler.DeleteMe)
	}
}
