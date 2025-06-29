package interceptors

import (
	"context"

	"github.com/niksmo/gophkeeper/pkg/logger"
	pb "github.com/niksmo/gophkeeper/proto/usersdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidToken = status.Error(codes.Unauthenticated, "invalid token")

type Interceptor interface {
	Intercept(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
}

func WithUser(i Interceptor) grpc.UnaryServerInterceptor {
	return i.Intercept
}

type key int8

const UserIDKey key = 0

type UserID int

func (id *UserID) Int() int {
	return int(*id)
}

type UsersTokenVerifier interface {
	Verify(token string) (int, error)
}

type UserIDInterceptor struct {
	log      logger.Logger
	verifier UsersTokenVerifier
}

func NewUseIDInterceptor(l logger.Logger, v UsersTokenVerifier) UserIDInterceptor {
	return UserIDInterceptor{l, v}
}

func (e UserIDInterceptor) Intercept(ctx context.Context,
	req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	const op = "UserIDInterceptor.Intercept"
	log := e.log.With().Str("op", op).Str("method", info.FullMethod).Logger()

	var token string
	switch r := req.(type) {
	default:
		return handler(ctx, req)
	case *pb.GetComparableRequest:
		token = r.Token
	case *pb.GetAllRequest:
		token = r.Token
	case *pb.GetSliceRequest:
		token = r.Token
	case *pb.UpdateSliceRequest:
		token = r.Token
	case *pb.InsertSliceRequest:
		token = r.Token
	}

	userID, err := e.getUserID(token)
	if err != nil {
		log.Warn().Err(err).Msg("invalid token")
		return nil, ErrInvalidToken
	}
	log.Debug().Int("userID", userID).Send()

	ctx = e.updateContext(ctx, userID)

	return handler(ctx, req)
}

func (e UserIDInterceptor) getUserID(token string) (int, error) {
	return e.verifier.Verify(token)
}

func (e UserIDInterceptor) updateContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}
