package genservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	updateRepo interface {
		Update(ctx context.Context, id int, name string, data []byte) error
	}
)

type EditService[T any] struct {
	l         logger.Logger
	r         updateRepo
	encoder   encoder
	encrypter encrypter
}

func NewEdit[T any](
	logger logger.Logger,
	repository updateRepo,
	encoder encoder,
	encrypter encrypter,
) *EditService[T] {
	return &EditService[T]{
		l:         logger,
		r:         repository,
		encoder:   encoder,
		encrypter: encrypter,
	}
}

func (s *EditService[T]) Edit(
	ctx context.Context, key string, entryNum int, name string, dto T,
) error {
	const op = "EditService.Update"
	log := s.l.With().Str("op", op).Logger()

	b, err := s.encoder.Encode(dto)
	if err != nil {
		log.Debug().Err(err).Msg("failed to encode object to bytes")
		return fmt.Errorf("%s: %w", op, err)
	}

	s.encrypter.SetKey(key)
	data := s.encrypter.Encrypt(b)

	if err = s.r.Update(ctx, entryNum, name, data); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			log.Debug().Str("name", name).Msg("object already exists")
			return service.ErrAlreadyExists
		}
		if errors.Is(err, repository.ErrNotExists) {
			log.Debug().Str("name", name).Msg("object not exists")
			return service.ErrNotExists
		}
		log.Debug().Err(err).Msg("failed to save updated object to repository")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
