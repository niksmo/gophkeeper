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
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL 
	);

	CREATE TABLE IF NOT EXISTS passwords (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS cards (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS texts (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN DEFAULT FALSE
	);
	
	CREATE TABLE IF NOT EXISTS binaries (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	data BLOB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted BOOLEAN DEFAULT FALSE
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
