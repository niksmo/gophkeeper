package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/niksmo/gophkeeper/internal/server/storage"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type usersSuite struct {
	t       *testing.T
	storage Storage
	repo    *UsersRepository
}

func newUsersSuite(t *testing.T) *usersSuite {
	t.Helper()

	dsn := os.Getenv("GOPHKEEPER_TEST_DB")
	if dsn == "" {
		t.Skip("Env GOPHKEEPER_TEST_DB not set")
	}
	logger := logger.NewPretty("debug")
	storage := storage.New(logger, dsn)
	repo := NewUsersRepository(logger, storage)
	st := &usersSuite{t, storage, repo}

	t.Cleanup(st.cleanup)
	return st
}

func (st *usersSuite) cleanup() {
	st.t.Helper()
	_, err := st.storage.ExecContext(
		context.Background(),
		"DELETE FROM users;",
	)
	if err != nil {
		st.t.Fatal(err)
	}
}

func TestUsers(t *testing.T) {

	t.Run("Create", func(t *testing.T) {

		t.Run("Ordinary", func(t *testing.T) {
			st := newUsersSuite(t)
			expectedID := 1
			expectedLogin := "testLogin"
			expectedPassword := "testPassword"
			expectedCreatedAt := time.Now()
			expectedDisable := false

			user, err := st.repo.Create(
				t.Context(), expectedLogin, expectedPassword,
			)
			require.NoError(t, err)
			assert.Equal(t, expectedID, user.ID)
			assert.Equal(t, expectedLogin, user.Login)
			assert.Equal(t, expectedPassword, user.Password)
			assert.True(
				t, user.CreatedAt.Sub(expectedCreatedAt) < time.Second*5,
			)
			assert.False(t, expectedDisable)
		})

		t.Run("AlreadyExists", func(t *testing.T) {
			st := newUsersSuite(t)
			login := "testLogin"
			password := "testPassword"

			_, err := st.repo.Create(t.Context(), login, password)
			require.NoError(t, err)

			_, err = st.repo.Create(t.Context(), login, password)
			require.ErrorIs(t, err, ErrAlreadyExists)
		})
	})

	t.Run("Read", func(t *testing.T) {

		t.Run("Ordinary", func(t *testing.T) {
			st := newUsersSuite(t)

			expectedUser, err := st.repo.Create(
				t.Context(), "testLogin", "testPassword",
			)
			require.NoError(t, err)

			user, err := st.repo.Read(t.Context(), expectedUser.ID)
			require.NoError(t, err)
			assert.Equal(t, expectedUser, user)
		})

		t.Run("NotExists", func(t *testing.T) {
			st := newUsersSuite(t)

			_, err := st.repo.Read(t.Context(), 1234)
			require.ErrorIs(t, err, ErrNotExists)
		})
	})
}
