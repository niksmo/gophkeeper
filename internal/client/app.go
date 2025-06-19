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
	"github.com/niksmo/gophkeeper/internal/client/service/removeservice"
	"github.com/niksmo/gophkeeper/internal/client/storage"
	"github.com/niksmo/gophkeeper/pkg/cipher"
	"github.com/niksmo/gophkeeper/pkg/encode"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type App struct {
	log       logger.Logger
	cmd       *command.Command
	storage   *storage.Storage
	encoder   *encode.Encoder
	decoder   *encode.Decoder
	encrypter *cipher.Encrypter
	decrypter *cipher.Decrypter
}

func New(logLevel, dsn string) *App {
	logger := logger.NewPretty(logLevel)
	storage := storage.New(logger, dsn)
	encoder := encode.NewEncoder()
	decoder := encode.NewDecoder()
	encrypter := cipher.NewEncrypter()
	decrypter := cipher.NewDecrypter()

	cmd := command.NewRootCommand()

	app := &App{
		logger,
		cmd,
		storage,
		encoder, decoder,
		encrypter, decrypter,
	}

	app.registerCommands()

	return app
}

func (a *App) Run(ctx context.Context) {
	a.storage.MustRun(ctx)
	if err := a.cmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func (a *App) registerCommands() {
	a.cmd.AddCommand(
		a.getPasswordCommand(),
	)
}

func (a *App) getPasswordCommand() *command.Command {
	repo := repository.NewPwdRepository(a.log, a.storage)

	addS := addservice.New[dto.PWD](a.log, repo, a.encoder, a.encrypter)
	addH := pwdhandler.NewAddHandler(a.log, addS, os.Stdout)

	readS := readservice.New[dto.PWD](a.log, repo, a.decoder, a.decrypter)
	readH := pwdhandler.NewReadHandler(a.log, readS, os.Stdout)

	listS := listservice.New(a.log, repo)
	listH := pwdhandler.NewListHandler(a.log, listS, os.Stdout)

	editS := editservice.New[dto.PWD](a.log, repo, a.encoder, a.encrypter)
	editH := pwdhandler.NewEditHandler(a.log, editS, os.Stdout)

	removeS := removeservice.New(a.log, repo)
	removeH := pwdhandler.NewRemoveHandler(a.log, removeS, os.Stdout)

	addC := pwdcommand.NewPwdAddCommand(addH)
	readC := pwdcommand.NewPwdReadCommand(readH)
	listC := pwdcommand.NewPwdListCommand(listH)
	editC := pwdcommand.NewPwdEditCommand(editH)
	removeC := pwdcommand.NewPwdRemoveCommand(removeH)

	passwordC := pwdcommand.NewPwdCommand()
	passwordC.AddCommand(addC, readC, listC, editC, removeC)
	return passwordC
}
