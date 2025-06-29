package syncservice

import (
	"context"
	"fmt"
	"time"

	"github.com/niksmo/gophkeeper/internal/model"
	"github.com/niksmo/gophkeeper/pkg/logger"
	usersdatapb "github.com/niksmo/gophkeeper/proto/usersdata"
)

type gRPCSyncClient struct {
	logger logger.Logger
	client usersdatapb.UsersDataClient
	entity string
	token  string
}

func NewGRPCSyncClientPwd(
	l logger.Logger, c usersdatapb.UsersDataClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "passwords"}
}

func NewGRPCSyncClientCard(
	l logger.Logger, c usersdatapb.UsersDataClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "cards"}
}

func NewGRPCSyncClientBin(
	l logger.Logger, c usersdatapb.UsersDataClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "binaries"}
}

func NewGRPCSyncClientText(
	l logger.Logger, c usersdatapb.UsersDataClient,
) ServerClient {
	return &gRPCSyncClient{logger: l, client: c, entity: "texts"}
}

func (c *gRPCSyncClient) SetToken(token string) {
	c.token = token
}

func (c *gRPCSyncClient) GetComparable(
	ctx context.Context,
) ([]model.SyncComparable, error) {
	const op = "gRPCSyncClient.GetComparable"
	log := c.logger.With().Str("op", op).Str("intity", c.entity).Logger()

	log.Debug().Str("token", c.token).Str("entity", c.entity).Msg("start request")

	req := &usersdatapb.GetComparableRequest{Token: c.token, Entity: c.entity}
	res, err := c.client.GetComparable(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get comparable objects")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c.pbToSyncComprable(res.Data), nil
}

func (c *gRPCSyncClient) GetAll(
	ctx context.Context,
) ([]model.SyncPayload, error) {
	const op = "gRPCSyncClient.GetAll"
	log := c.logger.With().Str("op", op).Str("intity", c.entity).Logger()
	req := &usersdatapb.GetAllRequest{Token: c.token, Entity: c.entity}
	res, err := c.client.GetAll(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all objects")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c.pbToSyncPayload(res.Data), nil
}

func (c *gRPCSyncClient) GetSliceByIDs(
	ctx context.Context, IDs []int64,
) ([]model.SyncPayload, error) {
	const op = "gRPCSyncClient.GetSliceByIDs"
	log := c.logger.With().Str("op", op).Str("intity", c.entity).Logger()
	req := &usersdatapb.GetSliceRequest{
		Token:  c.token,
		Entity: c.entity,
		IDs:    IDs,
	}

	res, err := c.client.GetSlice(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get slice of objects")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c.pbToSyncPayload(res.Data), nil
}

func (c *gRPCSyncClient) UpdateSliceByIDs(
	ctx context.Context, data []model.SyncPayload,
) error {
	const op = "gRPCSyncClient.UpdateSliceByIDs"
	log := c.logger.With().Str("op", op).Str("intity", c.entity).Logger()
	req := &usersdatapb.UpdateSliceRequest{
		Token:  c.token,
		Entity: c.entity,
		Data:   c.syncToPBPayload(data),
	}
	_, err := c.client.UpdateSlice(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to update slice of objects")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (c *gRPCSyncClient) InsertSlice(
	ctx context.Context, data []model.LocalPayload,
) ([]int64, error) {
	const op = "gRPCSyncClient.InsertSlice"
	log := c.logger.With().Str("op", op).Str("intity", c.entity).Logger()
	req := &usersdatapb.InsertSliceRequest{
		Token:  c.token,
		Entity: c.entity,
		Data:   c.localToPBPayload(data),
	}
	res, err := c.client.InsertSlice(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert sclice of objects")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res.IDs, nil
}

func (c *gRPCSyncClient) pbToSyncComprable(
	data []*usersdatapb.Comparable,
) []model.SyncComparable {
	s := make([]model.SyncComparable, 0, len(data))
	for _, o := range data {
		cvt := model.SyncComparable{
			ID:        o.ID,
			Name:      o.Name,
			UpdatedAt: time.UnixMilli(o.UpdatedAt),
		}
		s = append(s, cvt)
	}
	return s
}

func (c *gRPCSyncClient) pbToSyncPayload(
	data []*usersdatapb.Payload,
) []model.SyncPayload {
	s := make([]model.SyncPayload, 0, len(data))
	for _, o := range data {
		cvt := model.SyncPayload{
			ID:        o.ID,
			Name:      o.Name,
			Data:      o.Data,
			CreatedAt: time.UnixMilli(o.CreatedAt),
			UpdatedAt: time.UnixMilli(o.UpdatedAt),
			Deleted:   o.Deleted,
		}
		s = append(s, cvt)
	}
	return s
}

func (c *gRPCSyncClient) syncToPBPayload(
	data []model.SyncPayload,
) []*usersdatapb.Payload {
	s := make([]*usersdatapb.Payload, 0, len(data))
	for _, o := range data {
		cvt := &usersdatapb.Payload{
			ID:        o.ID,
			Name:      o.Name,
			Data:      o.Data,
			CreatedAt: o.CreatedAt.UnixMilli(),
			UpdatedAt: o.UpdatedAt.UnixMilli(),
			Deleted:   o.Deleted,
		}
		s = append(s, cvt)
	}
	return s
}

func (c *gRPCSyncClient) localToPBPayload(data []model.LocalPayload) []*usersdatapb.Payload {
	s := make([]*usersdatapb.Payload, 0, len(data))
	for _, o := range data {
		cvt := &usersdatapb.Payload{
			ID:        o.SyncID,
			Name:      o.Name,
			Data:      o.Data,
			CreatedAt: o.CreatedAt.UnixMilli(),
			UpdatedAt: o.UpdatedAt.UnixMilli(),
			Deleted:   o.Deleted,
		}
		s = append(s, cvt)
	}
	return s
}
