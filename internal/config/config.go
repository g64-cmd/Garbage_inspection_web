package config

import (
	"errors"
	"os"
)

// Config 保存了应用的所有配置
// 字段标签 `env` 用于指定对应的环境变量名
type Config struct {
	PGDsn          string
	EMQXHost       string
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	LLMApiKey      string
	LLMBaseURL     string
	ONNXModelPath  string
	JWTSecret      string
	WebsocketAllowedOrigins string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	cfg := &Config{
		PGDsn:          os.Getenv("PG_DSN"),
		EMQXHost:       os.Getenv("EMQX_HOST"),
		MinIOEndpoint:  os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		LLMApiKey:      os.Getenv("LLM_API_KEY"),
		LLMBaseURL:     os.Getenv("LLM_BASE_URL"),
		ONNXModelPath:  os.Getenv("ONNX_MODEL_PATH"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		WebsocketAllowedOrigins: os.Getenv("WEBSOCKET_ALLOWED_ORIGINS"),
	}

	// 验证必须的配置项
	if cfg.PGDsn == "" {
		return nil, errors.New("missing required environment variable: PG_DSN")
	}
	if cfg.EMQXHost == "" {
		return nil, errors.New("missing required environment variable: EMQX_HOST")
	}
	if cfg.MinIOEndpoint == "" {
		return nil, errors.New("missing required environment variable: MINIO_ENDPOINT")
	}
	if cfg.MinIOAccessKey == "" {
		return nil, errors.New("missing required environment variable: MINIO_ACCESS_KEY")
	}
	if cfg.MinIOSecretKey == "" {
		return nil, errors.New("missing required environment variable: MINIO_SECRET_KEY")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("missing required environment variable: JWT_SECRET")
	}
	if cfg.WebsocketAllowedOrigins == "" {
		return nil, errors.New("missing required environment variable: WEBSOCKET_ALLOWED_ORIGINS")
	}
	// LLM 服务是核心功能，所以 API Key 和 URL 也应该是必须的
	if cfg.LLMApiKey == "" {
		return nil, errors.New("missing required environment variable: LLM_API_KEY")
	}
	if cfg.LLMBaseURL == "" {
		return nil, errors.New("missing required environment variable: LLM_BASE_URL")
	}

	return cfg, nil
}
