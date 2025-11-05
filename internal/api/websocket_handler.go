package api

import (
	"log"
	"net/http"
	"patrol-cloud/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler 遵循 4.2.2 的设计
type WebSocketHandler struct {
	telemetryHub *services.TelemetryHub
	authSvc      *services.AuthService
	upgrader     websocket.Upgrader
}

func NewWebSocketHandler(hub *services.TelemetryHub, authSvc *services.AuthService, allowedOrigins string) *WebSocketHandler {
	allowedOriginMap := make(map[string]bool)
	for _, origin := range strings.Split(allowedOrigins, ",") {
		if origin != "" {
			allowedOriginMap[strings.TrimSpace(origin)] = true
		}
	}

	h := &WebSocketHandler{
		telemetryHub: hub,
		authSvc:      authSvc,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if allowedOriginMap[origin] {
					return true
				}
				log.Printf("WARN: WebSocket connection denied for origin: %s", origin)
				return false
			},
		},
	}
	return h
}

// HandleTelemetry 升级连接并将其注册到 Hub
func (h *WebSocketHandler) HandleTelemetry(c *gin.Context) {
	// 1. 从查询参数中获取 token 并验证
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}

	token, err := h.authSvc.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 2. 升级连接
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// CheckOrigin 失败时，Upgrade 会自动发送 403 Forbidden，这里只记录日志
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	// 3. 将连接注册到 Hub (遵循 4.2.7)
	client := services.NewTelemetryClient(h.telemetryHub, conn)
	h.telemetryHub.RegisterClient(client)

	// 4. 启动读写 Goroutine
	go client.WriteLoop()
	go client.ReadLoop()
}
