package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/niksmo/gophkeeper/internal/client"
)

const (
	syncTick    = 5 * time.Second
	authTimeout = 10 * time.Second
)

var (
	logLevel   = "debug"
	dsn        = ".gophkeeper.db"
	serverAddr = "127.0.0.1:8000"
	version    = "N/A"
	buildDate  = "N/A"
)

func main() {
	stopCtx, stopFn := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stopFn()
	app := client.New(
		client.Opt{
			LogLevel:    logLevel,
			DSN:         dsn,
			ServerAddr:  serverAddr,
			Version:     version,
			BuildDate:   buildDate,
			SyncTick:    syncTick,
			AuthTimeout: authTimeout,
		},
	)
	app.Run(stopCtx)
}
