package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// main 是应用程序的入口点。
// 它负责初始化所有组件、注入依赖并启动服务。
func main() {
	// --- 1. 初始化配置 (从环境变量加载) ---
	// 在生产环境中，这些值应该来自环境变量或配置文件
	pgDsn := os.Getenv("PG_DSN")
	if pgDsn == "" {
		pgDsn = "postgres://user:password@localhost:5432/garbage_db?sslmode=disable" // 默认值
	}
	emqxHost := os.Getenv("EMQX_HOST")
	if emqxHost == "" {
		emqxHost = "tcp://localhost:1883"
	}
	// ... 其他配置，如 MinIO, LLM API Key 等

	log.Println("Starting application...")

	// --- 2. 依赖初始化 (按照 design.md 的顺序) ---
	// db_pool = db.InitPostgreSQL(PG_DSN)
	// repo = db.NewRepository(db_pool)
	// storage = storage.InitMinIO(...)
	// mqtt_client = background.InitMQTTClient(EMQX_HOST, ...)
	// ...
	// 实际项目中，这些函数会位于各自的包中
	// 例如: db.InitPostgreSQL(), services.NewTelemetryHub()
	// 这里我们使用占位符
	dbPool, err := InitPostgreSQL(pgDsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close() // 确保数据库连接最后被关闭
	log.Println("Database connection established.")

	// 初始化 MQTT 客户端
	mqttClient, err := InitMQTTClient(emqxHost)
	if err != nil {
		log.Fatalf("Failed to initialize MQTT client: %v", err)
	}
	defer mqttClient.Disconnect(250) // 优雅停机时断开连接
	log.Println("MQTT client initialized.")

	// --- 3. 服务和后台任务初始化 ---
	// telemetry_hub = services.NewTelemetryHub()
	// go telemetry_hub.Run()
	telemetryHub := NewTelemetryHub()
	go telemetryHub.Run()
	log.Println("Telemetry hub is running.")

	// ai_service = services.NewAIService(ONNX_MODEL_PATH)
	// llm_service = services.NewLLMService(LLM_API_KEY)
	// ... (初始化其他 services, 注入依赖)
	// 同样，这里使用占位符
	authService := NewAuthService()
	commandService := NewCommandService(mqttClient)
	decisionService := NewDecisionService()

	// listener = background.NewMQTTListener(mqtt_client, telemetry_hub.BroadcastChannel)
	// go listener.StartListening()
	mqttListener := NewMQTTListener(mqttClient, telemetryHub.BroadcastChannel)
	go mqttListener.StartListening()
	log.Println("MQTT listener started.")

	// --- 4. HTTP 服务启动 ---
	// router = api.SetupRouter(auth_service, command_service, ...)
	router := SetupRouter(authService, commandService, decisionService, telemetryHub)

	// router.Run(":8080")
	// 采用优雅停机的方式启动服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("HTTP server starting on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// --- 5. 实现优雅停机 (Graceful Shutdown) ---
	quit := make(chan os.Signal, 1)
	// 监听 syscall.SIGINT 和 syscall.SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞，直到接收到信号

	log.Println("Shutting down server...")

	// 创建一个有超时的上下文，用于通知服务器在 5 秒内完成现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting.")
}

// --- 以下是占位符，用于模拟真实模块 ---

// 模拟数据库连接池
type DBPool struct{}

func (p *DBPool) Close() {}
func InitPostgreSQL(dsn string) (*DBPool, error) {
	// 实际会使用 pgx 或 database/sql
	return &DBPool{}, nil
}

// 模拟 MQTT 客户端
type MQTTClient struct{}

func (c *MQTTClient) Disconnect(quiesce uint) {}
func InitMQTTClient(host string) (*MQTTClient, error) {
	// 实际会使用 paho.mqtt.golang
	return &MQTTClient{}, nil
}

// 模拟 TelemetryHub
type TelemetryHub struct {
	BroadcastChannel chan []byte
}

func NewTelemetryHub() *TelemetryHub {
	return &TelemetryHub{BroadcastChannel: make(chan []byte)}
}
func (h *TelemetryHub) Run() {
	// 模拟后台运行
	for {
		time.Sleep(1 * time.Second)
	}
}

// 模拟 MQTTListener
type MQTTListener struct{}

func NewMQTTListener(client *MQTTClient, channel chan<- []byte) *MQTTListener {
	return &MQTTListener{}
}
func (l *MQTTListener) StartListening() {
	log.Println("Mock MQTTListener is 'listening' to topics.")
}

// 模拟业务服务
type AuthService struct{}

func NewAuthService() *AuthService { return &AuthService{} }

type CommandService struct{}

func NewCommandService(client *MQTTClient) *CommandService { return &CommandService{} }

type DecisionService struct{}

func NewDecisionService() *DecisionService { return &DecisionService{} }

// 模拟 Gin 路由设置
func SetupRouter(auth *AuthService, cmd *CommandService, dec *DecisionService, hub *TelemetryHub) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// 实际会在这里根据 design.md 定义所有 API 路由
	// e.g., api.SetupRouter(r, auth, cmd, dec, hub)
	log.Println("Router setup complete.")
	return r
}
