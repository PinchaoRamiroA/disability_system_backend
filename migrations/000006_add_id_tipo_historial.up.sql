-- Add id_tipo_historial column to historial
ALTER TABLE historial ADD COLUMN IF NOT EXISTS id_tipo_historial BIGINT;

-- Create index
CREATE INDEX IF NOT EXISTS idx_historial_id_tipo ON historial(id_tipo_historial);