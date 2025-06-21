package repository_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/storage"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type syncRepoSuite struct {
	ctx context.Context
	r   *repository.SyncRepository
	s   *storage.Storage
}

func newSyncSuite(t *testing.T) *syncRepoSuite {
	log := logger.NewPretty("debug")
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()

	t.Cleanup(func() {
		os.Remove(dsn)
	})

	s := storage.New(log, dsn)
	s.MustRun(ctx)
	r := repository.NewSync(log, s)
	return &syncRepoSuite{ctx, r, s}
}

func TestSyncCreate(t *testing.T) {
	st := newSyncSuite(t)
	expectedPID := 12345
	expectedStartedAt := time.Now()
	obj, err := st.r.Create(st.ctx, 12345, time.Now())
	require.NoError(t, err)
	assert.Equal(t, expectedPID, obj.PID)
	assert.Zero(t, expectedStartedAt.Compare(obj.StartedAt))
}

func TestSyncReadLast(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newSyncSuite(t)
		_, err := st.r.Create(st.ctx, 12345, time.Now())
		require.NoError(t, err)

		expectedPID := 777
		expectedStartedAt := time.Now()
		_, err = st.r.Create(st.ctx, expectedPID, expectedStartedAt)
		require.NoError(t, err)

		last, err := st.r.ReadLast(st.ctx)
		require.NoError(t, err)
		assert.Equal(t, expectedPID, last.PID)
		assert.Zero(t, expectedStartedAt.Compare(last.StartedAt))

	})
	t.Run("EmptyTable", func(t *testing.T) {
		st := newSyncSuite(t)
		_, err := st.r.ReadLast(st.ctx)
		require.ErrorIs(t, err, repository.ErrNotExists)
	})
}

func TestSyncUpdate(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newSyncSuite(t)
		stoppedTime := time.Now()
		var obj dto.SyncDTO
		obj.ID = 2
		obj.PID = 12345
		obj.StartedAt = time.Now().Add(-time.Hour)
		obj.StoppedAt = &stoppedTime
		err := st.r.Update(st.ctx, obj)
		assert.ErrorIs(t, err, repository.ErrNotExists)
	})

	t.Run("NotExists", func(t *testing.T) {
		st := newSyncSuite(t)
		pid := 12345
		startedAt := time.Now().Add(-time.Hour)
		_, err := st.r.Create(st.ctx, pid, startedAt)
		require.NoError(t, err)

		obj, err := st.r.ReadLast(st.ctx)
		require.NoError(t, err)
		assert.Equal(t, pid, obj.PID)

		stoppedTime := time.Now()
		obj.StoppedAt = &stoppedTime
		err = st.r.Update(st.ctx, obj)
		require.NoError(t, err)

		updated, err := st.r.ReadLast(st.ctx)
		require.NoError(t, err)
		assert.Equal(t, pid, updated.PID)
		assert.Zero(t, obj.StartedAt.Compare(updated.StartedAt))
		assert.Zero(t, obj.StoppedAt.Compare(*updated.StoppedAt))
	})
}
