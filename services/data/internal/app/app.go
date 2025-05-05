package app

import (
	grpcserv "data/internal/grpc"

	"github.com/sirupsen/logrus"
)

type App struct {
	l    *logrus.Logger
	s    grpcserv.Server
	port int
}

func New(logger *logrus.Logger, port int, dataProvider grpcserv.DataProvider) *App {
	return &App{
		l:    logger,
		s:    *grpcserv.NewServer(logger, dataProvider),
		port: port,
	}
}

func (a App) MustRun() {
	err := a.s.RunServer(a.port)

	if err != nil {
		panic(err)
	}
}
