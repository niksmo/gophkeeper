package pwdservice

import (
	"context"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	encoder interface {
		Encode(any) ([]byte, error)
	}

	keySetter interface {
		SetKey(string)
	}

	encrypter interface {
		keySetter
		Encrypt([]byte) []byte
	}

	decrypter interface {
		keySetter
		Decrypt([]byte) ([]byte, error)
	}

	pwdAddRepository interface {
		Add(ctx context.Context, name string, data []byte) (int, error)
	}

	pwdAddService struct {
		l         logger.Logger
		r         pwdAddRepository
		encoder   encoder
		encrypter encrypter
	}
)

func NewAddService(
	log logger.Logger,
	repo pwdAddRepository,
	encoder encoder,
	encrypter encrypter,
) *pwdAddService {
	return &pwdAddService{log, repo, encoder, encrypter}
}

func (s *pwdAddService) Add(ctx context.Context, key string, o objects.PWD) (int, error) {
	const op = "pwdservice.Add"
	log := s.l.With().Str("op", op).Logger()

	b, err := s.encoder.Encode(o)
	if err != nil {
		log.Error().Err(err).Msg("failed to encode object to bytes")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	s.encrypter.SetKey(key)
	data := s.encrypter.Encrypt(b)

	no, err := s.r.Add(ctx, o.Name, data)
	if err != nil {
		log.Error().Err(err).Msg("failed to add password to repository")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return no, nil
}
