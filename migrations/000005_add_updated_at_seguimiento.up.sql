-- Add updated_at column to seguimiento_cobro
ALTER TABLE seguimiento_cobro ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();