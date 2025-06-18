package client

import (
	"context"
	"os"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler/pwdhandler"
	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service/addservice"
	"github.com/niksmo/gophkeeper/internal/client/service/editservice"
	"github.com/niksmo/gophkeeper/internal/client/service/listservice"
	"github.com/niksmo/gophkeeper/internal/client/service/readservice"
	"github.com/niksmo/gophkeeper/internal/client/service/removeservice.go"
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

	pwdRepository := repository.NewPwdRepository(logger, storage)

	pwdAddService := addservice.New[dto.PWD](
		logger, pwdRepository, encoder, encrypter,
	)
	pwdAddHandler := pwdhandler.NewAddHandler(
		logger, pwdAddService, os.Stdout,
	)

	pwdReadService := readservice.New[dto.PWD](
		logger, pwdRepository, decoder, decrypter,
	)
	pwdReadHandler := pwdhandler.NewReadHandler(
		logger, pwdReadService, os.Stdout,
	)

	pwdListService := listservice.New(logger, pwdRepository)
	pwdListHandler := pwdhandler.NewListHandler(
		logger, pwdListService, os.Stdout,
	)

	pwdEditService := editservice.New[dto.PWD](
		logger, pwdRepository, encoder, encrypter,
	)
	pwdEditHandler := pwdhandler.NewEditHandler(
		logger, pwdEditService, os.Stdout,
	)

	pwdRemoveService := removeservice.New(logger, pwdRepository)
	pwdRemoveHandler := pwdhandler.NewRemoveHandler(
		logger, pwdRemoveService, os.Stdout,
	)

	pwdAddCommand := pwdcommand.NewPwdAddCommand(pwdAddHandler)
	pwdReadCommand := pwdcommand.NewPwdReadCommand(pwdReadHandler)
	pwdListCommand := pwdcommand.NewPwdListCommand(pwdListHandler)
	pwdEditCommand := pwdcommand.NewPwdEditCommand(pwdEditHandler)
	pwdRemoveCommand := pwdcommand.NewPwdRemoveCommand(pwdRemoveHandler)

	pwdCommand := pwdcommand.NewPwdCommand()
	pwdCommand.AddCommand(
		pwdAddCommand,
		pwdReadCommand,
		pwdListCommand,
		pwdEditCommand,
		pwdRemoveCommand,
	)

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
