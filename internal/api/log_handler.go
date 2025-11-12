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
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'page' parameter: must be an integer"})
		return
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'pageSize' parameter: must be an integer"})
		return
	}

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
		"logs":  logs, // Changed from data to logs for consistency
		"total": total,
	})
}

// HandleListAllDecisionLogs 处理获取所有决策日志的请求
func (h *LogHandler) HandleListAllDecisionLogs(c *gin.Context) {
	logs, err := h.repo.ListAllDecisionLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list all decision logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": len(logs),
	})
}
