-- 000007_add_origen_to_tipo_incapacidad.down.sql

ALTER TABLE tipo_incapacidad DROP COLUMN IF EXISTS origen;
