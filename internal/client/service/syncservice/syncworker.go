package syncservice

import (
	"context"

	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type LocalRepo interface {
	GetComparable(context.Context) ([]dto.LocalComparable, error)
	GetPayloadAll(context.Context) ([]dto.LocalPayload, error)
	GetPayload(ctx context.Context, sID []int) ([]dto.LocalPayload, error)
	InsertPayload(ctx context.Context, data []dto.LocalPayload) error
	UpdateSyncID(ctx context.Context, IDSyncIDPairs [][2]int) error
}

type ServerClient interface {
	GetComparable(context.Context) ([]dto.ServerComparable, error)
	GetPayloadAll(context.Context) ([]dto.ServerPayload, error)
	GetPayload(ctx context.Context, sID []int) ([]dto.ServerPayload, error)
	InsertPayload(ctx context.Context, data []dto.ServerPayload) error
}

type Worker struct {
	logger  logger.Logger
	clientR LocalRepo
	serverR ServerClient
}

func NewWorker(l logger.Logger, clR LocalRepo, srvR ServerClient) *Worker {
	return &Worker{l, clR, srvR}
}

func (w *Worker) DoJob(ctx context.Context) {

}
