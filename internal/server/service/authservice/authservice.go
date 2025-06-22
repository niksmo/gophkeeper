package authservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/server/dto"
	"github.com/niksmo/gophkeeper/internal/server/repository"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var (
	ErrAlreadyExists      = errors.New("the login is busy")
	ErrInvalidCredentials = errors.New("the login or password is incorrect")
)

type UserTokenProvider interface {
	GetTokenString(userID int) (string, error)
}

type UserCreator interface {
	Create(ctx context.Context, login, password string) (dto.User, error)
}

type UserProvider interface {
	Read(ctx context.Context, login string) (dto.User, error)
}

type AuthService struct {
	logger        logger.Logger
	userCreator   UserCreator
	userProvider  UserProvider
	tokenProvider UserTokenProvider
}

func NewAuthService(
	logger logger.Logger,
	userCreator UserCreator,
	userProvider UserProvider,
	tokenProvider UserTokenProvider,
) *AuthService {
	return &AuthService{
		logger, userCreator, userProvider, tokenProvider,
	}
}

func (s *AuthService) RegisterNewUser(
	ctx context.Context, login, password string,
) (string, error) {
	const op = "AuthService.RegisterNewUser"
	log := s.logger.WithOp(op)

	userObj, err := s.userCreator.Create(ctx, login, password)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			log.Debug().Str("login", login).Msg("already exists")
			return "", fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		}
		log.Error().Err(err).Msg("failed to create new user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := s.getUserToken(userObj.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *AuthService) AuthorizeUser(
	ctx context.Context, login, password string,
) (string, error) {
	const op = "AuthService.AuthorizeUser"
	log := s.logger.WithOp(op)

	userObj, err := s.userProvider.Read(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			log.Debug().Str("userLogin", login).Msg("not exists")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error().Err(err).Msg("failed to get users data")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := s.getUserToken(userObj.ID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *AuthService) getUserToken(userID int) (string, error) {
	const op = "AuthService.getUserToken"
	log := s.logger.WithOp(op)
	token, err := s.tokenProvider.GetTokenString(userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user token")
		return "", err
	}
	return token, err
}
