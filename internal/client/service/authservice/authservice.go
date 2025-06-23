package authservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/niksmo/gophkeeper/pkg/logger"
	authbp "github.com/niksmo/gophkeeper/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAlreadyExists = errors.New("user already exists")
)

const (
	timeout = 5 * time.Second
)

type ClientConnector interface {
	ConnClient() (authbp.AuthClient, error)
	Close() error
}

type AuthService struct {
	logger     logger.Logger
	clientConn ClientConnector
}

func New(logger logger.Logger, clientConn ClientConnector) *AuthService {
	return &AuthService{logger, clientConn}
}

func (s *AuthService) RegisterNewUser(
	ctx context.Context, login, password string,
) (token string, err error) {
	const op = "AuthService.RegisterNewUser"

	log := s.logger.WithOp(op)

	authClient, err := s.clientConn.ConnClient()
	if err != nil {
		log.Debug().Err(err).Msg("failed to establish connection")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer s.clientConn.Close()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reqData := &authbp.RegUserRequest{
		Login:    login,
		Password: []byte(password),
	}

	resData, err := authClient.RegisterUser(ctx, reqData)
	if err != nil {
		switch status.Code(err) {
		case codes.AlreadyExists:
			log.Debug().Msg("user already exists")
			return "", fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		case codes.InvalidArgument:
			log.Debug().Err(err).Msg("invalid provided data")
			return "", fmt.Errorf("%s: %w", op, err)
		default:
			log.Debug().Err(err).Msg("failed to register new user")
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return resData.Token, nil
}

func (s *AuthService) AuthorizeUser(
	ctx context.Context, login, password string,
) (token string, err error) {
	const op = "AuthService.AuthorizeUser"

	log := s.logger.WithOp(op)

	authClient, err := s.clientConn.ConnClient()
	if err != nil {
		log.Error().Err(err).Msg("failed to establish connection")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer s.clientConn.Close()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reqData := &authbp.AuthUserRequest{
		Login:    login,
		Password: []byte(password),
	}

	resData, err := authClient.AuthorizeUser(ctx, reqData)
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			log.Debug().Err(err).Msg("invalid provided data")
			return "", fmt.Errorf("%s: %w", op, err)
		default:
			log.Error().Err(err).Msg("failed to register new user")
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return resData.Token, nil
}
