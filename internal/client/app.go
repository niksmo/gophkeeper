package client

import (
	"context"
	"os"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/bincommand"
	"github.com/niksmo/gophkeeper/internal/client/command/cardcommand"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/command/synccommand"
	"github.com/niksmo/gophkeeper/internal/client/command/textcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler/binhandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/cardhandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/pwdhandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/synchandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/texthandler"
	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service/addservice"
	"github.com/niksmo/gophkeeper/internal/client/service/editservice"
	"github.com/niksmo/gophkeeper/internal/client/service/listservice"
	"github.com/niksmo/gophkeeper/internal/client/service/readservice"
	"github.com/niksmo/gophkeeper/internal/client/service/removeservice"
	"github.com/niksmo/gophkeeper/internal/client/service/syncservice"
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
		a.getCardCommand(),
		a.getTextCommand(),
		a.getSyncCommand(),
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

func (a *App) getCardCommand() *command.Command {
	repo := repository.NewCard(a.log, a.storage)

	addS := addservice.New[dto.BankCard](a.log, repo, a.encoder, a.encrypter)
	addH := cardhandler.NewAdd(a.log, addS, os.Stdout)
	addC := cardcommand.NewAdd(addH)

	readS := readservice.New[dto.BankCard](a.log, repo, a.decoder, a.decrypter)
	readH := cardhandler.NewRead(a.log, readS, os.Stdout)
	readC := cardcommand.NewRead(readH)

	listS := listservice.New(a.log, repo)
	listH := cardhandler.NewList(a.log, listS, os.Stdout)
	listC := cardcommand.NewList(listH)

	editS := editservice.New[dto.BankCard](a.log, repo, a.encoder, a.encrypter)
	editH := cardhandler.NewEdit(a.log, editS, os.Stdout)
	editC := cardcommand.NewEdit(editH)

	removeS := removeservice.New(a.log, repo)
	removeH := cardhandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := cardcommand.NewRemove(removeH)

	cardC := cardcommand.New()
	cardC.AddCommand(addC, readC, listC, editC, removeC)
	return cardC
}

func (a *App) getTextCommand() *command.Command {
	repo := repository.NewText(a.log, a.storage)

	addS := addservice.New[dto.Text](a.log, repo, a.encoder, a.encrypter)
	addH := texthandler.NewAdd(a.log, addS, os.Stdout)
	addC := textcommand.NewAdd(addH)

	readS := readservice.New[dto.Text](a.log, repo, a.decoder, a.decrypter)
	readH := texthandler.NewRead(a.log, readS, os.Stdout)
	readC := textcommand.NewRead(readH)

	listS := listservice.New(a.log, repo)
	listH := texthandler.NewList(a.log, listS, os.Stdout)
	listC := textcommand.NewList(listH)

	editS := editservice.New[dto.Text](a.log, repo, a.encoder, a.encrypter)
	editH := texthandler.NewEdit(a.log, editS, os.Stdout)
	editC := textcommand.NewEdit(editH)

	removeS := removeservice.New(a.log, repo)
	removeH := texthandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := textcommand.NewRemove(removeH)

	textC := textcommand.New()
	textC.AddCommand(addC, readC, listC, editC, removeC)
	return textC
}

func (a *App) getSyncCommand() *command.Command {
	const syncTick = time.Second * 5

	syncRepo := repository.NewSync(a.log, a.storage)
	signupS := syncservice.NewSignup(a.log, syncRepo)
	signupH := synchandler.NewSignup(a.log, signupS, os.Stdout)
	signupC := synccommand.NewSignup(signupH)

	signinS := syncservice.NewSignin(a.log, syncRepo)
	signinH := synchandler.NewSignin(a.log, signinS, os.Stdout)
	signinC := synccommand.NewSignin(signinH)

	logoutS := syncservice.NewLogout(a.log, syncRepo)
	logoutH := synchandler.NewLogout(a.log, logoutS, os.Stdout)
	logoutC := synccommand.NewLogout(logoutH)

	syncHS := syncservice.New(a.log, syncRepo, syncTick)
	startC := synccommand.NewStart(syncHS)

	syncC := synccommand.New()
	syncC.AddCommand(signupC, signinC, logoutC, startC)
	return syncC
}
