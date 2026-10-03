package handlers

import (
	"errors"
	"net/http"

	"ghbn/gin-study/internal/auth"
	"ghbn/gin-study/internal/services"

	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshHandler struct {
	refreshService *services.RefreshService
	jwtSecret      string
}

func NewRefreshHandler(refreshService *services.RefreshService, jwtSecret string) *RefreshHandler {
	return &RefreshHandler{
		refreshService: refreshService,
		jwtSecret:      jwtSecret,
	}
}

func (h *RefreshHandler) Refresh(c *gin.Context) {
	var input RefreshRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRefresh, userID, err := h.refreshService.Rotate(c.Request.Context(), input.RefreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshToken) ||
			errors.Is(err, services.ErrExpiredRefreshToken) ||
			errors.Is(err, services.ErrRevokedRefreshToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	newAccess, err := auth.GenerateToken(userID, "", h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	})
}

func (h *RefreshHandler) Logout(c *gin.Context) {
	var input RefreshRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.refreshService.Validate(c.Request.Context(), input.RefreshToken)
	if err != nil {
		c.Status(http.StatusNoContent) // idempotent — token invalid or already revoked
		return
	}

	if err := h.refreshService.RevokeAll(c.Request.Context(), token.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}
