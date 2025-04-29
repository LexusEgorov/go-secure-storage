package app

import (
	grpcserv "auth/internal/grpc"

	"github.com/sirupsen/logrus"
)

type App struct {
	l    *logrus.Logger
	s    grpcserv.Server
	port int
}

func New(logger *logrus.Logger, port int, authProvider grpcserv.AuthProvider) *App {
	return &App{
		l:    logger,
		s:    *grpcserv.NewServer(logger, authProvider),
		port: port,
	}
}

func (a App) MustRun() {
	err := a.s.RunServer(a.port)

	if err != nil {
		panic(err)
	}
}
