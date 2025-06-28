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
	"github.com/niksmo/gophkeeper/internal/client/handler/authhandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/binhandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/cardhandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/pwdhandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/synchandler"
	"github.com/niksmo/gophkeeper/internal/client/handler/texthandler"
	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service/authservice"
	"github.com/niksmo/gophkeeper/internal/client/service/genservice"
	"github.com/niksmo/gophkeeper/internal/client/service/syncservice"
	"github.com/niksmo/gophkeeper/internal/client/storage"
	"github.com/niksmo/gophkeeper/pkg/cipher"
	"github.com/niksmo/gophkeeper/pkg/encode"
	"github.com/niksmo/gophkeeper/pkg/logger"
	authbp "github.com/niksmo/gophkeeper/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Opt struct {
	LogLevel    string
	DSN         string
	ServerAddr  string
	Version     string
	BuildDate   string
	SyncTick    time.Duration
	AuthTimeout time.Duration
}

type App struct {
	log         logger.Logger
	cmd         *command.Command
	storage     *storage.Storage
	encoder     *encode.Encoder
	decoder     *encode.Decoder
	encrypter   *cipher.Encrypter
	decrypter   *cipher.Decrypter
	serverAddr  string
	syncTick    time.Duration
	authTimeout time.Duration
	conn        *grpc.ClientConn
}

func New(opt Opt) *App {
	log := logger.NewPretty(opt.LogLevel)

	app := &App{
		log:         log,
		cmd:         command.NewRootCommand(opt.Version, opt.BuildDate),
		storage:     storage.New(log, opt.DSN),
		encoder:     encode.NewEncoder(),
		decoder:     encode.NewDecoder(),
		encrypter:   cipher.NewEncrypter(),
		decrypter:   cipher.NewDecrypter(),
		serverAddr:  opt.ServerAddr,
		syncTick:    opt.SyncTick,
		authTimeout: opt.AuthTimeout,
	}

	app.initGRPCConn()
	app.registerCommands()
	return app
}

