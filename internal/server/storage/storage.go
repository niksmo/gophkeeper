package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type Storage struct {
	*sql.DB
	log logger.Logger
}

func New(logger logger.Logger, dsn string) *Storage {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open sql db")
	}

	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Str("dsn", dsn).Msg("failed to ping sql db")
	}
	logger.Info().Msg("database opens successfully")

	return &Storage{db, logger}
}
