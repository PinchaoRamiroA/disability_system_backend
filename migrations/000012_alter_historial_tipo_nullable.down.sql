-- Rollback migration 000012
DROP TABLE IF EXISTS tipo_historial CASCADE;
ALTER TABLE historial ALTER COLUMN tipo_historial DROP DEFAULT;
ALTER TABLE historial ALTER COLUMN tipo_historial SET NOT NULL;
