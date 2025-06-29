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

const (
	storagePathFlag   = "storage-path"
	migrationPathFlag = "migrations-path"
)

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

func main() {
	logger := NewMigrationLogger()
	storagePath, migrationsPath := getFlagsValues()
	validateFlags(logger, storagePath, migrationsPath)
	makeMigrations(logger, storagePath, migrationsPath)
}

func getFlagsValues() (storage, migrations string) {
	storagePath := pflag.StringP(storagePathFlag, "s", "", "")
	migrationsPath := pflag.StringP(migrationPathFlag, "m", "", "")
	pflag.Parse()
	return *storagePath, *migrationsPath
}

func validateFlags(logger *MigrationLogger, storagePath, migrationsPath string) {
	var errs []error

	if storagePath == "" {
		errs = append(errs, fmt.Errorf("--%s flag: required", storagePathFlag))
	}

	if migrationsPath == "" {
		errs = append(errs, fmt.Errorf("--%s flag: required", migrationPathFlag))
	}

	if len(errs) != 0 {
		logger.logger.Fatal().Err(errors.Join(errs...)).Send()
	}
}

func makeMigrations(
	logger *MigrationLogger, storagePath, migrationsPath string,
) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("sqlite3://%s?x-no-tx-wrap=true", storagePath),
	)
	if err != nil {
		logger.logger.Fatal().Err(err).Send()
	}

	m.Log = NewMigrationLogger()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.Log.Printf("no migrations to apply")
			return
		}
		logger.logger.Fatal().Err(err)
	}
	m.Log.Printf("migration applied\n")
}
