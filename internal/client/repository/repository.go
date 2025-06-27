package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/niksmo/gophkeeper/internal/client/model"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotExists     = errors.New("not exists")
)

const (
	passwords = "passwords"
	cards     = "cards"
	texts     = "texts"
	binaries  = "binaries"
)

type Storage interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Repository struct {
	log   logger.Logger
	db    Storage
	table string
}

func NewPwd(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, passwords}
}

func NewCard(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, cards}
}

func NewText(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, texts}
}

func NewBin(l logger.Logger, db Storage) *Repository {
	return &Repository{l, db, binaries}
}

func (r *Repository) Create(
	ctx context.Context, name string, data []byte,
) (int, error) {
	const op = "Repository.Create"
	log := r.log.With().Str("op", op).Logger()

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
	log := r.log.With().Str("op", op).Logger()

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
	log := r.log.With().Str("op", op).Logger()

	stmt := fmt.Sprintf(
		`
		SELECT id, name FROM %s WHERE deleted=FALSE
		ORDER BY name ASC, created_at ASC;
		`,
		r.table,
	)

	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

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
	log := r.log.With().Str("op", op).Logger()

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
		if IsSQLiteEniqueErr(err) {
			log.Debug().Err(err).Msg("object already exists")
			return fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		}
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("object is not exists")
			return fmt.Errorf("%s: %w", op, ErrNotExists)
		}
		log.Error().Err(err).Msg("failed to update")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, entryNum int) error {
	const op = "Repository.Delete"
	log := r.log.With().Str("op", op).Logger()

	stmt := fmt.Sprintf(`
	UPDATE %s SET
	  name=NULL, data=NULL, updated_at=?, deleted=TRUE
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

type SyncEntityRepository struct {
	logger logger.Logger
	db     Storage
	table  string
}

func NewPwdSync(l logger.Logger, s Storage) *SyncEntityRepository {
	return &SyncEntityRepository{l, s, passwords}
}

func NewCardSync(l logger.Logger, s Storage) *SyncEntityRepository {
	return &SyncEntityRepository{l, s, cards}
}

func NewTextSync(l logger.Logger, s Storage) *SyncEntityRepository {
	return &SyncEntityRepository{l, s, texts}
}

func NewBinSync(l logger.Logger, s Storage) *SyncEntityRepository {
	return &SyncEntityRepository{l, s, binaries}
}

func (r *SyncEntityRepository) GetComparable(
	ctx context.Context,
) ([]model.LocalComparable, error) {
	const op = "SyncEntityRepository.GetComparable"
	log := r.logger.WithOp(op)

	stmt := fmt.Sprintf(
		"SELECT id, name, updated_at, sync_id FROM %s;",
		r.table,
	)

	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	data := make([]model.LocalComparable, 0)

	for rows.Next() {
		var m model.LocalComparable

		if err := m.ScanRow(rows); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		data = append(data, m)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (r *SyncEntityRepository) GetAll(
	ctx context.Context,
) ([]model.LocalPayload, error) {
	const op = "SyncEntityRepository.GetAll"
	log := r.logger.WithOp(op)

	stmt := fmt.Sprintf(`
		SELECT id, name, data, created_at, updated_at, deleted, sync_id
		FROM %s;`,
		r.table,
	)

	return r.querySlice(ctx, log, op, stmt)
}

func (r *SyncEntityRepository) GetSliceByIDs(
	ctx context.Context, sID []int64,
) ([]model.LocalPayload, error) {
	const op = "SyncEntityRepository.GetSliceByIDs"
	log := r.logger.WithOp(op)

	stmt := fmt.Sprintf(`
		SELECT id, name, data, created_at, updated_at, deleted, sync_id
		FROM %s
		WHERE id IN (%s);`,
		r.table, r.makeStrIDList(sID),
	)

	return r.querySlice(ctx, log, op, stmt)
}

func (r *SyncEntityRepository) querySlice(
	ctx context.Context, log logger.Logger, op string, stmt string,
) ([]model.LocalPayload, error) {
	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	data := make([]model.LocalPayload, 0)

	for rows.Next() {
		var m model.LocalPayload

		if err := m.ScanRow(rows); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		data = append(data, m)
	}

	if err := rows.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("failed to select rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (r *SyncEntityRepository) makeStrIDList(sID []int64) string {
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

func (r *SyncEntityRepository) UpdateBySyncIDs(
	ctx context.Context, data []model.SyncPayload,
) error {
	const op = "SyncEntityRepository.UpdateBySyncIDs"
	log := r.logger.WithOp(op)

	q := fmt.Sprintf(`
		UPDATE %s
		SET name=?, data=?, created_at=?, updated_at=?, deleted=?
		WHERE sync_id=?;
		`,
		r.table,
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
			if IsSQLiteEniqueErr(err) {
				log.Error().Err(err).Int("index", i).Msg("unexpected name")
				continue
			}
			log.Error().Err(err).Int("index", i).Msg("failed to exec upsert")
			return tx.Rollback()
		}
	}
	return tx.Commit()

}

func (r *SyncEntityRepository) Insert(
	ctx context.Context, data []model.LocalPayload,
) error {
	const op = "SyncEntityRepository.Insert"
	log := r.logger.WithOp(op)

	q := fmt.Sprintf(`
		INSERT INTO %s
		  (name, data, created_at, updated_at, deleted, sync_id)
		VALUES (?, ?, ?, ?, ?, ?);
		`,
		r.table,
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
			o.CreatedAt, o.UpdatedAt, o.Deleted, o.SyncID)
		if execErr != nil {
			if IsSQLiteEniqueErr(err) {
				log.Error().Err(err).Int("index", i).Msg("unexpected name")
				continue
			}
			log.Error().Err(err).Int("index", i).Msg("failed to exec upsert")
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

func (r *SyncEntityRepository) UpdateSyncID(
	ctx context.Context, IDSyncIDPairs [][2]int,
) error {
	const op = "SyncEntityRepository.UpdateSyncID"
	log := r.logger.WithOp(op)

	q := fmt.Sprintf("UPDATE %s SET sync_id=? WHERE id=?;", r.table)

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

	for i, p := range IDSyncIDPairs {
		_, execErr := stmt.ExecContext(ctx, p[1], p[0])
		if execErr != nil {
			log.Error().Err(err).Int("index", i).Msg("failed to exec upsert")
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

func IsSQLiteEniqueErr(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) &&
		sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}
