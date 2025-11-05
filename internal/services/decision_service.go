package services

import (
	"context"
	"log"
	"patrol-cloud/internal/db"
	"patrol-cloud/internal/models"
	"patrol-cloud/internal/storage"
	"patrol-cloud/internal/tasks"

	"github.com/google/uuid"
)

// DecisionService 遵循 4.2.3 的设计
type DecisionService struct {
	aiSvc      *AIService
	repo       db.Repository
	uploader   *storage.MinIOClient
	taskQueue *tasks.FileQueue
}

func NewDecisionService(ai *AIService, r db.Repository, s *storage.MinIOClient, tq *tasks.FileQueue) *DecisionService {
	return &DecisionService{
		aiSvc:    ai,
		repo:     r,
		uploader: s,
		taskQueue: tq,
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
	go s.logAndUploadAsync(result, image, metadata)

	// 3. (同步) 立即返回 AI 结果
	return result, nil
}

// FailedDecisionLogTask 定义了写入文件队列的任务结构
type FailedDecisionLogTask struct {
	Result   *models.DecisionResult        `json:"result"`
	ImageURL string                        `json:"image_url"`
	Metadata models.DecisionRequestMetadata `json:"metadata"`
}

// logAndUploadAsync 在后台处理图片上传和数据库日志记录
func (s *DecisionService) logAndUploadAsync(result *models.DecisionResult, image []byte, metadata models.DecisionRequestMetadata) {
	// 使用一个新的 background context，因为原始的 API 请求可能已经结束
	bgCtx := context.Background()

	// 1. 上传图片到 MinIO
	fileName := result.ImageID + ".jpg"
	imageURL, err := s.uploader.Upload(bgCtx, "decisions", fileName, image, "image/jpeg")
	if err != nil {
		// 使用结构化日志记录后台任务的失败
		log.Printf(
			"level=error msg=\"background task failed: image upload\" image_id=%s vehicle_id=%s error=\"%v\"",
			result.ImageID,
			metadata.VehicleID,
			err,
		)
		// 即使上传失败，我们仍然尝试记录日志，imageURL 会是空的
	}

	// 2. 记录日志到数据库
	if err := s.repo.LogDecision(bgCtx, result, imageURL, metadata); err != nil {
		log.Printf(
			"level=error msg=\"background task failed: log decision\" image_id=%s vehicle_id=%s error=\"%v\"",
			result.ImageID,
			metadata.VehicleID,
			err,
		)

		// 将失败的任务写入文件队列
		failedTask := FailedDecisionLogTask{
			Result:   result,
			ImageURL: imageURL,
			Metadata: metadata,
		}
		if qErr := s.taskQueue.LogFailedTask(failedTask); qErr != nil {
			log.Printf(
				"level=critical msg=\"FATAL: could not write failed task to queue\" image_id=%s error=\"%v\"",
				result.ImageID,
				qErr,
			)
		}
		return // 日志记录失败后，终止后续操作
	}

	log.Printf("level=info msg=\"background task complete: decision logged\" image_id=%s", result.ImageID)
}
