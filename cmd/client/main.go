package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/niksmo/gophkeeper/internal/client"
)

const (
	syncTick    = 5 * time.Minute
	authTimeout = 10 * time.Second
)

var (
	LogLevel   = "debug"
	DSN        = ".gophkeeper.db"
	ServerAddr = "127.0.0.1:8000"
	Version    = "N/A"
	BuildDate  = "N/A"
)

func main() {
	stopCtx, stopFn := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stopFn()

	opt := client.Opt{
		LogLevel:    LogLevel,
		DSN:         DSN,
		ServerAddr:  ServerAddr,
		Version:     Version,
		BuildDate:   BuildDate,
		SyncTick:    syncTick,
		AuthTimeout: authTimeout,
	}

	app := client.New(opt)
	app.Run(stopCtx)
}
