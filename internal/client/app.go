package client

import (
	"context"
	"os"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/bincommand"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler/binhandler"
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
	log := logger.NewPretty(logLevel)
	app := &App{
		log,
		command.NewRootCommand(),
		storage.New(log, dsn),
		encode.NewEncoder(),
		encode.NewDecoder(),
		cipher.NewEncrypter(),
		cipher.NewDecrypter(),
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
		a.getBinaryCommand(),
	)
}

func (a *App) getPasswordCommand() *command.Command {
	repo := repository.NewPwd(a.log, a.storage)

	addS := addservice.New[dto.PWD](a.log, repo, a.encoder, a.encrypter)
	addH := pwdhandler.NewAdd(a.log, addS, os.Stdout)
	addC := pwdcommand.NewAdd(addH)

	readS := readservice.New[dto.PWD](a.log, repo, a.decoder, a.decrypter)
	readH := pwdhandler.NewRead(a.log, readS, os.Stdout)
	readC := pwdcommand.NewRead(readH)

	listS := listservice.New(a.log, repo)
	listH := pwdhandler.NewList(a.log, listS, os.Stdout)
	listC := pwdcommand.NewList(listH)

	editS := editservice.New[dto.PWD](a.log, repo, a.encoder, a.encrypter)
	editH := pwdhandler.NewEdit(a.log, editS, os.Stdout)
	editC := pwdcommand.NewEdit(editH)

	removeS := removeservice.New(a.log, repo)
	removeH := pwdhandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := pwdcommand.NewRemove(removeH)

	passwordC := pwdcommand.New()
	passwordC.AddCommand(addC, readC, listC, editC, removeC)
	return passwordC
}

func (a *App) getBinaryCommand() *command.Command {
	repo := repository.NewBin(a.log, a.storage)

	addS := addservice.New[dto.BIN](a.log, repo, a.encoder, a.encrypter)
	addH := binhandler.NewAdd(a.log, addS, os.Stdout)
	addC := bincommand.NewAdd(addH)

	readS := readservice.New[dto.BIN](a.log, repo, a.decoder, a.decrypter)
	readH := binhandler.NewRead(a.log, readS, os.Stdout)
	readC := bincommand.NewRead(readH)

	listS := listservice.New(a.log, repo)
	listH := binhandler.NewList(a.log, listS, os.Stdout)
	listC := bincommand.NewList(listH)

	editS := editservice.New[dto.BIN](a.log, repo, a.encoder, a.encrypter)
	editH := binhandler.NewEdit(a.log, editS, os.Stdout)
	editC := bincommand.NewEdit(editH)

	removeS := removeservice.New(a.log, repo)
	removeH := binhandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := bincommand.NewRemove(removeH)

	binaryC := bincommand.New()
	binaryC.AddCommand(addC, readC, listC, editC, removeC)
	return binaryC
}
