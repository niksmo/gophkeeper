package migrations

import (
	"context"
	"database/sql"
	"fmt"
)

var Seq = []func(Storage) error{
	init0,
}

type Storage interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func GetLastID(s Storage) (int, error) {
	const op = "migrations.GetLast"
	stmt := `
	SELECT id FROM migrations
	ORDER BY id DESC
	LIMIT 1;
	`

	r := s.QueryRowContext(context.Background(), stmt)
	var id int
	err := r.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
