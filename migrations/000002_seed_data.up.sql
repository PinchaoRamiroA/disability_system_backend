-- Seed data para datos de referencia

-- =====================================================
-- PERMISOS (TABLA)
-- =====================================================

INSERT INTO permisos (nombre, descripcion) VALUES
('crear_incapacidad', 'Permiso para crear incapacidades'),
('editar_incapacidad', 'Permiso para editar incapacidades'),
('consultar_incapacidad', 'Permiso para consultar incapacidades'),
('validar_documentos', 'Permiso para validar documentos'),
('rechazar_documentos', 'Permiso para rechazar documentos'),
('registrar_pago', 'Permiso para registrar pagos'),
('consultar_reportes', 'Permiso para consultar reportes'),
('gestionar_usuarios', 'Permiso para gestionar usuarios'),
('gestionar_roles', 'Permiso para gestionar roles'),
('generar_alertas', 'Permiso para generar alertas'),
('consultar_historial', 'Permiso para consultar historial'),
('realizar_conciliacion', 'Permiso para realizar conciliaciones'),
('gestionar_cobro_persuasivo', 'Permiso para gestionar cobro persuasivo'),
('gestionar_cobro_juridico', 'Permiso para gestionar cobro jurídico'),
('archivar_incapacidad', 'Permiso para archivar incapacidades');

-- =====================================================
-- ROLES (TABLA)
-- =====================================================

INSERT INTO rol (nombre, permisos) VALUES
('Administrador', '["crear_incapacidad","editar_incapacidad","consultar_incapacidad","validar_documentos","rechazar_documentos","registrar_pago","consultar_reportes","gestionar_usuarios","gestionar_roles","generar_alertas","consultar_historial","realizar_conciliacion","gestionar_cobro_persuasivo","gestionar_cobro_juridico","archivar_incapacidad"]'),
('Gestión Humana', '["crear_incapacidad","editar_incapacidad","consultar_incapacidad","validar_documentos","rechazar_documentos","consultar_reportes","consultar_historial","generar_alertas"]'),
('Gerencia', '["consultar_incapacidad","consultar_reportes","consultar_historial"]'),
('Empleado', '["consultar_incapacidad"]'),
('Contabilidad', '["consultar_incapacidad","registrar_pago","consultar_reportes","realizar_conciliacion"]'),
('Jurídica', '["consultar_incapacidad","consultar_reportes","consultar_historial","gestionar_cobro_juridico"]'),
('Recepcionista', '["crear_incapacidad","consultar_incapacidad"]'),
('Cartera', '["consultar_incapacidad","registrar_pago","consultar_reportes","realizar_conciliacion","gestionar_cobro_persuasivo"]'),
('SG-SST', '["crear_incapacidad","editar_incapacidad","consultar_incapacidad","validar_documentos","consultar_reportes","consultar_historial","generar_alertas"]');

-- =====================================================
-- TIPOS DE DOCUMENTO (TABLA)
-- =====================================================

INSERT INTO tipo_documento (nombre, descripcion, requerido) VALUES
('Certificado de incapacidad', 'Certificado médico de incapacidad', true),
('Epicrisis', 'Epicrisis médica', false),
('FURIPS', 'Formato único de reportabilidad de instituciones de salud', false),
('Historia clínica', 'Historia clínica completa', true),
('Certificado de nacido vivo', 'Certificado de nacido vivo', false),
('Registro civil', 'Registro civil del nacido', false),
('Documento de identidad', 'Cédula de ciudadanía', true),
('Soporte de atención médica', 'Soporte de la atención médica', true),
('Concepto de rehabilitación', 'Concepto de rehabilitación', false),
('Evidencia de radicación', 'Evidencia de radicación', false),
('Soporte de pago', 'Soporte del pago', false),
('Formato de seguimiento', 'Formato de seguimiento', false);

-- =====================================================
-- ESTADOS DE DOCUMENTO (TABLA)
-- =====================================================

INSERT INTO estado_documento (nombre, descripcion, color) VALUES
('Pendiente', 'Documento pendiente de validación', 'yellow'),
('Validado', 'Documento validado correctamente', 'green'),
('Rechazado', 'Documento rechazado', 'red'),
('Incompleto', 'Documento incompleto', 'orange'),
('Vencido', 'Documento vencido', 'gray'),
('Archivado', 'Documento archivado', 'gray');

-- =====================================================
-- TIPOS DE ENTIDAD (TABLA)
-- =====================================================

INSERT INTO tipo_entidad (nombre, descripcion) VALUES
('EPS', 'Entidad Promotora de Salud'),
('ARL', 'Administradora de Riesgos Laborales'),
('AFP', 'Administradora de Fondo de Pensiones');

-- =====================================================
-- CANALES DE RECEPCIÓN (TABLA)
-- =====================================================

INSERT INTO canal_recepcion (nombre, descripcion, activo) VALUES
('Recepción física', 'Recibido en recepción física', true),
('Oficina principal', 'Recibido en oficina principal', true),
('Líder comercial', 'Recibido por líder comercial', true),
('Correo electrónico', 'Recibido por correo electrónico', true),
('Portal EPS', 'Recibido desde portal EPS', true),
('Portal ARL', 'Recibido desde portal ARL', true);

