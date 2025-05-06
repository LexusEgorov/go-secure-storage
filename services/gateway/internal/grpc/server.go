package grpcsrv

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/LexusEgorov/go-secure-storage-protos/gen/golang/authpb"
	"github.com/LexusEgorov/go-secure-storage-protos/gen/golang/datapb"
	"github.com/LexusEgorov/go-secure-storage-protos/gen/golang/gatewaypb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	authClient authpb.AuthClient
	dataClient datapb.DataClient
	l          *logrus.Logger
	s          *grpc.Server
}

func (s Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	s.l.Info("register: ", req)
	return s.authClient.Register(ctx, req)
}

func (s Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.l.Info("login: ", req)
	return s.authClient.Login(ctx, req)
}

func (s Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	s.l.Info("validate: ", req)
	return s.authClient.Refresh(ctx, req)
}

func (s Server) Add(ctx context.Context, req *datapb.AddRequest) (*datapb.AddResponse, error) {
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values, found := md["authorization"]; found {
			token = values[0]
		}
	}

	if token == "" {
		return nil, errors.New("unauthorized")
	}

	check, err := s.authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	if !check.Valid {
		return nil, errors.New("unauthorized")
	}

	return s.dataClient.Add(ctx, req)
}

func (s Server) Get(ctx context.Context, req *datapb.GetRequest) (*datapb.GetResponse, error) {
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values, found := md["authorization"]; found {
			token = values[0]
		}
	}

	if token == "" {
		return nil, errors.New("unauthorized")
	}

	check, err := s.authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	if !check.Valid {
		return nil, errors.New("unauthorized")
	}

	return s.dataClient.Get(ctx, req)
}

func NewServer(l *logrus.Logger, authConn, dataConn string) *Server {
	grpcServer := grpc.NewServer()

	connAuth, err := grpc.NewClient(authConn)

	if err != nil {
		panic(err)
	}

	connData, err := grpc.NewClient(dataConn)

	if err != nil {
		panic(err)
	}

	server := Server{
		authClient: authpb.NewAuthClient(connAuth),
		dataClient: datapb.NewDataClient(connData),
		l:          l,
		s:          grpcServer,
	}

	gatewaypb.RegisterGatewayServer(grpcServer, server)

	return &server
}

func (s Server) RunServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		s.l.Panic(err)
		return err
	}

	s.l.Info("server is running on ", port, " port")

	if err := s.s.Serve(lis); err != nil {
		s.l.Panic(err)
		return err
	}

	return nil
}
