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

// Listar godoc
// @Summary Listar usuarios
// @Description Lista todos los usuarios con paginación
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Límite de resultados" default(20)
// @Param estado query bool false "Filtrar por estado"
// @Param id_rol query int false "Filtrar por rol"
// @Param search query string false "Buscar por nombre, correo o documento"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /usuarios [get]
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

	usuarios, roles, total, err := h.usecase.Listar(c.Request.Context(), query.Page, query.Limit, query.Estado, query.IDRol, query.Search)
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

// Obtener godoc
// @Summary Obtener usuario por ID
// @Description Obtiene los detalles de un usuario específico
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /usuarios/{id} [get]
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

// Crear godoc
// @Summary Crear usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usuariosdto.CrearUsuarioRequest true "Datos del usuario"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /usuarios [post]
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

// Actualizar godoc
// @Summary Actualizar usuario
// @Description Actualiza los datos de un usuario existente
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Param request body usuariosdto.ActualizarUsuarioRequest true "Datos a actualizar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /usuarios/{id} [put]
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
		IDRol         *uint64
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

// CambiarEstado godoc
// @Summary Cambiar estado de usuario
// @Description Activa o desactiva un usuario
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Param request body usuariosdto.CambiarEstadoUsuarioRequest true "Nuevo estado"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /usuarios/{id}/estado [patch]
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

// CambiarPassword godoc
// @Summary Cambiar contraseña
// @Description Cambia la contraseña del usuario autenticado
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usuariosdto.CambiarPasswordRequest true "Contraseñas"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /usuarios/{id}/password [post]
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

// AsignarRol godoc
// @Summary Asignar rol a usuario
// @Description Asigna un rol diferente a un usuario
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Param request body usuariosdto.AsignarRolRequest true "ID del rol"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /usuarios/{id}/rol [post]
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

// Eliminar godoc
// @Summary Eliminar usuario
// @Description Elimina un usuario del sistema (soft delete)
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /usuarios/{id} [delete]
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