-- =====================================================
-- CANALES DE ATENCIÓN ENTIDAD (TABLA)
-- =====================================================

INSERT INTO canal_atencion_entidad (nombre, descripcion) VALUES
('Portal web', 'Atención por portal web'),
('Correo electrónico', 'Atención por correo electrónico'),
('Atención presencial', 'Atención presencial'),
('Línea telefónica', 'Atención por línea telefónica'),
('Mesa de ayuda', 'Atención por mesa de ayuda');

-- =====================================================
-- TIPOS DE PAGO (TABLA)
-- =====================================================

INSERT INTO tipo_pago (nombre, descripcion) VALUES
('Transferencia bancaria', 'Pago por transferencia bancaria'),
('Consignación', 'Consignación bancaria'),
('Pago parcial', 'Pago parcial'),
('Pago total', 'Pago total'),
('Reintegro', 'Reintegro de pago');

-- =====================================================
-- ESTADOS DE PAGO (TABLA)
-- =====================================================

INSERT INTO estado_pago (nombre, descripcion, color) VALUES
('Pendiente', 'Pago pendiente', 'yellow'),
('En proceso', 'Pago en proceso', 'orange'),
('Pagado', 'Pago realizado', 'green'),
('Conciliado', 'Pago conciliado', 'blue'),
('Rechazado', 'Pago rechazado', 'red'),
('Parcial', 'Pago parcial', 'orange'),
('Anulado', 'Pago anulado', 'gray');

-- =====================================================
-- PERIODICIDAD DE REPORTES (TABLA)
-- =====================================================

INSERT INTO periodicidad_reporte (nombre, dias, descripcion) VALUES
('Diario', 1, 'Reporte diario'),
('Semanal', 7, 'Reporte semanal'),
('Mensual', 30, 'Reporte mensual'),
('Trimestral', 90, 'Reporte trimestral'),
('Anual', 365, 'Reporte anual');

-- =====================================================
-- ESTADOS DE INCAPACIDAD (TABLA)
-- =====================================================

INSERT INTO estado_incapacidad (nombre, descripcion, permite_transicion) VALUES
('Recibida', 'Incapacidad recibida en el sistema', true),
('En validación documental', 'Documentos siendo validados', true),
('Documentación incompleta', 'Faltan documentos requeridos', true),
('Pendiente transcripción', 'Por transcribir a la EPS', true),
('Transcrita', 'Ya transcrita a la EPS', true),
('En verificación EPS', 'EPS verificando la incapacidad', true),
('Aprobada', 'Incapacidad aprobada por EPS', true),
('Cobrada', 'Valor cobrado a la EPS', true),
('Rechazada', 'Incapacidad rechazada', false),
('Pendiente pago', 'Esperando pago de EPS', true),
('Pagada', 'Pago realizado', true),
('En conciliación', 'En proceso de conciliación', true),
('Conciliada', 'Conciliación completada', true),
('Cobro persuasivo', 'En etapa de cobro persuasivo', true),
('Cobro jurídico', 'En etapa de cobro jurídico', true),
('Archivada', 'Incapacidad archivada', false),
('Cerrada', 'Incapacidad cerrada', false);

-- =====================================================
-- TIPOS DE INCAPACIDAD (TABLA)
-- =====================================================

INSERT INTO tipo_incapacidad (nombre, documentos_requeridos) VALUES
('Enfermedad general', '["certificado_incapacidad","historia_clinica"]'),
('Accidente laboral', '["certificado_incapacidad","furips","historia_clinica"]'),
('Accidente de tránsito', '["certificado_incapacidad","furips","historia_clinica"]'),
('Licencia de maternidad', '["certificado_incapacidad","certificado_nacido_vivo","registro_civil"]'),
('Licencia de paternidad', '["certificado_incapacidad","registro_civil"]'),
('Enfermedad laboral', '["certificado_incapacidad","historia_clinica","concepto_rehabilitacion"]');

-- =====================================================
-- ENTIDADES (TABLA)
-- =====================================================

INSERT INTO entidad (nombre, tipo, plazo_transcripcion_dias, tiempo_maximo_pago_dias, canal_atencion, requiere_transcripcion) VALUES
('Salud Total', 'EPS', 5, 30, 'Portal web,Correo electrónico,Atención presencial', true),
('Nueva EPS', 'EPS', 5, 30, 'Portal web,Correo electrónico,Atención presencial,Línea telefónica', true),
('SOS', 'EPS', 5, 30, 'Portal web,Correo electrónico,Atención presencial', true),
('Sanitas', 'EPS', 5, 30, 'Portal web,Correo electrónico,Atención presencial,Mesa de ayuda', true),
('SURA EPS', 'EPS', 5, 30, 'Portal web,Correo electrónico,Atención presencial', true),
('Asmet Salud', 'EPS', 5, 30, 'Correo electrónico,Atención presencial', true),
('ARL SURA', 'ARL', 3, 45, 'Portal web,Correo electrónico,Atención presencial', true);