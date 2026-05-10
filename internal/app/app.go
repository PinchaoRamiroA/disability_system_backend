package app

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/shared/config"
	"disability_system_backend/internal/shared/logger"
	"disability_system_backend/internal/shared/router"

	"gorm.io/gorm"
)

type App struct {
	Config      *config.Config
	Logger      *logger.Logger
	DB          *gorm.DB
	Router      *router.Router
	AuthHandler *authhttp.AuthHandler
}

func New(cfg *config.Config, log *logger.Logger) *App {
	return &App{
		Config: cfg,
		Logger: log,
	}
}

func (a *App) WithDatabase(db *gorm.DB) *App {
	a.DB = db
	return a
}

func (a *App) WithRouter(r *router.Router) *App {
	a.Router = r
	return a
}