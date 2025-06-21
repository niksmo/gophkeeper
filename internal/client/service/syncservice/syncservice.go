package syncservice

import (
	"context"

	"github.com/niksmo/gophkeeper/pkg/logger"
)

type SignupService struct {
	logger logger.Logger
}

func NewSignup(logger logger.Logger) *SignupService {
	return &SignupService{logger}
}

func (s *SignupService) Signup(
	ctx context.Context, login, password string,
) error {
	return nil
}

type SigninService struct {
	logger logger.Logger
}

func NewSignin(logger logger.Logger) *SigninService {
	return &SigninService{logger}
}

func (s *SigninService) Signin(
	ctx context.Context, login, password string,
) error {
	return nil
}

type LogoutService struct {
	logger logger.Logger
}

func NewLogout(logger logger.Logger) *LogoutService {
	return &LogoutService{logger}
}

func (s *LogoutService) Logout(ctx context.Context) error {
	return nil
}
