package client

import (
	"context"
	"os"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/handler/pwdhandler"
	"github.com/niksmo/gophkeeper/internal/client/repository/pwdrepository"
	"github.com/niksmo/gophkeeper/internal/client/service/pwdservice"
	"github.com/niksmo/gophkeeper/internal/client/storage"
	"github.com/niksmo/gophkeeper/pkg/cipher"
	"github.com/niksmo/gophkeeper/pkg/encode"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type App struct {
	c *command.Command
	s *storage.Storage
}

func New(logLevel, dsn string) *App {
	logger := logger.NewPretty(logLevel)
	storage := storage.New(logger, dsn)
	encoder := encode.NewEncoder()
	decoder := encode.NewDecoder()
	encrypter := cipher.NewEncrypter()
	decrypter := cipher.NewDecrypter()

	pwdRepository := pwdrepository.New(logger, storage)
	pwdService := pwdservice.New(
		logger,
		pwdRepository,
		pwdRepository,
		pwdRepository,
		encoder,
		decoder,
		encrypter,
		decrypter,
	)
	pwdAddHandler := pwdhandler.NewAddHandler(
		logger, pwdService, os.Stdout,
	)
	pwdReadHandler := pwdhandler.NewReadHandler(
		logger, pwdService, os.Stdout,
	)
	pwdListHandler := pwdhandler.NewListHandler(
		logger, pwdService, os.Stdout,
	)
	pwdAddCommand := pwdcommand.NewPwdAddCommand(pwdAddHandler)
	pwdReadCommand := pwdcommand.NewPwdReadCommand(pwdReadHandler)
	pwdListCommand := pwdcommand.NewPwdListCommand(pwdListHandler)

	pwdCommand := pwdcommand.NewPwdCommand()
	pwdCommand.AddCommand(pwdAddCommand, pwdReadCommand, pwdListCommand)

	cmdRoot := command.NewRootCommand()
	cmdRoot.AddCommand(pwdCommand)

	return &App{c: cmdRoot, s: storage}
}

func (app *App) Run(ctx context.Context) {
	app.s.MustRun(ctx)
	if err := app.c.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
