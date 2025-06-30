package genservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/service/genservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCreater struct {
	mock.Mock
}

func (r *MockCreater) Create(
	ctx context.Context, name string, data []byte,
) (int, error) {
	args := r.Called(ctx, name, data)
	return args.Int(0), args.Error(1)
}

type AddSuite struct {
	t         *testing.T
	ctx       context.Context
	log       logger.Logger
	repo      *MockCreater
	encoder   *encoder
	encrypter *encrypter
	service   *genservice.AddService[any]
}

func newAddSuiteAdd(t *testing.T) *AddSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	encoder := &encoder{}
	encrypter := &encrypter{}
	repo := &MockCreater{}
	service := genservice.NewAdd[any](log, repo, encoder, encrypter)
	st := &AddSuite{
		t, ctx, log,
		repo,
		encoder,
		encrypter,
		service,
	}
	return st
}

func (st *AddSuite) PrettyPanic() {
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
		st := newAddSuiteAdd(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr error
		expected := 1

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData, nil)
		st.repo.On(
			Create, st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj.Name, obj)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("EncodeFailed", func(t *testing.T) {
		st := newAddSuiteAdd(t)
		defer st.PrettyPanic()

		var repoAddErr error
		var encodeErr = errors.New("encode failed")
		expected := 0

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData, nil)
		st.repo.On(
			Create, st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj.Name, obj)
		require.ErrorIs(t, err, encodeErr)
		assert.Equal(t, expected, actual)
	})

	t.Run("RepoAddFailed", func(t *testing.T) {
		st := newAddSuiteAdd(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr = errors.New("repo add failed")
		expected := 0

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData, nil)
		st.repo.On(
			Create, st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj.Name, obj)
		require.ErrorIs(t, err, repoAddErr)
		assert.Equal(t, expected, actual)
	})
}
