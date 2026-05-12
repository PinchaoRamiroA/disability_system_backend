package dto

type SubirDocumentoRequest struct {
	IDIncapacidad uint64 `json:"id_incapacidad" binding:"required"`
	Nombre        string `json:"nombre" binding:"required"`
	Tipo          string `json:"tipo" binding:"required"`
	URL           string `json:"url" binding:"required"`
	Formato       string `json:"formato" binding:"required"`
}

type ValidarDocumentoRequest struct {
	Estado      string  `json:"estado" binding:"required"`
	Comentario  *string `json:"comentario,omitempty"`
}

type ListarDocumentosQuery struct {
	IDIncapacidad uint64 `form:"id_incapacidad"`
	Estado        string `form:"estado"`
	Tipo          string `form:"tipo"`
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
}