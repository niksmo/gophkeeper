package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/gophkeeper/internal/client"
)

const (
	logLevel   = "error"
	dsn        = ".gophkeeper.db"
	serverAddr = "127.0.0.1:8000"
)

func main() {
	stopCtx, stopFn := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stopFn()
	app := client.New(logLevel, dsn, serverAddr)
	app.Run(stopCtx)
}
