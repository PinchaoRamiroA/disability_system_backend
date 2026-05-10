package app

import (
	"disability_system_backend/internal/shared/database"
	"disability_system_backend/internal/shared/router"
	"disability_system_backend/internal/shared/middleware"
)

func (a *App) InitRouter() error {
	if a.DB == nil {
		a.Logger.Error("database not initialized, cannot create router")
		return ErrDatabaseNotInitialized
	}

	r := router.New(a.Config, a.Logger)
	r.Use(middleware.Recovery(a.Logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(a.Logger))

	r.SetupHealth(a.DB)

	a.Router = r
	return nil
}

func (a *App) InitDependencies() error {
	db, err := database.NewConnection(a.Config)
	if err != nil {
		a.Logger.Error("failed to connect to database", "error", err)
		return err
	}

	if err := database.HealthCheck(db); err != nil {
		a.Logger.Error("database health check failed", "error", err)
		return err
	}

	a.DB = db
	return nil
}