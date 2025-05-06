package main

import (
	"auth/config"
	"auth/internal/app"
	"auth/internal/auth"
	"auth/internal/logger"
	"auth/internal/storage"
)

func main() {

	cfg := config.MustLoad()
	logger := logger.Init(cfg.Env)
	DBStorage := storage.NewDB(cfg.DBConnect)
	authProvider := auth.NewAuth(DBStorage)

	application := app.New(logger, cfg.GRPC.Port, authProvider)

	application.MustRun()
}
