package api

import (
	"net/http"
	"patrol-cloud/internal/db"
	"time"

	"github.com/gin-gonic/gin"
)

// TelemetryHandler 负责处理历史遥测数据相关的 API 请求
type TelemetryHandler struct {
	repo db.Repository
}

// NewTelemetryHandler 创建一个新的 TelemetryHandler
func NewTelemetryHandler(repo db.Repository) *TelemetryHandler {
	return &TelemetryHandler{repo: repo}
}

// HandleGetTelemetry 处理获取车辆历史轨迹的请求
func (h *TelemetryHandler) HandleGetTelemetry(c *gin.Context) {
	vehicleID := c.Param("id")

	// 解析时间范围查询参数
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time and end_time query parameters are required"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format; use RFC3339"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format; use RFC3339"})
		return
	}

	telemetry, err := h.repo.GetTelemetryByVehicleID(c.Request.Context(), vehicleID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get telemetry data"})
		return
	}

	c.JSON(http.StatusOK, telemetry)
}
