package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotExists     = errors.New("not exists")
)

type (
	Storage interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	}

	Repository struct {
		l     logger.Logger
		db    Storage
		table string
	}
)

func NewPwdRepository(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, "passwords"}
}

func NewCardRepository(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, "cards"}
}

func NewTextRepository(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, "texts"}
}

func NewBinRepository(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, "binaries"}
}

func (r *Repository) Add(
	ctx context.Context, name string, data []byte,
) (int, error) {
	const op = "Repository.Add"
	log := r.l.With().Str("op", op).Logger()

	stmt := fmt.Sprintf(`
	INSERT INTO %s (name, data, created_at, updated_at)
	VALUES (?, ?, ?, ?) RETURNING id;`,
		r.table,
	)

	var id int
	t := time.Now()
	err := r.db.QueryRowContext(ctx, stmt, name, data, t, t).Scan(&id)
	if err != nil {
		if IsSQLiteEniqueErr(err) {
			log.Debug().Err(err).Msg("object already exists")
			return 0, fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		}
		log.Error().Err(err).Msg("failed to insert")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *Repository) ReadByID(ctx context.Context, id int) ([]byte, error) {
	const op = "Repository.ReadByID"
	log := r.l.With().Str("op", op).Logger()

	stmt := fmt.Sprintf(
		`SELECT data FROM %s WHERE id=? AND deleted=FALSE;`, r.table,
	)

	var data []byte
	err := r.db.QueryRowContext(ctx, stmt, id).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("object is not exists")
			return nil, fmt.Errorf("%s: %w", op, ErrNotExists)
		}

		log.Error().Err(err).Msg("failed to select row")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (r *Repository) ListNames(ctx context.Context) ([][2]string, error) {
	const op = "Repository.ListNames"
	log := r.l.With().Str("op", op).Logger()

	stmt := fmt.Sprintf(
		`SELECT id, name FROM %s WHERE deleted=FALSE ORDER BY name ASC;`,
		r.table,
	)

	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data := make([][2]string, 0)
	var id int
	var name string
	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		item := [2]string{strconv.Itoa(id), name}
		data = append(data, item)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("failed while iterate rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (r *Repository) Update(
	ctx context.Context, entryNum int, name string, data []byte,
) error {
	const op = "Repository.Update"
	log := r.l.With().Str("op", op).Logger()

	stmt := fmt.Sprintf(`
	UPDATE %s SET
	  name=?, data=?, updated_at=?
	WHERE id=? RETURNING id;`,
		r.table,
	)

	var id int
	err := r.db.QueryRowContext(
		ctx, stmt, name, data, time.Now(), entryNum,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("object is not exists")
			return fmt.Errorf("%s: %w", op, ErrNotExists)
		}
		log.Error().Err(err).Msg("failed to update")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r Repository) Delete(ctx context.Context, entryNum int) error {
	const op = "Repository.Delete"
	log := r.l.With().Str("op", op).Logger()

	stmt := fmt.Sprintf(`
	UPDATE %s SET
	  name='', data=NULL, updated_at=?, deleted=TRUE
	WHERE id=? RETURNING id;`,
		r.table,
	)

	var id int
	err := r.db.QueryRowContext(ctx, stmt, time.Now(), entryNum).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("table is not exists")
			return fmt.Errorf("%s: %w", op, ErrNotExists)
		}
		log.Error().Err(err).Msg("failed to delete")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func IsSQLiteEniqueErr(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) &&
		sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}
