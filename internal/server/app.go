package server

import (
	"net"

	"github.com/niksmo/gophkeeper/internal/server/api"
	"github.com/niksmo/gophkeeper/internal/server/config"
	"github.com/niksmo/gophkeeper/internal/server/repository"
	"github.com/niksmo/gophkeeper/internal/server/service"
	"github.com/niksmo/gophkeeper/internal/server/storage"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	logger     logger.Logger
	config     *config.Config
	gRPCServer *grpc.Server
	storage    *storage.Storage
}

func New(config *config.Config) *App {
	app := &App{config: config}
	app.initLogger()
	app.initStorage()
	app.initGRPCServer()
	app.registerServices()
	return app
}

func (app *App) initLogger() {
	app.logger = logger.New(app.config.LogLevel)
	app.logger.Info().Str("init", "logger").Str("level", app.config.LogLevel).Send()
}

func (app *App) initStorage() {
	app.storage = storage.New(app.logger, app.config.DSN)
	app.logger.Info().Str("init", "storage").Str("dns", app.config.DSN).Send()
}

func (app *App) initGRPCServer() {
	app.gRPCServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
	)
	app.logger.Info().Str(
		"init", "gRPCServer",
	).Str("addr", app.config.TCPAddr.String()).Send()
}

func (app *App) registerServices() {
	usersR := repository.NewUsersRepository(app.logger, app.storage)
	authS := service.NewAuthService(app.logger, usersR, usersR)
	api.RegisterAuthAPI(app.logger, app.gRPCServer, authS)
	app.logger.Info().Str("register", "AuthService").Send()
}

func (app *App) MustRun() {
	const op = "App.MustRun"

	log := app.logger.With().Str("op", op).Logger()

	lis, err := net.ListenTCP("tcp", app.config.TCPAddr)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen tcp address")
		app.Stop()
	}

	log.Info().Msg("gRPC server started")

	err = app.gRPCServer.Serve(lis)
	if err != nil {
		log.Error().Err(err).Msg("gRPC server serve error")
		app.Stop()
	}
}

func (app *App) Stop() {
	const op = "App.Stop"
	log := app.logger.With().Str("op", op).Logger()

	log.Info().Msg("stopping application")

	log.Info().Msg("gRPC server closing")
	app.gRPCServer.GracefulStop()
	log.Info().Msg("gRPC server closed successfully")

	log.Info().Msg("storage closing")
	if err := app.storage.Close(); err != nil {
		log.Info().Err(err).Msg("storage closed with errors")
	} else {
		log.Info().Msg("storage closed successfully")
	}
	log.Info().Msg("application stopped")
}
