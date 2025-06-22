package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/spf13/pflag"
)

func main() {
	const (
		storagePathFlag   = "storage-path"
		migrationPathFlag = "migrations-path"
	)

	var errs []error

	storagePath := pflag.StringP(storagePathFlag, "s", "", "")
	migrationsPath := pflag.StringP(migrationPathFlag, "m", "", "")
	pflag.Parse()

	if *storagePath == "" {
		errs = append(errs, fmt.Errorf("--%s flag: required", storagePathFlag))
	}

	if *migrationsPath == "" {
		errs = append(errs, fmt.Errorf("--%s flag: required", migrationPathFlag))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", *migrationsPath),
		fmt.Sprintf("sqlite3://%s?x-no-tx-wrap=true", *storagePath),
	)
	if err != nil {
		panic(err)
	}

	m.Log = NewMigrationLogger()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.Log.Printf("no migrations to apply")
			return
		}
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}
	m.Log.Printf("migration applied\n")
}

type MigrationLogger struct {
	logger  logger.Logger
	verbose bool
}

func NewMigrationLogger() *MigrationLogger {

	return &MigrationLogger{logger.NewPretty("debug"), true}
}

func (ml *MigrationLogger) Printf(format string, v ...any) {
	ml.logger.Info().Msgf(format, v...)
}

func (ml *MigrationLogger) Verbose() bool {
	return ml.verbose
}
