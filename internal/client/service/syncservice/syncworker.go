package syncservice

import (
	"context"

	"github.com/niksmo/gophkeeper/internal/client/model"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type LocalRepo interface {
	GetComparable(context.Context) ([]model.LocalComparable, error)
	GetAll(context.Context) ([]model.LocalPayload, error)
	GetSliceByIDs(ctx context.Context, sID []int64) ([]model.LocalPayload, error)
	UpdateBySyncIDs(ctx context.Context, data []model.SyncPayload) error
	Insert(ctx context.Context, data []model.LocalPayload) error
	UpdateSyncID(ctx context.Context, IDSyncIDPairs [][2]int64) error
}

type ServerClient interface {
	GetComparable(context.Context) ([]model.SyncComparable, error)
	GetAll(context.Context) ([]model.SyncPayload, error)
	GetSliceByIDs(context.Context) ([]model.SyncPayload, error)
	UpdateByIDs(ctx context.Context, data []model.SyncPayload) error
	Insert(ctx context.Context, data []model.SyncPayload) ([]int64, error)
}

type Worker struct {
	logger  logger.Logger
	localR  LocalRepo
	serverR ServerClient
}

func NewWorker(l logger.Logger, clR LocalRepo, srvR ServerClient) *Worker {
	return &Worker{l, clR, srvR}
}

func (w *Worker) DoJob(ctx context.Context) {

}
