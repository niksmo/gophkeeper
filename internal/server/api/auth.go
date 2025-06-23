package api

import (
	"context"
	"errors"

	"github.com/niksmo/gophkeeper/internal/server/service/authservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	authpb "github.com/niksmo/gophkeeper/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	RegisterNewUser(
		ctx context.Context, login string, password []byte,
	) (token string, err error)

	AuthorizeUser(
		ctx context.Context, login string, password []byte,
	) (token string, err error)
}

type authHandler struct {
	authpb.UnimplementedAuthServer
	service AuthService
	logger  logger.Logger
}

func RegisterAuthAPI(
	logger logger.Logger, gRPCServer *grpc.Server, service AuthService,
) {
	authpb.RegisterAuthServer(
		gRPCServer, &authHandler{service: service, logger: logger},
	)
}

func (h *authHandler) RegisterUser(
	ctx context.Context, in *authpb.RegUserRequest,
) (*authpb.RegUserResponse, error) {
	const op = "authAPI.RegisterUser"
	log := h.logger.WithOp(op)

	// TODO: verify on pattern login and password

	token, err := h.service.RegisterNewUser(
		ctx, in.GetLogin(), in.GetPassword(),
	)
	if err != nil {
		if errors.Is(err, authservice.ErrAlreadyExists) {
			log.Debug().Str("login", in.Login).Msg("user already exists")
			return nil, status.Errorf(
				codes.AlreadyExists,
				"user with login %s already exists",
				in.Login,
			)
		}
		log.Error().Err(err).Msg("internal errror")
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authpb.RegUserResponse{Token: token}, nil
}

func (h *authHandler) AuthorizeUser(
	ctx context.Context, in *authpb.AuthUserRequest,
) (*authpb.AuthUserResponse, error) {
	// TODO: verify on pattern login and password
	token, err := h.service.AuthorizeUser(
		ctx, in.GetLogin(), in.GetPassword(),
	)
	if err != nil {
		// TODO: handle errors
	}

	return &authpb.AuthUserResponse{Token: token}, nil
}
