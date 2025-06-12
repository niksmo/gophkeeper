package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/gophkeeper/internal/client"
	"github.com/niksmo/gophkeeper/internal/client/cmd"
	cmdpassword "github.com/niksmo/gophkeeper/internal/client/cmd/password"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

func main() {
	stopCtx, stopFn := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stopFn()
	l := logger.NewPretty("debug")
	cmdRoot := cmd.New()
	cmdRoot.AddCommand(
		cmdpassword.New(),
	)
	app := client.New(l, cmdRoot)
	app.Run(stopCtx)
}
