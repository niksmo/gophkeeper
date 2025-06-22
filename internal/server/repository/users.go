package repository

import (
	"context"

	"github.com/niksmo/gophkeeper/internal/server/dto"
	"github.com/niksmo/gophkeeper/internal/server/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var _ service.UserCreator = (*UsersRepository)(nil)
var _ service.UserProvider = (*UsersRepository)(nil)

type UsersRepository struct {
	logger logger.Logger
	db     Storage
}

func NewUsersRepository(logger logger.Logger, storage Storage) *UsersRepository {
	return &UsersRepository{logger, storage}
}

func (r *UsersRepository) Create(
	ctx context.Context, login, password string,
) (dto.User, error) {
	return dto.User{}, nil
}

func (r *UsersRepository) Read(ctx context.Context, id int) (dto.User, error) {
	return dto.User{}, nil
}
