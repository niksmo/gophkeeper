package editservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/service/editservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type encoder struct {
	mock.Mock
}

func (e *encoder) Encode(src any) ([]byte, error) {
	args := e.Called(src)
	return args.Get(0).([]byte), args.Error(1)
}

type encrypter struct {
	mock.Mock
}

func (e *encrypter) SetKey(k string) {
	e.Called(k)
}

func (e *encrypter) Encrypt(data []byte) []byte {
	args := e.Called(data)
	return args.Get(0).([]byte)
}

type repo struct {
	mock.Mock
}

func (r *repo) Update(
	ctx context.Context, entryNum int, name string, data []byte,
) error {
	args := r.Called(ctx, entryNum, name, data)
	return args.Error(0)
}

type suite struct {
	t         *testing.T
	ctx       context.Context
	log       logger.Logger
	repo      *repo
	encoder   *encoder
	encrypter *encrypter
	service   *editservice.EditService[any]
}

func newSuite(t *testing.T) *suite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	encoder := &encoder{}
	encrypter := &encrypter{}
	repo := &repo{}
	service := editservice.New[any](log, repo, encoder, encrypter)
	st := &suite{
		t, ctx, log,
		repo,
		encoder,
		encrypter,
		service,
	}
	return st
}

func (st *suite) PrettyPanic() {
	st.t.Helper()
	if r := recover(); r != nil {
		st.t.Log(r)
		st.t.FailNow()
	}
}

func TestEditService(t *testing.T) {
	const (
		Update  = "Update"
		Encode  = "Encode"
		SetKey  = "SetKey"
		Encrypt = "Encrypt"
	)
	key := "testMasterKey"
	entryNum := 1
	obj := struct{ Name string }{
		Name: "testName",
	}
	encodedData := []byte("encodedData")
	encryptedData := []byte("encryptedData")

	t.Run("Ordinary", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr error

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData)
		st.repo.On(
			Update, st.ctx, entryNum, obj.Name, encryptedData,
		).Return(repoAddErr)

		err := st.service.Edit(st.ctx, key, entryNum, obj.Name, obj)
		require.NoError(t, err)
	})

	t.Run("EncodeFailed", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		var repoAddErr error
		var encodeErr = errors.New("encode failed")

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData)
		st.repo.On(
			Update, st.ctx, entryNum, obj.Name, encryptedData,
		).Return(repoAddErr)

		err := st.service.Edit(st.ctx, key, entryNum, obj.Name, obj)
		require.ErrorIs(t, err, encodeErr)
	})

	t.Run("RepoUpdateFailed", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr = errors.New("repo add failed")

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData)
		st.repo.On(
			Update, st.ctx, entryNum, obj.Name, encryptedData,
		).Return(repoAddErr)

		err := st.service.Edit(st.ctx, key, entryNum, obj.Name, obj)
		require.ErrorIs(t, err, repoAddErr)
	})
}
