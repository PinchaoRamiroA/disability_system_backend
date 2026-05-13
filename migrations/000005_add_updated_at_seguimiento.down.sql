-- Remove updated_at column from seguimiento_cobro
ALTER TABLE seguimiento_cobro DROP COLUMN IF EXISTS updated_at;