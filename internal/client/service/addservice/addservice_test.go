package addservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/service/addservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
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

func (r *repo) Create(
	ctx context.Context, name string, data []byte,
) (int, error) {
	args := r.Called(ctx, name, data)
	return args.Int(0), args.Error(1)
}

type suite struct {
	t         *testing.T
	ctx       context.Context
	log       logger.Logger
	repo      *repo
	encoder   *encoder
	encrypter *encrypter
	service   *addservice.AddService[any]
}

func newSuite(t *testing.T) *suite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	encoder := &encoder{}
	encrypter := &encrypter{}
	repo := &repo{}
	service := addservice.New[any](log, repo, encoder, encrypter)
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

func TestAddService(t *testing.T) {
	const (
		Add     = "Add"
		Create  = "Create"
		Encode  = "Encode"
		SetKey  = "SetKey"
		Encrypt = "Encrypt"
	)

	key := "testMasterKey"
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
		expected := 1

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData)
		st.repo.On(
			Create, st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj.Name, obj)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("EncodeFailed", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		var repoAddErr error
		var encodeErr = errors.New("encode failed")
		expected := 0

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData)
		st.repo.On(
			Create, st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj.Name, obj)
		require.ErrorIs(t, err, encodeErr)
		assert.Equal(t, expected, actual)
	})

	t.Run("RepoAddFailed", func(t *testing.T) {
		st := newSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr = errors.New("repo add failed")
		expected := 0

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData)
		st.repo.On(
			Create, st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj.Name, obj)
		require.ErrorIs(t, err, repoAddErr)
		assert.Equal(t, expected, actual)
	})
}
