-- database.sql - Esquema completo del sistema de incapacidades

-- =====================================================
-- TABLAS DE USUARIOS Y ROLES
-- =====================================================

CREATE TABLE rol (
    id_rol BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    permisos JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE permisos (
    id_permiso BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
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

-- =====================================================
-- TABLAS DE REFERENCIA (TABLAS)
-- =====================================================

CREATE TABLE tipo_documento (
    id_tipo_documento BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    descripcion TEXT,
    requerido BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE estado_documento (
    id_estado_documento BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL,
    descripcion TEXT,
    color VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tipo_entidad (
    id_tipo_entidad BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE canal_recepcion (
    id_canal_recepcion BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    descripcion TEXT,
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE canal_atencion_entidad (
    id_canal_atencion BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tipo_pago (
    id_tipo_pago BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE estado_pago (
    id_estado_pago BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL,
    descripcion TEXT,
    color VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE periodicidad_reporte (
    id_periodicidad BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL,
    dias INTEGER,
    descripcion TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- =====================================================
-- ENTIDADES, ESTADOS Y TIPOS
-- =====================================================

CREATE TABLE entidad (
    id_entidad BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(150) UNIQUE NOT NULL,
    tipo VARCHAR(50) NOT NULL,
    plazo_transcripcion_dias INTEGER,
    tiempo_maximo_pago_dias INTEGER,
    canal_atencion VARCHAR(150),
    canales_atencion JSONB,
    requiere_transcripcion BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
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
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- =====================================================
-- TABLA PRINCIPAL: INCAPACIDAD
-- =====================================================

CREATE TABLE incapacidad (
    id_incapacidad BIGSERIAL PRIMARY KEY,
    id_usuario BIGINT NOT NULL,
    id_estado BIGINT NOT NULL,
    id_tipo BIGINT NOT NULL,
    id_entidad BIGINT NOT NULL,
    canal_recepcion VARCHAR(100),
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

-- =====================================================
-- DOCUMENTOS
-- =====================================================

CREATE TABLE documento (
    id_documento BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    tipo_documento VARCHAR(100) NOT NULL,
    estado_documento VARCHAR(50) NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    formato VARCHAR(20) NOT NULL,
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

-- =====================================================
-- PAGOS Y SEGUIMIENTO
-- =====================================================

CREATE TABLE pago (
    id_pago BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    id_entidad BIGINT NOT NULL,
    tipo_pago VARCHAR(100) NOT NULL,
    estado_pago VARCHAR(50) NOT NULL,
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

CREATE TABLE seguimiento_cobro (
    id_seguimiento BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    tipo_seguimiento VARCHAR(50) NOT NULL,
    resultado_seguimiento VARCHAR(100),
    descripcion TEXT,
    fecha TIMESTAMP DEFAULT NOW(),
    gestionado_por BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE CASCADE,
    FOREIGN KEY (gestionado_por) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL
);

-- =====================================================
-- HISTORIAL Y NOTIFICACIONES
-- =====================================================

CREATE TABLE historial (
    id_historial BIGSERIAL PRIMARY KEY,
    id_incapacidad BIGINT NOT NULL,
    tipo_historial VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL,
    fecha TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    gestor_id BIGINT,
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE CASCADE,
    FOREIGN KEY (gestor_id) REFERENCES usuario(id_usuario) ON DELETE SET NULL
);

CREATE TABLE notificacion (
    id_notificacion BIGSERIAL PRIMARY KEY,
    id_usuario BIGINT NOT NULL,
    id_incapacidad BIGINT,
    tipo_notificacion VARCHAR(100) NOT NULL,
    mensaje TEXT NOT NULL,
    fecha TIMESTAMP DEFAULT NOW(),
    leida BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    FOREIGN KEY (id_incapacidad) REFERENCES incapacidad(id_incapacidad) ON DELETE SET NULL
);

-- =====================================================
-- ÍNDICES
-- =====================================================

CREATE INDEX idx_incapacidad_usuario ON incapacidad(id_usuario);
CREATE INDEX idx_incapacidad_estado ON incapacidad(id_estado);
CREATE INDEX idx_incapacidad_entidad ON incapacidad(id_entidad);
CREATE INDEX idx_incapacidad_tipo ON incapacidad(id_tipo);
CREATE INDEX idx_incapacidad_fecha_inicio ON incapacidad(fecha_inicio);
CREATE INDEX idx_incapacidad_fecha_fin ON incapacidad(fecha_fin);
CREATE INDEX idx_incapacidad_origen ON incapacidad(origen);
CREATE INDEX idx_incapacidad_canal_recepcion ON incapacidad(canal_recepcion);
CREATE INDEX idx_incapacidad_created_by ON incapacidad(created_by);
CREATE INDEX idx_incapacidad_fecha_radicacion ON incapacidad(fecha_radicacion);

CREATE INDEX idx_documento_incapacidad ON documento(id_incapacidad);
CREATE INDEX idx_documento_tipo ON documento(tipo_documento);
CREATE INDEX idx_documento_estado ON documento(estado_documento);
CREATE INDEX idx_documento_validado_por ON documento(validado_por);

CREATE INDEX idx_pago_incapacidad ON pago(id_incapacidad);
CREATE INDEX idx_pago_entidad ON pago(id_entidad);
CREATE INDEX idx_pago_tipo_pago ON pago(tipo_pago);
CREATE INDEX idx_pago_estado_pago ON pago(estado_pago);
CREATE INDEX idx_pago_fecha_pago ON pago(fecha_pago);

CREATE INDEX idx_seguimiento_incapacidad ON seguimiento_cobro(id_incapacidad);
CREATE INDEX idx_seguimiento_tipo ON seguimiento_cobro(tipo_seguimiento);
CREATE INDEX idx_seguimiento_gestionado_por ON seguimiento_cobro(gestionado_por);

CREATE INDEX idx_historial_incapacidad ON historial(id_incapacidad);
CREATE INDEX idx_historial_tipo ON historial(tipo_historial);
CREATE INDEX idx_historial_gestor ON historial(gestor_id);

CREATE INDEX idx_notificacion_usuario ON notificacion(id_usuario);
CREATE INDEX idx_notificacion_incapacidad ON notificacion(id_incapacidad);
CREATE INDEX idx_notificacion_leida ON notificacion(leida);

CREATE INDEX idx_usuario_rol ON usuario(id_rol);
CREATE INDEX idx_usuario_correo ON usuario(correo);
CREATE INDEX idx_usuario_documento ON usuario(numero_documento);

CREATE INDEX idx_entidad_tipo ON entidad(tipo);