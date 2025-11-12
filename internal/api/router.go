package api

import (
	"patrol-cloud/internal/api/middleware"
	"patrol-cloud/internal/db"
	"patrol-cloud/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRouter 组装所有 API 路由，如 4.2.1 所述
func SetupRouter(
	repo db.Repository,
	authSvc *services.AuthService,
	cmdSvc *services.CommandService,
	decisionSvc *services.DecisionService,
	llmSvc *services.LLMService,
	telemetryHub *services.TelemetryHub,
	jwtSecret []byte,
	websocketAllowedOrigins string,
) *gin.Engine {

	router := gin.Default()

	// 实例化 Handlers
	authHandler := NewAuthHandler(authSvc)
	llmHandler := NewLLMHandler(llmSvc)
	commandHandler := NewCommandHandler(cmdSvc)
	decisionHandler := NewDecisionHandler(decisionSvc)
	wsHandler := NewWebSocketHandler(telemetryHub, authSvc, websocketAllowedOrigins)
	vehicleHandler := NewVehicleHandler(repo)
	telemetryHandler := NewTelemetryHandler(repo)
	logHandler := NewLogHandler(repo)

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 3.3.1 用户认证 (公开路由)
		v1.POST("/auth/login", authHandler.HandleLogin)

		// 创建需要认证的路由组
		authRequired := v1.Group("/")
		authRequired.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// 指令
			authRequired.POST("/commands/send", commandHandler.HandleSendCommand)

			// LLM
			authRequired.POST("/llm/plan", llmHandler.HandlePlan)

			// 同步决策
			authRequired.POST("/decisions/recognize", decisionHandler.HandleDecision)

			// 车辆
			authRequired.GET("/vehicles", vehicleHandler.HandleListVehicles)
			authRequired.GET("/vehicles/:id", vehicleHandler.HandleGetVehicleByID)

			// 遥测
			authRequired.GET("/vehicles/:id/telemetry", telemetryHandler.HandleGetTelemetry)

			// 日志
			authRequired.GET("/decision-logs", logHandler.HandleListAllDecisionLogs) // New global log route
			authRequired.GET("/vehicles/:id/decision-logs", logHandler.HandleListDecisionLogs)
		}
	}

	// WebSocket 实时遥测
	router.GET("/ws/telemetry", wsHandler.HandleTelemetry)

	return router
}
