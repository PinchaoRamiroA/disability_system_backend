-- 000004_add_transcripcion_fields.down.sql - Eliminar campos de transcripción

ALTER TABLE incapacidad
DROP CONSTRAINT IF EXISTS fk_incapacidad_transcrito_por,
DROP COLUMN IF EXISTS fecha_transcripcion,
DROP COLUMN IF EXISTS transcrito_por,
DROP COLUMN IF EXISTS observaciones_transcripcion,
DROP COLUMN IF EXISTS fecha_limite_transcripcion,
DROP COLUMN IF EXISTS estado_transcripcion;

DROP INDEX IF EXISTS idx_incapacidad_estado_transcripcion;
DROP INDEX IF EXISTS idx_incapacidad_fecha_limite_transcripcion;
