package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/server/interceptors"
	"github.com/niksmo/gophkeeper/internal/server/service/usersdataservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	usrdatapb "github.com/niksmo/gophkeeper/proto/usersdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidEntity = status.Error(codes.InvalidArgument, "invalid entity")

type UsersDataService interface {
	GetComparable(ctx context.Context,
		userID int, entity string) ([]*usrdatapb.Comparable, error)

	GetAll(ctx context.Context,
		userID int, entity string) ([]*usrdatapb.Payload, error)

	GetSliceByIDs(ctx context.Context,
		userID int, entity string, IDs []int64) ([]*usrdatapb.Payload, error)

	UpdateSliceByIDs(ctx context.Context,
		entity string, data []*usrdatapb.Payload) error

	InsertSlice(ctx context.Context,
		userID int, entity string, data []*usrdatapb.Payload) ([]int64, error)
}

type usersDataSyncHandler struct {
	usrdatapb.UnimplementedUsersDataServer
	service UsersDataService
	logger  logger.Logger
}

func RegisterUsersDataSyncAPI(
	logger logger.Logger, gRPCServer *grpc.Server, service UsersDataService,
) {
	usrdatapb.RegisterUsersDataServer(
		gRPCServer, &usersDataSyncHandler{service: service, logger: logger},
	)
}

func (h *usersDataSyncHandler) GetComparable(
	ctx context.Context, in *usrdatapb.GetComparableRequest,
) (*usrdatapb.GetComparableResponse, error) {
	const op = "usersDataSyncHandler.GetComparable"
	log := h.logger.WithOp(op)

	userID, err := h.getUserID(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, ErrInternal
	}

	data, err := h.service.GetComparable(ctx, userID, in.Entity)
	if err != nil {
		if errors.Is(err, usersdataservice.ErrInvalidEntity) {
			log.Warn().Str("entity", in.Entity).Msg("invalid entity")
			return nil, ErrInvalidEntity
		}
		log.Error().Err(err).Msg("internal error")
		return nil, ErrInternal
	}

	return &usrdatapb.GetComparableResponse{Data: data}, nil
}

func (h *usersDataSyncHandler) GetAll(
	ctx context.Context, in *usrdatapb.GetAllRequest,
) (*usrdatapb.GetAllResponse, error) {
	const op = "usersDataSyncHandler.GetAll"
	log := h.logger.WithOp(op)

	userID, err := h.getUserID(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, ErrInternal
	}

	data, err := h.service.GetAll(ctx, userID, in.Entity)
	if err != nil {
		if errors.Is(err, usersdataservice.ErrInvalidEntity) {
			return nil, ErrInvalidEntity
		}
		log.Error().Err(err).Msg("internal error")
		return nil, ErrInternal
	}

	return &usrdatapb.GetAllResponse{Data: data}, nil
}

func (h *usersDataSyncHandler) GetSlice(
	ctx context.Context, in *usrdatapb.GetSliceRequest,
) (*usrdatapb.GetSliceResponse, error) {
	const op = "usersDataSyncHandler.GetSlice"
	log := h.logger.WithOp(op)

	userID, err := h.getUserID(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, ErrInternal
	}

	data, err := h.service.GetAll(ctx, userID, in.Entity)
	if err != nil {
		if errors.Is(err, usersdataservice.ErrInvalidEntity) {
			log.Warn().Str("entity", in.Entity).Msg("invalid entity")
			return nil, ErrInvalidEntity
		}
		log.Error().Err(err).Msg("internal error")
		return nil, ErrInternal
	}

	return &usrdatapb.GetSliceResponse{Data: data}, nil
}

func (h *usersDataSyncHandler) UpdateSlice(
	ctx context.Context, in *usrdatapb.UpdateSliceRequest,
) (*usrdatapb.UpdateSliceResponse, error) {
	const op = "usersDataSyncHandler.UpdateSlice"
	log := h.logger.WithOp(op)

	err := h.service.UpdateSliceByIDs(ctx, in.Entity, in.Data)
	if err != nil {
		if errors.Is(err, usersdataservice.ErrInvalidEntity) {
			log.Warn().Str("entity", in.Entity).Msg("invalid entity")
			return nil, ErrInvalidEntity
		}
		log.Error().Err(err).Msg("internal error")
		return nil, ErrInternal
	}

	return &usrdatapb.UpdateSliceResponse{Ok: true}, nil
}

func (h *usersDataSyncHandler) InsertSlice(
	ctx context.Context, in *usrdatapb.InsertSliceRequest,
) (*usrdatapb.InsertSliceResponse, error) {
	const op = "usersDataSyncHandler.InsertSlice"
	log := h.logger.WithOp(op)

	userID, err := h.getUserID(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, ErrInternal
	}

	IDs, err := h.service.InsertSlice(ctx, userID, in.Entity, in.Data)
	if err != nil {
		if errors.Is(err, usersdataservice.ErrInvalidEntity) {
			log.Warn().Str("entity", in.Entity).Msg("invalid entity")
			return nil, ErrInvalidEntity
		}
		log.Error().Err(err).Msg("internal error")
		return nil, ErrInternal
	}

	return &usrdatapb.InsertSliceResponse{IDs: IDs}, nil
}

func (h *usersDataSyncHandler) getUserID(ctx context.Context) (int, error) {
	const op = "usersDataSyncHandler.getUserID"
	userID, ok := ctx.Value(interceptors.UserIDKey).(interceptors.UserID)
	if !ok {
		return 0, fmt.Errorf(
			"%s: %w", op, errors.New("expected userID not provided"),
		)
	}
	return userID.Int(), nil
}
