-- 000004_add_transcripcion_fields.up.sql - Agregar campos de transcripción EPS/ARL

ALTER TABLE incapacidad
ADD COLUMN fecha_transcripcion TIMESTAMP,
ADD COLUMN transcrito_por BIGINT,
ADD COLUMN observaciones_transcripcion TEXT,
ADD COLUMN fecha_limite_transcripcion DATE,
ADD COLUMN estado_transcripcion VARCHAR(50) DEFAULT 'pendiente';

COMMENT ON COLUMN incapacidad.fecha_transcripcion IS 'Fecha y hora en que se realizó la transcripción';
COMMENT ON COLUMN incapacidad.transcrito_por IS 'ID del usuario que realizó la transcripción';
COMMENT ON COLUMN incapacidad.observaciones_transcripcion IS 'Observaciones adicionales del proceso de transcripción';
COMMENT ON COLUMN incapacidad.fecha_limite_transcripcion IS 'Fecha límite para realizar la transcripción (3 días hábiles después de creación)';
COMMENT ON COLUMN incapacidad.estado_transcripcion IS 'Estado de transcripción: pendiente, en_proceso, completado, vencida';

ALTER TABLE incapacidad
ADD CONSTRAINT fk_incapacidad_transcrito_por
FOREIGN KEY (transcrito_por) REFERENCES usuario(id_usuario) ON DELETE SET NULL;

CREATE INDEX idx_incapacidad_estado_transcripcion ON incapacidad(estado_transcripcion);
CREATE INDEX idx_incapacidad_fecha_limite_transcripcion ON incapacidad(fecha_limite_transcripcion);
