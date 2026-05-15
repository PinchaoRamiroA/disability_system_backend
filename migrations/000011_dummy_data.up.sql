-- Insertar usuarios administrativos
INSERT INTO usuario (id_rol, nombre, correo, numero_documento, password_hash, estado) VALUES
(1, 'Admin Principal', 'admin@empresa.com', '100000001', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(2, 'Gestor Humano', 'gh@empresa.com', '100000002', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(3, 'Gerente General', 'gerencia@empresa.com', '100000003', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(5, 'Contador', 'contabilidad@empresa.com', '100000005', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(6, 'Asesor Jurídico', 'juridica@empresa.com', '100000006', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(7, 'Recepcionista', 'recepcion@empresa.com', '100000007', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(8, 'Analista de Cartera', 'cartera@empresa.com', '100000008', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(9, 'Coordinador SG-SST', 'sgsst@empresa.com', '100000009', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true);

-- Insertar empleados (asumiendo que id_rol = 4 es Empleado según el seed anterior)
INSERT INTO usuario (id_rol, nombre, correo, numero_documento, password_hash, estado) VALUES
(4, 'Juan Perez', 'juan.perez@empleado.com', '200000001', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Maria Gomez', 'maria.gomez@empleado.com', '200000002', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Carlos Rodriguez', 'carlos.rodriguez@empleado.com', '200000003', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Ana Martinez', 'ana.martinez@empleado.com', '200000004', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Luis Fernandez', 'luis.fernandez@empleado.com', '200000005', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Laura Sanchez', 'laura.sanchez@empleado.com', '200000006', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'David Torres', 'david.torres@empleado.com', '200000007', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Elena Ruiz', 'elena.ruiz@empleado.com', '200000008', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Jorge Diaz', 'jorge.diaz@empleado.com', '200000009', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true),
(4, 'Sofia Morales', 'sofia.morales@empleado.com', '200000010', '$2a$10$5JGzZsF/VZbcpMPaoFcWr.Ur.923/5d8dR8O54FJv02zbvh2brTim', true);

-- Las incapacidades se asociarán buscando el ID del usuario insertado.
INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 1, 1, 1, CURRENT_DATE - INTERVAL '10 days', CURRENT_DATE - INTERVAL '5 days', 'Común', 'Reposo absoluto - Gripe Fuerte', 'Incapacidad Gripe', (SELECT id_usuario FROM usuario WHERE correo = 'admin@empresa.com') FROM usuario WHERE correo = 'juan.perez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 1, 7, 1, CURRENT_DATE - INTERVAL '30 days', CURRENT_DATE - INTERVAL '25 days', 'Común', 'Reposo y dieta - Infección Intestinal', 'Incapacidad Infección', (SELECT id_usuario FROM usuario WHERE correo = 'admin@empresa.com') FROM usuario WHERE correo = 'juan.perez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 7, 8, 2, CURRENT_DATE - INTERVAL '45 days', CURRENT_DATE - INTERVAL '30 days', 'Laboral', 'Accidente en oficina - Fractura de Brazo', 'Accidente Brazo', (SELECT id_usuario FROM usuario WHERE correo = 'gh@empresa.com') FROM usuario WHERE correo = 'maria.gomez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 7, 10, 2, CURRENT_DATE - INTERVAL '15 days', CURRENT_DATE - INTERVAL '5 days', 'Laboral', 'Continuación - Prórroga Fractura', 'Prorroga Fractura', (SELECT id_usuario FROM usuario WHERE correo = 'gh@empresa.com') FROM usuario WHERE correo = 'maria.gomez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 2, 11, 1, CURRENT_DATE - INTERVAL '60 days', CURRENT_DATE - INTERVAL '55 days', 'Común', 'Reposo - Migraña', 'Incapacidad Migraña', (SELECT id_usuario FROM usuario WHERE correo = 'admin@empresa.com') FROM usuario WHERE correo = 'carlos.rodriguez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 2, 12, 1, CURRENT_DATE - INTERVAL '90 days', CURRENT_DATE - INTERVAL '80 days', 'Común', 'Reposo y fisioterapia - Lumbago', 'Incapacidad Lumbago', (SELECT id_usuario FROM usuario WHERE correo = 'admin@empresa.com') FROM usuario WHERE correo = 'carlos.rodriguez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 4, 14, 1, CURRENT_DATE - INTERVAL '120 days', CURRENT_DATE - INTERVAL '100 days', 'Común', 'Hospitalización - Neumonía', 'Neumonía Severa', (SELECT id_usuario FROM usuario WHERE correo = 'gh@empresa.com') FROM usuario WHERE correo = 'ana.martinez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 4, 15, 4, CURRENT_DATE - INTERVAL '200 days', CURRENT_DATE - INTERVAL '80 days', 'Común', 'Licencia Maternidad - Parto normal', 'Licencia Ana M.', (SELECT id_usuario FROM usuario WHERE correo = 'gh@empresa.com') FROM usuario WHERE correo = 'ana.martinez@empleado.com';

INSERT INTO incapacidad (id_usuario, id_entidad, id_estado, id_tipo, fecha_inicio, fecha_fin, origen, observaciones, titulo, created_by)
SELECT id_usuario, 5, 2, 1, CURRENT_DATE - INTERVAL '2 days', CURRENT_DATE + INTERVAL '3 days', 'Común', 'Reposo e hidratación - Dengue', 'Incapacidad Dengue', (SELECT id_usuario FROM usuario WHERE correo = 'admin@empresa.com') FROM usuario WHERE correo = 'luis.fernandez@empleado.com';

-- Pagos
INSERT INTO pago (id_incapacidad, id_entidad, tipo_pago, estado_pago, descripcion, valor, fecha_pago, periodo_contable, conciliado, registrado_por)
SELECT id_incapacidad, 2, 'Transferencia bancaria', 'Pagado', 'Pago de EPS por 5 días', 250000.00, CURRENT_DATE - INTERVAL '40 days', '2026-04', false, (SELECT id_usuario FROM usuario WHERE correo = 'contabilidad@empresa.com')
FROM incapacidad WHERE titulo = 'Incapacidad Migraña';

INSERT INTO pago (id_incapacidad, id_entidad, tipo_pago, estado_pago, descripcion, valor, fecha_pago, periodo_contable, conciliado, registrado_por)
SELECT id_incapacidad, 2, 'Consignación', 'Parcial', 'Pago incompleto de EPS por lumbago', 300000.00, CURRENT_DATE - INTERVAL '60 days', '2026-03', false, (SELECT id_usuario FROM usuario WHERE correo = 'cartera@empresa.com')
FROM incapacidad WHERE titulo = 'Incapacidad Lumbago';

-- Seguimientos
INSERT INTO seguimiento_cobro (id_incapacidad, tipo_seguimiento, descripcion, fecha, resultado_seguimiento, gestionado_por)
SELECT id_incapacidad, 'Cobro persuasivo', 'Se envió correo de cobro a EPS Sanitas.', CURRENT_DATE - INTERVAL '30 days', 'Pendiente respuesta', (SELECT id_usuario FROM usuario WHERE correo = 'cartera@empresa.com')
FROM incapacidad WHERE titulo = 'Neumonía Severa';

INSERT INTO seguimiento_cobro (id_incapacidad, tipo_seguimiento, descripcion, fecha, resultado_seguimiento, gestionado_por)
SELECT id_incapacidad, 'Cobro persuasivo', 'Se realizó llamada a EPS Sanitas, indican que está en trámite.', CURRENT_DATE - INTERVAL '15 days', 'En revisión', (SELECT id_usuario FROM usuario WHERE correo = 'cartera@empresa.com')
FROM incapacidad WHERE titulo = 'Neumonía Severa';

INSERT INTO seguimiento_cobro (id_incapacidad, tipo_seguimiento, descripcion, fecha, resultado_seguimiento, gestionado_por)
SELECT id_incapacidad, 'Cobro jurídico', 'Se inicia radicación de demanda ante la Superintendencia.', CURRENT_DATE - INTERVAL '60 days', 'En revisión', (SELECT id_usuario FROM usuario WHERE correo = 'juridica@empresa.com')
FROM incapacidad WHERE titulo = 'Licencia Ana M.';

INSERT INTO seguimiento_cobro (id_incapacidad, tipo_seguimiento, descripcion, fecha, resultado_seguimiento, gestionado_por)
SELECT id_incapacidad, 'Cobro jurídico', 'Respuesta de la Superintendencia a favor de la empresa.', CURRENT_DATE - INTERVAL '10 days', 'Favorable', (SELECT id_usuario FROM usuario WHERE correo = 'juridica@empresa.com')
FROM incapacidad WHERE titulo = 'Licencia Ana M.';

-- Auditoria
INSERT INTO auditoria (id_usuario, id_incapacidad, tipo_accion, modulo, descripcion, cambio_anterior, cambio_nuevo, created_at)
SELECT (SELECT id_usuario FROM usuario WHERE correo = 'admin@empresa.com'), id_incapacidad, 'crear', 'incapacidad', 'El admin creó la incapacidad', null, 'Recibida', CURRENT_DATE - INTERVAL '10 days'
FROM incapacidad WHERE titulo = 'Incapacidad Gripe';

INSERT INTO auditoria (id_usuario, id_incapacidad, tipo_accion, modulo, descripcion, cambio_anterior, cambio_nuevo, created_at)
SELECT (SELECT id_usuario FROM usuario WHERE correo = 'cartera@empresa.com'), id_incapacidad, 'cambiar_estado', 'cobro', 'Se inició cobro persuasivo', 'Pendiente pago', 'Cobro persuasivo', CURRENT_DATE - INTERVAL '30 days'
FROM incapacidad WHERE titulo = 'Neumonía Severa';
