-- 000005_create_vehicle_telemetry_table.up.sql

CREATE TABLE IF NOT EXISTS vehicle_telemetry (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id VARCHAR(255) NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    "timestamp" TIMESTAMPTZ NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    battery DOUBLE PRECISION NOT NULL,
    state VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create an index on vehicle_id and timestamp for faster queries
CREATE INDEX IF NOT EXISTS idx_vehicle_telemetry_vehicle_id_timestamp ON vehicle_telemetry(vehicle_id, "timestamp" DESC);
