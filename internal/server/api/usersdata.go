package api

import (
	"context"

	"github.com/niksmo/gophkeeper/pkg/logger"
	usersdatapb "github.com/niksmo/gophkeeper/proto/usersdata"
	"google.golang.org/grpc"
)

type UsersDataSyncService interface {
}

type usersDataSyncHandler struct {
	usersdatapb.UnimplementedUsersDataSynchronizationServer
	service UsersDataSyncService
	logger  logger.Logger
}

func RegisterUsersDataSyncAPI(
	logger logger.Logger, gRPCServer *grpc.Server, service AuthService,
) {
	usersdatapb.RegisterUsersDataSynchronizationServer(
		gRPCServer, &usersDataSyncHandler{service: service, logger: logger},
	)
}

func (h *usersDataSyncHandler) GetAll(
	ctx context.Context, in *usersdatapb.GetAllRequest,
) (*usersdatapb.GetAllResponse, error) {

}

func (h *usersDataSyncHandler) GetComparable(
	ctx context.Context, in *usersdatapb.GetComparableRequest,
) (*usersdatapb.GetComparableResponse, error) {

}

func (h *usersDataSyncHandler) GetSlice(
	ctx context.Context, in *usersdatapb.GetSliceRequest,
) (*usersdatapb.GetSliceResponse, error) {

}

func (h *usersDataSyncHandler) InsertSlice(
	ctx context.Context, in *usersdatapb.InsertSliceRequest,
) (*usersdatapb.InsertSliceResponse, error) {

}

func (h *usersDataSyncHandler) UpdateSlice(
	ctx context.Context, in *usersdatapb.UpdateSliceRequest,
) (*usersdatapb.UpdateSliceResponse, error) {

}
