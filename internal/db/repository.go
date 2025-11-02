package db

import (
	"context"
	"encoding/json"
	"log"
	"patrol-cloud/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository 定义了数据库操作接口
type Repository interface {
	LogDecision(ctx context.Context, result *models.DecisionResult, imageURL string, metadata models.DecisionRequestMetadata) error
	// (其他方法... e.g., GetUser, CreateVehicle)
}

// postgresRepository 是 Repository 的 PG 实现
type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository (4.2.1) 创建一个新的仓库实例
func NewRepository(dsn string) (Repository, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	// (可以添加 Ping() 检查)
	return &postgresRepository{pool: pool}, nil
}

// LogDecision (4.2.3) 将决策写入 decision_logs 表
func (r *postgresRepository) LogDecision(ctx context.Context, result *models.DecisionResult, imageURL string, metadata models.DecisionRequestMetadata) error {

	// (将 metadata 序列化为 jsonb)
	metadataBytes, _ := json.Marshal(metadata)

	query := `
		INSERT INTO decision_logs (
			id, vehicle_id, timestamp, image_url, server_decision, request_metadata
		) VALUES (
			$1, $2, NOW(), $3, $4, $5
		)
	`
	// (我们使用 result.ImageID 作为 PK，如 3.2.1 所述)
	// (server_decision 可以是 result action 或完整的 result json)
	decisionBytes, _ := json.Marshal(result)

	_, err := r.pool.Exec(ctx, query,
		result.ImageID,
		metadata.VehicleID,
		imageURL,
		decisionBytes,  // (存为 jsonb)
		metadataBytes,  // (存为 jsonb)
	)

	if err != nil {
		log.Printf("ERROR: Failed to execute LogDecision query: %v", err)
	}
	return err
}
