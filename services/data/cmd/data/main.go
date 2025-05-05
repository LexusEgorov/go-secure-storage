package main

import (
	"data/config"
	"data/internal/app"
	"data/internal/data"
	"data/internal/logger"
	"data/internal/storage"
)

func main() {

	cfg := config.MustLoad()
	logger := logger.Init(cfg.Env)
	DBStorage := storage.NewDB(cfg.DBConnect)
	dataProvider := data.NewData(DBStorage)

	application := app.New(logger, cfg.GRPC.Port, dataProvider)

	application.MustRun()
}
