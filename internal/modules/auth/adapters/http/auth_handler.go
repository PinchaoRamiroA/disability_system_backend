package authhttp

import (
	"errors"

	_ "disability_system_backend/docs"
	"disability_system_backend/internal/modules/auth/dto"
	"disability_system_backend/internal/modules/auth/mapper"
	"disability_system_backend/internal/modules/auth/usecase"
	"disability_system_backend/internal/shared/response"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	loginUseCase    *usecase.LoginUseCase
	registerUseCase *usecase.RegisterUseCase
	refreshUseCase  *usecase.RefreshTokenUseCase
}

func NewAuthHandler(
	loginUseCase *usecase.LoginUseCase,
	registerUseCase *usecase.RegisterUseCase,
	refreshUseCase *usecase.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		loginUseCase:    loginUseCase,
		registerUseCase: registerUseCase,
		refreshUseCase:  refreshUseCase,
	}
}

// Login godoc
// @Summary Iniciar sesión
// @Description Autentica un usuario y devuelve tokens de acceso
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciales de login"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	tokens, user, role, err := h.loginUseCase.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := mapper.ToLoginResponse(user, role, tokens.AccessToken, tokens.RefreshToken, h.loginUseCase.GetExpirationSeconds())
	response.Success(c, resp, "login exitoso")
}

// Register godoc
// @Summary Registrar usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Datos del usuario"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	user, err := h.registerUseCase.Execute(
		c.Request.Context(),
		req.Nombre,
		req.Email,
		req.Password,
		req.NumeroDocumento,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := mapper.ToUserResponse(user, nil)
	response.Created(c, resp, "usuario registrado exitosamente")
}

// RefreshToken godoc
// @Summary Renovar token
// @Description Obtiene nuevos tokens usando un refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response{data=dto.TokenResponse}
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "token requerido", err.Error())
		return
	}

	tokens, err := h.refreshUseCase.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := mapper.ToTokenResponse(tokens.AccessToken, tokens.RefreshToken, 0)
	response.Success(c, resp, "token renovado")
}

func handleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	response.InternalError(c, "error interno", "INTERNAL_ERROR")
}