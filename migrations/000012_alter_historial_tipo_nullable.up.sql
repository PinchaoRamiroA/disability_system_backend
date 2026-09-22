-- Allow NULL for legacy tipo_historial column and set default empty string
ALTER TABLE historial ALTER COLUMN tipo_historial DROP NOT NULL;
ALTER TABLE historial ALTER COLUMN tipo_historial SET DEFAULT '';

-- Create tipo_historial catalog table if it does not exist
CREATE TABLE IF NOT EXISTS tipo_historial (
    id_tipo_historial BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    descripcion TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Seed basic categories
INSERT INTO tipo_historial (id_tipo_historial, nombre, descripcion) VALUES
(1, 'Creación', 'Creación inicial de la incapacidad o registro'),
(2, 'Cambio de estado', 'Actualización del estado de la incapacidad'),
(3, 'Carga de documento', 'Documento adjunto cargado'),
(4, 'Validación de documento', 'Documento aprobado y validado'),
(5, 'Rechazo de documento', 'Rechazo de documento')
ON CONFLICT (id_tipo_historial) DO NOTHING;

SELECT setval('tipo_historial_id_tipo_historial_seq', COALESCE((SELECT MAX(id_tipo_historial) FROM tipo_historial), 1));