func (a *App) Run(ctx context.Context) {
	a.storage.MustRun(ctx)
	if err := a.cmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func (a *App) initGRPCConn() {
	dialOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(a.serverAddr, dialOpt)
	if err != nil {
		a.log.Fatal().Err(err).Msg("failed to init gRPC conn")
	}
	a.conn = conn
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

	addS := genservice.NewAdd[dto.PWD](a.log, repo, a.encoder, a.encrypter)
	addH := pwdhandler.NewAdd(a.log, addS, os.Stdout)
	addC := pwdcommand.NewAdd(addH)

	readS := genservice.NewRead[dto.PWD](a.log, repo, a.decoder, a.decrypter)
	readH := pwdhandler.NewRead(a.log, readS, os.Stdout)
	readC := pwdcommand.NewRead(readH)

	listS := genservice.NewList(a.log, repo)
	listH := pwdhandler.NewList(a.log, listS, os.Stdout)
	listC := pwdcommand.NewList(listH)

	editS := genservice.NewEdit[dto.PWD](a.log, repo, a.encoder, a.encrypter)
	editH := pwdhandler.NewEdit(a.log, editS, os.Stdout)
	editC := pwdcommand.NewEdit(editH)

	removeS := genservice.NewRemove(a.log, repo)
	removeH := pwdhandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := pwdcommand.NewRemove(removeH)

	passwordC := pwdcommand.New()
	passwordC.AddCommand(addC, readC, listC, editC, removeC)
	return passwordC
}

func (a *App) getBinaryCommand() *command.Command {
	repo := repository.NewBin(a.log, a.storage)

	addS := genservice.NewAdd[dto.BIN](a.log, repo, a.encoder, a.encrypter)
	addH := binhandler.NewAdd(a.log, addS, os.Stdout)
	addC := bincommand.NewAdd(addH)

	readS := genservice.NewRead[dto.BIN](a.log, repo, a.decoder, a.decrypter)
	readH := binhandler.NewRead(a.log, readS, os.Stdout)
	readC := bincommand.NewRead(readH)

	listS := genservice.NewList(a.log, repo)
	listH := binhandler.NewList(a.log, listS, os.Stdout)
	listC := bincommand.NewList(listH)

	editS := genservice.NewEdit[dto.BIN](a.log, repo, a.encoder, a.encrypter)
	editH := binhandler.NewEdit(a.log, editS, os.Stdout)
	editC := bincommand.NewEdit(editH)

	removeS := genservice.NewRemove(a.log, repo)
	removeH := binhandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := bincommand.NewRemove(removeH)

	binaryC := bincommand.New()
	binaryC.AddCommand(addC, readC, listC, editC, removeC)
	return binaryC
}

func (a *App) getCardCommand() *command.Command {
	repo := repository.NewCard(a.log, a.storage)

	addS := genservice.NewAdd[dto.BankCard](a.log, repo, a.encoder, a.encrypter)
	addH := cardhandler.NewAdd(a.log, addS, os.Stdout)
	addC := cardcommand.NewAdd(addH)

	readS := genservice.NewRead[dto.BankCard](a.log, repo, a.decoder, a.decrypter)
	readH := cardhandler.NewRead(a.log, readS, os.Stdout)
	readC := cardcommand.NewRead(readH)

	listS := genservice.NewList(a.log, repo)
	listH := cardhandler.NewList(a.log, listS, os.Stdout)
	listC := cardcommand.NewList(listH)

	editS := genservice.NewEdit[dto.BankCard](a.log, repo, a.encoder, a.encrypter)
	editH := cardhandler.NewEdit(a.log, editS, os.Stdout)
	editC := cardcommand.NewEdit(editH)

	removeS := genservice.NewRemove(a.log, repo)
	removeH := cardhandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := cardcommand.NewRemove(removeH)

	cardC := cardcommand.New()
	cardC.AddCommand(addC, readC, listC, editC, removeC)
	return cardC
}

func (a *App) getTextCommand() *command.Command {
	repo := repository.NewText(a.log, a.storage)

	addS := genservice.NewAdd[dto.Text](a.log, repo, a.encoder, a.encrypter)
	addH := texthandler.NewAdd(a.log, addS, os.Stdout)
	addC := textcommand.NewAdd(addH)

	readS := genservice.NewRead[dto.Text](a.log, repo, a.decoder, a.decrypter)
	readH := texthandler.NewRead(a.log, readS, os.Stdout)
	readC := textcommand.NewRead(readH)

	listS := genservice.NewList(a.log, repo)
	listH := texthandler.NewList(a.log, listS, os.Stdout)
	listC := textcommand.NewList(listH)

	editS := genservice.NewEdit[dto.Text](a.log, repo, a.encoder, a.encrypter)
	editH := texthandler.NewEdit(a.log, editS, os.Stdout)
	editC := textcommand.NewEdit(editH)

	removeS := genservice.NewRemove(a.log, repo)
	removeH := texthandler.NewRemove(a.log, removeS, os.Stdout)
	removeC := textcommand.NewRemove(removeH)

	textC := textcommand.New()
	textC.AddCommand(addC, readC, listC, editC, removeC)
	return textC
}

func (a *App) getSyncCommand() *command.Command {

	authClient := authservice.NewGRPCAuthClient(
		a.log, authbp.NewAuthClient(a.conn), a.authTimeout,
	)
	syncRepo := repository.NewSync(a.log, a.storage)

	syncStarter := syncservice.NewSyncRunner(a.log, syncRepo)
	userRegistrar := authservice.NewUserRegistrar(a.log, authClient, syncStarter)
	userAuthorizer := authservice.NewUserAuthorizer(a.log, authClient, syncStarter)

	signupH := authhandler.NewSignup(a.log, userRegistrar, os.Stdout)
	signupC := synccommand.NewSignup(signupH)

	signinH := authhandler.NewSignin(a.log, userAuthorizer, os.Stdout)
	signinC := synccommand.NewSignin(signinH)

	syncCloser := syncservice.NewSyncCloser(a.log, syncRepo)
	logoutH := authhandler.NewLogout(a.log, syncCloser, os.Stdout)
	logoutC := synccommand.NewLogout(logoutH)

	syncRunner := syncservice.NewWorkerPool(a.log, syncRepo, nil, a.syncTick)
	startH := synchandler.NewStart(a.log, syncRunner, os.Stdout)
	startC := synccommand.NewStart(startH)

	syncC := synccommand.New()
	syncC.AddCommand(signupC, signinC, logoutC, startC)
	return syncC
}
