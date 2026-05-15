CREATE TABLE IF NOT EXISTS auditoria (
    id_auditoria SERIAL PRIMARY KEY,
    id_usuario INT REFERENCES usuario(id_usuario) ON DELETE SET NULL,
    id_incapacidad INT REFERENCES incapacidad(id_incapacidad) ON DELETE SET NULL,
    tipo_accion VARCHAR(50) NOT NULL,
    modulo VARCHAR(50) NOT NULL,
    descripcion TEXT NOT NULL,
    cambio_anterior TEXT,
    cambio_nuevo TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_auditoria_id_usuario ON auditoria(id_usuario);
CREATE INDEX idx_auditoria_id_incapacidad ON auditoria(id_incapacidad);
CREATE INDEX idx_auditoria_modulo ON auditoria(modulo);
CREATE INDEX idx_auditoria_tipo_accion ON auditoria(tipo_accion);
