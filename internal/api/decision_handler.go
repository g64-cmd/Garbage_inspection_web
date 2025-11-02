package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"patrol-cloud/internal/models"
	"patrol-cloud/internal/services"
	"github.com/gin-gonic/gin"
)

// DecisionHandler 封装了 4.2.2 中描述的决策逻辑
type DecisionHandler struct {
	decisionSvc *services.DecisionService
}

func NewDecisionHandler(svc *services.DecisionService) *DecisionHandler {
	return &DecisionHandler{decisionSvc: svc}
}

// HandleDecision 严格遵循 3.2.1 和 4.2.2 的设计
func (h *DecisionHandler) HandleDecision(c *gin.Context) {
	// 1. 解析 multipart/form-data
	// 1a. 解析 metadata
	metadataStr := c.PostForm("metadata")
	if metadataStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metadata is required"})
		return
	}

	var metadata models.DecisionRequestMetadata
	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metadata JSON"})
		return
	}

	// 1b. 解析 image
	imageFile, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer imageFile.Close()

	imageBytes, err := io.ReadAll(imageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image file"})
		return
	}

	// 2. 调用 Service 层处理
	result, err := h.decisionSvc.ProcessDecision(c.Request.Context(), imageBytes, metadata)
	if err != nil {
		// (根据错误类型返回 500 或其他)
		log.Printf("ERROR: ProcessDecision failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process decision"})
		return
	}

	// 3. 响应 (格式必须符合 3.2.1)
	c.JSON(http.StatusOK, result)
}
