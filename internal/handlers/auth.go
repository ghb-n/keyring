package handlers

import (
	"net/http"

	"ghbn/gin-study/internal/auth"
	"ghbn/gin-study/internal/services"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type AuthHandler struct {
	userService    *services.UserService
	refreshService *services.RefreshService
	jwtSecret      string
}

func NewAuthHandler(userService *services.UserService, refreshService *services.RefreshService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		refreshService: refreshService,
		jwtSecret:      jwtSecret,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.refreshService.Generate(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.refreshService.Generate(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	})
}
