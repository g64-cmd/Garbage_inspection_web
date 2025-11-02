package api

import (
	"patrol-cloud/internal/services"
	"github.com/gin-gonic/gin"
)

// SetupRouter 组装所有 API 路由，如 4.2.1 所述
func SetupRouter(
	authSvc *services.AuthService,
	cmdSvc *services.CommandService,
	decisionSvc *services.DecisionService,
	llmSvc *services.LLMService,
	telemetryHub *services.TelemetryHub,
) *gin.Engine {

	router := gin.Default()

	// 实例化 Handlers
	// (AuthHandler 和 LLMHandler 需要实现)
	// authHandler := NewAuthHandler(authSvc)
	// llmHandler := NewLLMHandler(llmSvc)
	commandHandler := NewCommandHandler(cmdSvc)
	decisionHandler := NewDecisionHandler(decisionSvc)
	wsHandler := NewWebSocketHandler(telemetryHub) // (4.2.2 WebSocketHandler)

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 3.3.1 用户认证
		// v1.POST("/auth/login", authHandler.Login)

		// 3.3.1 宏观指令 (交互 A)
		// (假设有 authMiddleware)
		v1.POST("/commands/send", commandHandler.HandleSendCommand)

		// 3.3.1 LLM 智能交互
		// v1.POST("/llm/plan", llmHandler.HandlePlan)

		// 3.2.1 同步决策 (交互 C)
		// 注意：这个接口没有 /auth/ 前缀，可能不需要 JWT？
		// 如果需要，应添加 authMiddleware
		v1.POST("/decisions/recognize", decisionHandler.HandleDecision)
	}

	// 3.3.2 WebSocket 实时遥测
	router.GET("/ws/telemetry", wsHandler.HandleTelemetry)

	return router
}
