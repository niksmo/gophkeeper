package listservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/service/listservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type repo struct {
	mock.Mock
}

func (r *repo) ListNames(ctx context.Context) ([][2]string, error) {
	args := r.Called(ctx)
	return args.Get(0).([][2]string), args.Error(1)
}

type suite struct {
	t       *testing.T
	ctx     context.Context
	log     logger.Logger
	repo    *repo
	service *listservice.ListService
}

func newListSuite(t *testing.T) *suite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	repo := &repo{}
	service := listservice.New(log, repo)
	return &suite{
		t, ctx, log,
		repo,
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

func TestList(t *testing.T) {
	const ListNames = "ListNames"
	t.Run("Ordinary", func(t *testing.T) {
		st := newListSuite(t)
		defer st.PrettyPanic()

		expectedData := [][2]string{
			{"1", "testName1"},
			{"2", "testName2"},
		}

		st.repo.On(ListNames, st.ctx).Return(expectedData, nil)
		actual, err := st.service.List(st.ctx)
		require.NoError(t, err)
		assert.Equal(t, expectedData, actual)
	})

	t.Run("EmptyList", func(t *testing.T) {
		st := newListSuite(t)
		defer st.PrettyPanic()

		expectedData := [][2]string{}

		st.repo.On(ListNames, st.ctx).Return(expectedData, nil)
		actual, err := st.service.List(st.ctx)
		require.NoError(t, err)
		assert.Equal(t, expectedData, actual)
	})

	t.Run("RepoFailed", func(t *testing.T) {
		st := newListSuite(t)
		defer st.PrettyPanic()

		data := [][2]string{
			{"1", "testName1"},
			{"2", "testName2"},
		}

		listNamesErr := errors.New("something happened with repo")

		st.repo.On(ListNames, st.ctx).Return(data, listNamesErr)
		actual, err := st.service.List(st.ctx)
		require.Error(t, err)
		assert.Nil(t, actual)
	})
}
