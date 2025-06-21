package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type SyncRepository struct {
	log logger.Logger
	db  Storage
}

func NewSync(l logger.Logger, db Storage) *SyncRepository {
	return &SyncRepository{l, db}
}

func (r *SyncRepository) Create(
	ctx context.Context, pid int, startedAt time.Time,
) (dto.SyncDTO, error) {
	const op = "SyncRepository.Create"

	stmt := `
	INSERT INTO synchronizations (pid, started_at)
	VALUES (?, ?)
	RETURNING id, pid, started_at;
	`
	var obj dto.SyncDTO
	err := r.db.QueryRowContext(ctx, stmt, pid, startedAt).Scan(
		&obj.ID, &obj.PID, &obj.StartedAt,
	)
	if err != nil {
		r.log.Debug().Str("op", op).Err(err).Msg("failed to create sync entry")
		return obj, fmt.Errorf("%s: %w", op, err)
	}

	return obj, nil
}

func (r *SyncRepository) ReadLast(ctx context.Context) (dto.SyncDTO, error) {
	const op = "SyncRepository.ReadLast"

	log := r.log.With().Str("op", op).Logger()

	stmt := `
	SELECT id, pid, started_at, stopped_at
	FROM synchronizations
	ORDER BY id DESC
	LIMIT 1;
	`
	var obj dto.SyncDTO
	err := r.db.QueryRowContext(ctx, stmt).Scan(
		&obj.ID, &obj.PID, &obj.StartedAt, &obj.StoppedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		log.Debug().Err(err).Msg("object is not exists")
		return obj, fmt.Errorf("%s: %w", op, ErrNotExists)
	}

	if err != nil {
		log.Debug().Err(err).Msg("failed to read sync entry")
		return obj, fmt.Errorf("%s: %w", op, err)
	}

	return obj, nil
}

func (r *SyncRepository) Update(ctx context.Context, dto dto.SyncDTO) error {
	const op = "SyncRepository.Update"

	log := r.log.With().Str("op", op).Logger()

	stmt := `
	UPDATE synchronizations
	SET pid=?, started_at=?, stopped_at=?
	WHERE id=?
	RETURNING id;
	`

	var id int
	err := r.db.QueryRowContext(
		ctx, stmt, dto.PID, dto.StartedAt, *dto.StoppedAt, dto.ID,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Debug().Err(err).Msg("object is not exists")
		return fmt.Errorf("%s: %w", op, ErrNotExists)
	}

	if err != nil {
		log.Debug().Err(err).Msg("failed to update")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
