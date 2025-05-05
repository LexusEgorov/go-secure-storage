package grpcserv

import (
	"auth/internal/models"
	"context"
	"fmt"
	"net"

	"github.com/LexusEgorov/go-secure-storage-protos/gen/golang/authpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type AuthProvider interface {
	Register(email string, password string) (*models.Credentials, error)
	Auth(email, password string) (*models.Credentials, error)
	Refresh(token string) (*models.Credentials, error)
	Validate(token string) bool
}

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
	credentials, err := s.auth.Register(req.GetEmail(), req.GetPassword())

	if err != nil {
		s.l.Error(err)
		return &authpb.RegisterResponse{
			Ok: false,
			Response: &authpb.RegisterResponse_Bad{
				Bad: &authpb.BadResponse{
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &authpb.RegisterResponse{
		Ok: true,
		Response: &authpb.RegisterResponse_Success{
			Success: &authpb.SuccessTokenResponse{
				Jwt:     credentials.JWT,
				Refresh: credentials.Refresh,
			},
		},
	}, nil
}

func (s Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.l.Info("login: ", req)

	credentials, err := s.auth.Auth(req.GetEmail(), req.GetPassword())

	if err != nil {
		s.l.Error(err)
		return &authpb.LoginResponse{
			Ok: false,
			Response: &authpb.LoginResponse_Bad{
				Bad: &authpb.BadResponse{
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &authpb.LoginResponse{
		Ok: true,
		Response: &authpb.LoginResponse_Success{
			Success: &authpb.SuccessTokenResponse{
				Jwt:     credentials.JWT,
				Refresh: credentials.Refresh,
			},
		},
	}, nil
}

func (s Server) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	s.l.Info("validate: ", req)

	isValid := s.auth.Validate(req.GetToken())

	return &authpb.ValidateTokenResponse{
		Valid: isValid,
	}, nil
}
func (s Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	s.l.Info("refresh: ", req)

	credentials, err := s.auth.Refresh(req.GetRefresh())

	if err != nil {
		return &authpb.RefreshResponse{
			Ok: false,
			Response: &authpb.RefreshResponse_Bad{
				Bad: &authpb.BadResponse{
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &authpb.RefreshResponse{
		Ok: false,
		Response: &authpb.RefreshResponse_Success{
			Success: &authpb.SuccessTokenResponse{
				Jwt:     credentials.JWT,
				Refresh: credentials.Refresh,
			},
		},
	}, nil
}
