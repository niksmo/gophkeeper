package interceptors

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithRecovery(log logger.Logger) grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
		log.Error().Any("panic", p).Msg("panic recovered")
		return status.Errorf(codes.Internal, "internal error")
	}))
}
