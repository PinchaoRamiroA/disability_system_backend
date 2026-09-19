// @title Sistema de Gestión de Incapacidades
// @version 1.0
// @description API REST para la gestión de incapacidades médicas, cobro a EPS/ARL y seguimiento documental
// @BasePath /api/v1
// @schemes http https

package main

import (
	"fmt"
	"os"

	"disability_system_backend/internal/app"
	"disability_system_backend/internal/shared/config"
	"disability_system_backend/internal/shared/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.Env)

	log.Info("Starting server",
		"name", cfg.App.Name,
		"env", cfg.App.Env,
		"version", cfg.App.Version,
	)

	a := app.New(cfg, log)

	if err := a.InitDependencies(); err != nil {
		log.Error("failed to initialize dependencies", "error", err)
		os.Exit(1)
	}

	if err := a.InitRouter(); err != nil {
		log.Error("failed to initialize router", "error", err)
		os.Exit(1)
	}

	jwtService := a.InitAuth()
	_ = jwtService

	port := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Info("server listening", "port", cfg.Server.Port)

	if err := a.Router.Run(port); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}