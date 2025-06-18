package removeservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type deleteRepo interface {
	Delete(ctx context.Context, id int) error
}

type RemoveService struct {
	l logger.Logger
	r deleteRepo
}

func New(
	logger logger.Logger,
	repository deleteRepo,
) *RemoveService {
	return &RemoveService{
		l: logger, r: repository,
	}
}

func (s *RemoveService) Remove(ctx context.Context, entryNum int) error {
	const op = "RemoveService.Remove"
	log := s.l.With().Str("op", op).Logger()

	err := s.r.Delete(ctx, entryNum)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			log.Debug().Msg("object not exists")
			return service.ErrNotExists
		}
		log.Debug().Err(err).Msg("failed to remove object from repository")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
