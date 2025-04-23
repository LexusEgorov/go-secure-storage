package main

import (
	"auth/config"
	"auth/internal/app"
	"auth/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.Init(cfg.Env)

	application := app.New(logger, cfg.GRPC.Port)

	application.GRPCServer.MustRun()
}
