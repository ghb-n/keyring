package main

import (
	"log"

	"ghbn/gin-study/internal/config"
	"ghbn/gin-study/internal/db"
	"ghbn/gin-study/internal/handlers"
	"ghbn/gin-study/internal/repository"
	"ghbn/gin-study/internal/routes"
	"ghbn/gin-study/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 2. DB
	pool, err := db.CreateConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// 3. Repository
	userRepo := repository.NewUserRepository(pool)
	refreshRepo := repository.NewRefreshRepository(pool)
	// 4. Service
	userService := services.NewUserService(userRepo)
	refreshService := services.NewRefreshService(refreshRepo)
	// 5. Handlers
	authHandler := handlers.NewAuthHandler(userService, refreshService, cfg.JWTSecret)

	userHandler := handlers.NewUserHandler(userService)

	// 6. Router
	r := gin.Default()

	// 7. Routes
	routes.Setup(r, authHandler, userHandler, cfg.JWTSecret)

	// 8. Run
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
