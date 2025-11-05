package api

import (
	"net/http"
	"patrol-cloud/internal/services"

	"github.com/gin-gonic/gin"
)

// LLMHandler 负责处理与 LLM 相关的 API 请求
type LLMHandler struct {
	llmSvc *services.LLMService
}

// NewLLMHandler 创建一个新的 LLMHandler
func NewLLMHandler(svc *services.LLMService) *LLMHandler {
	return &LLMHandler{llmSvc: svc}
}

// PlanRequest 定义了 LLM 规划请求的 JSON 结构
type PlanRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// HandlePlan 处理 /llm/plan 请求 (对应 3.3.1)
func (h *LLMHandler) HandlePlan(c *gin.Context) {
	var req PlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: prompt is required"})
		return
	}

	plan, err := h.llmSvc.PlanMission(c.Request.Context(), req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate plan from LLM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan": plan})
}
