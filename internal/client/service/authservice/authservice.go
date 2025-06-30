package authservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/internal/client/service/syncservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	authbp "github.com/niksmo/gophkeeper/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAlreadyExists         = service.ErrAlreadyExists
	ErrCredentials           = errors.New("invalid credentials")
	ErrTimeoutExpired        = errors.New("deadline exceeded")
	ErrAuthServerUnavailable = errors.New("authorization service unavailable")
	ErrSyncAlreadyRunning    = errors.New("synchronization is already running")
)

type (
	AuthClient interface {
		RegisterUser(ctx context.Context, login, password string) (token string, err error)
		AuthorizeUser(ctx context.Context, login, password string) (token string, err error)
	}

	SyncExecuter interface {
		ExecSynchronization(ctx context.Context, token string) error
	}
)

type gRPCAuthClient struct {
	logger  logger.Logger
	client  authbp.AuthClient
	timeout time.Duration
}

func NewGRPCAuthClient(
	logger logger.Logger, client authbp.AuthClient, timeout time.Duration,
) AuthClient {

	return &gRPCAuthClient{logger, client, timeout}
}

func (c *gRPCAuthClient) RegisterUser(
	ctx context.Context, login, password string,
) (token string, err error) {
	ctx, cancel := c.setTimeout(ctx)
	defer cancel()

	reqData := &authbp.RegUserRequest{
		Login:    login,
		Password: []byte(password),
	}

	resData, err := c.client.RegisterUser(ctx, reqData)
	if err != nil {
		return "", c.handleRegistrationErr(err)
	}
	return resData.Token, nil
}

func (c *gRPCAuthClient) AuthorizeUser(
	ctx context.Context, login, password string,
) (token string, err error) {
	ctx, cancel := c.setTimeout(ctx)
	defer cancel()

	reqData := &authbp.AuthUserRequest{
		Login:    login,
		Password: []byte(password),
	}

	resData, err := c.client.AuthorizeUser(ctx, reqData)
	if err != nil {
		return "", c.handleAuthorizationErr(err)
	}

	return resData.Token, nil
}

func (c *gRPCAuthClient) setTimeout(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, c.timeout)
}

func (c *gRPCAuthClient) handleRegistrationErr(err error) error {
	if err == nil {
		return nil
	}

	const op = "gRPCAuthClient.handleRegistrationErr"
	log := c.logger.WithOp(op)

	switch status.Code(err) {
	case codes.Unavailable:
		log.Debug().Msg("unavailable")
		return ErrAuthServerUnavailable
	case codes.DeadlineExceeded:
		log.Debug().Msg("timeout expired")
		return ErrTimeoutExpired
	case codes.AlreadyExists:
		log.Debug().Msg("user already exists")
		return ErrAlreadyExists
	case codes.InvalidArgument:
		log.Debug().Err(err).Msg("invalid provided data")
		return err
	default:
		log.Debug().Err(err).Msg("failed to register new user")
		return err
	}
}

func (c *gRPCAuthClient) handleAuthorizationErr(err error) error {
	if err == nil {
		return nil
	}

	const op = "gRPCAuthClient.handleAuthorizationErr"
	log := c.logger.WithOp(op)

	switch status.Code(err) {
	case codes.Unavailable:
		log.Debug().Msg("unavailable")
		return ErrAuthServerUnavailable
	case codes.DeadlineExceeded:
		log.Debug().Msg("timeout expired")
		return ErrTimeoutExpired
	case codes.Unauthenticated:
		log.Debug().Err(err).Msg("invalid login or password")
		return ErrCredentials
	case codes.InvalidArgument:
		log.Debug().Err(err).Msg("invalid provided data")
		return err
	default:
		log.Error().Err(err).Msg("failed to authorize user")
		return err
	}
}

type UserRegistrar struct {
	logger      logger.Logger
	authClient  AuthClient
	syncStarter SyncExecuter
}

func NewUserRegistrar(
	logger logger.Logger, authClient AuthClient, syncStarter SyncExecuter,
) *UserRegistrar {
	return &UserRegistrar{logger, authClient, syncStarter}
}

func (r *UserRegistrar) RegisterUser(
	ctx context.Context, login, password string,
) error {
	const op = "UserRegistrar.RegisterUser"

	token, err := r.registerUser(ctx, login, password)
	if err != nil {
		return r.error(op, err)
	}

	if err := r.startSynchronization(ctx, token); err != nil {
		return r.error(op, err)
	}

	return nil
}

func (r *UserRegistrar) registerUser(
	ctx context.Context, login, password string,
) (string, error) {
	token, err := r.authClient.RegisterUser(ctx, login, password)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (r *UserRegistrar) startSynchronization(
	ctx context.Context, token string,
) error {
	const op = "UserRegistrar.startSynchronization"
	log := r.logger.WithOp(op)
	return startSynchronization(ctx, log, r.syncStarter, token)
}

func (r *UserRegistrar) error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

type UserAuthorizer struct {
	logger      logger.Logger
	authClient  AuthClient
	syncStarter SyncExecuter
}

func NewUserAuthorizer(
	logger logger.Logger, authClient AuthClient, syncStarter SyncExecuter,
) *UserAuthorizer {
	return &UserAuthorizer{logger, authClient, syncStarter}
}

func (a *UserAuthorizer) AuthorizeUser(
	ctx context.Context, login, password string,
) error {
	const op = "AuthService.AuthorizeUser"

	token, err := a.authorizeUser(ctx, login, password)
	if err != nil {
		return a.error(op, err)
	}

	if err := a.startSynchronization(ctx, token); err != nil {
		return a.error(op, err)
	}

	return nil
}

func (a *UserAuthorizer) authorizeUser(
	ctx context.Context, login, password string,
) (string, error) {
	token, err := a.authClient.AuthorizeUser(ctx, login, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *UserAuthorizer) startSynchronization(ctx context.Context, token string) error {
	const op = "UserAuthorizer.startSynchronization"
	log := r.logger.WithOp(op)
	return startSynchronization(ctx, log, r.syncStarter, token)
}

func (a *UserAuthorizer) error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

func startSynchronization(
	ctx context.Context, log logger.Logger, ss SyncExecuter, token string,
) error {
	err := ss.ExecSynchronization(ctx, token)
	if err != nil {
		if errors.Is(err, syncservice.ErrPIDConflict) {
			log.Debug().Err(err).Msg("synchronization is already running")
			return ErrSyncAlreadyRunning
		}
		log.Debug().Err(err).Msg("failed to start synchronization")
		return err
	}
	return nil
}
