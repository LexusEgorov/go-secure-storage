package grpcapp

import (
	authgrpc "auth/internal/grpc"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	log        *logrus.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *logrus.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "auth.Run"

	a.log.Info(
		fmt.Printf("%s: %s\n", op, "starting gRPC server"),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%d", a.port))

	if err != nil {
		a.log.Error(fmt.Sprintf("%s: %s", op, err))
		return err
	}

	a.log.Info(
		fmt.Printf("%s: %s %s\n", op, "gRPC server is running on", listener.Addr().String()),
	)

	if err := a.gRPCServer.Serve(listener); err != nil {
		a.log.Error(fmt.Sprintf("%s: %s", op, err))
		return err
	}

	return nil
}

func (a *App) Stop() {
	const op = "auth.Stop"

	a.log.Info(
		fmt.Sprintf("%s: %s", op, "stopping gRPC server"),
	)

	a.gRPCServer.GracefulStop()
}
