package app

import (
	grpcapp "auth/internal/app/grpc"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *logrus.Logger, port int) *App {
	gRPCServer := grpcapp.New(log, port)

	return &App{
		GRPCServer: gRPCServer,
	}
}
