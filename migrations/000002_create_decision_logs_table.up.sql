-- 000002_create_decision_logs_table.up.sql

CREATE TABLE IF NOT EXISTS decision_logs (
    id VARCHAR(255) PRIMARY KEY,
    vehicle_id VARCHAR(255) NOT NULL,
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    image_url VARCHAR(255),
    server_decision JSONB,
    request_metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
