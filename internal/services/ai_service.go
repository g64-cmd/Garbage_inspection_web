package services

// #cgo CFLAGS: -I/path/to/onnxruntime/include
// #cgo LDFLAGS: -L/path/to/onnxruntime/lib -lonnxruntime
// #include <onnxruntime_c_api.h>
//
// (这里需要 CGo 包装函数来调用 run_onnx_inference)
import "C"

import (
	"log"
	"patrol-cloud/internal/models"
)

// AIService 遵循 4.2.4，封装 ONNX CGo 调用
type AIService struct {
	// onnx_session C.OrtSession (在 NewAIService 中初始化)
	modelPath string
}

func NewAIService(modelPath string) (*AIService, error) {
	// (此处应包含 CGo/ONNX 的真实初始化逻辑)
	// C.InitORTEnv()
	// C.CreateSession(modelPath)
	log.Printf("INFO: (STUB) AIService 'initialized' with model %s", modelPath)
	return &AIService{modelPath: modelPath}, nil
}

// Recognize 模拟 CGo 推理
func (s *AIService) Recognize(image []byte) (*models.DecisionResult, error) {
	log.Println("INFO: (STUB) AIService.Recognize called")

	// 1. PreProcess(image) -> tensor_input
	// 2. tensor_output = C.run_onnx_inference(...)
	// 3. result = PostProcess(tensor_output)

	// --- 模拟实现 ---
	// 模拟一个高置信度的 "pickup" 决策
	result := &models.DecisionResult{
		ImageID:    "", // DecisionService 将填充此项
		Action:     "pickup",
		Confidence: 0.95,
		Reason:     "is_trash_type_A (stubbed)",
	}
	// --- 结束模拟 ---

	return result, nil
}
