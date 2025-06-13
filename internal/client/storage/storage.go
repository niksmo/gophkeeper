package storage

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/niksmo/gophkeeper/internal/client/storage/migrations"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type Storage struct {
	*sql.DB
	l logger.Logger
}

func New(l logger.Logger, dsn string) *Storage {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to open sql db")
	}

	if err := db.Ping(); err != nil {
		l.Fatal().Err(err).Str("dsn", dsn).Msg("failed to ping sql db")
	}
	l.Debug().Msg("database opens successfully")

	return &Storage{db, l}
}

func (s *Storage) MustRun(ctx context.Context) {
	s.migrate(ctx)
}

func (s *Storage) migrate(ctx context.Context) {
	id := s.lastMigrationID(ctx)
	s.makeMigrations(ctx, id)
}

func (s *Storage) lastMigrationID(ctx context.Context) int {
	const op = "storage.lastMigrationID"
	log := s.l.With().Str("op", op).Logger()
	id, err := migrations.GetLastID(ctx, s)
	if err != nil {
		log.Debug().Err(err).Send()
		return 0
	}
	log.Debug().Int("lastID", id).Send()
	return id
}

func (s *Storage) makeMigrations(ctx context.Context, lastID int) {
	const op = "storage.makeMigrations"
	for i, m := range migrations.Seq[lastID:] {
		log := s.l.With().Str("op", op).Int("migrationID", i).Logger()

		log.Debug().Msg("start migration")
		if err := m(ctx, s); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate")
		}
		log.Debug().Msg("complete migration")
	}
}
