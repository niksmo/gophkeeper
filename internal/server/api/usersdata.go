package api

import (
	"github.com/niksmo/gophkeeper/pkg/logger"
	usersdatapb "github.com/niksmo/gophkeeper/proto/usersdata"
	"google.golang.org/grpc"
)

type UsersDataSyncService interface {
	// TODO: add methods
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

// func (h *usersDataSyncHandler) GetAll(
// 	ctx context.Context, in *usersdatapb.GetAllRequest,
// ) (*usersdatapb.GetAllResponse, error) {
// 	return nil, nil
// }

// func (h *usersDataSyncHandler) GetComparable(
// 	ctx context.Context, in *usersdatapb.GetComparableRequest,
// ) (*usersdatapb.GetComparableResponse, error) {
// 	return nil, nil

// }

// func (h *usersDataSyncHandler) GetSlice(
// 	ctx context.Context, in *usersdatapb.GetSliceRequest,
// ) (*usersdatapb.GetSliceResponse, error) {
// 	return nil, nil

// }

// func (h *usersDataSyncHandler) InsertSlice(
// 	ctx context.Context, in *usersdatapb.InsertSliceRequest,
// ) (*usersdatapb.InsertSliceResponse, error) {
// 	return nil, nil

// }

// func (h *usersDataSyncHandler) UpdateSlice(
// 	ctx context.Context, in *usersdatapb.UpdateSliceRequest,
// ) (*usersdatapb.UpdateSliceResponse, error) {
// 	return nil, nil

// }
