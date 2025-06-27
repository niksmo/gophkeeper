package syncservice

import (
	"context"

	"github.com/niksmo/gophkeeper/internal/client/model"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type LocalRepo interface {
	GetComparable(context.Context) ([]model.LocalComparable, error)
	GetPayloadAll(context.Context) ([]model.LocalPayload, error)
	GetPayloadSlice(ctx context.Context, sID []int64) ([]model.LocalPayload, error)
	InsertPayload(ctx context.Context, data []model.LocalPayload) error
	UpdateSyncID(ctx context.Context, IDSyncIDPairs [][2]int64) error
}

type ServerClient interface {
	GetComparable(context.Context) ([]model.SyncComparable, error)
	GetPayloadAll(context.Context) ([]model.SyncPayload, error)
	GetPayloadSlice(ctx context.Context, sID []int64) ([]model.SyncPayload, error)
	InsertPayload(ctx context.Context, data []model.SyncPayload) error
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
