package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/niksmo/gophkeeper/internal/client"
)

const (
	logLevel   = "debug"
	dsn        = ".gophkeeper.db"
	serverAddr = "127.0.0.1:8000"
	syncTick   = 5 * time.Second
)

func main() {
	stopCtx, stopFn := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stopFn()
	app := client.New(logLevel, dsn, serverAddr, syncTick)
	app.Run(stopCtx)
}
