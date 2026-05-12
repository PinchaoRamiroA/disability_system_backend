-- Revertir: Volver FKs a gestion_humana

-- =====================================================
-- INCAPACIDAD
-- =====================================================
ALTER TABLE incapacidad DROP CONSTRAINT IF EXISTS incapacidad_created_by_fkey;
ALTER TABLE incapacidad ADD CONSTRAINT incapacidad_created_by_fkey FOREIGN KEY (created_by) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL;

-- =====================================================
-- DOCUMENTO
-- =====================================================
ALTER TABLE documento DROP CONSTRAINT IF EXISTS documento_validado_por_fkey;
ALTER TABLE documento ADD CONSTRAINT documento_validado_por_fkey FOREIGN KEY (validado_por) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL;

-- =====================================================
-- PAGO
-- =====================================================
ALTER TABLE pago DROP CONSTRAINT IF EXISTS pago_registrado_por_fkey;
ALTER TABLE pago ADD CONSTRAINT pago_registrado_por_fkey FOREIGN KEY (registrado_por) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL;

-- =====================================================
-- SEGUIMIENTO_COBRO
-- =====================================================
ALTER TABLE seguimiento_cobro DROP CONSTRAINT IF EXISTS seguimiento_cobro_gestionado_por_fkey;
ALTER TABLE seguimiento_cobro ADD CONSTRAINT seguimiento_cobro_gestionado_por_fkey FOREIGN KEY (gestionado_por) REFERENCES gestion_humana(id_usuario) ON DELETE SET NULL;
