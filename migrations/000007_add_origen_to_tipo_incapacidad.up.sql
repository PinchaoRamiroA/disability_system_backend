-- 000007_add_origen_to_tipo_incapacidad.up.sql

ALTER TABLE tipo_incapacidad ADD COLUMN origen VARCHAR(50);

-- Actualizar datos existentes basándose en el nombre
UPDATE tipo_incapacidad SET origen = 'Común' WHERE nombre = 'Enfermedad general';
UPDATE tipo_incapacidad SET origen = 'Laboral' WHERE nombre = 'Accidente laboral';
UPDATE tipo_incapacidad SET origen = 'Tránsito' WHERE nombre = 'Accidente de tránsito';
UPDATE tipo_incapacidad SET origen = 'Maternidad' WHERE nombre = 'Licencia de maternidad';
UPDATE tipo_incapacidad SET origen = 'Paternidad' WHERE nombre = 'Licencia de paternidad';
UPDATE tipo_incapacidad SET origen = 'Laboral' WHERE nombre = 'Enfermedad laboral';

-- Establecer como NOT NULL después de poblar los datos
ALTER TABLE tipo_incapacidad ALTER COLUMN origen SET NOT NULL;
