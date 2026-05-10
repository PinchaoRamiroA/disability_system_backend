-- database.sql

CREATE TABLE rol (
    id_rol BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    permisos JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE usuario (
    id_usuario BIGSERIAL PRIMARY KEY,
    id_rol BIGINT NOT NULL,
    nombre VARCHAR(150) NOT NULL,
    correo VARCHAR(150) UNIQUE NOT NULL,
    numero_celular VARCHAR(20),
    direccion VARCHAR(255),
    password_hash TEXT NOT NULL,
    numero_documento VARCHAR(30) UNIQUE NOT NULL,
    numero_acudiente VARCHAR(20),
    estado BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_rol) REFERENCES rol(id_rol) ON DELETE RESTRICT
);

CREATE TABLE empleado (
    id_usuario BIGINT PRIMARY KEY,
    puesto_trabajo VARCHAR(150) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),    
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON DELETE CASCADE
);

CREATE TABLE gerencia (
    id_usuario BIGINT PRIMARY KEY,
    puesto_trabajo VARCHAR(150) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON DELETE CASCADE
);

CREATE TABLE gestion_humana (
    id_usuario BIGINT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON DELETE CASCADE
);

CREATE TABLE estado_incapacidad (
    id_estado BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    descripcion TEXT,
    permite_transicion BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tipo_incapacidad (
    id_tipo BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    documentos_requeridos JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
);

CREATE TABLE entidad (
    id_entidad BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(150) UNIQUE NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    plazo_transcripcion_dias INTEGER,
    tiempo_maximo_pago_dias INTEGER,
    canal_atencion VARCHAR(150),
    requiere_transcripcion BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
);

CREATE TABLE incapacidad (
    id_incapacidad BIGSERIAL PRIMARY KEY,
    id_usuario BIGINT NOT NULL,
    id_estado BIGINT NOT NULL,
    id_tipo BIGINT NOT NULL,
    id_entidad BIGINT NOT NULL,
    titulo VARCHAR(200) NOT NULL,
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE,
    origen VARCHAR(100) NOT NULL,
    fecha_radicacion DATE,
    fecha_pago DATE,
    observaciones TEXT,
    created_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON DELETE RESTRICT,
    FOREIGN KEY (id_estado) REFERENCES estado_incapacidad(id_estado) ON DELETE RESTRICT,
    FOREIGN KEY (id_tipo) REFERENCES tipo_incapacidad(id_tipo) ON DELETE RESTRICT,
    FOREIGN KEY (id_entidad) REFERENCES entidad(id_entidad) ON DELETE RESTRICT,
    FOREIGN KEY (created_by) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL,
    CHECK (fecha_fin IS NULL OR fecha_fin >= fecha_inicio)
);

CREATE TABLE documento (
    id_documento BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    formato VARCHAR(20) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    comentario TEXT,
    fecha_carga TIMESTAMP DEFAULT NOW(),
    validado_por BIGINT,
    fecha_validacion TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE CASCADE,
    FOREIGN KEY (validado_por) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL
);

CREATE TABLE pago (
    id_pago BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    id_entidad BIGINT NOT NULL,
    descripcion TEXT,
    valor NUMERIC(14,2) NOT NULL,
    fecha_pago DATE NOT NULL,
    periodo_contable VARCHAR(20),
    conciliado BOOLEAN DEFAULT FALSE,
    registrado_por BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE CASCADE,
    FOREIGN KEY (id_entidad) REFERENCES entidad(id_entidad) ON DELETE RESTRICT,
    FOREIGN KEY (registrado_por) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL,
    CHECK (valor >= 0)
);

CREATE TABLE tipo_seguimiento (
    id_tipo_seguimiento BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    descripcion TEXT
);

CREATE TABLE seguimiento_cobro (
    id_seguimiento BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    id_tipo_seguimiento BIGINT NOT NULL,
    descripcion TEXT,
    fecha TIMESTAMP DEFAULT NOW(),
    resultado TEXT,
    gestionado_por BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE CASCADE,
    FOREIGN KEY (id_tipo_seguimiento) REFERENCES tipo_seguimiento(id_tipo_seguimiento) ON DELETE RESTRICT,
    FOREIGN KEY (gestionado_por) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL
);

CREATE TABLE tipo_historial (
    id_tipo_historial BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    descripcion TEXT
);

CREATE TABLE historial (
    id_historial BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    id_tipo_historial BIGINT NOT NULL,
    descripcion TEXT NOT NULL,
    fecha TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    gestor_id BIGINT,
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE CASCADE,
    FOREIGN KEY (id_tipo_historial) REFERENCES tipo_historial(id_tipo_historial) ON DELETE RESTRICT,
    FOREIGN KEY (gestor_id) REFERENCES usuario(id_usuario) ON DELETE SET NULL
);

CREATE TABLE tipo_notificacion (
    id_tipo_notificacion BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    descripcion TEXT
);

CREATE TABLE notificacion (
    id_notificacion BIGSERIAL PRIMARY KEY,
    id_usuario BIGINT NOT NULL,
    id_incapacidad BIGINT,
    id_tipo_notificacion BIGINT NOT NULL,
    mensaje TEXT NOT NULL,
    fecha TIMESTAMP DEFAULT NOW(),
    leida BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE SET NULL,
    FOREIGN KEY (id_tipo_notificacion) REFERENCES tipo_notificacion(id_tipo_notificacion) ON DELETE RESTRICT
);

-- Indices para optimizar consultas
CREATE INDEX idx_incapacidad_usuario ON incapacidad(id_usuario);
CREATE INDEX idx_incapacidad_estado ON incapacidad(id_estado);
CREATE INDEX idx_incapacidad_entidad ON incapacidad(id_entidad);
CREATE INDEX idx_incapacidad_fecha_inicio ON incapacidad(fecha_inicio);
CREATE INDEX idx_documento_incapacidad ON documento(id_incapacidad);
CREATE INDEX idx_historial_incapacidad ON historial(id_incapacidad);
