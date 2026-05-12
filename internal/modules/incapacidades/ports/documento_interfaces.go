package ports

type DocumentoRepository interface {
	Create(ctx context.Context, documento *domain.Documento) error
	FindByID(ctx context.Context, id uint64) (*domain.Documento, error)
	List(ctx context.Context, incapacidadID uint64, estado, tipo string, page, limit int) ([]domain.Documento, int64, error)
	Update(ctx context.Context, documento *domain.Documento) error
	Delete(ctx context.Context, id uint64) error
	ExistsIncapacidad(ctx context.Context, id uint64) (bool, error)
}

type HistorialRepository interface {
	Create(ctx context.Context, historial *domain.Historial) error
	List(ctx context.Context, incapacidadID uint64, tipoID *uint64, page, limit int) ([]domain.Historial, int64, error)
	FindByID(ctx context.Context, id uint64) (*domain.Historial, error)
	FindTipoByID(ctx context.Context, id uint64) (*domain.TipoHistorial, error)
}