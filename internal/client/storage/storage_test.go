package storage_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/storage"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Suite struct {
	S *storage.Storage
}

func NewSuite(t *testing.T) Suite {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.db")
	t.Cleanup(func() {
		t.Helper()
		os.Remove(dsn)
	})
	return Suite{storage.New(logger.NewPretty("debug"), dsn)}
}

func TestNew(t *testing.T) {
	if os.Getenv("GOPHKEEPER_STORAGETEST_NEW") == "1" {
		NewSuite(t)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNew")
	cmd.Env = append(os.Environ(), "GOPHKEEPER_STORAGETEST_NEW=1")
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	require.NoError(t, err)
	assert.True(t, cmd.ProcessState.Success())
}

func TestMustRun(t *testing.T) {
	st := NewSuite(t)
	if os.Getenv("GOPHKEEPER_STORAGETEST_MUSTRUN") == "1" {
		st.S.MustRun()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMustRun")
	cmd.Env = append(os.Environ(), "GOPHKEEPER_STORAGETEST_MUSTRUN=1")
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	require.NoError(t, err)
	assert.True(t, cmd.ProcessState.Success())

}
