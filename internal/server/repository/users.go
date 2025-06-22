package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/niksmo/gophkeeper/internal/server/dto"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

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
	const op = "UsersRepository.Create"

	log := r.logger.WithOp(op)

	stmt := `
	INSERT INTO users (login, password, created_at)
	VALUES (?, ?, ?)
	RETURNING id, login, password, created_at, disabled;
	`

	var obj dto.User
	err := r.db.QueryRowContext(
		ctx, stmt, login, password, time.Now(),
	).Scan(
		&obj.ID, &obj.Login, &obj.Password, &obj.CreatedAt, &obj.Disabled,
	)

	if err != nil {
		if r.uniqueConstraintErr(err) {
			log.Debug().Str("login", login).Msg("login already exists")
			return dto.User{}, fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		}
		log.Error().Err(err).Msg("failed to create user")
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return obj, nil
}

func (r *UsersRepository) Read(ctx context.Context, login string) (dto.User, error) {
	const op = "UsersRepository.Read"

	log := r.logger.WithOp(op)

	stmt := `
	SELECT id, login, password, created_at, disabled
	FROM users
	WHERE login=?;
	`

	var obj dto.User
	err := r.db.QueryRowContext(ctx, stmt, login).Scan(
		&obj.ID, &obj.Login, &obj.Password, &obj.CreatedAt, &obj.Disabled,
	)
	if err != nil {
		if r.noRowErr(err) {
			log.Debug().Str("userLogin", login).Msg("user not exists")
			return dto.User{}, fmt.Errorf("%s: %w", op, ErrNotExists)
		}
		log.Error().Err(err).Msg("failed to read user")
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return obj, nil
}

func (r *UsersRepository) uniqueConstraintErr(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) &&
		sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}

func (r *UsersRepository) noRowErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
