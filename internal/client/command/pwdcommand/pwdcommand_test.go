package pwdcommand_test

import (
	"context"
	"os"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/stretchr/testify/mock"
)

type addHandler struct {
	mock.Mock
}

func (h *addHandler) Handle(ctx context.Context, v command.ValueGetter) {
	h.Called(ctx, v)
}

type addSuite struct {
	ctx context.Context
	h   *addHandler
	cmd *command.Command
}

func newAddSuite(t *testing.T) *addSuite {
	ctx := context.Background()
	h := new(addHandler)
	cmd := pwdcommand.NewPwdAddCommand(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &addSuite{ctx, h, cmd}
	return st
}

func (st *addSuite) SetArgs(args []string) {
	os.Args = append(os.Args[:1], args...)
}

func TestAdd(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newAddSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.MasterKeyFlag, "testKey",
			"--" + pwdcommand.NameFlag, "testName",
			"--" + pwdcommand.PasswordFlag, "testPassword",
			"--" + pwdcommand.LoginFlag, "testLogin",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("OnlyRequiredFlags", func(t *testing.T) {
		st := newAddSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.MasterKeyFlag, "testKey",
			"--" + pwdcommand.NameFlag, "testName",
			"--" + pwdcommand.PasswordFlag, "testPassword",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlag", func(t *testing.T) {
		st := newAddSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.MasterKeyFlag, "testKey",
			"--" + pwdcommand.NameFlag, "testName",
			"--" + pwdcommand.LoginFlag, "testLogin",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}
