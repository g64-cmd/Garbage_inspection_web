package models

import (
	"encoding/json"
	"time"
)

// 基于 design.md 3.1.1 的状态上报 (边缘端 -> 云端)
type VehicleStatus struct {
	Timestamp int64    `json:"timestamp"`
	Position  Position `json:"position"`
	Battery   float64  `json:"battery"`
	State     string   `json:"state"`
}

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// 基于 design.md 3.3.2 的 WebSocket 遥测 (云端 -> 客户端)
// 注意：这个结构在发送给客户端时，由云端动态添加了 vehicle_id
type TelemetryUpdate struct {
	VehicleID string `json:"vehicle_id"`
	VehicleStatus
}

// 基于 design.md 3.1.2 的宏观指令 (云端 -> 边缘端)
type Command struct {
	CommandID string `json:"command_id"`
	Command   string `json:"command"`
	TaskID    string `json:"task_id,omitempty"`
}

// 基于 design.md 3.2.1 的决策请求元数据
type DecisionRequestMetadata struct {
	VehicleID string `json:"vehicle_id"`
	Timestamp int64  `json:"timestamp"`
}

// 基于 design.md 3.2.1 的决策响应
type DecisionResult struct {
	ImageID    string  `json:"image_id"`
	Action     string  `json:"action"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason,omitempty"`
}

// 基于 design.md 3.3.1 的客户端指令请求
type SendCommandRequest struct {
	VehicleID string `json:"vehicle_id" binding:"required"`
	Command   string `json:"command" binding:"required"`
}

// User 对应于数据库中的 'users' 表 (design.md 5.)
type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"` // (密码哈希不应被序列化到 JSON 中)
}

// Vehicle 对应于数据库中的 'vehicles' 表
type Vehicle struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Model         string         `json:"model"`
	CurrentStatus *VehicleStatus `json:"current_status"` // 使用指针以允许 null
}

// VehicleTelemetry 对应于 'vehicle_telemetry' 表，用于存储历史轨迹点
type VehicleTelemetry struct {
	ID        int64     `json:"id"`
	VehicleID string    `json:"vehicle_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Battery   float64   `json:"battery"`
	State     string    `json:"state"`
}

// DecisionLog 对应于数据库中的 'decision_logs' 表
type DecisionLog struct {
	ID              string          `json:"id"`
	VehicleID       string          `json:"vehicle_id"`
	Timestamp       time.Time       `json:"timestamp"`
	ImageURL        string          `json:"image_url"`
	ServerDecision  json.RawMessage `json:"server_decision"`
	RequestMetadata json.RawMessage `json:"request_metadata"`
}
