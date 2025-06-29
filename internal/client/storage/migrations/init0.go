package migrations

import (
	"context"
	"time"
)

func init0(ctx context.Context, s Storage) error {
	stmt := `
	BEGIN;
	CREATE TABLE IF NOT EXISTS migrations (
	id INTEGER PRIMARY KEY,
	name TEXT,
	created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS synchronizations (
	id INTEGER PRIMARY KEY,
	pid INTEGER NOT NULL,
	started_at TIMESTAMP NOT NULL,
	stopped_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS passwords (
	id INTEGER PRIMARY KEY,
	name TEXT,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	sync_id INTEGER UNIQUE
	);

	CREATE TABLE IF NOT EXISTS cards (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	sync_id INTEGER UNIQUE
	);

	CREATE TABLE IF NOT EXISTS texts (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	sync_id INTEGER UNIQUE
	);
	
	CREATE TABLE IF NOT EXISTS binaries (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	sync_id INTEGER UNIQUE
	);

	INSERT INTO migrations (name, created_at) VALUES (?, ?);
	COMMIT;
	`
	_, err := s.ExecContext(ctx, stmt, "init0", time.Now())
	if err != nil {
		return err
	}

	return nil
}
