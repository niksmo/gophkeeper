package readservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	decoder interface {
		Decode(dst any, src []byte) error
	}

	decrypter interface {
		SetKey(string)
		Decrypt([]byte) ([]byte, error)
	}

	readRepo interface {
		ReadByID(ctx context.Context, id int) ([]byte, error)
	}
)

type ReadService[T any] struct {
	l         logger.Logger
	r         readRepo
	decoder   decoder
	decrypter decrypter
	dto       T
}

func New[T any](
	logger logger.Logger,
	repository readRepo,
	decoder decoder,
	decrypter decrypter,
) *ReadService[T] {
	return &ReadService[T]{
		l:         logger,
		r:         repository,
		decoder:   decoder,
		decrypter: decrypter,
	}
}

func (s *ReadService[T]) Read(
	ctx context.Context, key string, entryNum int,
) (T, error) {
	const op = "ReadService.Read"
	log := s.l.With().Str("op", op).Logger()

	data, err := s.r.ReadByID(ctx, entryNum)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			log.Debug().Msg("object not exists")
			return s.dto, service.ErrNotExists
		}
		log.Debug().Err(err).Msg("failed to read object from repository")
		return s.dto, fmt.Errorf("%s: %w", op, err)
	}

	s.decrypter.SetKey(key)
	b, err := s.decrypter.Decrypt(data)
	if err != nil {
		log.Debug().Err(err).Msg("failed to decrypt")
		return s.dto, service.ErrInvalidKey
	}

	if err := s.decoder.Decode(&s.dto, b); err != nil {
		log.Error().Err(err).Msg("failed to decode bytes to object")
		return s.dto, fmt.Errorf("%s: %w", op, err)
	}

	return s.dto, nil
}
