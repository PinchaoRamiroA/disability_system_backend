package usecase

import (
	"context"
	"strings"
	"time"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

const (
	DefaultPlazoEntregaDias = 3
)

type DocumentosRequeridosResult struct {
	Requeridos         []domain.TipoDocumento
	Faltantes          []domain.TipoDocumento
	TodosPresentes     bool
	FechaLimiteEntrega time.Time
	DiasRestantes      int
}

type IncapacidadDocumentosService struct {
	repo ports.IncapacidadRepository
}

func NewIncapacidadDocumentosService(repo ports.IncapacidadRepository) *IncapacidadDocumentosService {
	return &IncapacidadDocumentosService{repo: repo}
}

func (s *IncapacidadDocumentosService) ObtenerDocumentosRequeridos(ctx context.Context, tipoID uint64) ([]domain.TipoDocumento, error) {
	tipo, err := s.repo.FindTipoByID(ctx, tipoID)
	if err != nil {
		return nil, err
	}

	if len(tipo.DocumentosRequeridos) == 0 {
		return []domain.TipoDocumento{}, nil
	}

	return s.repo.FindTiposDocumentoByCodigo(ctx, tipo.DocumentosRequeridos)
}

func (s *IncapacidadDocumentosService) ValidarDocumentosRequeridos(ctx context.Context, tipoID uint64, documentosCargados []domain.Documento) (*DocumentosRequeridosResult, error) {
	tipo, err := s.repo.FindTipoByID(ctx, tipoID)
	if err != nil {
		return nil, err
	}

	requeridos, err := s.repo.FindTiposDocumentoByCodigo(ctx, tipo.DocumentosRequeridos)
	if err != nil {
		return nil, err
	}

	cargadosMap := make(map[string]bool)
	for _, doc := range documentosCargados {
		cargadosMap[doc.Tipo] = true
	}

	var faltantes []domain.TipoDocumento
	for _, req := range requeridos {
		if req.Requerido && !cargadosMap[req.Codigo] {
			faltantes = append(faltantes, req)
		}
	}

	result := &DocumentosRequeridosResult{
		Requeridos:     requeridos,
		Faltantes:      faltantes,
		TodosPresentes: len(faltantes) == 0,
	}

	return result, nil
}

func (s *IncapacidadDocumentosService) ObtenerPlazoEntrega(ctx context.Context, tipoIncapacidadID uint64, canalRecepcion string) (int, error) {
	return DefaultPlazoEntregaDias, nil
}

func (s *IncapacidadDocumentosService) ObtenerFechaLimiteEntrega(ctx context.Context, fechaCreacion time.Time, tipoIncapacidadID uint64, canalRecepcion string) (time.Time, error) {
	plazo, err := s.ObtenerPlazoEntrega(ctx, tipoIncapacidadID, canalRecepcion)
	if err != nil {
		return time.Time{}, err
	}

	return fechaCreacion.AddDate(0, 0, plazo), nil
}

func (s *IncapacidadDocumentosService) ObtenerPlazoTranscripcion(ctx context.Context, entidadID uint64) (int, error) {
	entidad, err := s.repo.FindEntidadByID(ctx, entidadID)
	if err != nil {
		return 0, err
	}

	if entidad.PlazoTranscripcionDias == nil || *entidad.PlazoTranscripcionDias == 0 {
		return 5, nil
	}
	return *entidad.PlazoTranscripcionDias, nil
}

func (s *IncapacidadDocumentosService) ObtenerFechaLimiteTranscripcion(ctx context.Context, fechaRecepcion time.Time, entidadID uint64) (time.Time, error) {
	plazo, err := s.ObtenerPlazoTranscripcion(ctx, entidadID)
	if err != nil {
		return time.Time{}, err
	}

	return fechaRecepcion.AddDate(0, 0, plazo), nil
}

func (s *IncapacidadDocumentosService) ObtenerTiempoMaximoPago(ctx context.Context, entidadID uint64) (int, error) {
	entidad, err := s.repo.FindEntidadByID(ctx, entidadID)
	if err != nil {
		return 0, err
	}

	if entidad.TiempoMaximoPagoDias == nil || *entidad.TiempoMaximoPagoDias == 0 {
		return 30, nil
	}
	return *entidad.TiempoMaximoPagoDias, nil
}

func (s *IncapacidadDocumentosService) VerificarEstadoTransicion(ctx context.Context, incapacidad *domain.Incapacidad, nuevoEstado *domain.EstadoIncapacidad) error {
	if incapacidad.Estado != nil && !incapacidad.Estado.PermiteTransicion {
		return apperrors.ErrConflict.WithMessage("el estado actual no permite transiciones")
	}

	if nuevoEstado.Nombre == "Rechazada" || nuevoEstado.Nombre == "Archivada" || nuevoEstado.Nombre == "Cerrada" {
		return nil
	}

	validTransitions := map[string][]string{
		"Recibida":                 {"En validación documental", "Pendiente transcripción", "Documentación incompleta", "Archivada"},
		"En validación documental": {"Documentación incompleta", "Pendiente transcripción", "Rechazada", "Archivada"},
		"Documentación incompleta": {"En validación documental", "Pendiente transcripción", "Rechazada", "Archivada"},
		"Pendiente transcripción":  {"Transcrita", "Rechazada", "Archivada"},
		"Transcrita":               {"En verificación EPS", "Rechazada", "Archivada"},
		"En verificación EPS":      {"Aprobada", "Rechazada", "Archivada"},
		"Aprobada":                 {"Cobrada", "Rechazada", "Archivada"},
		"Cobrada":                  {"Pendiente pago", "Rechazada", "Archivada"},
		"Pendiente pago":           {"Pagada", "Cobro persuasivo", "Cobro jurídico", "Rechazada", "Archivada"},
		"Pagada":                   {"En conciliación", "Archivada"},
		"En conciliación":          {"Conciliada", "Archivada"},
		"Conciliada":               {"Archivada", "Cerrada"},
		"Cobro persuasivo":         {"Cobro jurídico", "Pagada", "Archivada"},
		"Cobro jurídico":           {"Pagada", "Rechazada", "Archivada"},
	}

	if currentState, ok := validTransitions[incapacidad.Estado.Nombre]; ok {
		for _, valid := range currentState {
			if valid == nuevoEstado.Nombre {
				return nil
			}
		}
	}

	return apperrors.ErrConflict.WithMessage("transición de estado no válida: de " + incapacidad.Estado.Nombre + " a " + nuevoEstado.Nombre)
}

func (s *IncapacidadDocumentosService) ObtenerDiasTranscurridos(fechaInicio time.Time) int {
	duracion := time.Since(fechaInicio)
	return int(duracion.Hours() / 24)
}

func (s *IncapacidadDocumentosService) EsIncapacidadVencida(fechaInicio time.Time, diasMaximos int) bool {
	dias := s.ObtenerDiasTranscurridos(fechaInicio)
	return dias > diasMaximos
}

func (s *IncapacidadDocumentosService) ObtenerAlertasVencimientos(fechaInicio time.Time) []string {
	var alertas []string
	dias := s.ObtenerDiasTranscurridos(fechaInicio)

	limites := []struct {
		dias    int
		mensaje string
	}{
		{90, "Alerta: Incapacidad supera 90 días"},
		{120, "Alerta: Incapacidad supera 120 días"},
		{150, "Alerta: Incapacidad supera 150 días"},
		{180, "Alerta: Incapacidad supera 180 días - Revisar estado"},
	}

	for _, limite := range limites {
		if dias >= limite.dias {
			alertas = append(alertas, limite.mensaje)
		}
	}

	return alertas
}

func normalizeDocumentKey(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	replacements := map[string]string{
		"á": "a",
		"é": "e",
		"í": "i",
		"ó": "o",
		"ú": "u",
		"ñ": "n",
	}
	for old, newValue := range replacements {
		normalized = strings.ReplaceAll(normalized, old, newValue)
	}
	normalized = strings.ReplaceAll(normalized, " de ", " ")
	normalized = strings.ReplaceAll(normalized, "-", " ")
	normalized = strings.ReplaceAll(normalized, "_", " ")
	normalized = strings.Join(strings.Fields(normalized), "_")
	return normalized
}
