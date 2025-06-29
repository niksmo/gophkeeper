package interceptors

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"google.golang.org/grpc"
)

func WithLog(logger logger.Logger) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(
		interceptorLogger(logger),
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	)
}

func interceptorLogger(logger logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		log := logger.With().Fields(fields).Logger()
		switch lvl {
		case logging.LevelDebug:
			log.Debug().Msg(msg)
		case logging.LevelInfo:
			log.Info().Msg(msg)
		case logging.LevelWarn:
			log.Warn().Msg(msg)
		case logging.LevelError:
			log.Error().Msg(msg)
		}
	})
}
