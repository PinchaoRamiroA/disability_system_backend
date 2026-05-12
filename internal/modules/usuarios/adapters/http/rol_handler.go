package http

import (
	"net/http"
	"strconv"

	usuariosdto "disability_system_backend/internal/modules/usuarios/dto"
	"disability_system_backend/internal/modules/usuarios/mapper"
	"disability_system_backend/internal/modules/usuarios/usecase"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type RolHandler struct {
	usecase *usecase.RolUseCase
}

func NewRolHandler(usecase *usecase.RolUseCase) *RolHandler {
	return &RolHandler{usecase: usecase}
}

func (h *RolHandler) Listar(c *gin.Context) {
	query := usuariosdto.ListarRolesQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Parámetros inválidos", "INVALID_PARAMS", err.Error())
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}

	roles, total, err := h.usecase.Listar(c.Request.Context(), query.Page, query.Limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al listar roles", "INTERNAL_ERROR", err.Error())
		return
	}

	items := make([]usuariosdto.RolResponse, len(roles))
	for i, r := range roles {
		items[i] = mapper.ToRolResponse(r)
	}

	response.Paginated(c, items, total, int64(query.Page), int64(query.Limit))
}

func (h *RolHandler) Obtener(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	rol, err := h.usecase.Obtener(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al obtener rol", "INTERNAL_ERROR", err.Error())
		return
	}

	resp := mapper.ToRolResponse(*rol)
	response.Success(c, resp, "Rol obtenido correctamente")
}

func (h *RolHandler) Crear(c *gin.Context) {
	var req struct {
		Nombre   string   `json:"nombre" binding:"required,min=2,max=100"`
		Permisos []string `json:"permisos" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "VALIDATION_ERROR", err.Error())
		return
	}

	rol, err := h.usecase.Crear(c.Request.Context(), req.Nombre, req.Permisos)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al crear rol", "INTERNAL_ERROR", err.Error())
		return
	}

	resp := mapper.ToRolResponse(*rol)
	response.Created(c, resp, "Rol creado correctamente")
}

func (h *RolHandler) Actualizar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	var req struct {
		Nombre   string   `json:"nombre" binding:"required,min=2,max=100"`
		Permisos []string `json:"permisos" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "VALIDATION_ERROR", err.Error())
		return
	}

	rol, err := h.usecase.Actualizar(c.Request.Context(), id, req.Nombre, req.Permisos)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al actualizar rol", "INTERNAL_ERROR", err.Error())
		return
	}

	resp := mapper.ToRolResponse(*rol)
	response.Success(c, resp, "Rol actualizado correctamente")
}

func (h *RolHandler) Eliminar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	if err := h.usecase.Eliminar(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al eliminar rol", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, nil, "Rol eliminado correctamente")
}
