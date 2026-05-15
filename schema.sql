--
-- PostgreSQL database dump
--

\restrict TjsT8TI2gXisPpKNM5UX8WQlLJgb0OidrrzTHAFy9Te3hojJmkYRJH2FzSNhQdV

-- Dumped from database version 18.3 (Debian 18.3-1.pgdg13+1)
-- Dumped by pg_dump version 18.3 (Debian 18.3-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: canal_atencion_entidad; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.canal_atencion_entidad (
    id_canal_atencion bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    descripcion text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.canal_atencion_entidad OWNER TO root;

--
-- Name: canal_atencion_entidad_id_canal_atencion_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.canal_atencion_entidad_id_canal_atencion_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.canal_atencion_entidad_id_canal_atencion_seq OWNER TO root;

--
-- Name: canal_atencion_entidad_id_canal_atencion_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.canal_atencion_entidad_id_canal_atencion_seq OWNED BY public.canal_atencion_entidad.id_canal_atencion;


--
-- Name: canal_recepcion; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.canal_recepcion (
    id_canal_recepcion bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    descripcion text,
    activo boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.canal_recepcion OWNER TO root;

--
-- Name: canal_recepcion_id_canal_recepcion_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.canal_recepcion_id_canal_recepcion_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.canal_recepcion_id_canal_recepcion_seq OWNER TO root;

--
-- Name: canal_recepcion_id_canal_recepcion_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.canal_recepcion_id_canal_recepcion_seq OWNED BY public.canal_recepcion.id_canal_recepcion;


--
-- Name: documento; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.documento (
    id_documento bigint NOT NULL,
    id_incapacidad bigint NOT NULL,
    tipo_documento character varying(100) NOT NULL,
    estado_documento character varying(50) NOT NULL,
    nombre character varying(255) NOT NULL,
    url text NOT NULL,
    formato character varying(20) NOT NULL,
    comentario text,
    fecha_carga timestamp without time zone DEFAULT now(),
    validado_por bigint,
    fecha_validacion timestamp without time zone,
    is_deleted boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.documento OWNER TO root;

--
-- Name: documento_id_documento_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.documento_id_documento_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.documento_id_documento_seq OWNER TO root;

--
-- Name: documento_id_documento_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.documento_id_documento_seq OWNED BY public.documento.id_documento;


--
-- Name: empleado; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.empleado (
    id_usuario bigint NOT NULL,
    puesto_trabajo character varying(150) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.empleado OWNER TO root;

--
-- Name: entidad; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.entidad (
    id_entidad bigint NOT NULL,
    nombre character varying(150) NOT NULL,
    tipo character varying(50) NOT NULL,
    plazo_transcripcion_dias integer,
    tiempo_maximo_pago_dias integer,
    canal_atencion character varying(150),
    canales_atencion jsonb,
    requiere_transcripcion boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.entidad OWNER TO root;

--
-- Name: entidad_id_entidad_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.entidad_id_entidad_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.entidad_id_entidad_seq OWNER TO root;

--
-- Name: entidad_id_entidad_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.entidad_id_entidad_seq OWNED BY public.entidad.id_entidad;


--
-- Name: estado_documento; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.estado_documento (
    id_estado_documento bigint NOT NULL,
    nombre character varying(50) NOT NULL,
    descripcion text,
    color character varying(20),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.estado_documento OWNER TO root;

--
-- Name: estado_documento_id_estado_documento_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.estado_documento_id_estado_documento_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.estado_documento_id_estado_documento_seq OWNER TO root;

--
-- Name: estado_documento_id_estado_documento_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.estado_documento_id_estado_documento_seq OWNED BY public.estado_documento.id_estado_documento;


--
-- Name: estado_incapacidad; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.estado_incapacidad (
    id_estado bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    descripcion text,
    permite_transicion boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.estado_incapacidad OWNER TO root;

--
-- Name: estado_incapacidad_id_estado_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.estado_incapacidad_id_estado_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.estado_incapacidad_id_estado_seq OWNER TO root;

--
-- Name: estado_incapacidad_id_estado_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.estado_incapacidad_id_estado_seq OWNED BY public.estado_incapacidad.id_estado;


--
-- Name: estado_pago; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.estado_pago (
    id_estado_pago bigint NOT NULL,
    nombre character varying(50) NOT NULL,
    descripcion text,
    color character varying(20),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.estado_pago OWNER TO root;

--
-- Name: estado_pago_id_estado_pago_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.estado_pago_id_estado_pago_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.estado_pago_id_estado_pago_seq OWNER TO root;

--
-- Name: estado_pago_id_estado_pago_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.estado_pago_id_estado_pago_seq OWNED BY public.estado_pago.id_estado_pago;


--
-- Name: gerencia; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.gerencia (
    id_usuario bigint NOT NULL,
    puesto_trabajo character varying(150) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.gerencia OWNER TO root;

--
-- Name: gestion_humana; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.gestion_humana (
    id_usuario bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.gestion_humana OWNER TO root;

--
-- Name: historial; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.historial (
    id_historial bigint NOT NULL,
    id_incapacidad bigint NOT NULL,
    tipo_historial character varying(100) NOT NULL,
    descripcion text NOT NULL,
    fecha timestamp without time zone DEFAULT now(),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    gestor_id bigint,
    id_tipo_historial bigint,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.historial OWNER TO root;

--
-- Name: historial_id_historial_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.historial_id_historial_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.historial_id_historial_seq OWNER TO root;

--
-- Name: historial_id_historial_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.historial_id_historial_seq OWNED BY public.historial.id_historial;


--
-- Name: incapacidad; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.incapacidad (
    id_incapacidad bigint NOT NULL,
    id_usuario bigint NOT NULL,
    id_estado bigint NOT NULL,
    id_tipo bigint NOT NULL,
    id_entidad bigint NOT NULL,
    canal_recepcion character varying(100),
    titulo character varying(200) NOT NULL,
    fecha_inicio date NOT NULL,
    fecha_fin date,
    origen character varying(100) NOT NULL,
    fecha_radicacion date,
    fecha_pago date,
    observaciones text,
    created_by bigint,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false,
    fecha_transcripcion timestamp without time zone,
    transcrito_por bigint,
    observaciones_transcripcion text,
    fecha_limite_transcripcion date,
    estado_transcripcion character varying(50) DEFAULT 'pendiente'::character varying,
    CONSTRAINT incapacidad_check CHECK (((fecha_fin IS NULL) OR (fecha_fin >= fecha_inicio)))
);


ALTER TABLE public.incapacidad OWNER TO root;

--
-- Name: COLUMN incapacidad.fecha_transcripcion; Type: COMMENT; Schema: public; Owner: root
--

COMMENT ON COLUMN public.incapacidad.fecha_transcripcion IS 'Fecha y hora en que se realizó la transcripción';


--
-- Name: COLUMN incapacidad.transcrito_por; Type: COMMENT; Schema: public; Owner: root
--

COMMENT ON COLUMN public.incapacidad.transcrito_por IS 'ID del usuario que realizó la transcripción';


--
-- Name: COLUMN incapacidad.observaciones_transcripcion; Type: COMMENT; Schema: public; Owner: root
--

COMMENT ON COLUMN public.incapacidad.observaciones_transcripcion IS 'Observaciones adicionales del proceso de transcripción';


--
-- Name: COLUMN incapacidad.fecha_limite_transcripcion; Type: COMMENT; Schema: public; Owner: root
--

COMMENT ON COLUMN public.incapacidad.fecha_limite_transcripcion IS 'Fecha límite para realizar la transcripción (3 días hábiles después de creación)';


--
-- Name: COLUMN incapacidad.estado_transcripcion; Type: COMMENT; Schema: public; Owner: root
--

COMMENT ON COLUMN public.incapacidad.estado_transcripcion IS 'Estado de transcripción: pendiente, en_proceso, completado, vencida';


--
-- Name: incapacidad_id_incapacidad_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.incapacidad_id_incapacidad_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.incapacidad_id_incapacidad_seq OWNER TO root;

--
-- Name: incapacidad_id_incapacidad_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.incapacidad_id_incapacidad_seq OWNED BY public.incapacidad.id_incapacidad;


--
-- Name: notificacion; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.notificacion (
    id_notificacion bigint NOT NULL,
    id_usuario bigint NOT NULL,
    id_incapacidad bigint,
    tipo_notificacion character varying(100) NOT NULL,
    mensaje text NOT NULL,
    fecha timestamp without time zone DEFAULT now(),
    leida boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.notificacion OWNER TO root;

--
-- Name: notificacion_id_notificacion_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.notificacion_id_notificacion_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.notificacion_id_notificacion_seq OWNER TO root;

--
-- Name: notificacion_id_notificacion_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.notificacion_id_notificacion_seq OWNED BY public.notificacion.id_notificacion;


--
-- Name: pago; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.pago (
    id_pago bigint NOT NULL,
    id_incapacidad bigint NOT NULL,
    id_entidad bigint NOT NULL,
    tipo_pago character varying(100) NOT NULL,
    estado_pago character varying(50) NOT NULL,
    descripcion text,
    valor numeric(14,2) NOT NULL,
    fecha_pago date NOT NULL,
    periodo_contable character varying(20),
    conciliado boolean DEFAULT false,
    registrado_por bigint,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false,
    CONSTRAINT pago_valor_check CHECK ((valor >= (0)::numeric))
);


ALTER TABLE public.pago OWNER TO root;

--
-- Name: pago_id_pago_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.pago_id_pago_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pago_id_pago_seq OWNER TO root;

--
-- Name: pago_id_pago_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.pago_id_pago_seq OWNED BY public.pago.id_pago;


--
-- Name: periodicidad_reporte; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.periodicidad_reporte (
    id_periodicidad bigint NOT NULL,
    nombre character varying(50) NOT NULL,
    dias integer,
    descripcion text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.periodicidad_reporte OWNER TO root;

--
-- Name: periodicidad_reporte_id_periodicidad_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.periodicidad_reporte_id_periodicidad_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.periodicidad_reporte_id_periodicidad_seq OWNER TO root;

--
-- Name: periodicidad_reporte_id_periodicidad_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.periodicidad_reporte_id_periodicidad_seq OWNED BY public.periodicidad_reporte.id_periodicidad;


--
-- Name: permisos; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.permisos (
    id_permiso bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    descripcion text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.permisos OWNER TO root;

--
-- Name: permisos_id_permiso_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.permisos_id_permiso_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permisos_id_permiso_seq OWNER TO root;

--
-- Name: permisos_id_permiso_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.permisos_id_permiso_seq OWNED BY public.permisos.id_permiso;


--
-- Name: rol; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.rol (
    id_rol bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    permisos jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false
);


ALTER TABLE public.rol OWNER TO root;

--
-- Name: rol_id_rol_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.rol_id_rol_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.rol_id_rol_seq OWNER TO root;

--
-- Name: rol_id_rol_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.rol_id_rol_seq OWNED BY public.rol.id_rol;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO root;

--
-- Name: seguimiento_cobro; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.seguimiento_cobro (
    id_seguimiento bigint NOT NULL,
    id_incapacidad bigint NOT NULL,
    tipo_seguimiento character varying(50) NOT NULL,
    resultado_seguimiento character varying(100),
    descripcion text,
    fecha timestamp without time zone DEFAULT now(),
    gestionado_por bigint,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.seguimiento_cobro OWNER TO root;

--
-- Name: seguimiento_cobro_id_seguimiento_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.seguimiento_cobro_id_seguimiento_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seguimiento_cobro_id_seguimiento_seq OWNER TO root;

--
-- Name: seguimiento_cobro_id_seguimiento_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.seguimiento_cobro_id_seguimiento_seq OWNED BY public.seguimiento_cobro.id_seguimiento;


--
-- Name: tipo_documento; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tipo_documento (
    id_tipo_documento bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    descripcion text,
    requerido boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    codigo character varying(100) NOT NULL
);


ALTER TABLE public.tipo_documento OWNER TO root;

--
-- Name: tipo_documento_id_tipo_documento_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.tipo_documento_id_tipo_documento_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tipo_documento_id_tipo_documento_seq OWNER TO root;

--
-- Name: tipo_documento_id_tipo_documento_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.tipo_documento_id_tipo_documento_seq OWNED BY public.tipo_documento.id_tipo_documento;


--
-- Name: tipo_entidad; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tipo_entidad (
    id_tipo_entidad bigint NOT NULL,
    nombre character varying(50) NOT NULL,
    descripcion text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.tipo_entidad OWNER TO root;

--
-- Name: tipo_entidad_id_tipo_entidad_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.tipo_entidad_id_tipo_entidad_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tipo_entidad_id_tipo_entidad_seq OWNER TO root;

--
-- Name: tipo_entidad_id_tipo_entidad_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.tipo_entidad_id_tipo_entidad_seq OWNED BY public.tipo_entidad.id_tipo_entidad;


--
-- Name: tipo_incapacidad; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tipo_incapacidad (
    id_tipo bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    documentos_requeridos jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    origen character varying(50) NOT NULL
);


ALTER TABLE public.tipo_incapacidad OWNER TO root;

--
-- Name: tipo_incapacidad_id_tipo_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.tipo_incapacidad_id_tipo_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tipo_incapacidad_id_tipo_seq OWNER TO root;

--
-- Name: tipo_incapacidad_id_tipo_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.tipo_incapacidad_id_tipo_seq OWNED BY public.tipo_incapacidad.id_tipo;


--
-- Name: tipo_pago; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tipo_pago (
    id_tipo_pago bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    descripcion text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.tipo_pago OWNER TO root;

--
-- Name: tipo_pago_id_tipo_pago_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.tipo_pago_id_tipo_pago_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tipo_pago_id_tipo_pago_seq OWNER TO root;

--
-- Name: tipo_pago_id_tipo_pago_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.tipo_pago_id_tipo_pago_seq OWNED BY public.tipo_pago.id_tipo_pago;


--
-- Name: usuario; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.usuario (
    id_usuario bigint NOT NULL,
    id_rol bigint NOT NULL,
    nombre character varying(150) NOT NULL,
    correo character varying(150) NOT NULL,
    numero_celular character varying(20),
    direccion character varying(255),
    password_hash text NOT NULL,
    numero_documento character varying(30) NOT NULL,
    numero_acudiente character varying(20),
    estado boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false
);


ALTER TABLE public.usuario OWNER TO root;

--
-- Name: usuario_id_usuario_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.usuario_id_usuario_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.usuario_id_usuario_seq OWNER TO root;

--
-- Name: usuario_id_usuario_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.usuario_id_usuario_seq OWNED BY public.usuario.id_usuario;


--
-- Name: canal_atencion_entidad id_canal_atencion; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.canal_atencion_entidad ALTER COLUMN id_canal_atencion SET DEFAULT nextval('public.canal_atencion_entidad_id_canal_atencion_seq'::regclass);


--
-- Name: canal_recepcion id_canal_recepcion; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.canal_recepcion ALTER COLUMN id_canal_recepcion SET DEFAULT nextval('public.canal_recepcion_id_canal_recepcion_seq'::regclass);


--
-- Name: documento id_documento; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.documento ALTER COLUMN id_documento SET DEFAULT nextval('public.documento_id_documento_seq'::regclass);


--
-- Name: entidad id_entidad; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.entidad ALTER COLUMN id_entidad SET DEFAULT nextval('public.entidad_id_entidad_seq'::regclass);


--
-- Name: estado_documento id_estado_documento; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_documento ALTER COLUMN id_estado_documento SET DEFAULT nextval('public.estado_documento_id_estado_documento_seq'::regclass);


--
-- Name: estado_incapacidad id_estado; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_incapacidad ALTER COLUMN id_estado SET DEFAULT nextval('public.estado_incapacidad_id_estado_seq'::regclass);


--
-- Name: estado_pago id_estado_pago; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_pago ALTER COLUMN id_estado_pago SET DEFAULT nextval('public.estado_pago_id_estado_pago_seq'::regclass);


--
-- Name: historial id_historial; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.historial ALTER COLUMN id_historial SET DEFAULT nextval('public.historial_id_historial_seq'::regclass);


--
-- Name: incapacidad id_incapacidad; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad ALTER COLUMN id_incapacidad SET DEFAULT nextval('public.incapacidad_id_incapacidad_seq'::regclass);


--
-- Name: notificacion id_notificacion; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.notificacion ALTER COLUMN id_notificacion SET DEFAULT nextval('public.notificacion_id_notificacion_seq'::regclass);


--
-- Name: pago id_pago; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.pago ALTER COLUMN id_pago SET DEFAULT nextval('public.pago_id_pago_seq'::regclass);


--
-- Name: periodicidad_reporte id_periodicidad; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.periodicidad_reporte ALTER COLUMN id_periodicidad SET DEFAULT nextval('public.periodicidad_reporte_id_periodicidad_seq'::regclass);


--
-- Name: permisos id_permiso; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.permisos ALTER COLUMN id_permiso SET DEFAULT nextval('public.permisos_id_permiso_seq'::regclass);


--
-- Name: rol id_rol; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.rol ALTER COLUMN id_rol SET DEFAULT nextval('public.rol_id_rol_seq'::regclass);


--
-- Name: seguimiento_cobro id_seguimiento; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.seguimiento_cobro ALTER COLUMN id_seguimiento SET DEFAULT nextval('public.seguimiento_cobro_id_seguimiento_seq'::regclass);


--
-- Name: tipo_documento id_tipo_documento; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_documento ALTER COLUMN id_tipo_documento SET DEFAULT nextval('public.tipo_documento_id_tipo_documento_seq'::regclass);


--
-- Name: tipo_entidad id_tipo_entidad; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_entidad ALTER COLUMN id_tipo_entidad SET DEFAULT nextval('public.tipo_entidad_id_tipo_entidad_seq'::regclass);


--
-- Name: tipo_incapacidad id_tipo; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_incapacidad ALTER COLUMN id_tipo SET DEFAULT nextval('public.tipo_incapacidad_id_tipo_seq'::regclass);


--
-- Name: tipo_pago id_tipo_pago; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_pago ALTER COLUMN id_tipo_pago SET DEFAULT nextval('public.tipo_pago_id_tipo_pago_seq'::regclass);


--
-- Name: usuario id_usuario; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.usuario ALTER COLUMN id_usuario SET DEFAULT nextval('public.usuario_id_usuario_seq'::regclass);


--
-- Name: canal_atencion_entidad canal_atencion_entidad_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.canal_atencion_entidad
    ADD CONSTRAINT canal_atencion_entidad_nombre_key UNIQUE (nombre);


--
-- Name: canal_atencion_entidad canal_atencion_entidad_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.canal_atencion_entidad
    ADD CONSTRAINT canal_atencion_entidad_pkey PRIMARY KEY (id_canal_atencion);


--
-- Name: canal_recepcion canal_recepcion_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.canal_recepcion
    ADD CONSTRAINT canal_recepcion_nombre_key UNIQUE (nombre);


--
-- Name: canal_recepcion canal_recepcion_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.canal_recepcion
    ADD CONSTRAINT canal_recepcion_pkey PRIMARY KEY (id_canal_recepcion);


--
-- Name: documento documento_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.documento
    ADD CONSTRAINT documento_pkey PRIMARY KEY (id_documento);


--
-- Name: empleado empleado_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.empleado
    ADD CONSTRAINT empleado_pkey PRIMARY KEY (id_usuario);


--
-- Name: entidad entidad_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.entidad
    ADD CONSTRAINT entidad_nombre_key UNIQUE (nombre);


--
-- Name: entidad entidad_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.entidad
    ADD CONSTRAINT entidad_pkey PRIMARY KEY (id_entidad);


--
-- Name: estado_documento estado_documento_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_documento
    ADD CONSTRAINT estado_documento_nombre_key UNIQUE (nombre);


--
-- Name: estado_documento estado_documento_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_documento
    ADD CONSTRAINT estado_documento_pkey PRIMARY KEY (id_estado_documento);


--
-- Name: estado_incapacidad estado_incapacidad_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_incapacidad
    ADD CONSTRAINT estado_incapacidad_nombre_key UNIQUE (nombre);


--
-- Name: estado_incapacidad estado_incapacidad_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_incapacidad
    ADD CONSTRAINT estado_incapacidad_pkey PRIMARY KEY (id_estado);


--
-- Name: estado_pago estado_pago_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_pago
    ADD CONSTRAINT estado_pago_nombre_key UNIQUE (nombre);


--
-- Name: estado_pago estado_pago_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.estado_pago
    ADD CONSTRAINT estado_pago_pkey PRIMARY KEY (id_estado_pago);


--
-- Name: gerencia gerencia_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.gerencia
    ADD CONSTRAINT gerencia_pkey PRIMARY KEY (id_usuario);


--
-- Name: gestion_humana gestion_humana_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.gestion_humana
    ADD CONSTRAINT gestion_humana_pkey PRIMARY KEY (id_usuario);


--
-- Name: historial historial_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.historial
    ADD CONSTRAINT historial_pkey PRIMARY KEY (id_historial);


--
-- Name: incapacidad incapacidad_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad
    ADD CONSTRAINT incapacidad_pkey PRIMARY KEY (id_incapacidad);


--
-- Name: notificacion notificacion_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.notificacion
    ADD CONSTRAINT notificacion_pkey PRIMARY KEY (id_notificacion);


--
-- Name: pago pago_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.pago
    ADD CONSTRAINT pago_pkey PRIMARY KEY (id_pago);


--
-- Name: periodicidad_reporte periodicidad_reporte_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.periodicidad_reporte
    ADD CONSTRAINT periodicidad_reporte_nombre_key UNIQUE (nombre);


--
-- Name: periodicidad_reporte periodicidad_reporte_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.periodicidad_reporte
    ADD CONSTRAINT periodicidad_reporte_pkey PRIMARY KEY (id_periodicidad);


--
-- Name: permisos permisos_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.permisos
    ADD CONSTRAINT permisos_nombre_key UNIQUE (nombre);


--
-- Name: permisos permisos_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.permisos
    ADD CONSTRAINT permisos_pkey PRIMARY KEY (id_permiso);


--
-- Name: rol rol_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.rol
    ADD CONSTRAINT rol_nombre_key UNIQUE (nombre);


--
-- Name: rol rol_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.rol
    ADD CONSTRAINT rol_pkey PRIMARY KEY (id_rol);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: seguimiento_cobro seguimiento_cobro_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.seguimiento_cobro
    ADD CONSTRAINT seguimiento_cobro_pkey PRIMARY KEY (id_seguimiento);


--
-- Name: tipo_documento tipo_documento_codigo_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_documento
    ADD CONSTRAINT tipo_documento_codigo_key UNIQUE (codigo);


--
-- Name: tipo_documento tipo_documento_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_documento
    ADD CONSTRAINT tipo_documento_nombre_key UNIQUE (nombre);


--
-- Name: tipo_documento tipo_documento_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_documento
    ADD CONSTRAINT tipo_documento_pkey PRIMARY KEY (id_tipo_documento);


--
-- Name: tipo_entidad tipo_entidad_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_entidad
    ADD CONSTRAINT tipo_entidad_nombre_key UNIQUE (nombre);


--
-- Name: tipo_entidad tipo_entidad_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_entidad
    ADD CONSTRAINT tipo_entidad_pkey PRIMARY KEY (id_tipo_entidad);


--
-- Name: tipo_incapacidad tipo_incapacidad_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_incapacidad
    ADD CONSTRAINT tipo_incapacidad_nombre_key UNIQUE (nombre);


--
-- Name: tipo_incapacidad tipo_incapacidad_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_incapacidad
    ADD CONSTRAINT tipo_incapacidad_pkey PRIMARY KEY (id_tipo);


--
-- Name: tipo_pago tipo_pago_nombre_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_pago
    ADD CONSTRAINT tipo_pago_nombre_key UNIQUE (nombre);


--
-- Name: tipo_pago tipo_pago_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tipo_pago
    ADD CONSTRAINT tipo_pago_pkey PRIMARY KEY (id_tipo_pago);


--
-- Name: usuario usuario_correo_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.usuario
    ADD CONSTRAINT usuario_correo_key UNIQUE (correo);


--
-- Name: usuario usuario_numero_documento_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.usuario
    ADD CONSTRAINT usuario_numero_documento_key UNIQUE (numero_documento);


--
-- Name: usuario usuario_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.usuario
    ADD CONSTRAINT usuario_pkey PRIMARY KEY (id_usuario);


--
-- Name: idx_documento_estado; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_documento_estado ON public.documento USING btree (estado_documento);


--
-- Name: idx_documento_incapacidad; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_documento_incapacidad ON public.documento USING btree (id_incapacidad);


--
-- Name: idx_documento_tipo; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_documento_tipo ON public.documento USING btree (tipo_documento);


--
-- Name: idx_documento_validado_por; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_documento_validado_por ON public.documento USING btree (validado_por);


--
-- Name: idx_entidad_tipo; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_entidad_tipo ON public.entidad USING btree (tipo);


--
-- Name: idx_historial_gestor; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_historial_gestor ON public.historial USING btree (gestor_id);


--
-- Name: idx_historial_id_tipo; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_historial_id_tipo ON public.historial USING btree (id_tipo_historial);


--
-- Name: idx_historial_incapacidad; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_historial_incapacidad ON public.historial USING btree (id_incapacidad);


--
-- Name: idx_historial_tipo; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_historial_tipo ON public.historial USING btree (tipo_historial);


--
-- Name: idx_incapacidad_canal_recepcion; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_canal_recepcion ON public.incapacidad USING btree (canal_recepcion);


--
-- Name: idx_incapacidad_created_by; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_created_by ON public.incapacidad USING btree (created_by);


--
-- Name: idx_incapacidad_entidad; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_entidad ON public.incapacidad USING btree (id_entidad);


--
-- Name: idx_incapacidad_estado; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_estado ON public.incapacidad USING btree (id_estado);


--
-- Name: idx_incapacidad_estado_transcripcion; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_estado_transcripcion ON public.incapacidad USING btree (estado_transcripcion);


--
-- Name: idx_incapacidad_fecha_fin; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_fecha_fin ON public.incapacidad USING btree (fecha_fin);


--
-- Name: idx_incapacidad_fecha_inicio; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_fecha_inicio ON public.incapacidad USING btree (fecha_inicio);


--
-- Name: idx_incapacidad_fecha_limite_transcripcion; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_fecha_limite_transcripcion ON public.incapacidad USING btree (fecha_limite_transcripcion);


--
-- Name: idx_incapacidad_fecha_radicacion; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_fecha_radicacion ON public.incapacidad USING btree (fecha_radicacion);


--
-- Name: idx_incapacidad_origen; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_origen ON public.incapacidad USING btree (origen);


--
-- Name: idx_incapacidad_tipo; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_tipo ON public.incapacidad USING btree (id_tipo);


--
-- Name: idx_incapacidad_usuario; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_incapacidad_usuario ON public.incapacidad USING btree (id_usuario);


--
-- Name: idx_notificacion_incapacidad; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_notificacion_incapacidad ON public.notificacion USING btree (id_incapacidad);


--
-- Name: idx_notificacion_leida; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_notificacion_leida ON public.notificacion USING btree (leida);


--
-- Name: idx_notificacion_usuario; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_notificacion_usuario ON public.notificacion USING btree (id_usuario);


--
-- Name: idx_pago_entidad; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_pago_entidad ON public.pago USING btree (id_entidad);


--
-- Name: idx_pago_estado_pago; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_pago_estado_pago ON public.pago USING btree (estado_pago);


--
-- Name: idx_pago_fecha_pago; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_pago_fecha_pago ON public.pago USING btree (fecha_pago);


--
-- Name: idx_pago_incapacidad; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_pago_incapacidad ON public.pago USING btree (id_incapacidad);


--
-- Name: idx_pago_tipo_pago; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_pago_tipo_pago ON public.pago USING btree (tipo_pago);


--
-- Name: idx_seguimiento_gestionado_por; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_seguimiento_gestionado_por ON public.seguimiento_cobro USING btree (gestionado_por);


--
-- Name: idx_seguimiento_incapacidad; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_seguimiento_incapacidad ON public.seguimiento_cobro USING btree (id_incapacidad);


--
-- Name: idx_seguimiento_tipo; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_seguimiento_tipo ON public.seguimiento_cobro USING btree (tipo_seguimiento);


--
-- Name: idx_usuario_correo; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_usuario_correo ON public.usuario USING btree (correo);


--
-- Name: idx_usuario_documento; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_usuario_documento ON public.usuario USING btree (numero_documento);


--
-- Name: idx_usuario_rol; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_usuario_rol ON public.usuario USING btree (id_rol);


--
-- Name: documento documento_id_incapacidad_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.documento
    ADD CONSTRAINT documento_id_incapacidad_fkey FOREIGN KEY (id_incapacidad) REFERENCES public.incapacidad(id_incapacidad) ON DELETE CASCADE;


--
-- Name: documento documento_validado_por_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.documento
    ADD CONSTRAINT documento_validado_por_fkey FOREIGN KEY (validado_por) REFERENCES public.usuario(id_usuario) ON DELETE SET NULL;


--
-- Name: empleado empleado_id_usuario_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.empleado
    ADD CONSTRAINT empleado_id_usuario_fkey FOREIGN KEY (id_usuario) REFERENCES public.usuario(id_usuario) ON DELETE CASCADE;


--
-- Name: incapacidad fk_incapacidad_transcrito_por; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad
    ADD CONSTRAINT fk_incapacidad_transcrito_por FOREIGN KEY (transcrito_por) REFERENCES public.usuario(id_usuario) ON DELETE SET NULL;


--
-- Name: gerencia gerencia_id_usuario_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.gerencia
    ADD CONSTRAINT gerencia_id_usuario_fkey FOREIGN KEY (id_usuario) REFERENCES public.usuario(id_usuario) ON DELETE CASCADE;


--
-- Name: gestion_humana gestion_humana_id_usuario_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.gestion_humana
    ADD CONSTRAINT gestion_humana_id_usuario_fkey FOREIGN KEY (id_usuario) REFERENCES public.usuario(id_usuario) ON DELETE CASCADE;


--
-- Name: historial historial_gestor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.historial
    ADD CONSTRAINT historial_gestor_id_fkey FOREIGN KEY (gestor_id) REFERENCES public.usuario(id_usuario) ON DELETE SET NULL;


--
-- Name: historial historial_id_incapacidad_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.historial
    ADD CONSTRAINT historial_id_incapacidad_fkey FOREIGN KEY (id_incapacidad) REFERENCES public.incapacidad(id_incapacidad) ON DELETE CASCADE;


--
-- Name: incapacidad incapacidad_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad
    ADD CONSTRAINT incapacidad_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.usuario(id_usuario) ON DELETE SET NULL;


--
-- Name: incapacidad incapacidad_id_entidad_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad
    ADD CONSTRAINT incapacidad_id_entidad_fkey FOREIGN KEY (id_entidad) REFERENCES public.entidad(id_entidad) ON DELETE RESTRICT;


--
-- Name: incapacidad incapacidad_id_estado_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad
    ADD CONSTRAINT incapacidad_id_estado_fkey FOREIGN KEY (id_estado) REFERENCES public.estado_incapacidad(id_estado) ON DELETE RESTRICT;


--
-- Name: incapacidad incapacidad_id_tipo_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad
    ADD CONSTRAINT incapacidad_id_tipo_fkey FOREIGN KEY (id_tipo) REFERENCES public.tipo_incapacidad(id_tipo) ON DELETE RESTRICT;


--
-- Name: incapacidad incapacidad_id_usuario_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.incapacidad
    ADD CONSTRAINT incapacidad_id_usuario_fkey FOREIGN KEY (id_usuario) REFERENCES public.usuario(id_usuario) ON DELETE RESTRICT;


--
-- Name: notificacion notificacion_id_incapacidad_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.notificacion
    ADD CONSTRAINT notificacion_id_incapacidad_fkey FOREIGN KEY (id_incapacidad) REFERENCES public.incapacidad(id_incapacidad) ON DELETE SET NULL;


--
-- Name: notificacion notificacion_id_usuario_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.notificacion
    ADD CONSTRAINT notificacion_id_usuario_fkey FOREIGN KEY (id_usuario) REFERENCES public.usuario(id_usuario) ON DELETE CASCADE;


--
-- Name: pago pago_id_entidad_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.pago
    ADD CONSTRAINT pago_id_entidad_fkey FOREIGN KEY (id_entidad) REFERENCES public.entidad(id_entidad) ON DELETE RESTRICT;


--
-- Name: pago pago_id_incapacidad_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.pago
    ADD CONSTRAINT pago_id_incapacidad_fkey FOREIGN KEY (id_incapacidad) REFERENCES public.incapacidad(id_incapacidad) ON DELETE CASCADE;


--
-- Name: pago pago_registrado_por_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.pago
    ADD CONSTRAINT pago_registrado_por_fkey FOREIGN KEY (registrado_por) REFERENCES public.usuario(id_usuario) ON DELETE SET NULL;


--
-- Name: seguimiento_cobro seguimiento_cobro_gestionado_por_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.seguimiento_cobro
    ADD CONSTRAINT seguimiento_cobro_gestionado_por_fkey FOREIGN KEY (gestionado_por) REFERENCES public.usuario(id_usuario) ON DELETE SET NULL;


--
-- Name: seguimiento_cobro seguimiento_cobro_id_incapacidad_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.seguimiento_cobro
    ADD CONSTRAINT seguimiento_cobro_id_incapacidad_fkey FOREIGN KEY (id_incapacidad) REFERENCES public.incapacidad(id_incapacidad) ON DELETE CASCADE;


--
-- Name: usuario usuario_id_rol_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.usuario
    ADD CONSTRAINT usuario_id_rol_fkey FOREIGN KEY (id_rol) REFERENCES public.rol(id_rol) ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

\unrestrict TjsT8TI2gXisPpKNM5UX8WQlLJgb0OidrrzTHAFy9Te3hojJmkYRJH2FzSNhQdV

