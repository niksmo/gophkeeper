package repository

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotExists     = errors.New("not exist")
	ErrAlreadyExists = errors.New("already exists")
)

type Storage interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
