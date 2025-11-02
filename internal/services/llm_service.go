package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"patrol-cloud/internal/models" // (需要定义 models.LLMPlan)
	"time"
)

const LLM_API_URL = "https://api.example.com/v1/chat/completions" // 示例 URL

// LLMService 遵循 4.2.5 设计
type LLMService struct {
	httpClient *http.Client
	apiKey     string
}

func NewLLMService(apiKey string) *LLMService {
	return &LLMService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     apiKey,
	}
}

// PlanMission 调用外部 LLM API
func (s *LLMService) PlanMission(prompt string) (/* *models.LLMPlan */ interface{}, error) {
	// 1. 构建请求体
	reqBodyMap := map[string]interface{}{
		"model": "gpt-4", // 示例模型
		"messages": []map[string]string{
			{"role": "system", "content": "You are a patrol route planner."},
			{"role": "user", "content": prompt},
		},
	}
	reqBodyBytes, _ := json.Marshal(reqBodyMap)

	// 2. 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", LLM_API_URL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 3. 发送请求
	httpResp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// 4. 解析响应
	var respBody interface{} // (应定义
	if err := json.NewDecoder(httpResp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	// 5. (将 respBody 转换为 models.LLMPlan)
	return respBody, nil
}
