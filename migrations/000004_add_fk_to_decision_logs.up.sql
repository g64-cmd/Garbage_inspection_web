-- 000004_add_fk_to_decision_logs.up.sql

ALTER TABLE decision_logs
ADD CONSTRAINT fk_vehicles
FOREIGN KEY (vehicle_id)
REFERENCES vehicles(id)
ON DELETE CASCADE;
