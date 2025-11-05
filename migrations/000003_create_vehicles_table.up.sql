-- 000003_create_vehicles_table.up.sql

CREATE TABLE IF NOT EXISTS vehicles (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    model VARCHAR(255),
    current_status JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
