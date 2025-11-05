package api

import (
	"net/http"
	"patrol-cloud/internal/db"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list vehicles"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// HandleGetVehicleByID 处理获取单个车辆详情的请求
func (h *VehicleHandler) HandleGetVehicleByID(c *gin.Context) {
	vehicleID := c.Param("id")

	vehicle, err := h.repo.GetVehicleByID(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get vehicle"})
		return
	}

	if vehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}
