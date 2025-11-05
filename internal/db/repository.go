package db

import (
	"context"
	"encoding/json"
	"log"
	"patrol-cloud/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository 定义了数据库操作接口
type Repository interface {
	// User methods
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error

	// Decision Log methods
	LogDecision(ctx context.Context, result *models.DecisionResult, imageURL string, metadata models.DecisionRequestMetadata) error
	ListDecisionLogsByVehicleID(ctx context.Context, vehicleID string, page, pageSize int) ([]*models.DecisionLog, int, error)

	// Vehicle methods
	CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error
	GetVehicleByID(ctx context.Context, id string) (*models.Vehicle, error)
	ListVehicles(ctx context.Context) ([]*models.Vehicle, error)
	UpdateVehicleStatus(ctx context.Context, vehicleID string, status *models.VehicleStatus) error

	// Telemetry methods
	CreateTelemetryEntry(ctx context.Context, telemetry *models.VehicleTelemetry) error
	GetTelemetryByVehicleID(ctx context.Context, vehicleID string, startTime, endTime time.Time) ([]*models.VehicleTelemetry, error)
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

// --- User Methods ---

func (r *postgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, hashed_password) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Username, user.HashedPassword)
	if err != nil {
		log.Printf("ERROR: Failed to create user: %v", err)
	}
	return err
}

func (r *postgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, hashed_password FROM users WHERE username = $1`
	row := r.pool.QueryRow(ctx, query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.HashedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // 用户不存在，不视为错误
		}
		log.Printf("ERROR: Failed to query user by username: %v", err)
		return nil, err
	}

	return &user, nil
}

// --- Decision Log Methods ---

func (r *postgresRepository) LogDecision(ctx context.Context, result *models.DecisionResult, imageURL string, metadata models.DecisionRequestMetadata) error {
	metadataBytes, _ := json.Marshal(metadata)
	decisionBytes, _ := json.Marshal(result)

	query := `
		INSERT INTO decision_logs (id, vehicle_id, image_url, server_decision, request_metadata)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		result.ImageID,
		metadata.VehicleID,
		imageURL,
		decisionBytes,
		metadataBytes,
	)

	if err != nil {
		log.Printf("ERROR: Failed to execute LogDecision query: %v", err)
	}
	return err
}

func (r *postgresRepository) ListDecisionLogsByVehicleID(ctx context.Context, vehicleID string, page, pageSize int) ([]*models.DecisionLog, int, error) {
	// 1. Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM decision_logs WHERE vehicle_id = $1`
	if err := r.pool.QueryRow(ctx, countQuery, vehicleID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2. Get paginated results
	query := `
		SELECT id, vehicle_id, timestamp, image_url, server_decision, request_metadata
		FROM decision_logs
		WHERE vehicle_id = $1
		ORDER BY "timestamp" DESC
		LIMIT $2 OFFSET $3
	`
	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx, query, vehicleID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*models.DecisionLog
	for rows.Next() {
		var log models.DecisionLog
		if err := rows.Scan(&log.ID, &log.VehicleID, &log.Timestamp, &log.ImageURL, &log.ServerDecision, &log.RequestMetadata); err != nil {
			return nil, 0, err
		}
		logs = append(logs, &log)
	}

	return logs, total, nil
}

// --- Vehicle Methods ---

func (r *postgresRepository) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	query := `INSERT INTO vehicles (id, name, model) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, vehicle.ID, vehicle.Name, vehicle.Model)
	if err != nil {
		log.Printf("ERROR: Failed to create vehicle: %v", err)
	}
	return err
}

func (r *postgresRepository) GetVehicleByID(ctx context.Context, id string) (*models.Vehicle, error) {
	query := `SELECT id, name, model, current_status FROM vehicles WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	var v models.Vehicle
	err := row.Scan(&v.ID, &v.Name, &v.Model, &v.CurrentStatus)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *postgresRepository) ListVehicles(ctx context.Context) ([]*models.Vehicle, error) {
	query := `SELECT id, name, model, current_status FROM vehicles ORDER BY name`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []*models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		if err := rows.Scan(&v.ID, &v.Name, &v.Model, &v.CurrentStatus); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, &v)
	}
	return vehicles, nil
}

func (r *postgresRepository) UpdateVehicleStatus(ctx context.Context, vehicleID string, status *models.VehicleStatus) error {
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return err
	}
	query := `UPDATE vehicles SET current_status = $1 WHERE id = $2`
	_, err = r.pool.Exec(ctx, query, statusJSON, vehicleID)
	return err
}

// --- Telemetry Methods ---

func (r *postgresRepository) CreateTelemetryEntry(ctx context.Context, telemetry *models.VehicleTelemetry) error {
	query := `
		INSERT INTO vehicle_telemetry (vehicle_id, "timestamp", latitude, longitude, battery, state)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		telemetry.VehicleID,
		telemetry.Timestamp,
		telemetry.Latitude,
		telemetry.Longitude,
		telemetry.Battery,
		telemetry.State,
	)
	return err
}

func (r *postgresRepository) GetTelemetryByVehicleID(ctx context.Context, vehicleID string, startTime, endTime time.Time) ([]*models.VehicleTelemetry, error) {
	query := `
		SELECT id, vehicle_id, "timestamp", latitude, longitude, battery, state
		FROM vehicle_telemetry
		WHERE vehicle_id = $1 AND "timestamp" >= $2 AND "timestamp" <= $3
		ORDER BY "timestamp" ASC
	`
	rows, err := r.pool.Query(ctx, query, vehicleID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var telemetryEntries []*models.VehicleTelemetry
	for rows.Next() {
		var entry models.VehicleTelemetry
		if err := rows.Scan(&entry.ID, &entry.VehicleID, &entry.Timestamp, &entry.Latitude, &entry.Longitude, &entry.Battery, &entry.State); err != nil {
			return nil, err
		}
		telemetryEntries = append(telemetryEntries, &entry)
	}
	return telemetryEntries, nil
}
