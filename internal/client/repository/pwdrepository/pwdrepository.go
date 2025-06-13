package pwdrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	storage interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	}

	PwdRepository struct {
		l  logger.Logger
		db storage
	}
)

func New(l logger.Logger, db storage) *PwdRepository {
	return &PwdRepository{l, db}
}

func (r *PwdRepository) Add(
	ctx context.Context, name string, data []byte,
) (int, error) {
	const op = "pwdrepository.Add"
	log := r.l.With().Str("op", op).Logger()

	stmt := `
	INSERT INTO passwords (name, data, created_at, updated_at)
	VALUES (?, ?, ?, ?) RETURNING id;
	`

	var id int
	t := time.Now()
	err := r.db.QueryRowContext(ctx, stmt, name, data, t, t).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
