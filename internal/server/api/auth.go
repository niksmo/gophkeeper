package api

import (
	"context"

	"github.com/niksmo/gophkeeper/pkg/logger"
	authpb "github.com/niksmo/gophkeeper/proto/auth"
	"google.golang.org/grpc"
)

type AuthService interface {
	RegisterNewUser(
		ctx context.Context, login, password string,
	) (token string, err error)

	AuthorizeUser(
		ctx context.Context, login, password string,
	) (token string, err error)
}

type authAPI struct {
	authpb.UnimplementedAuthServer
	service AuthService
	logger  logger.Logger
}

func RegisterAuthAPI(
	logger logger.Logger, gRPCServer *grpc.Server, service AuthService,
) {
	authpb.RegisterAuthServer(
		gRPCServer, &authAPI{service: service, logger: logger},
	)
}

func (auth *authAPI) RegisterUser(
	ctx context.Context, in *authpb.RegUserRequest,
) (*authpb.RegUserResponse, error) {
	// TODO: check login and password
	token, err := auth.service.RegisterNewUser(
		ctx, in.GetLogin(), in.GetPassword(),
	)
	if err != nil {
		// TODO: handle errors
	}

	return &authpb.RegUserResponse{Token: token}, nil
}

func (auth *authAPI) AuthorizeUser(
	ctx context.Context, in *authpb.AuthUserRequest,
) (*authpb.AuthUserResponse, error) {
	// TODO: check login and password
	token, err := auth.service.AuthorizeUser(
		ctx, in.GetLogin(), in.GetPassword(),
	)
	if err != nil {
		// TODO: handle errors
	}

	return &authpb.AuthUserResponse{Token: token}, nil
}
