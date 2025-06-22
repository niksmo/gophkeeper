package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(level string) Logger {
	setLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	return Logger{
		zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

func NewPretty(level string) Logger {
	setLevel(level)
	zl := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly},
	).With().Timestamp().Logger()
	return Logger{zl}
}

func setLevel(level string) {
	const op = "logger.setLevel"
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		fmt.Printf("%s: %s\n", op, "unknown level")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(lvl)
}
