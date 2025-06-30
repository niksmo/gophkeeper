package genservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/internal/client/service/genservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDeleter struct {
	mock.Mock
}

func (r *MockDeleter) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

type RemoveSuite struct {
	t       *testing.T
	ctx     context.Context
	log     logger.Logger
	repo    *MockDeleter
	service *genservice.RemoveService
}

func newRemoveSuite(t *testing.T) *RemoveSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	repo := &MockDeleter{}
	service := genservice.NewRemove(log, repo)
	return &RemoveSuite{
		t, ctx, log,
		repo,
		service,
	}
}

func (st *RemoveSuite) PrettyPanic() {
	st.t.Helper()
	if r := recover(); r != nil {
		st.t.Log(r)
		st.t.FailNow()
	}
}

func TestRemove(t *testing.T) {
	const Delete = "Delete"
	id := 1
	t.Run("Ordinary", func(t *testing.T) {
		st := newRemoveSuite(t)
		defer st.PrettyPanic()

		st.repo.On(Delete, st.ctx, id).Return(nil)
		err := st.service.Remove(st.ctx, id)
		assert.NoError(t, err)
	})

	t.Run("NotExists", func(t *testing.T) {
		st := newRemoveSuite(t)
		defer st.PrettyPanic()

		st.repo.On(Delete, st.ctx, id).Return(repository.ErrNotExists)
		err := st.service.Remove(st.ctx, id)
		assert.ErrorIs(t, err, service.ErrNotExists)
	})

	t.Run("RepoRemoveFailed", func(t *testing.T) {
		st := newRemoveSuite(t)
		defer st.PrettyPanic()

		repoErr := errors.New("something happened with repo")

		st.repo.On(Delete, st.ctx, id).Return(repoErr)
		err := st.service.Remove(st.ctx, id)
		assert.ErrorIs(t, err, repoErr)
	})

}
