package usersdataservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/niksmo/gophkeeper/internal/model"
	"github.com/niksmo/gophkeeper/internal/server/repository"
	"github.com/niksmo/gophkeeper/pkg/logger"
	usrdatapb "github.com/niksmo/gophkeeper/proto/usersdata"
)

var ErrInvalidEntity = errors.New("invalid entity")

type DataProvider interface {
	GetComparable(
		ctx context.Context, t repository.Table, userID int,
	) ([]model.SyncComparable, error)

	GetAll(
		ctx context.Context, t repository.Table, userID int,
	) ([]model.SyncPayload, error)

	GetSliceByIDs(
		ctx context.Context, t repository.Table, userID int, IDs []int64,
	) ([]model.SyncPayload, error)

	UpdateSliceByIDs(
		ctx context.Context, t repository.Table, data []model.SyncPayload,
	) error

	InsertSlice(
		ctx context.Context, t repository.Table,
		userID int, data []model.SyncPayload,
	) ([]int64, error)
}

type UsersDataService struct {
	logger       logger.Logger
	dataProvider DataProvider
}

func New(l logger.Logger, p DataProvider) *UsersDataService {
	return &UsersDataService{l, p}
}

func (s *UsersDataService) GetComparable(ctx context.Context,
	userID int, entity string) ([]*usrdatapb.Comparable, error) {
	const op = "UsersDataService.GetComparable"
	log := s.logger.WithOp(op)

	table, err := s.parseEntity(entity)
	if err != nil {
		log.Warn().Err(err).Send()
		return nil, err
	}

	compData, err := s.dataProvider.GetComparable(ctx, table, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get comparable")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return s.comparableToPB(compData), nil
}

func (s *UsersDataService) GetAll(ctx context.Context,
	userID int, entity string) ([]*usrdatapb.Payload, error) {
	const op = "UsersDataService.GetAll"
	log := s.logger.WithOp(op)

	table, err := s.parseEntity(entity)
	if err != nil {
		log.Warn().Err(err).Send()
		return nil, err
	}

	payloadData, err := s.dataProvider.GetAll(ctx, table, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return s.payloadToPB(payloadData), nil
}

func (s *UsersDataService) GetSliceByIDs(ctx context.Context,
	userID int, entity string, IDs []int64) ([]*usrdatapb.Payload, error) {
	const op = "UsersDataService.GetSliceByIDs"
	log := s.logger.WithOp(op)

	table, err := s.parseEntity(entity)
	if err != nil {
		log.Warn().Err(err).Send()
		return nil, err
	}

	payloadData, err := s.dataProvider.GetSliceByIDs(ctx, table, userID, IDs)
	if err != nil {
		log.Error().Err(err).Msg("failed to get slice by IDs")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return s.payloadToPB(payloadData), nil
}

func (s *UsersDataService) UpdateSliceByIDs(ctx context.Context,
	entity string, data []*usrdatapb.Payload) error {
	const op = "UsersDataService.UpdateSliceByIDs"
	log := s.logger.WithOp(op)

	table, err := s.parseEntity(entity)
	if err != nil {
		log.Warn().Err(err).Send()
		return err
	}

	err = s.dataProvider.UpdateSliceByIDs(ctx, table, s.pbToPayload(data))
	if err != nil {
		log.Error().Err(err).Msg("failed to get update slice by IDs")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *UsersDataService) InsertSlice(ctx context.Context,
	userID int, entity string, data []*usrdatapb.Payload) ([]int64, error) {
	const op = "UsersDataService.InsertSlice"
	log := s.logger.WithOp(op)

	table, err := s.parseEntity(entity)
	if err != nil {
		log.Warn().Err(err).Send()
		return nil, err
	}

	IDs, err := s.dataProvider.InsertSlice(
		ctx, table, userID, s.pbToPayload(data),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert slice")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return IDs, nil
}

func (s *UsersDataService) parseEntity(
	entity string,
) (repository.Table, error) {
	switch {
	case strings.EqualFold(entity, repository.Passwords.String()):
		return repository.Passwords, nil
	case strings.EqualFold(entity, repository.Cards.String()):
		return repository.Cards, nil
	case strings.EqualFold(entity, repository.Binaries.String()):
		return repository.Binaries, nil
	case strings.EqualFold(entity, repository.Texts.String()):
		return repository.Texts, nil
	}
	return repository.Table(-1), ErrInvalidEntity
}

func (s *UsersDataService) comparableToPB(
	compData []model.SyncComparable,
) []*usrdatapb.Comparable {
	data := make([]*usrdatapb.Comparable, 0, len(compData))
	for _, o := range compData {
		pb := &usrdatapb.Comparable{
			ID:        o.ID,
			Name:      o.Name,
			UpdatedAt: o.UpdatedAt.UnixMilli(),
		}
		data = append(data, pb)
	}
	return data
}

func (s *UsersDataService) payloadToPB(
	payloadData []model.SyncPayload,
) []*usrdatapb.Payload {
	data := make([]*usrdatapb.Payload, 0, len(payloadData))
	for _, o := range payloadData {
		pb := &usrdatapb.Payload{
			ID:        o.ID,
			Name:      o.Name,
			Data:      o.Data,
			CreatedAt: o.CreatedAt.UnixMilli(),
			UpdatedAt: o.UpdatedAt.UnixMilli(),
			Deleted:   o.Deleted,
		}
		data = append(data, pb)
	}
	return data
}

func (s *UsersDataService) pbToPayload(
	pbData []*usrdatapb.Payload,
) []model.SyncPayload {
	data := make([]model.SyncPayload, 0, len(pbData))
	for _, o := range pbData {
		pb := model.SyncPayload{
			ID:        o.ID,
			Name:      o.Name,
			Data:      o.Data,
			CreatedAt: time.UnixMilli(o.CreatedAt),
			UpdatedAt: time.UnixMilli(o.UpdatedAt),
			Deleted:   o.Deleted,
		}
		data = append(data, pb)
	}
	return data
}
