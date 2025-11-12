package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"patrol-cloud/internal/api"
	"patrol-cloud/internal/background"
	"patrol-cloud/internal/config"
	"patrol-cloud/internal/db"
	//"patrol-cloud/internal/models"
	"patrol-cloud/internal/services"
	"patrol-cloud/internal/storage"
	"patrol-cloud/internal/tasks"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	//"github.com/gin-gonic/gin"
	//"github.com/google/uuid"
	//"golang.org/x/crypto/bcrypt"
)

// main 是应用程序的入口点。
// 它负责初始化所有组件、注入依赖并启动服务。
func main() {
	// --- 1. 初始化配置 ---
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded.")

	log.Println("Starting application...")

	// --- 2. 依赖初始化 ---
	repo, err := db.NewRepository(cfg.PGDsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database repository initialized.")

	// (为测试添加一个临时用户)
	createTempUser(repo)

	minioClient, err := storage.NewMinIOClient(cfg.MinIOEndpoint, cfg.MinIOAccessKey, cfg.MinIOSecretKey, false)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}
	log.Println("MinIO client initialized.")

	opts := mqtt.NewClientOptions().AddBroker(cfg.EMQXHost).SetClientID("patrol-cloud-server")
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer mqttClient.Disconnect(250)
	log.Println("MQTT client initialized.")

	// 初始化失败任务队列
	failedTaskQueue := tasks.NewFileQueue("failed_tasks.log")
	log.Println("Failed task queue initialized.")

	// --- 3. 服务和后台任务初始化 ---
	telemetryHub := services.NewTelemetryHub()
	go telemetryHub.Run()
	log.Println("Telemetry hub is running.")

	aiService, err := services.NewAIService(cfg.ONNXModelPath)
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	llmService := services.NewLLMService(cfg.LLMApiKey, cfg.LLMBaseURL)
	authService := services.NewAuthService(repo, []byte(cfg.JWTSecret))
	commandService := services.NewCommandService(mqttClient)
	decisionService := services.NewDecisionService(aiService, repo, minioClient, failedTaskQueue)

	log.Println("All services initialized.")

	// 启动 MQTT 监听器
	mqttListener := background.NewMQTTListener(mqttClient, telemetryHub.BroadcastChannel, repo)
	mqttListener.StartListening()
	log.Println("MQTT listener started.")

	// --- 4. HTTP 服务启动 ---
	router := api.SetupRouter(repo, authService, commandService, decisionService, llmService, telemetryHub, []byte(cfg.JWTSecret), cfg.WebsocketAllowedOrigins)

	server := &http.Server{
		Addr:    ":8888",
		Handler: router,
	}

	go func() {
		log.Println("HTTP server starting on port 8888...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// --- 5. 实现优雅停机 (Graceful Shutdown) ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting.")
}
