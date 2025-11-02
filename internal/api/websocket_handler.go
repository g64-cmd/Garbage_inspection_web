package api

import (
	"log"
	"net/http"
	"patrol-cloud/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader 负责将 HTTP 连接升级到 WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// (在生产中，这里应该检查 Origin 头部，只允许受信任的域)
		return true
	},
}

// WebSocketHandler 遵循 4.2.2 的设计
type WebSocketHandler struct {
	telemetryHub *services.TelemetryHub
}

func NewWebSocketHandler(hub *services.TelemetryHub) *WebSocketHandler {
	return &WebSocketHandler{telemetryHub: hub}
}

// HandleTelemetry 升级连接并将其注册到 Hub
func (h *WebSocketHandler) HandleTelemetry(c *gin.Context) {
	// 1. 升级连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	// 2. 将连接注册到 Hub (遵循 4.2.7)
	// NewClient 封装了连接和 Hub 的引用
	client := services.NewTelemetryClient(h.telemetryHub, conn)
	h.telemetryHub.RegisterClient(client)

	// 3. 启动读写 Goroutine
	// (ReadLoop 也在 4.2.2 中被提及，用于检测断开连接)
	go client.WriteLoop()
	go client.ReadLoop()
}
