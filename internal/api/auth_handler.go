package api

import (
	"log"
	"net/http"
	"patrol-cloud/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler 负责处理用户认证相关的 HTTP 请求
type AuthHandler struct {
	authSvc *services.AuthService
}

// NewAuthHandler 创建一个新的 AuthHandler
func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: svc}
}

// LoginRequest 定义了登录请求的 JSON 结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HandleLogin 处理用户登录请求 (对应 3.3.1)
func (h *AuthHandler) HandleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		// Add detailed logging for debugging
		log.Printf("[AUTH_DEBUG] Login failed for user '%s'. Reason: %v", req.Username, err)

		switch err {
		case services.ErrUserNotFound, services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
