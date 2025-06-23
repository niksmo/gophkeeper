package authservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/service"
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
)

const (
	timeout = 5 * time.Second
)

type (
	ClientConnector interface {
		ConnClient() (authbp.AuthClient, error)
		Close() error
	}

	SyncStarter interface {
		StartSynchronization(ctx context.Context, token string) error
	}
)

type UserRegistrar struct {
	logger      logger.Logger
	clientConn  ClientConnector
	syncStarter SyncStarter
}

func NewUserRegistrar(
	logger logger.Logger, clientConn ClientConnector, syncStarter SyncStarter,
) *UserRegistrar {
	return &UserRegistrar{logger, clientConn, syncStarter}
}

func (r *UserRegistrar) RegisterUser(
	ctx context.Context, login, password string,
) error {
	const op = "UserRegistrar.RegisterUser"

	log := r.logger.WithOp(op)

	token, err := r.registerUser(ctx, login, password)
	if err != nil {
		return r.error(op, err)
	}

	err = r.syncStarter.StartSynchronization(ctx, token)
	if err != nil {
		log.Debug().Err(err).Msg("failed to start synchronization")
		return r.error(op, err)
	}

	return nil
}

func (r *UserRegistrar) registerUser(ctx context.Context, login, password string) (string, error) {
	const op = "UserRegistrar.registerUser"

	log := r.logger.WithOp(op)

	authClient, err := r.clientConn.ConnClient()
	if err != nil {
		log.Debug().Err(err).Msg("failed create auth client")
		return "", err
	}
	defer r.clientConn.Close()

	ctxWT, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reqData := &authbp.RegUserRequest{
		Login:    login,
		Password: []byte(password),
	}

	resData, err := authClient.RegisterUser(ctxWT, reqData)
	if err := r.handleRegistrationErr(err); err != nil {
		return "", err
	}

	return resData.Token, nil

}

func (r *UserRegistrar) handleRegistrationErr(err error) error {
	if err == nil {
		return nil
	}

	const op = "UserRegistrar.handleRegistrationErr"
	log := r.logger.WithOp(op)

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

func (r *UserRegistrar) error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

type UserAuthorizer struct {
	logger      logger.Logger
	clientConn  ClientConnector
	syncStarter SyncStarter
}

func NewUserAuthorizer(
	logger logger.Logger, clientConn ClientConnector, syncStarter SyncStarter,
) *UserAuthorizer {
	return &UserAuthorizer{logger, clientConn, syncStarter}
}

func (a *UserAuthorizer) AuthorizeUser(
	ctx context.Context, login, password string,
) error {
	const op = "AuthService.AuthorizeUser"

	log := a.logger.WithOp(op)

	token, err := a.authorizeUser(ctx, login, password)
	if err != nil {
		return a.error(op, err)
	}

	err = a.syncStarter.StartSynchronization(ctx, token)
	if err != nil {
		log.Debug().Err(err).Msg("failed to start synchronization")
		return a.error(op, err)
	}

	return nil
}

func (a *UserAuthorizer) authorizeUser(
	ctx context.Context, login, password string,
) (string, error) {
	const op = "AuthService.authorizeUser"

	log := a.logger.WithOp(op)

	authClient, err := a.clientConn.ConnClient()
	if err != nil {
		log.Debug().Err(err).Msg("failed create auth client")
		return "", err
	}
	defer a.clientConn.Close()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reqData := &authbp.AuthUserRequest{
		Login:    login,
		Password: []byte(password),
	}

	resData, err := authClient.AuthorizeUser(ctx, reqData)
	if err := a.handleAuthorizationErr(err); err != nil {
		return "", err
	}

	return resData.Token, nil
}

func (a *UserAuthorizer) handleAuthorizationErr(err error) error {
	if err == nil {
		return nil
	}

	const op = "UserAuthorizer.handleAuthorizationErr"
	log := a.logger.WithOp(op)

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

func (a *UserAuthorizer) error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
