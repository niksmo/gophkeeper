package pwdservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var (
	ErrPwdNotExists = errors.New("password is not exists")
	ErrInvalidKey   = errors.New("invalid key provided")
	ErrEmptyList    = errors.New("there are no saved passwords")
)

type (
	encoder interface {
		Encode(src any) ([]byte, error)
	}

	decoder interface {
		Decode(dst any, src []byte) error
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

	pwdReadRepository interface {
		ReadByID(ctx context.Context, id int) ([]byte, error)
	}

	pwdListRepository interface {
		ListNames(ctx context.Context) ([][2]string, error)
	}

	PwdService struct {
		l         logger.Logger
		addR      pwdAddRepository
		readR     pwdReadRepository
		listR     pwdListRepository
		encoder   encoder
		decoder   decoder
		encrypter encrypter
		decrypter decrypter
	}
)

func New(
	log logger.Logger,
	addR pwdAddRepository,
	readR pwdReadRepository,
	listR pwdListRepository,
	encoder encoder,
	decoder decoder,
	encrypter encrypter,
	decrypter decrypter,
) *PwdService {
	return &PwdService{
		log,
		addR, readR, listR,
		encoder, decoder,
		encrypter, decrypter,
	}
}

func (s *PwdService) Add(ctx context.Context, key string, o objects.PWD) (int, error) {
	const op = "pwdService.Add"
	log := s.l.With().Str("op", op).Logger()

	b, err := s.encoder.Encode(o)
	if err != nil {
		log.Error().Err(err).Msg("failed to encode object to bytes")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	s.encrypter.SetKey(key)
	data := s.encrypter.Encrypt(b)

	id, err := s.addR.Add(ctx, o.Name, data)
	if err != nil {
		log.Error().Err(err).Msg("failed to add password to repository")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *PwdService) Read(ctx context.Context, key string, id int) (objects.PWD, error) {
	const op = "pwdService.Read"
	log := s.l.With().Str("op", op).Logger()

	data, err := s.readR.ReadByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			log.Debug().Msg("password not exists")
			return objects.PWD{}, ErrPwdNotExists
		}
		log.Debug().Err(err).Msg("failed to read password from repository")
		return objects.PWD{}, fmt.Errorf("%s: %w", op, err)
	}

	s.decrypter.SetKey(key)
	b, err := s.decrypter.Decrypt(data)
	if err != nil {
		log.Debug().Err(err).Msg("failed to decrypt")
		return objects.PWD{}, ErrInvalidKey
	}

	var obj objects.PWD
	if err := s.decoder.Decode(&obj, b); err != nil {
		log.Error().Err(err).Msg("failed to decode bytes to object")
		return objects.PWD{}, fmt.Errorf("%s: %w", op, err)
	}

	return obj, nil
}

func (s *PwdService) List(ctx context.Context) ([][2]string, error) {
	const op = "pwdService.List"
	log := s.l.With().Str("op", op).Logger()

	data, err := s.listR.ListNames(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list passwords names")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(data) == 0 {
		log.Debug().Msg("empty list")
		return nil, ErrEmptyList
	}

	return data, nil
}
