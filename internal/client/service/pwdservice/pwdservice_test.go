package pwdservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/internal/client/service/pwdservice"
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

type pwdAddRepository struct {
	mock.Mock
}

func (r *pwdAddRepository) Add(
	ctx context.Context, name string, data []byte,
) (int, error) {
	args := r.Called(ctx, name, data)
	return args.Int(0), args.Error(1)
}

type addSuite struct {
	ctx       context.Context
	log       logger.Logger
	repo      *pwdAddRepository
	encoder   *encoder
	encrypter *encrypter
	service   *pwdservice.PwdAddService
	t         *testing.T
}

func newAddSuite(t *testing.T) *addSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	encoder := &encoder{}
	encrypter := &encrypter{}
	repo := &pwdAddRepository{}
	service := pwdservice.NewAddService(log, repo, encoder, encrypter)
	st := &addSuite{ctx, log, repo, encoder, encrypter, service, t}
	return st
}

func (st *addSuite) PrettyPanic() {
	prettyPanic(st.t)
}

func TestAdd(t *testing.T) {

	key := "testMasterKey"

	obj := objects.PWD{
		Name:     "testName",
		Login:    "testLogin",
		Password: "TestPassword",
	}

	encodedData := []byte("encodedData")
	encryptedData := []byte("encryptedData")

	t.Run("Ordinary", func(t *testing.T) {
		st := newAddSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr error
		expected := 1

		st.encoder.On("Encode", obj).Return(encodedData, encodeErr)
		st.encrypter.On("SetKey", key)
		st.encrypter.On("Encrypt", encodedData).Return(encryptedData)
		st.repo.On(
			"Add", st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("EncodeFailed", func(t *testing.T) {
		st := newAddSuite(t)
		defer st.PrettyPanic()

		var repoAddErr error
		var encodeErr = errors.New("encode failed")
		expected := 0

		st.encoder.On("Encode", obj).Return(encodedData, encodeErr)
		st.encrypter.On("SetKey", key)
		st.encrypter.On("Encrypt", encodedData).Return(encryptedData)
		st.repo.On(
			"Add", st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj)
		require.ErrorIs(t, err, encodeErr)
		assert.Equal(t, expected, actual)
	})

	t.Run("RepoAddFailed", func(t *testing.T) {
		st := newAddSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr = errors.New("repo add failed")
		expected := 0

		st.encoder.On("Encode", obj).Return(encodedData, encodeErr)
		st.encrypter.On("SetKey", key)
		st.encrypter.On("Encrypt", encodedData).Return(encryptedData)
		st.repo.On(
			"Add", st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj)
		require.ErrorIs(t, err, repoAddErr)
		assert.Equal(t, expected, actual)
	})
}

func prettyPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Log(r)
		t.FailNow()
	}
}
