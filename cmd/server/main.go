package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/gophkeeper/internal/server"
	"github.com/niksmo/gophkeeper/internal/server/config"
)

func main() {
	ctxStop, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stop()

	app := server.New(config.MustLoad())

	go app.MustRun()

	<-ctxStop.Done()

	app.Stop()
}
