--
-- PostgreSQL database dump
--

\restrict 7tW6PzhroU44F9PIeVQbEmRw7j9pgs3NOcTIi35V8cTC1L4oePuW6dW0XtlEcNw

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
-- Data for Name: canal_atencion_entidad; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.canal_atencion_entidad (id_canal_atencion, nombre, descripcion, created_at, updated_at) FROM stdin;
1	Portal web	Atención por portal web	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	Correo electrónico	Atención por correo electrónico	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	Atención presencial	Atención presencial	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Línea telefónica	Atención por línea telefónica	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	Mesa de ayuda	Atención por mesa de ayuda	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: canal_recepcion; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.canal_recepcion (id_canal_recepcion, nombre, descripcion, activo, created_at, updated_at) FROM stdin;
1	Recepción física	Recibido en recepción física	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	Oficina principal	Recibido en oficina principal	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	Líder comercial	Recibido por líder comercial	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Correo electrónico	Recibido por correo electrónico	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	Portal EPS	Recibido desde portal EPS	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
6	Portal ARL	Recibido desde portal ARL	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: documento; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.documento (id_documento, id_incapacidad, tipo_documento, estado_documento, nombre, url, formato, comentario, fecha_carga, validado_por, fecha_validacion, is_deleted, created_at, updated_at) FROM stdin;
1	1	certificado	Validado	Proyecto_Incapacidades (3).pdf	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/1/1778633794_proyecto_incapacidades_(3).pdf	.pdf		2026-05-12 19:56:35.21817	3	2026-05-12 21:49:37.282389	f	2026-05-12 19:56:35.218679	2026-05-12 21:49:37.282619
2	1	certificado	Pendiente	Diagrama de clases - incapacidades.drawio.png	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/1/1778721779_diagrama_de_clases_-_incapacidades.drawio.png	.png	\N	2026-05-13 20:23:01.929118	\N	\N	f	2026-05-13 20:23:01.929437	2026-05-13 20:23:01.929437
3	1	certificado	Pendiente	primera evaluación IS CB 26-1 .pdf	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/1/1778727300_primera_evaluación_is_cb_26-1_.pdf	.pdf	\N	2026-05-13 21:55:01.076519	\N	\N	f	2026-05-13 21:55:01.076671	2026-05-13 21:55:01.076671
4	4	Certificado de incapacidad	Pendiente	Diccionario Clases Postgresql Incapacidades.pdf	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/4/1778730637_diccionario_clases_postgresql_incapacidades.pdf	.pdf	\N	2026-05-13 22:50:38.080957	\N	\N	f	2026-05-13 22:50:38.081622	2026-05-13 22:50:38.081622
5	4	Historia clínica	Pendiente	4. Redes neuronales.pdf	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/4/1778730651_4._redes_neuronales.pdf	.pdf	\N	2026-05-13 22:50:54.671456	\N	\N	f	2026-05-13 22:50:54.681797	2026-05-13 22:50:54.681797
6	4	Concepto de rehabilitación	Pendiente	Diagrama de calses - incapacidades.drawio.pdf	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/4/1778730663_diagrama_de_calses_-_incapacidades.drawio.pdf	.pdf	\N	2026-05-13 22:51:04.177136	\N	\N	f	2026-05-13 22:51:04.17763	2026-05-13 22:51:04.17763
7	1	certificado_incapacidad	Pendiente	te vas de mi.jpeg	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/1/1778795786_te_vas_de_mi.jpeg	.jpeg	\N	2026-05-14 16:56:27.111488	\N	\N	f	2026-05-14 16:56:27.112781	2026-05-14 16:56:27.112781
8	1	historia_clinica	Pendiente	Black Sabbath.jpeg	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/1/1778795796_black_sabbath.jpeg	.jpeg	\N	2026-05-14 16:56:36.720507	\N	\N	f	2026-05-14 16:56:36.720899	2026-05-14 16:56:36.720899
9	10	certificado_incapacidad	Pendiente	Black Sabbath.jpeg	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/10/1778796997_black_sabbath.jpeg	.jpeg	\N	2026-05-14 17:16:37.751424	\N	\N	f	2026-05-14 17:16:37.751755	2026-05-14 17:16:37.751755
10	10	historia_clinica	Pendiente	radiohead.jpeg	https://pub-1f217c1a1b0549878cd7a1492f6f2771.r2.dev/incapacidad/10/1778797008_radiohead.jpeg	.jpeg	\N	2026-05-14 17:16:49.025116	\N	\N	f	2026-05-14 17:16:49.025378	2026-05-14 17:16:49.025378
\.


