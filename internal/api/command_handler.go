package api

import (
	"log"
	"net/http"
	"patrol-cloud/internal/models"
	"patrol-cloud/internal/services"
	"github.com/gin-gonic/gin"
)

// CommandHandler 负责处理 3.3.1 中定义的宏观指令
type CommandHandler struct {
	cmdSvc *services.CommandService
}

func NewCommandHandler(svc *services.CommandService) *CommandHandler {
	return &CommandHandler{cmdSvc: svc}
}

// HandleSendCommand 接收来自客户端的指令并将其转发到 CommandService
func (h *CommandHandler) HandleSendCommand(c *gin.Context) {
	var req models.SendCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层发布 MQTT 消息
	commandID, err := h.cmdSvc.SendCommand(req.VehicleID, req.Command, "") // task_id 可选
	if err != nil {
		log.Printf("ERROR: Failed to send command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue command"})
		return
	}

	// 响应 (遵循 3.3.1 的 202 Accepted)
	c.JSON(http.StatusAccepted, gin.H{
		"status":     "queued",
		"command_id": commandID,
	})
}
