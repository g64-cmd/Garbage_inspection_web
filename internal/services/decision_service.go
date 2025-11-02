package services

import (
	"context"
	"log"
	"patrol-cloud/internal/db"
	"patrol-cloud/internal/models"
	"patrol-cloud/internal/storage"

	"github.com/google/uuid"
)

// DecisionService 遵循 4.2.3 的设计
type DecisionService struct {
	aiSvc    *AIService
	repo     db.Repository
	uploader *storage.MinIOClient
}

func NewDecisionService(ai *AIService, r db.Repository, s *storage.MinIOClient) *DecisionService {
	return &DecisionService{
		aiSvc:    ai,
		repo:     r,
		uploader: s,
	}
}

// ProcessDecision 编排同步 AI 决策和异步日志记录
func (s *DecisionService) ProcessDecision(ctx context.Context, image []byte, metadata models.DecisionRequestMetadata) (*models.DecisionResult, error) {

	// 1. (同步) 调用 AI 服务进行识别
	result, err := s.aiSvc.Recognize(image)
	if err != nil {
		log.Printf("ERROR: AIService.Recognize failed: %v", err)
		return nil, err
	}

	// 确保 ImageID 已生成
	if result.ImageID == "" {
		result.ImageID = uuid.NewString()
	}

	// 2. (异步) 启动 Goroutine 上传图片和记录日志
	// 这种“即发即忘”的设计遵循 4.2.3，确保了 API 的低延迟
	go func() {
		// (使用 background context，因为原始请求可能已结束)
		bgCtx := context.Background()

		// 2a. 上传图片到 MinIO
		// (文件名可以使用 image_id)
		fileName := result.ImageID + ".jpg"
		imageURL, err := s.uploader.Upload(bgCtx, "decisions", fileName, image, "image/jpeg")
		if err != nil {
			log.Printf("ERROR: Failed to upload image %s: %v", fileName, err)
			// 即使上传失败，我们仍然尝试记录日志
		}

		// 2b. 记录日志到数据库
		err = s.repo.LogDecision(bgCtx, result, imageURL, metadata)
		if err != nil {
			log.Printf("ERROR: Failed to log decision %s: %v", result.ImageID, err)
		}

		log.Printf("INFO: Decision %s processed and logged.", result.ImageID)
	}()

	// 3. (同步) 立即返回 AI 结果
	return result, nil
}
