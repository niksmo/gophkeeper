package pwdrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	storage interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	}

	PwdRepository struct {
		l  logger.Logger
		db storage
	}

	listItem [2]string
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

func (r *PwdRepository) ReadByID(ctx context.Context, id int) ([]byte, error) {
	const op = "pwdrepository.ReadByID"
	log := r.l.With().Str("op", op).Logger()

	stmt := `SELECT data FROM passwords WHERE id=? AND deleted=FALSE;`

	var data []byte
	err := r.db.QueryRowContext(ctx, stmt, id).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("password is not exists")
			return nil, fmt.Errorf("%s: %w", op, repository.ErrNotExists)
		}

		log.Error().Err(err).Msg("failed to select row")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (r *PwdRepository) ListNames(ctx context.Context) ([]listItem, error) {
	const op = "pwdrepository.ListNames"
	log := r.l.With().Str("op", op).Logger()

	stmt := `
	SELECT id, name FROM passwords 
	WHERE deleted=FALSE
	ORDER BY name ASC;
	`

	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data := make([]listItem, 0)
	var id int
	var name string
	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		item := listItem{strconv.Itoa(id), name}
		data = append(data, item)
	}

	if rows.Err() != nil {
		log.Error().Err(err).Msg("failed to select row while iterate rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}
