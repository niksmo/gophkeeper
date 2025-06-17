package readservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/internal/client/service/readservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type dto struct {
	Name string
	Data string
}

type decoder struct {
	mock.Mock
}

func (d *decoder) Decode(dst any, src []byte) error {
	args := d.Called(dst, src)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	if obj, ok := dst.(*dto); ok {
		obj.Name = "testName"
		obj.Data = "decodedData"
	}
	return args.Error(0)
}

type decrypter struct {
	mock.Mock
}

func (d *decrypter) SetKey(k string) {
	d.Called(k)
}

func (d *decrypter) Decrypt(data []byte) ([]byte, error) {
	args := d.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

type repo struct {
	mock.Mock
}

func (r *repo) ReadByID(ctx context.Context, id int) ([]byte, error) {
	args := r.Called(ctx, id)
	return args.Get(0).([]byte), args.Error(1)
}

type suite struct {
	t         *testing.T
	ctx       context.Context
	log       logger.Logger
	repo      *repo
	decoder   *decoder
	decrypter *decrypter
	service   *readservice.ReadService[dto]
}

func newSuite(t *testing.T) *suite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	decoder := &decoder{}
	decrypter := &decrypter{}
	repo := &repo{}
	service := readservice.New[dto](log, repo, decoder, decrypter)
	return &suite{
		t, ctx, log,
		repo,
		decoder,
		decrypter,
		service,
	}
}

func (st *suite) PrettyPanic() {
	st.t.Helper()
	if r := recover(); r != nil {
		st.t.Log(r)
		st.t.FailNow()
	}
}

func TestRead(t *testing.T) {
	key := "testMasterKey"
	id := 1
	encryptedData := []byte("encryptedData")
	encodedData := []byte("encodedData")
	var obj dto
	t.Run("Ordinary", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		st.repo.On("ReadByID", st.ctx, id).Return(encryptedData, nil)
		st.decrypter.On("SetKey", key)
		st.decrypter.On("Decrypt", encryptedData).Return(encodedData, nil)
		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.NoError(t, err)
		expectedObj := dto{
			Name: "testName",
			Data: "decodedData",
		}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("NotExists", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		st.repo.On("ReadByID", st.ctx, id).Return(
			encryptedData, repository.ErrNotExists,
		)

		st.decrypter.On("SetKey", key)
		st.decrypter.On("Decrypt", encryptedData).Return(encodedData, nil)
		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, service.ErrNotExists)
		expectedObj := dto{}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("RepoReadFailed", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		repoErr := errors.New("something happened with repo")

		st.repo.On("ReadByID", st.ctx, id).Return(encryptedData, repoErr)

		st.decrypter.On("SetKey", key)
		st.decrypter.On("Decrypt", encryptedData).Return(encodedData, nil)
		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, repoErr)
		expectedObj := dto{}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("DecryptFailedInvalidKey", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		st.repo.On("ReadByID", st.ctx, id).Return(encryptedData, nil)
		st.decrypter.On("SetKey", key)

		st.decrypter.On(
			"Decrypt", encryptedData,
		).Return(encodedData, errors.New("invalid key"))

		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, service.ErrInvalidKey)
		expectedObj := dto{}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("DecodeFailed", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		decodeErr := errors.New("failed to decode")

		st.repo.On("ReadByID", st.ctx, id).Return(encryptedData, nil)
		st.decrypter.On("SetKey", key)
		st.decrypter.On("Decrypt", encryptedData).Return(encodedData, nil)
		st.decoder.On("Decode", &obj, encodedData).Return(decodeErr)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, decodeErr)
		expectedObj := dto{}
		assert.Equal(t, expectedObj, obj)
	})
}
