-- Up
ALTER TABLE tipo_documento ADD COLUMN codigo VARCHAR(100) UNIQUE;

-- Poblar datos iniciales basados en los nombres existentes
UPDATE tipo_documento SET codigo = 'certificado_incapacidad' WHERE nombre = 'Certificado de incapacidad';
UPDATE tipo_documento SET codigo = 'epicrisis' WHERE nombre = 'Epicrisis';
UPDATE tipo_documento SET codigo = 'furips' WHERE nombre = 'FURIPS';
UPDATE tipo_documento SET codigo = 'historia_clinica' WHERE nombre = 'Historia clínica';
UPDATE tipo_documento SET codigo = 'certificado_nacido_vivo' WHERE nombre = 'Certificado de nacido vivo';
UPDATE tipo_documento SET codigo = 'registro_civil' WHERE nombre = 'Registro civil';
UPDATE tipo_documento SET codigo = 'documento_identidad' WHERE nombre = 'Documento de identidad';
UPDATE tipo_documento SET codigo = 'soporte_atencion_medica' WHERE nombre = 'Soporte de atención médica';
UPDATE tipo_documento SET codigo = 'concepto_rehabilitacion' WHERE nombre = 'Concepto de rehabilitación';
UPDATE tipo_documento SET codigo = 'evidencia_radicacion' WHERE nombre = 'Evidencia de radicación';
UPDATE tipo_documento SET codigo = 'soporte_pago' WHERE nombre = 'Soporte de pago';
UPDATE tipo_documento SET codigo = 'formato_seguimiento' WHERE nombre = 'Formato de seguimiento';

-- Hacer el campo obligatorio después de poblarlo
ALTER TABLE tipo_documento ALTER COLUMN codigo SET NOT NULL;
