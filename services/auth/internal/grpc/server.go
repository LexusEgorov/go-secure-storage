package authgrpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/LexusEgorov/go-secure-storage-protos/gen/golang/authpb"
)

type serverAPI struct {
	authpb.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Register(context.Context, *authpb.RegisterRequest) (*authpb.RegisterResponse, error)
	Login(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error)
	ValidateToken(context.Context, *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error)
	Refresh(context.Context, *authpb.RefreshRequest) (*authpb.RefreshResponse, error)
}

func Register(grpcServer *grpc.Server) {
	authpb.RegisterAuthServer(grpcServer, serverAPI{})
}

func (s serverAPI) Register(
	ctx context.Context,
	reg *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {
	//TODO
	return nil, errors.New("unimplemented")
}

func (s serverAPI) Login(
	ctx context.Context,
	login *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	//TODO
	return nil, errors.New("unimplemented")
}

func (s serverAPI) ValidateToken(
	ctx context.Context,
	token *authpb.ValidateTokenRequest,
) (*authpb.ValidateTokenResponse, error) {
	//TODO
	return nil, errors.New("unimplemented")
}

func (s serverAPI) Refresh(
	ctx context.Context,
	refresh *authpb.RefreshRequest,
) (*authpb.RefreshResponse, error) {
	//TODO
	return nil, errors.New("unimplemented")
}
