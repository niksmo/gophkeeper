package dataservice

import (
	"context"

	"github.com/niksmo/gophkeeper/internal/model"
	"github.com/niksmo/gophkeeper/internal/server/repository"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

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
