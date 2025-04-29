package grpcserv

import (
	"context"
	"fmt"
	"net"

	"github.com/LexusEgorov/go-secure-storage-protos/gen/golang/authpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthProvider interface{}

type Server struct {
	authpb.UnimplementedAuthServer
	l    *logrus.Logger
	s    *grpc.Server
	auth AuthProvider
}

func NewServer(l *logrus.Logger, auth AuthProvider) *Server {
	grpcServer := grpc.NewServer()

	server := Server{
		l:    l,
		s:    grpcServer,
		auth: auth,
	}

	authpb.RegisterAuthServer(grpcServer, server)

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

func (s Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	s.l.Info("register: ", req)
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func (s Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.l.Info("login: ", req)
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (s Server) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	s.l.Info("validate: ", req)
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (s Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	s.l.Info("refresh: ", req)
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
