-- 000004_add_fk_to_decision_logs.down.sql

ALTER TABLE decision_logs
DROP CONSTRAINT IF EXISTS fk_vehicles;
