//go:build !onnx

package services

import (
	"log"
	"patrol-cloud/internal/models"
)

// AIService is a stub implementation used when the onnx build tag is not set.
type AIService struct {
	modelPath string
}

// NewAIService creates a new stub AIService.
func NewAIService(modelPath string) (*AIService, error) {
	log.Printf("INFO: (STUB) AIService 'initialized' with model %s", modelPath)
	return &AIService{modelPath: modelPath}, nil
}

// Recognize returns a mocked decision without performing any real inference.
func (s *AIService) Recognize(image []byte) (*models.DecisionResult, error) {
	log.Println("INFO: (STUB) AIService.Recognize called")

	// 模拟一个高置信度的 "pickup" 决策
	result := &models.DecisionResult{
		ImageID:    "", // DecisionService 将填充此项
		Action:     "pickup",
		Confidence: 0.95,
		Reason:     "is_trash_type_A (stubbed)",
	}

	return result, nil
}
