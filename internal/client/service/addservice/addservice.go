package addservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	encoder interface {
		Encode(src any) ([]byte, error)
	}

	encrypter interface {
		SetKey(string)
		Encrypt([]byte) []byte
	}

	addRepo interface {
		Add(ctx context.Context, name string, data []byte) (int, error)
	}
)

type AddService[T any] struct {
	l         logger.Logger
	r         addRepo
	encoder   encoder
	encrypter encrypter
}

func New[T any](
	logger logger.Logger,
	repository addRepo,
	encoder encoder,
	encrypter encrypter,
) *AddService[T] {
	return &AddService[T]{
		l:         logger,
		r:         repository,
		encoder:   encoder,
		encrypter: encrypter,
	}
}

func (s *AddService[T]) Add(
	ctx context.Context, key, name string, dto T,
) (int, error) {
	const op = "AddService.Add"
	log := s.l.With().Str("op", op).Logger()

	b, err := s.encoder.Encode(dto)
	if err != nil {
		log.Debug().Err(err).Msg("failed to encode object to bytes")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	s.encrypter.SetKey(key)
	data := s.encrypter.Encrypt(b)

	id, err := s.r.Add(ctx, name, data)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			log.Debug().Str("name", name).Msg("object already exists")
			return 0, service.ErrAlreadyExists
		}
		log.Debug().Err(err).Msg("failed to add object to repository")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
