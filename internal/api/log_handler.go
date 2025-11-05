package api

import (
	"net/http"
	"patrol-cloud/internal/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LogHandler 负责处理日志相关的 API 请求
type LogHandler struct {
	repo db.Repository
}

// NewLogHandler 创建一个新的 LogHandler
func NewLogHandler(repo db.Repository) *LogHandler {
	return &LogHandler{repo: repo}
}

// HandleListDecisionLogs 处理获取决策日志列表的请求
func (h *LogHandler) HandleListDecisionLogs(c *gin.Context) {
	vehicleID := c.Param("id")

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	logs, total, err := h.repo.ListDecisionLogsByVehicleID(c.Request.Context(), vehicleID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list decision logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + pageSize - 1) / pageSize,
		},
	})
}
