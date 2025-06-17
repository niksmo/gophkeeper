package pwdcommand_test

import (
	"context"
	"os"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/stretchr/testify/mock"
)

type handler struct {
	mock.Mock
}

func (h *handler) Handle(ctx context.Context, v command.ValueGetter) {
	h.Called(ctx, v)
}

type suite struct {
	ctx context.Context
	h   *handler
	cmd *command.Command
}

func (st *suite) SetArgs(args []string) {
	os.Args = append(os.Args[:1], args...)
}

// Test *AddPassword* command

func newAddSuite(t *testing.T) *suite {
	ctx := context.Background()
	h := new(handler)
	cmd := pwdcommand.NewPwdAddCommand(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &suite{ctx, h, cmd}
	return st
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

	t.Run("MissedRequiredFlags", func(t *testing.T) {
		st := newAddSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.LoginFlag, "testLogin",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}

// Test *ReadPassword* command

func newReadSuite(t *testing.T) *suite {
	ctx := context.Background()
	h := new(handler)
	cmd := pwdcommand.NewPwdReadCommand(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &suite{ctx, h, cmd}
	return st
}

func TestRead(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newReadSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.MasterKeyFlag, "testKey",
			"--" + pwdcommand.EntryNumFlag, "1",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlag", func(t *testing.T) {
		st := newReadSuite(t)
		st.SetArgs([]string{})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedEntryNumFlag", func(t *testing.T) {
		st := newReadSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.MasterKeyFlag, "testKey",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}

// Test *ListPassword* command

func newListSuite(t *testing.T) *suite {
	ctx := context.Background()
	h := new(handler)
	cmd := pwdcommand.NewPwdListCommand(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &suite{ctx, h, cmd}
	return st
}

func TestList(t *testing.T) {
	st := newListSuite(t)
	st.SetArgs([]string{})
	st.h.On("Handle", st.ctx, st.cmd.Flags())
	st.cmd.ExecuteContext(st.ctx)
	expectedCalls := 1
	st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
}

// Test *EditPassword* command

func newEditSuite(t *testing.T) *suite {
	ctx := context.Background()
	h := new(handler)
	cmd := pwdcommand.NewPwdEditCommand(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &suite{ctx, h, cmd}
	return st
}

func TestEdit(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newEditSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.MasterKeyFlag, "testKey",
			"--" + pwdcommand.NameFlag, "testName",
			"--" + pwdcommand.PasswordFlag, "testPassword",
			"--" + pwdcommand.EntryNumFlag, "1",
			"--" + pwdcommand.LoginFlag, "testLogin",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("OnlyRequiredFlags", func(t *testing.T) {
		st := newEditSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.MasterKeyFlag, "testKey",
			"--" + pwdcommand.NameFlag, "testName",
			"--" + pwdcommand.PasswordFlag, "testPassword",
			"--" + pwdcommand.EntryNumFlag, "1",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlags", func(t *testing.T) {
		st := newEditSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.LoginFlag, "testLogin",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}

// Test *DeletePassword* command

func newDeleteSuite(t *testing.T) *suite {
	ctx := context.Background()
	h := new(handler)
	cmd := pwdcommand.NewPwdDeleteCommand(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &suite{ctx, h, cmd}
	return st
}

func TestDelete(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newDeleteSuite(t)
		st.SetArgs([]string{
			"--" + pwdcommand.EntryNumFlag, "1",
		})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlags", func(t *testing.T) {
		st := newDeleteSuite(t)
		st.SetArgs([]string{})
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}
