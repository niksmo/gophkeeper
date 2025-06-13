package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/gophkeeper/internal/client"
)

const (
	logLevel = "info"
	dsn      = ".gophkeeper.db"
)

func main() {
	stopCtx, stopFn := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stopFn()
	app := client.New(logLevel, dsn)
	app.Run(stopCtx)
}
