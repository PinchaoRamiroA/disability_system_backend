-- Eliminar datos de prueba en el orden inverso a las dependencias
DELETE FROM auditoria WHERE id_usuario > 0;
DELETE FROM seguimiento_cobro WHERE id_incapacidad IN (SELECT id_incapacidad FROM incapacidad WHERE titulo IN ('Incapacidad Gripe', 'Incapacidad Infección', 'Accidente Brazo', 'Prorroga Fractura', 'Incapacidad Migraña', 'Incapacidad Lumbago', 'Neumonía Severa', 'Licencia Ana M.', 'Incapacidad Dengue'));
DELETE FROM pago WHERE id_incapacidad IN (SELECT id_incapacidad FROM incapacidad WHERE titulo IN ('Incapacidad Gripe', 'Incapacidad Infección', 'Accidente Brazo', 'Prorroga Fractura', 'Incapacidad Migraña', 'Incapacidad Lumbago', 'Neumonía Severa', 'Licencia Ana M.', 'Incapacidad Dengue'));
DELETE FROM incapacidad WHERE titulo IN ('Incapacidad Gripe', 'Incapacidad Infección', 'Accidente Brazo', 'Prorroga Fractura', 'Incapacidad Migraña', 'Incapacidad Lumbago', 'Neumonía Severa', 'Licencia Ana M.', 'Incapacidad Dengue');
DELETE FROM usuario WHERE correo LIKE '%@empresa.com' OR correo LIKE '%@empleado.com';
