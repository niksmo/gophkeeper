package syncservice

import (
	"context"

	"github.com/niksmo/gophkeeper/pkg/logger"
)

type LocalRepo interface {
}

type ServerClient interface {
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
