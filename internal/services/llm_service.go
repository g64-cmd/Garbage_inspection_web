package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// --- Qwen API Data Structures (OpenAI Compatible) ---

type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenRequest struct {
	Model    string        `json:"model"`
	Messages []QwenMessage `json:"messages"`
}

type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message ResponseMessage `json:"message"`
}

type QwenResponse struct {
	Choices []Choice `json:"choices"`
}

// --- LLM Service ---

// LLMService 遵循 4.2.5 设计，适配 Qwen API
type LLMService struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewLLMService(apiKey, baseURL string) *LLMService {
	return &LLMService{
		httpClient: &http.Client{Timeout: 60 * time.Second}, // (增加超时)
		apiKey:     apiKey,
		baseURL:    baseURL,
	}
}

// PlanMission 调用 Qwen API 生成巡检计划
func (s *LLMService) PlanMission(ctx context.Context, prompt string) (string, error) {
	// 1. 构建请求体
	reqPayload := QwenRequest{
		Model: "qwen-plus", // (模型可配置)
		Messages: []QwenMessage{
			{Role: "system", Content: "You are a highly intelligent patrol route planner for autonomous vehicles."},
			{Role: "user", Content: prompt},
		},
	}
	reqBodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", err
	}

	// 2. 创建 HTTP 请求
	endpoint := s.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 3. 发送请求
	httpResp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()

	// 4. 解析响应
	var respPayload QwenResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&respPayload); err != nil {
		return "", err
	}

	// 5. 提取并返回结果
	if len(respPayload.Choices) == 0 {
		return "", errors.New("no response choices from LLM")
	}

	return respPayload.Choices[0].Message.Content, nil
}
