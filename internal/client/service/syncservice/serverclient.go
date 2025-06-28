package syncservice

import (
	"context"

	"github.com/niksmo/gophkeeper/internal/model"
	"github.com/niksmo/gophkeeper/pkg/logger"
	usersdatapb "github.com/niksmo/gophkeeper/proto/usersdata"
)

type gRPCSyncClient struct {
	logger logger.Logger
	client usersdatapb.UsersDataSynchronizationClient
	entity string
	token  string
}

func NewGRPCSyncClientPwd(
	l logger.Logger, c usersdatapb.UsersDataSynchronizationClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "passwords"}
}

func NewGRPCSyncClientCard(
	l logger.Logger, c usersdatapb.UsersDataSynchronizationClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "cards"}
}

func NewGRPCSyncClientBin(
	l logger.Logger, c usersdatapb.UsersDataSynchronizationClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "binaries"}
}

func NewGRPCSyncClientText(
	l logger.Logger, c usersdatapb.UsersDataSynchronizationClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "texts"}
}

func (c *gRPCSyncClient) SetToken(token string) {
	c.token = token
}

func (c *gRPCSyncClient) GetComparable(
	ctx context.Context,
) ([]model.SyncComparable, error) {
	return nil, nil
}

func (c *gRPCSyncClient) GetAll(
	ctx context.Context,
) ([]model.SyncPayload, error) {
	return nil, nil
}

func (c *gRPCSyncClient) GetSliceByIDs(
	ctx context.Context, sID []int64,
) ([]model.SyncPayload, error) {
	return nil, nil
}

func (c *gRPCSyncClient) UpdateSliceByIDs(
	ctx context.Context, data []model.SyncPayload,
) error {
	return nil
}

func (c *gRPCSyncClient) InsertSlice(
	ctx context.Context, data []model.LocalPayload,
) ([]int64, error) {
	return nil, nil
}
