package gateway

import (
	"gateway/config"
	"gateway/internal/app"
	"gateway/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.Init(cfg.Env)

	application := app.New(logger, cfg.GRPC.Port, cfg.AuthConnect, cfg.DataConnect)

	application.MustRun()
}
