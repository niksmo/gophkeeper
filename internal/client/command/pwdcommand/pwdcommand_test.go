package pwdcommand_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/stretchr/testify/mock"
)

// Test *AddPassword* command

type MockAddCmdHandler struct {
	mock.Mock
}

func (h *MockAddCmdHandler) Handle(
	ctx context.Context, fv pwdcommand.AddCmdFlags,
) {
	h.Called(ctx, fv)
}

type AddCmdSuite struct {
	ctx context.Context
	h   *MockAddCmdHandler
	cmd *command.Command
}

func newAddSuite(t *testing.T) *AddCmdSuite {
	h := new(MockAddCmdHandler)
	cmd := pwdcommand.NewAdd(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &AddCmdSuite{t.Context(), h, cmd}
	return st
}

func (st *AddCmdSuite) SetArgs(args ...string) {
	os.Args = append(os.Args[:1], args...)
}

func TestAdd(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newAddSuite(t)
		fv := pwdcommand.AddCmdFlags{
			Key:      "testKey",
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		st.SetArgs(
			"--"+pwdcommand.SecretKeyFlag, fv.Key,
			"--"+pwdcommand.NameFlag, fv.Name,
			"--"+pwdcommand.PasswordFlag, fv.Password,
			"--"+pwdcommand.LoginFlag, fv.Login,
		)
		st.h.On("Handle", st.ctx, fv)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("OnlyRequiredFlags", func(t *testing.T) {
		st := newAddSuite(t)
		fv := pwdcommand.AddCmdFlags{
			Key:      "testKey",
			Name:     "testName",
			Password: "testPassword",
		}
		st.SetArgs(
			"--"+pwdcommand.SecretKeyFlag, fv.Key,
			"--"+pwdcommand.NameFlag, fv.Name,
			"--"+pwdcommand.PasswordFlag, fv.Password,
		)
		st.h.On("Handle", st.ctx, fv)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlags", func(t *testing.T) {
		st := newAddSuite(t)
		fv := pwdcommand.AddCmdFlags{
			Login: "testLogin",
		}
		st.SetArgs(
			"--"+pwdcommand.LoginFlag, fv.Login,
		)
		st.h.On("Handle", st.ctx, fv)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}

// Test *ReadPassword* command

type MockReadCmdHandler struct {
	mock.Mock
}

func (h *MockReadCmdHandler) Handle(
	ctx context.Context, masterKey string, entryNum int,
) {
	h.Called(ctx, masterKey, entryNum)
}

type ReadCmdSuite struct {
	ctx context.Context
	h   *MockReadCmdHandler
	cmd *command.Command
}

func newReadSuite(t *testing.T) *ReadCmdSuite {
	h := new(MockReadCmdHandler)
	cmd := pwdcommand.NewRead(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &ReadCmdSuite{t.Context(), h, cmd}
	return st
}

func (st *ReadCmdSuite) SetArgs(args ...string) {
	os.Args = append(os.Args[:1], args...)
}

func TestRead(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newReadSuite(t)
		masterKey := "testKey"
		entryNum := 1
		st.SetArgs(
			"--"+pwdcommand.SecretKeyFlag, masterKey,
			"--"+pwdcommand.EntryNumFlag, strconv.Itoa(entryNum),
		)
		st.h.On("Handle", st.ctx, masterKey, entryNum)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlag", func(t *testing.T) {
		st := newReadSuite(t)
		st.SetArgs()
		st.h.On("Handle", st.ctx, st.cmd.Flags())
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedEntryNumFlag", func(t *testing.T) {
		st := newReadSuite(t)
		masterKey := "testKey"
		st.SetArgs(
			"--"+pwdcommand.SecretKeyFlag, masterKey,
		)
		st.h.On("Handle", st.ctx, masterKey, 0)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}

// Test *ListPassword* command

type MockListCmdHandler struct {
	mock.Mock
}

func (h *MockListCmdHandler) Handle(ctx context.Context) {
	h.Called(ctx)
}

func TestList(t *testing.T) {
	osArgs := os.Args
	t.Cleanup(func() {
		os.Args = osArgs
	})

	mockListH := new(MockListCmdHandler)
	listCmd := pwdcommand.NewList(mockListH)
	os.Args = os.Args[:1]

	mockListH.On("Handle", t.Context())
	listCmd.ExecuteContext(t.Context())
	expectedCalls := 1
	mockListH.AssertNumberOfCalls(t, "Handle", expectedCalls)
}

// Test *EditPassword* command

type MockEditCmdHandler struct {
	mock.Mock
}

func (h *MockEditCmdHandler) Handle(
	ctx context.Context, fv pwdcommand.EditCmdFlags,
) {
	h.Called(ctx, fv)
}

type AddEditSuite struct {
	ctx context.Context
	h   *MockEditCmdHandler
	cmd *command.Command
}

func newEditSuite(t *testing.T) *AddEditSuite {
	h := new(MockEditCmdHandler)
	cmd := pwdcommand.NewEdit(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &AddEditSuite{t.Context(), h, cmd}
	return st
}

func (st *AddEditSuite) SetArgs(args ...string) {
	os.Args = append(os.Args[:1], args...)
}

func TestEdit(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newEditSuite(t)
		fv := pwdcommand.EditCmdFlags{
			Key:      "testKey",
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
			EntryNum: 1,
		}
		st.SetArgs(
			"--"+pwdcommand.SecretKeyFlag, fv.Key,
			"--"+pwdcommand.NameFlag, fv.Name,
			"--"+pwdcommand.PasswordFlag, fv.Password,
			"--"+pwdcommand.EntryNumFlag, strconv.Itoa(fv.EntryNum),
			"--"+pwdcommand.LoginFlag, fv.Login,
		)
		st.h.On("Handle", st.ctx, fv)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("OnlyRequiredFlags", func(t *testing.T) {
		st := newEditSuite(t)
		fv := pwdcommand.EditCmdFlags{
			Key:      "testKey",
			Name:     "testName",
			Password: "testPassword",
			EntryNum: 1,
		}
		st.SetArgs(
			"--"+pwdcommand.SecretKeyFlag, fv.Key,
			"--"+pwdcommand.NameFlag, fv.Name,
			"--"+pwdcommand.PasswordFlag, fv.Password,
			"--"+pwdcommand.EntryNumFlag, strconv.Itoa(fv.EntryNum),
		)
		st.h.On("Handle", st.ctx, fv)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlags", func(t *testing.T) {
		st := newEditSuite(t)
		st.SetArgs(
			"--"+pwdcommand.LoginFlag, "testLogin",
		)
		st.h.On("Handle", st.ctx, pwdcommand.EditCmdFlags{})
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}

// Test *RemovePassword* command

type MockRemoveCmdHandler struct {
	mock.Mock
}

func (h *MockRemoveCmdHandler) Handle(
	ctx context.Context, entryNum int,
) {
	h.Called(ctx, entryNum)
}

type AddRemoveSuite struct {
	ctx context.Context
	h   *MockRemoveCmdHandler
	cmd *command.Command
}

func newRemoveSuite(t *testing.T) *AddRemoveSuite {
	h := new(MockRemoveCmdHandler)
	cmd := pwdcommand.NewRemove(h)
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	st := &AddRemoveSuite{t.Context(), h, cmd}
	return st
}

func (st *AddRemoveSuite) SetArgs(args ...string) {
	os.Args = append(os.Args[:1], args...)
}

func TestRemove(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newRemoveSuite(t)
		entryNum := 1
		st.SetArgs(
			"--"+pwdcommand.EntryNumFlag, strconv.Itoa(entryNum),
		)
		st.h.On("Handle", st.ctx, entryNum)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 1
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})

	t.Run("MissedRequiredFlags", func(t *testing.T) {
		st := newRemoveSuite(t)
		st.SetArgs()
		st.h.On("Handle", st.ctx, 0)
		st.cmd.ExecuteContext(st.ctx)
		expectedCalls := 0
		st.h.AssertNumberOfCalls(t, "Handle", expectedCalls)
	})
}
