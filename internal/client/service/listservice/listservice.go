package listservice

import (
	"context"
	"fmt"

	"github.com/niksmo/gophkeeper/pkg/logger"
)

type listRepo interface {
	ListNames(ctx context.Context) ([][2]string, error)
}

type ListService struct {
	l logger.Logger
	r listRepo
}

func New(
	logger logger.Logger,
	repository listRepo,
) *ListService {
	return &ListService{
		l: logger,
		r: repository,
	}
}

func (s *ListService) List(ctx context.Context) ([][2]string, error) {
	const op = "ListService.List"
	log := s.l.With().Str("op", op).Logger()

	idNameSlice, err := s.r.ListNames(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get list of names")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return idNameSlice, nil
}
