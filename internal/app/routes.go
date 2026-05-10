package app

import (
	"disability_system_backend/internal/shared/router"
)

func (a *App) RegisterRoutes(modules ...func(*router.Router)) {
	if a.Router == nil {
		a.Logger.Error("router not initialized")
		return
	}

	for _, module := range modules {
		module(a.Router)
	}
}

func (a *App) RegisterModuleRoutes(moduleFn func(v1 *router.APIVersion)) {
	if a.Router == nil {
		a.Logger.Error("router not initialized")
		return
	}

	moduleFn(a.Router.V1())
}

func (a *App) InitAuthRoutes() {
	if a.AuthHandler == nil {
		a.Logger.Error("auth handler not initialized")
		return
	}

	// Register Routes
	v1 := a.Router.V1()
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", a.AuthHandler.Login)
		authGroup.POST("/register", a.AuthHandler.Register)
		authGroup.POST("/refresh", a.AuthHandler.RefreshToken)
	}
}