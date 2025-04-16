package main

import (
	"auth/config"
	"auth/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.Init(cfg.Env)

	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Fatal("test")
	//TODO: init app
	//TODO: run app
}
