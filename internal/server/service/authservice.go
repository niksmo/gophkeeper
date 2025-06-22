package service

import (
	"context"

	"github.com/niksmo/gophkeeper/internal/server/dto"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type UserCreator interface {
	Create(ctx context.Context, login, password string) (dto.User, error)
}

type UserProvider interface {
	Read(ctx context.Context, id int) (dto.User, error)
}

type AuthService struct {
	logger       logger.Logger
	userCreator  UserCreator
	userProvider UserProvider
}

func NewAuthService(
	logger logger.Logger, c UserCreator, p UserProvider,
) *AuthService {
	return &AuthService{logger, c, p}
}

func (s *AuthService) RegisterNewUser(
	ctx context.Context, login, password string,
) (string, error) {
	return "registeredUserToken", nil
}

func (s *AuthService) AuthorizeUser(
	ctx context.Context, login, password string,
) (token string, err error) {
	return "authorizedUserToken", nil
}
