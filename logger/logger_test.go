package logger_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/niksmo/gophkeeper/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	if os.Getenv("GOPHKEEPER_TEST_LOGGER") == "1" {
		logger.New("invalidLevel")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLogger")
	cmd.Env = append(os.Environ(), "GOPHKEEPER_TEST_LOGGER=1")
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	require.Error(t, err)
	assert.False(t, cmd.ProcessState.Success())
}
