package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/niksmo/gophkeeper/internal/model"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type table int8

const (
	Passwords table = iota
	Cards
	Texts
	Binaries
)

func (t table) String() string {
	switch t {
	case Passwords:
		return "passwords"
	case Cards:
		return "cards"
	case Texts:
		return "texts"
	case Binaries:
		return "binaries"
	}
	return ""
}

type UsersDataRepository struct {
	logger logger.Logger
	db     Storage
}

func NewUsersDataRepository(l logger.Logger, s Storage) *UsersDataRepository {
	return &UsersDataRepository{l, s}
}

func (r *UsersDataRepository) GetComparable(
	ctx context.Context, t table, userID int,
) ([]model.SyncComparable, error) {
	const op = "UsersDataRepository.GetComparable"
	log := r.logger.WithOp(op)

	stmt := fmt.Sprintf(
		`SELECT id, name, updated_at FROM %s WHERE user_id=?;`, t,
	)

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var s []model.SyncComparable
	if rows.Next() {
		var o model.SyncComparable
		if err := o.ScanRow(rows); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s = append(s, o)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("failed get comparable data")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s, nil
}

func (r *UsersDataRepository) GetAll(
	ctx context.Context, t table, userID int,
) ([]model.SyncPayload, error) {
	const op = "UsersDataRepository.GetAll"
	log := r.logger.WithOp(op)

	stmt := fmt.Sprintf(`
		SELECT
			id, name, data, created_at, updated_at, deleted
		FROM %s
		WHERE user_id=?;`,
		t,
	)

	return r.querySlice(ctx, log, op, stmt, userID)
}

func (r *UsersDataRepository) GetSliceByIDs(
	ctx context.Context, t table, userID int, sID []int64,
) ([]model.SyncPayload, error) {
	const op = "UsersDataRepository.GetSliceByIDs"
	log := r.logger.WithOp(op)

	stmt := fmt.Sprintf(`
		SELECT id, name, data, created_at, updated_at, deleted, sync_id
		FROM %s
		WHERE user_id=? AND id IN (%s);`,
		t, r.makeStrIDList(sID),
	)

	return r.querySlice(ctx, log, op, stmt, userID)
}

func (r *UsersDataRepository) querySlice(
	ctx context.Context, log logger.Logger, op string, stmt string, userID int,
) ([]model.SyncPayload, error) {
	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	data := make([]model.SyncPayload, 0)

	for rows.Next() {
		var m model.SyncPayload

		if err := m.ScanRow(rows); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		data = append(data, m)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("failed to get users data rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (r *UsersDataRepository) UpdateSliceByIDs(
	ctx context.Context, t table, data []model.SyncPayload,
) error {
	const op = "UsersDataRepository.UpdateSliceByIDs"
	log := r.logger.WithOp(op)

	q := fmt.Sprintf(`
		UPDATE %s
		SET name=?, data=?, created_at=?, updated_at=?, deleted=?
		WHERE id=?;
		`, t,
	)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare stmt")
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	for i, o := range data {
		_, execErr := stmt.ExecContext(ctx, o.Name, o.Data,
			o.CreatedAt, o.UpdatedAt, o.Deleted, o.ID)
		if execErr != nil {
			log.Error().Err(err).Int("index", i).Msg("failed to exec update")
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

func (r *UsersDataRepository) InsertSlice(
	ctx context.Context, t table, userID int, data []model.SyncPayload,
) ([]int64, error) {
	const op = "UsersDataRepository.InsertSlice"
	log := r.logger.WithOp(op)

	q := fmt.Sprintf(`
		INSERT INTO %s
		  (user_id, name, data, created_at, updated_at, deleted)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id;`, t,
	)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare stmt")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var s []int64

	for i, o := range data {
		var id int64
		err := stmt.QueryRowContext(ctx, userID, o.Name, o.Data,
			o.CreatedAt, o.UpdatedAt, o.Deleted).Scan(&id)
		if err != nil {
			log.Error().Err(err).Int("index", i).Msg(
				"failed to insert row while iterate")
			return nil, tx.Rollback()
		}
		s = append(s, id)
	}
	return s, tx.Commit()
}

func (r *UsersDataRepository) makeStrIDList(sID []int64) string {
	var b strings.Builder
	lastIdx := len(sID) - 1
	for idx, id := range sID {
		b.WriteString(strconv.FormatInt(id, 10))

		if idx != lastIdx {
			b.WriteString(", ")
		}
	}
	return b.String()
}
