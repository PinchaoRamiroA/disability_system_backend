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

type UsuarioHandler struct {
	usecase *usecase.UsuarioUseCase
}

func NewUsuarioHandler(usecase *usecase.UsuarioUseCase) *UsuarioHandler {
	return &UsuarioHandler{usecase: usecase}
}

func (h *UsuarioHandler) Listar(c *gin.Context) {
	query := usuariosdto.ListarUsuariosQuery{}
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

	usuarios, roles, total, err := h.usecase.Listar(c.Request.Context(), query.Page, query.Limit, query.Estado, query.IDRol)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al listar usuarios", "INTERNAL_ERROR", err.Error())
		return
	}

	rolesMap := make(map[uint64]string)
	for _, rol := range roles {
		rolesMap[rol.ID] = rol.Nombre
	}

	items := make([]usuariosdto.UsuarioResponse, len(usuarios))
	for i, u := range usuarios {
		nombreRol := rolesMap[u.IDRol]
		items[i] = mapper.ToUsuarioResponse(u, nombreRol)
	}

	response.Paginated(c, items, total, int64(query.Page), int64(query.Limit))
}

func (h *UsuarioHandler) Obtener(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	usuario, rol, err := h.usecase.Obtener(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al obtener usuario", "INTERNAL_ERROR", err.Error())
		return
	}

	resp := mapper.ToUsuarioResponse(*usuario, rol.Nombre)
	response.Success(c, resp, "Usuario obtenido correctamente")
}

func (h *UsuarioHandler) Crear(c *gin.Context) {
	var req usuariosdto.CrearUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "VALIDATION_ERROR", err.Error())
		return
	}

	usuario, _, err := h.usecase.Crear(c.Request.Context(), struct {
		IDRol           uint64
		Nombre          string
		Correo          string
		NumeroCelular   *string
		Direccion       *string
		Password        string
		NumeroDocumento string
		NumeroAcudiente *string
	}{
		IDRol:           req.IDRol,
		Nombre:          req.Nombre,
		Correo:          req.Correo,
		NumeroCelular:   req.NumeroCelular,
		Direccion:       req.Direccion,
		Password:        req.Password,
		NumeroDocumento: req.NumeroDocumento,
		NumeroAcudiente: req.NumeroAcudiente,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al crear usuario", "INTERNAL_ERROR", err.Error())
		return
	}

	_, rol, err := h.usecase.Obtener(c.Request.Context(), usuario.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al obtener rol", "INTERNAL_ERROR", err.Error())
		return
	}

	resp := mapper.ToUsuarioResponse(*usuario, rol.Nombre)
	response.Created(c, resp, "Usuario creado correctamente")
}

func (h *UsuarioHandler) Actualizar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	var req usuariosdto.ActualizarUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "VALIDATION_ERROR", err.Error())
		return
	}

	usuario, err := h.usecase.Actualizar(c.Request.Context(), id, struct {
		IDRol         uint64
		Nombre        string
		Correo        string
		NumeroCelular *string
		Direccion     *string
	}{
		IDRol:         req.IDRol,
		Nombre:        req.Nombre,
		Correo:        req.Correo,
		NumeroCelular: req.NumeroCelular,
		Direccion:     req.Direccion,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al actualizar usuario", "INTERNAL_ERROR", err.Error())
		return
	}

	_, rol, err := h.usecase.Obtener(c.Request.Context(), usuario.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al obtener rol", "INTERNAL_ERROR", err.Error())
		return
	}

	resp := mapper.ToUsuarioResponse(*usuario, rol.Nombre)
	response.Success(c, resp, "Usuario actualizado correctamente")
}

func (h *UsuarioHandler) CambiarEstado(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	var req usuariosdto.CambiarEstadoUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.usecase.CambiarEstado(c.Request.Context(), id, req.Estado); err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al cambiar estado", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, nil, "Estado actualizado correctamente")
}

func (h *UsuarioHandler) CambiarPassword(c *gin.Context) {
	var req usuariosdto.CambiarPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "VALIDATION_ERROR", err.Error())
		return
	}

	id := c.GetUint64("user_id")

	if err := h.usecase.CambiarPassword(c.Request.Context(), id, req.PasswordActual, req.PasswordNuevo); err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al cambiar contraseña", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, nil, "Contraseña actualizada correctamente")
}

func (h *UsuarioHandler) AsignarRol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	var req usuariosdto.AsignarRolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.usecase.AsignarRol(c.Request.Context(), id, req.IDRol); err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al asignar rol", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, nil, "Rol asignado correctamente")
}

func (h *UsuarioHandler) Eliminar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	if err := h.usecase.Eliminar(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al eliminar usuario", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, nil, "Usuario eliminado correctamente")
}
