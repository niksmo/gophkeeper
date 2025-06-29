package server

import (
	"net"

	"github.com/niksmo/gophkeeper/internal/server/api"
	"github.com/niksmo/gophkeeper/internal/server/config"
	"github.com/niksmo/gophkeeper/internal/server/interceptors"
	"github.com/niksmo/gophkeeper/internal/server/repository"
	"github.com/niksmo/gophkeeper/internal/server/service/authservice"
	"github.com/niksmo/gophkeeper/internal/server/service/tokenservice"
	"github.com/niksmo/gophkeeper/internal/server/service/usersdataservice"
	"github.com/niksmo/gophkeeper/internal/server/storage"
	"github.com/niksmo/gophkeeper/pkg/hasher"
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

func (a *App) initLogger() {
	a.logger = logger.New(a.config.LogLevel)
	a.logger.Info().Str("init", "logger").Str(
		"level", a.config.LogLevel,
	).Send()
}

func (a *App) initStorage() {
	a.storage = storage.New(a.logger, a.config.DSN)
	a.logger.Info().Str("init", "storage").Str(
		"dns", a.config.DSN,
	).Send()
}

func (a *App) initGRPCServer() {
	tokenVerifier := tokenservice.NewUsersTokenVerifier(
		a.logger, a.config.TokenSecret,
	)
	userIDInterceptor := interceptors.NewUseIDInterceptor(
		a.logger, tokenVerifier,
	)
	a.gRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.WithRecovery(a.logger),
		interceptors.WithLog(a.logger),
		interceptors.WithUser(userIDInterceptor),
	))
	a.logger.Info().Str("init", "gRPCServer").Str(
		"addr", a.config.TCPAddr.String(),
	).Send()
}

func (a *App) registerServices() {
	cryptoHasher := hasher.NewCryptoHasher(a.config.HashCost)
	userTP := tokenservice.NewUsersTokenProvider(
		a.logger, a.config.TokenSecret, a.config.TokenTTL,
	)
	usersR := repository.NewUsersRepository(a.logger, a.storage)
	authS := authservice.New(
		authservice.ServiceDeps{
			Logger:        a.logger,
			Hasher:        cryptoHasher,
			UserCreator:   usersR,
			UserProvider:  usersR,
			TokenProvider: userTP,
		},
	)
	api.RegisterAuthAPI(a.logger, a.gRPCServer, authS)
	a.logger.Info().Str("register", "AuthService").Send()

	usersDataR := repository.NewUsersDataRepository(a.logger, a.storage)
	usersDataS := usersdataservice.New(a.logger, usersDataR)

	api.RegisterUsersDataSyncAPI(a.logger, a.gRPCServer, usersDataS)
	a.logger.Info().Str("register", "UsersDataSynchronizationService").Send()
}

func (a *App) MustRun() {
	const op = "App.MustRun"

	log := a.logger.With().Str("op", op).Logger()

	lis, err := net.ListenTCP("tcp", a.config.TCPAddr)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen tcp address")
		a.Stop()
	}

	log.Info().Msg("gRPC server started")

	err = a.gRPCServer.Serve(lis)
	if err != nil {
		log.Error().Err(err).Msg("gRPC server serve error")
		a.Stop()
	}
}

func (a *App) Stop() {
	const op = "App.Stop"
	log := a.logger.With().Str("op", op).Logger()

	log.Info().Msg("stopping application")

	log.Info().Msg("gRPC server closing")
	a.gRPCServer.GracefulStop()
	log.Info().Msg("gRPC server closed successfully")

	log.Info().Msg("storage closing")
	if err := a.storage.Close(); err != nil {
		log.Info().Err(err).Msg("storage closed with errors")
	} else {
		log.Info().Msg("storage closed successfully")
	}
	log.Info().Msg("application stopped")
}
