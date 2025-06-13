package pwdrepository_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/repository/pwdrepository"
	"github.com/niksmo/gophkeeper/internal/client/storage"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Suite struct {
	r   *pwdrepository.PwdRepository
	s   *storage.Storage
	ctx context.Context
}

func NewSuite(t *testing.T) Suite {
	l := logger.NewPretty("debug")
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
		os.Remove(dsn)
	})

	s := storage.New(l, dsn)
	s.MustRun(ctx)
	return Suite{pwdrepository.New(l, s), s, ctx}
}

func TestAdd(t *testing.T) {
	st := NewSuite(t)
	expectedID := 1
	expectedName := "testName"
	expectedData := []byte("helloWorld")
	id, err := st.r.Add(st.ctx, expectedName, expectedData)
	require.NoError(t, err)
	assert.Equal(t, expectedID, id)

	stmt := `
	SELECT name, data FROM passwords;
	`
	rows, err := st.s.QueryContext(st.ctx, stmt)
	require.NoError(t, err)
	defer rows.Close()
	var nRows int
	var name string
	var data []byte
	for rows.Next() {
		nRows++
		err := rows.Scan(&name, &data)
		require.NoError(t, err)
	}
	err = rows.Err()
	require.NoError(t, err)
	require.Equal(t, 1, nRows)
	assert.Equal(t, expectedName, name)
	assert.Equal(t, expectedData, data)
}
