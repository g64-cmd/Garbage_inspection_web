package api

import (
	"errors"
	"net/http"
	"patrol-cloud/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// VehicleHandler 负责处理与车辆相关的 API 请求
type VehicleHandler struct {
	repo db.Repository
}

// NewVehicleHandler 创建一个新的 VehicleHandler
func NewVehicleHandler(repo db.Repository) *VehicleHandler {
	return &VehicleHandler{repo: repo}
}

// HandleListVehicles 处理获取车辆列表的请求
func (h *VehicleHandler) HandleListVehicles(c *gin.Context) {
	vehicles, err := h.repo.ListVehicles(c.Request.Context())
	if err != nil {
		// 在实际应用中，这里也应该记录日志
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve vehicle list"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// HandleGetVehicleByID 处理获取单个车辆详情的请求
func (h *VehicleHandler) HandleGetVehicleByID(c *gin.Context) {
	vehicleID := c.Param("id")

	vehicle, err := h.repo.GetVehicleByID(c.Request.Context(), vehicleID)
	if err != nil {
		// 检查是否是“未找到记录”的特定错误
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "vehicle with the specified ID was not found"})
			return
		}
		// 对于所有其他类型的错误，返回内部服务器错误
		// 在实际应用中，这里应该记录详细的错误日志
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an internal error occurred while retrieving vehicle details"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}
