package genservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/service/genservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUpdater struct {
	mock.Mock
}

func (r *MockUpdater) Update(
	ctx context.Context, entryNum int, name string, data []byte,
) error {
	args := r.Called(ctx, entryNum, name, data)
	return args.Error(0)
}

type EditSuite struct {
	t         *testing.T
	ctx       context.Context
	log       logger.Logger
	repo      *MockUpdater
	encoder   *encoder
	encrypter *encrypter
	service   *genservice.EditService[any]
}

func newEditSuite(t *testing.T) *EditSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	encoder := &encoder{}
	encrypter := &encrypter{}
	repo := &MockUpdater{}
	service := genservice.NewEdit[any](log, repo, encoder, encrypter)
	st := &EditSuite{
		t, ctx, log,
		repo,
		encoder,
		encrypter,
		service,
	}
	return st
}

func (st *EditSuite) PrettyPanic() {
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
		st := newEditSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr error

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData, nil)
		st.repo.On(
			Update, st.ctx, entryNum, obj.Name, encryptedData,
		).Return(repoAddErr)

		err := st.service.Edit(st.ctx, key, entryNum, obj.Name, obj)
		require.NoError(t, err)
	})

	t.Run("EncodeFailed", func(t *testing.T) {
		st := newEditSuite(t)
		defer st.PrettyPanic()

		var repoAddErr error
		var encodeErr = errors.New("encode failed")

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData, nil)
		st.repo.On(
			Update, st.ctx, entryNum, obj.Name, encryptedData,
		).Return(repoAddErr)

		err := st.service.Edit(st.ctx, key, entryNum, obj.Name, obj)
		require.ErrorIs(t, err, encodeErr)
	})

	t.Run("RepoUpdateFailed", func(t *testing.T) {
		st := newEditSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr = errors.New("repo add failed")

		st.encoder.On(Encode, obj).Return(encodedData, encodeErr)
		st.encrypter.On(SetKey, key)
		st.encrypter.On(Encrypt, encodedData).Return(encryptedData, nil)
		st.repo.On(
			Update, st.ctx, entryNum, obj.Name, encryptedData,
		).Return(repoAddErr)

		err := st.service.Edit(st.ctx, key, entryNum, obj.Name, obj)
		require.ErrorIs(t, err, repoAddErr)
	})
}