--
-- Data for Name: empleado; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.empleado (id_usuario, puesto_trabajo, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: entidad; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.entidad (id_entidad, nombre, tipo, plazo_transcripcion_dias, tiempo_maximo_pago_dias, canal_atencion, canales_atencion, requiere_transcripcion, created_at, updated_at) FROM stdin;
1	Salud Total	EPS	5	30	Portal web,Correo electrónico,Atención presencial	\N	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	Nueva EPS	EPS	5	30	Portal web,Correo electrónico,Atención presencial,Línea telefónica	\N	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	SOS	EPS	5	30	Portal web,Correo electrónico,Atención presencial	\N	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Sanitas	EPS	5	30	Portal web,Correo electrónico,Atención presencial,Mesa de ayuda	\N	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	SURA EPS	EPS	5	30	Portal web,Correo electrónico,Atención presencial	\N	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
6	Asmet Salud	EPS	5	30	Correo electrónico,Atención presencial	\N	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
7	ARL SURA	ARL	3	45	Portal web,Correo electrónico,Atención presencial	\N	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: estado_documento; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.estado_documento (id_estado_documento, nombre, descripcion, color, created_at, updated_at) FROM stdin;
1	Pendiente	Documento pendiente de validación	yellow	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	Validado	Documento validado correctamente	green	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	Rechazado	Documento rechazado	red	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Incompleto	Documento incompleto	orange	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	Vencido	Documento vencido	gray	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
6	Archivado	Documento archivado	gray	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: estado_incapacidad; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.estado_incapacidad (id_estado, nombre, descripcion, permite_transicion, created_at, updated_at) FROM stdin;
1	Recibida	Incapacidad recibida en el sistema	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	En validación documental	Documentos siendo validados	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	Documentación incompleta	Faltan documentos requeridos	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Pendiente transcripción	Por transcribir a la EPS	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	Transcrita	Ya transcrita a la EPS	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
6	En verificación EPS	EPS verificando la incapacidad	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
7	Aprobada	Incapacidad aprobada por EPS	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
8	Cobrada	Valor cobrado a la EPS	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
9	Rechazada	Incapacidad rechazada	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
10	Pendiente pago	Esperando pago de EPS	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
11	Pagada	Pago realizado	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
12	En conciliación	En proceso de conciliación	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
13	Conciliada	Conciliación completada	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
14	Cobro persuasivo	En etapa de cobro persuasivo	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
15	Cobro jurídico	En etapa de cobro jurídico	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
16	Archivada	Incapacidad archivada	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
17	Cerrada	Incapacidad cerrada	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: estado_pago; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.estado_pago (id_estado_pago, nombre, descripcion, color, created_at, updated_at) FROM stdin;
1	Pendiente	Pago pendiente	yellow	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	En proceso	Pago en proceso	orange	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	Pagado	Pago realizado	green	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Conciliado	Pago conciliado	blue	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	Rechazado	Pago rechazado	red	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
6	Parcial	Pago parcial	orange	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
7	Anulado	Pago anulado	gray	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: gerencia; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.gerencia (id_usuario, puesto_trabajo, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: gestion_humana; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.gestion_humana (id_usuario, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: historial; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.historial (id_historial, id_incapacidad, tipo_historial, descripcion, fecha, created_at, gestor_id, id_tipo_historial, updated_at) FROM stdin;
\.


--
-- Data for Name: incapacidad; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.incapacidad (id_incapacidad, id_usuario, id_estado, id_tipo, id_entidad, canal_recepcion, titulo, fecha_inicio, fecha_fin, origen, fecha_radicacion, fecha_pago, observaciones, created_by, created_at, updated_at, is_deleted, fecha_transcripcion, transcrito_por, observaciones_transcripcion, fecha_limite_transcripcion, estado_transcripcion) FROM stdin;
1	3	11	1	1	presencial	Incapacidad actualizada	2024-01-15	2024-01-30	laboral	2024-01-20	2024-02-15	Cambio de estado	3	2026-05-12 19:24:37.774972	2026-05-13 18:36:25.548182	f	2026-05-12 21:38:34.653702	3	\N	2026-05-17	en_proceso
2	2	1	1	1	presencial	Incapacidad por enfermedad general	2024-01-15	2024-01-30	Común	2024-01-20	2024-02-15	Notas adicionales	2	2026-05-13 19:42:45.132294	2026-05-13 19:42:45.132294	f	\N	\N	\N	2026-05-18	pendiente
3	3	1	3	4	presencial	Incapacidad por accidente de transito	2026-05-13	2026-05-16	Tránsito	2026-05-13	\N	Radicada la incapacidad queda pendiente subida de documentos	3	2026-05-13 19:46:35.310192	2026-05-13 19:46:35.310192	f	\N	\N	\N	2026-05-18	pendiente
4	3	1	6	2	virtual	Incapacidad de un usuario	2026-05-13	2026-05-16	Laboral	2026-05-13	\N	El empleado presenta sintomas de dolor constante	3	2026-05-13 22:49:33.037031	2026-05-13 22:49:33.037031	f	\N	\N	\N	2026-05-18	pendiente
5	3	1	5	5	email	Incapacidad de un usuario	2026-05-13	2026-05-16	Paternidad	2026-05-13	\N	\N	3	2026-05-13 22:58:07.186869	2026-05-13 22:58:07.186869	f	\N	\N	\N	2026-05-18	pendiente
6	3	1	1	7	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Común	2026-05-14	\N	\N	3	2026-05-14 17:09:53.172399	2026-05-14 17:09:53.172399	f	\N	\N	\N	2026-05-17	pendiente
7	3	1	6	2	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Laboral	2026-05-14	\N	\N	3	2026-05-14 17:10:31.21419	2026-05-14 17:10:31.21419	f	\N	\N	\N	2026-05-19	pendiente
8	3	1	1	6	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Común	2026-05-14	\N	\N	3	2026-05-14 17:13:37.34842	2026-05-14 17:13:37.34842	f	\N	\N	\N	2026-05-19	pendiente
9	3	1	6	1	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Laboral	2026-05-14	\N	\N	3	2026-05-14 17:15:17.670636	2026-05-14 17:15:17.670636	f	\N	\N	\N	2026-05-19	pendiente
10	3	1	3	3	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Tránsito	2026-05-14	\N	\N	3	2026-05-14 17:16:25.866232	2026-05-14 17:16:25.866232	f	\N	\N	\N	2026-05-19	pendiente
11	3	1	2	2	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Laboral	2026-05-14	\N	\N	3	2026-05-14 20:37:47.5469	2026-05-14 20:37:47.5469	f	\N	\N	\N	2026-05-19	pendiente
13	4	3	1	2	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Común	2026-05-14	\N	\N	3	2026-05-14 20:40:41.410512	2026-05-15 02:30:54.260307	f	\N	\N	\N	2026-05-19	pendiente
12	4	1	4	2	presencial	Incapacidad de un usuario	2026-05-14	2026-05-17	Maternidad	2026-05-14	\N	\N	3	2026-05-14 20:38:42.793933	2026-05-14 20:38:42.793933	f	\N	\N	\N	2026-05-19	pendiente
\.


--
-- Data for Name: notificacion; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.notificacion (id_notificacion, id_usuario, id_incapacidad, tipo_notificacion, mensaje, fecha, leida, created_at, updated_at) FROM stdin;
2	1	1	Recordatorio	El pago está próximo a vencer	2026-05-12 20:51:07.879457	f	2026-05-12 20:51:07.879457	2026-05-12 20:51:07.879457
1	3	1	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-12 19:24:37.789956	t	2026-05-12 19:24:37.789956	2026-05-12 20:52:13.471446
3	3	1	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-12 21:49:37.292071	f	2026-05-12 21:49:37.292071	2026-05-12 21:49:37.292071
4	2	2	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-13 19:42:45.154713	f	2026-05-13 19:42:45.154713	2026-05-13 19:42:45.154713
5	3	3	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-13 19:46:35.313749	f	2026-05-13 19:46:35.313749	2026-05-13 19:46:35.313749
6	3	4	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-13 22:49:33.074206	f	2026-05-13 22:49:33.074206	2026-05-13 22:49:33.074206
7	3	4	Documento faltante	Documento faltante: Historia clínica	2026-05-13 22:50:38.09342	f	2026-05-13 22:50:38.09342	2026-05-13 22:50:38.09342
8	3	5	Documento faltante	Documento faltante: Certificado de incapacidad	2026-05-13 22:58:07.205329	f	2026-05-13 22:58:07.205329	2026-05-13 22:58:07.205329
9	3	1	Documento faltante	Documento faltante: Historia clínica	2026-05-14 16:56:27.129481	f	2026-05-14 16:56:27.129481	2026-05-14 16:56:27.129481
10	3	6	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-14 17:09:53.182732	f	2026-05-14 17:09:53.182732	2026-05-14 17:09:53.182732
11	3	7	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-14 17:10:31.230648	f	2026-05-14 17:10:31.230648	2026-05-14 17:10:31.230648
12	3	8	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-14 17:13:37.365922	f	2026-05-14 17:13:37.365922	2026-05-14 17:13:37.365922
13	3	9	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-14 17:15:17.690318	f	2026-05-14 17:15:17.690318	2026-05-14 17:15:17.690318
14	3	10	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-14 17:16:25.888006	f	2026-05-14 17:16:25.888006	2026-05-14 17:16:25.888006
15	3	10	Documento faltante	Documento faltante: Historia clínica	2026-05-14 17:16:37.761515	f	2026-05-14 17:16:37.761515	2026-05-14 17:16:37.761515
16	3	11	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-14 20:37:47.560081	f	2026-05-14 20:37:47.560081	2026-05-14 20:37:47.560081
17	3	12	Documento faltante	Documento faltante: Certificado de incapacidad	2026-05-14 20:38:42.809617	f	2026-05-14 20:38:42.809617	2026-05-14 20:38:42.809617
18	4	13	Documento faltante	Documentos faltantes: Certificado de incapacidad, Historia clínica	2026-05-14 20:40:41.424705	f	2026-05-14 20:40:41.424705	2026-05-14 20:40:41.424705
\.


--
-- Data for Name: pago; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.pago (id_pago, id_incapacidad, id_entidad, tipo_pago, estado_pago, descripcion, valor, fecha_pago, periodo_contable, conciliado, registrado_por, created_at, updated_at, is_deleted) FROM stdin;
1	1	1	Consignación	Conciliado	Pago conciliado	1500000.00	2024-01-20	2024-01	t	3	2026-05-12 20:26:58.718527	2026-05-12 20:31:11.260231	f
2	1	1	Consignación	Pendiente	Pago de incapacidad	1500000.00	2024-01-20	2024-01	f	3	2026-05-12 21:49:59.190282	2026-05-12 21:49:59.190282	f
\.


--
-- Data for Name: periodicidad_reporte; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.periodicidad_reporte (id_periodicidad, nombre, dias, descripcion, created_at, updated_at) FROM stdin;
1	Diario	1	Reporte diario	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	Semanal	7	Reporte semanal	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	Mensual	30	Reporte mensual	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Trimestral	90	Reporte trimestral	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	Anual	365	Reporte anual	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: permisos; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.permisos (id_permiso, nombre, descripcion, created_at, updated_at) FROM stdin;
1	crear_incapacidad	Permiso para crear incapacidades	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	editar_incapacidad	Permiso para editar incapacidades	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	consultar_incapacidad	Permiso para consultar incapacidades	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	validar_documentos	Permiso para validar documentos	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	rechazar_documentos	Permiso para rechazar documentos	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
6	registrar_pago	Permiso para registrar pagos	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
7	consultar_reportes	Permiso para consultar reportes	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
8	gestionar_usuarios	Permiso para gestionar usuarios	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
9	gestionar_roles	Permiso para gestionar roles	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
10	generar_alertas	Permiso para generar alertas	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
11	consultar_historial	Permiso para consultar historial	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
12	realizar_conciliacion	Permiso para realizar conciliaciones	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
13	gestionar_cobro_persuasivo	Permiso para gestionar cobro persuasivo	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
14	gestionar_cobro_juridico	Permiso para gestionar cobro jurídico	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
15	archivar_incapacidad	Permiso para archivar incapacidades	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: rol; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.rol (id_rol, nombre, permisos, created_at, updated_at, is_deleted) FROM stdin;
1	Administrador	["crear_incapacidad", "editar_incapacidad", "consultar_incapacidad", "validar_documentos", "rechazar_documentos", "registrar_pago", "consultar_reportes", "gestionar_usuarios", "gestionar_roles", "generar_alertas", "consultar_historial", "realizar_conciliacion", "gestionar_cobro_persuasivo", "gestionar_cobro_juridico", "archivar_incapacidad"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
2	Gestión Humana	["crear_incapacidad", "editar_incapacidad", "consultar_incapacidad", "validar_documentos", "rechazar_documentos", "consultar_reportes", "consultar_historial", "generar_alertas"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
3	Gerencia	["consultar_incapacidad", "consultar_reportes", "consultar_historial"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
4	Empleado	["consultar_incapacidad"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
5	Contabilidad	["consultar_incapacidad", "registrar_pago", "consultar_reportes", "realizar_conciliacion"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
6	Jurídica	["consultar_incapacidad", "consultar_reportes", "consultar_historial", "gestionar_cobro_juridico"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
7	Recepcionista	["crear_incapacidad", "consultar_incapacidad"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
8	Cartera	["consultar_incapacidad", "registrar_pago", "consultar_reportes", "realizar_conciliacion", "gestionar_cobro_persuasivo"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
9	SG-SST	["crear_incapacidad", "editar_incapacidad", "consultar_incapacidad", "validar_documentos", "consultar_reportes", "consultar_historial", "generar_alertas"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	f
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.schema_migrations (version, dirty) FROM stdin;
9	f
\.


--
-- Data for Name: seguimiento_cobro; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.seguimiento_cobro (id_seguimiento, id_incapacidad, tipo_seguimiento, resultado_seguimiento, descripcion, fecha, gestionado_por, created_at, updated_at) FROM stdin;
1	1	Persuasivo	Pago realizado	Seguimiento actualizado	2026-05-12 20:34:45.596657	3	2026-05-12 20:34:45.596657	2026-05-12 20:35:53.754404
\.


--
-- Data for Name: tipo_documento; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.tipo_documento (id_tipo_documento, nombre, descripcion, requerido, created_at, updated_at, codigo) FROM stdin;
1	Certificado de incapacidad	Certificado médico de incapacidad	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	certificado_incapacidad
2	Epicrisis	Epicrisis médica	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	epicrisis
3	FURIPS	Formato único de reportabilidad de instituciones de salud	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	furips
4	Historia clínica	Historia clínica completa	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	historia_clinica
5	Certificado de nacido vivo	Certificado de nacido vivo	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	certificado_nacido_vivo
6	Registro civil	Registro civil del nacido	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	registro_civil
7	Documento de identidad	Cédula de ciudadanía	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	documento_identidad
8	Soporte de atención médica	Soporte de la atención médica	t	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	soporte_atencion_medica
9	Concepto de rehabilitación	Concepto de rehabilitación	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	concepto_rehabilitacion
10	Evidencia de radicación	Evidencia de radicación	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	evidencia_radicacion
11	Soporte de pago	Soporte del pago	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	soporte_pago
12	Formato de seguimiento	Formato de seguimiento	f	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	formato_seguimiento
\.


--
-- Data for Name: tipo_entidad; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.tipo_entidad (id_tipo_entidad, nombre, descripcion, created_at, updated_at) FROM stdin;
1	EPS	Entidad Promotora de Salud	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	ARL	Administradora de Riesgos Laborales	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	AFP	Administradora de Fondo de Pensiones	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: tipo_incapacidad; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.tipo_incapacidad (id_tipo, nombre, documentos_requeridos, created_at, updated_at, origen) FROM stdin;
1	Enfermedad general	["certificado_incapacidad", "historia_clinica"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	Común
2	Accidente laboral	["certificado_incapacidad", "furips", "historia_clinica"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	Laboral
3	Accidente de tránsito	["certificado_incapacidad", "furips", "historia_clinica"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	Tránsito
4	Licencia de maternidad	["certificado_incapacidad", "certificado_nacido_vivo", "registro_civil"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	Maternidad
5	Licencia de paternidad	["certificado_incapacidad", "registro_civil"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	Paternidad
6	Enfermedad laboral	["certificado_incapacidad", "historia_clinica", "concepto_rehabilitacion"]	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628	Laboral
\.


--
-- Data for Name: tipo_pago; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.tipo_pago (id_tipo_pago, nombre, descripcion, created_at, updated_at) FROM stdin;
1	Transferencia bancaria	Pago por transferencia bancaria	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
2	Consignación	Consignación bancaria	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
3	Pago parcial	Pago parcial	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
4	Pago total	Pago total	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
5	Reintegro	Reintegro de pago	2026-05-12 17:14:38.813628	2026-05-12 17:14:38.813628
\.


--
-- Data for Name: usuario; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.usuario (id_usuario, id_rol, nombre, correo, numero_celular, direccion, password_hash, numero_documento, numero_acudiente, estado, created_at, updated_at, is_deleted) FROM stdin;
2	2	Gestion Humana	GH@ejemplo.com	\N	\N	$2a$10$A0xFLxlWoSKZbOzD97CP9O.LRoJXW.w2vwvn51o2lRoNUg1jVSjVe	12345678	\N	t	2026-05-12 19:02:56.307435	2026-05-12 19:02:56.307931	f
3	1	Administrador R	r@ejemplo.com	\N	\N	$2a$10$ik/9OVydlL1B8FEDNlIjJOwTSXFR7/SZsQobOzSYP.ZTGbIYiGoaK	12345675	\N	t	2026-05-12 19:09:50.184872	2026-05-12 19:09:50.185243	f
1	1	Administrador	admin@example.com	+573009999999	\N	$2a$10$5aOqVlctjvrqRRauxXkKVuo8OHcLzeigAz6lE6UQhiOkg2l6FqJpe	string	\N	t	2026-05-12 12:21:20.568864	2026-05-12 21:31:42.558623	f
4	4	Empleado normal	empleado@ejemplo.com	\N	\N	$2a$10$jGdzYhRxWF74mTpwhNMif.3AE/PXaQkr8.qZBYuNxYftJ1p8a.Yki	12345671	\N	t	2026-05-13 22:09:13.98015	2026-05-13 22:09:13.980899	f
5	4	Kenneth	k@empleado.com	6549846	5154	$2a$10$TzlJ6YLSAZpkXgs5wX8EEeQOv3hCvkbU8kdcOx7SUNGCwJSlUexg6	123465741	\N	t	2026-05-15 02:48:26.128035	2026-05-15 02:48:26.128035	f
6	4	Kenneth Rodriguez	k3@empleado.com	6549846	5154	$2a$10$KGXf.821OtxKNaWxsbSEwutjx7p7zlg4znfCRSYvjpVWBTBb0LE5C	12346245345	\N	t	2026-05-15 02:52:42.764336	2026-05-15 02:52:42.764336	f
7	4	Kenneth Rodriguez	k2@empleado.com	6549846	5154	$2a$10$3fmiAxbdk1YJxHnhI3uVruIIMrirQdrTbej9m3/UVRAnK0GETGInq	1235345	\N	t	2026-05-15 02:53:29.57243	2026-05-15 02:53:29.57243	f
8	4	Kenneth Rodriguez	k4@empleado.com	6549846	5154	$2a$10$qQ4hAW0OCRB8hnHgCHZf9eDjXRZ7bbyO.6WFhTAXcni34R/Gb6ON.	1235556	\N	t	2026-05-15 02:55:05.893176	2026-05-15 02:55:05.893176	f
\.


--
-- Name: canal_atencion_entidad_id_canal_atencion_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.canal_atencion_entidad_id_canal_atencion_seq', 5, true);


--
-- Name: canal_recepcion_id_canal_recepcion_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.canal_recepcion_id_canal_recepcion_seq', 6, true);


--
-- Name: documento_id_documento_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.documento_id_documento_seq', 10, true);


--
-- Name: entidad_id_entidad_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.entidad_id_entidad_seq', 7, true);


--
-- Name: estado_documento_id_estado_documento_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.estado_documento_id_estado_documento_seq', 6, true);


--
-- Name: estado_incapacidad_id_estado_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.estado_incapacidad_id_estado_seq', 17, true);


--
-- Name: estado_pago_id_estado_pago_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.estado_pago_id_estado_pago_seq', 7, true);


--
-- Name: historial_id_historial_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.historial_id_historial_seq', 1, false);


--
-- Name: incapacidad_id_incapacidad_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.incapacidad_id_incapacidad_seq', 13, true);


--
-- Name: notificacion_id_notificacion_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.notificacion_id_notificacion_seq', 18, true);


--
-- Name: pago_id_pago_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.pago_id_pago_seq', 2, true);


--
-- Name: periodicidad_reporte_id_periodicidad_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.periodicidad_reporte_id_periodicidad_seq', 5, true);


--
-- Name: permisos_id_permiso_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.permisos_id_permiso_seq', 15, true);


--
-- Name: rol_id_rol_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.rol_id_rol_seq', 9, true);


--
-- Name: seguimiento_cobro_id_seguimiento_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.seguimiento_cobro_id_seguimiento_seq', 1, true);


--
-- Name: tipo_documento_id_tipo_documento_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.tipo_documento_id_tipo_documento_seq', 12, true);


--
-- Name: tipo_entidad_id_tipo_entidad_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.tipo_entidad_id_tipo_entidad_seq', 3, true);


--
-- Name: tipo_incapacidad_id_tipo_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.tipo_incapacidad_id_tipo_seq', 6, true);


--
-- Name: tipo_pago_id_tipo_pago_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.tipo_pago_id_tipo_pago_seq', 5, true);


--
-- Name: usuario_id_usuario_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.usuario_id_usuario_seq', 8, true);


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

\unrestrict 7tW6PzhroU44F9PIeVQbEmRw7j9pgs3NOcTIi35V8cTC1L4oePuW6dW0XtlEcNw

