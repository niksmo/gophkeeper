package dataservice

import (
	"context"
	"errors"
	"strings"

	"github.com/niksmo/gophkeeper/internal/model"
	"github.com/niksmo/gophkeeper/internal/server/repository"
	"github.com/niksmo/gophkeeper/pkg/logger"
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
		ctx context.Context, t repository.Table, userID int, sID []int64,
	) ([]model.SyncPayload, error)

	UpdateSliceByIDs(
		ctx context.Context, t repository.Table, data []model.SyncPayload,
	) error

	InsertSlice(
		ctx context.Context, t repository.Table,
		userID int, data []model.SyncPayload,
	) ([]int64, error)
}

type DataSyncService struct {
	logger       logger.Logger
	dataProvider DataProvider
}

func NewSync(l logger.Logger, p DataProvider) *DataSyncService {
	return &DataSyncService{l, p}
}

func (s *DataSyncService) parseEntity(
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
