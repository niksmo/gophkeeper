package pwdhandler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler/pwdhandler"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// *Tests Mocks*

type valueGetter struct {
	mock.Mock
}

func (v *valueGetter) GetString(name string) (string, error) {
	args := v.Called(name)
	return args.String(0), args.Error(1)
}

func (v *valueGetter) GetInt(name string) (int, error) {
	args := v.Called(name)
	return args.Int(0), args.Error(1)
}

type pwdAddService struct {
	mock.Mock
}

func (s *pwdAddService) Add(
	ctx context.Context, key, name string, dto dto.PWD,
) (int, error) {
	args := s.Called(ctx, key, name, dto)
	return args.Int(0), args.Error(1)
}

type pwdReadService struct {
	mock.Mock
}

func (s *pwdReadService) Read(
	ctx context.Context, key string, id int,
) (dto.PWD, error) {
	args := s.Called(ctx, key, id)
	return args.Get(0).(dto.PWD), args.Error(1)
}

type pwdListService struct {
	mock.Mock
}

func (s *pwdListService) List(ctx context.Context) ([][2]string, error) {
	args := s.Called(ctx)
	return args.Get(0).([][2]string), args.Error(1)
}

type pwdEditService struct {
	mock.Mock
}

func (s *pwdEditService) Edit(
	ctx context.Context, key string, entryNum int, name string, dto dto.PWD,
) error {
	args := s.Called(ctx, key, entryNum, name, dto)
	return args.Error(0)
}

type pwdRemoveService struct {
	mock.Mock
}

func (s *pwdRemoveService) Remove(ctx context.Context, entryNum int) error {
	args := s.Called(ctx, entryNum)
	return args.Error(0)
}

// Test *AddPassword* handler

type addSuite struct {
	ctx         context.Context
	log         logger.Logger
	service     *pwdAddService
	valueGetter *valueGetter
	handler     *pwdhandler.PwdAddHandler
	w           io.Writer
	t           *testing.T
}

func newAddSuite(t *testing.T, w io.Writer) *addSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	service := &pwdAddService{}
	valueGetter := &valueGetter{}
	handler := pwdhandler.NewAddHandler(log, service, w)
	return &addSuite{ctx, log, service, valueGetter, handler, w, t}
}

func (st *addSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestAdd(t *testing.T) {
	const subProcEnv = "GOPHKEEPER_TEST_PWDHDR_ADD"
	const serviceMethod = "Add"
	t.Run("Ordinary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		st := newAddSuite(t, buf)
		defer st.PrettyPanic()

		expectedEntryNo := 1
		var expectedErr error
		expectedOut := fmt.Sprintf(
			"the password is saved under the record number %d\n",
			expectedEntryNo,
		)

		masterKey := "testKey"

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		st.valueGetter.On(
			"GetString", pwdcommand.MasterKeyFlag,
		).Return(masterKey, nil)

		st.valueGetter.On(
			"GetString", pwdcommand.NameFlag,
		).Return(obj.Name, nil)

		st.valueGetter.On(
			"GetString", pwdcommand.PasswordFlag,
		).Return(obj.Password, nil)

		st.valueGetter.On(
			"GetString", pwdcommand.LoginFlag,
		).Return(obj.Login, nil)

		st.service.On(
			serviceMethod, st.ctx, masterKey, obj.Name, obj,
		).Return(expectedEntryNo, expectedErr)

		st.handler.Handle(st.ctx, st.valueGetter)
		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("RequiredFlagsNotSpecified", func(t *testing.T) {
		var expectedErr error
		expectedEntryNo := 0
		masterKey := "testKey"
		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newAddSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(masterKey, errors.New(""))

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, errors.New(""))

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, errors.New(""))

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, errors.New(""))

			st.service.On(
				serviceMethod, st.ctx, masterKey, obj.Name, obj,
			).Return(expectedEntryNo, expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newAddSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n--%s\n--%s\n",
			pwdcommand.MasterKeyFlag,
			pwdcommand.NameFlag,
			pwdcommand.PasswordFlag,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("RequiredFlagsSpaceStr", func(t *testing.T) {
		var expectedErr error
		expectedEntryNo := 0
		masterKey := "      "
		obj := dto.PWD{
			Name:     " ",
			Login:    "",
			Password: "      ",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newAddSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(masterKey, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, nil)

			st.service.On(
				serviceMethod, st.ctx, masterKey, obj.Name, obj,
			).Return(expectedEntryNo, expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newAddSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n--%s\n--%s\n",
			pwdcommand.MasterKeyFlag,
			pwdcommand.NameFlag,
			pwdcommand.PasswordFlag,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("FailedSave", func(t *testing.T) {
		expectedEntryNo := 0
		expectedErr := errors.New("something happened with database")
		masterKey := "testKey"
		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newAddSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(masterKey, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, nil)

			st.service.On(
				serviceMethod, st.ctx, masterKey, obj.Name, obj,
			).Return(expectedEntryNo, expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newAddSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"the application completed with an error: %s\n",
			expectedErr.Error(),
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})
}

// Test *ReadPassword* handler

type readSuite struct {
	ctx         context.Context
	log         logger.Logger
	service     *pwdReadService
	valueGetter *valueGetter
	handler     *pwdhandler.PwdReadHandler
	w           io.Writer
	t           *testing.T
}

func newReadSuite(t *testing.T, w io.Writer) *readSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	service := &pwdReadService{}
	valueGetter := &valueGetter{}
	handler := pwdhandler.NewReadHandler(log, service, w)
	return &readSuite{ctx, log, service, valueGetter, handler, w, t}
}

func (st *readSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestRead(t *testing.T) {
	const subProcEnv = "GOPHKEEPER_TEST_PWDHDR_READ"
	const serviceMethod = "Read"
	t.Run("Ordinary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		st := newReadSuite(t, buf)
		defer st.PrettyPanic()

		key := "testKey"
		id := 1
		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}
		expectedOut := fmt.Sprintf(
			"the password with entry %d: name=%q, login=%q, password=%q\n",
			id, obj.Name, obj.Login, obj.Password,
		)

		st.valueGetter.On(
			"GetString", pwdcommand.MasterKeyFlag,
		).Return(key, nil)

		st.valueGetter.On(
			"GetInt", pwdcommand.EntryNumFlag,
		).Return(id, nil)

		st.service.On(serviceMethod, st.ctx, key, id).Return(obj, nil)

		st.handler.Handle(st.ctx, st.valueGetter)
		assert.Equal(t, expectedOut, buf.String())
	})

	t.Run("RequiredFlagsNotSpecified", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newReadSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, errors.New(""))

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, errors.New(""))

			st.service.On(serviceMethod, st.ctx, key, id).Return(obj, nil)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newReadSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n--%s\n",
			pwdcommand.MasterKeyFlag,
			pwdcommand.EntryNumFlag,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("RequiredFlagsSpaceStr", func(t *testing.T) {
		key := "    "
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newReadSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.service.On(serviceMethod, st.ctx, key, id).Return(obj, nil)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newReadSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n",
			pwdcommand.MasterKeyFlag,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("InvalidKey", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newReadSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.service.On(
				serviceMethod, st.ctx, key, id,
			).Return(obj, service.ErrInvalidKey)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newReadSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintln(service.ErrInvalidKey)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("NotExistsPassword", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newReadSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.service.On(
				serviceMethod, st.ctx, key, id,
			).Return(obj, service.ErrNotExists)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newReadSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 0
		expectedOut := fmt.Sprintf(
			"the password with entry number %d is not exists\nPASS\n",
			id,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)
		require.Equal(t, expectedExitCode, cmd.ProcessState.ExitCode())

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("ReadFailedInternalErr", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		expectedErr := errors.New("something happened with service")

		if os.Getenv(subProcEnv) == "1" {
			st := newReadSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.service.On(
				serviceMethod, st.ctx, key, id,
			).Return(obj, expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newReadSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"the application completed with an error: %s\n",
			expectedErr.Error(),
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})
}

// Test *ListPassword* handler

type listSuite struct {
	ctx         context.Context
	log         logger.Logger
	service     *pwdListService
	valueGetter *valueGetter
	handler     *pwdhandler.PwdListHandler
	w           io.Writer
	t           *testing.T
}

func newListSuite(t *testing.T, w io.Writer) *listSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	service := &pwdListService{}
	valueGetter := &valueGetter{}
	handler := pwdhandler.NewListHandler(log, service, w)
	return &listSuite{ctx, log, service, valueGetter, handler, w, t}
}

func (st *listSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestList(t *testing.T) {
	const subProcEnv = "GOPHKEEPER_TEST_PWDHDR_LIST"
	const serviceMethod = "List"
	t.Run("Ordinary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		st := newListSuite(t, buf)
		defer st.PrettyPanic()

		namesSlice := [][2]string{
			{"1", "testName1"},
			{"2", "testName1"},
		}

		st.service.On(serviceMethod, st.ctx).Return(namesSlice, nil)
		st.handler.Handle(st.ctx, st.valueGetter)
		expectedOut := fmt.Sprintf(
			"saved passwords names:\n%s\n%s\n",
			fmt.Sprintf("%s: %s", namesSlice[0][0], namesSlice[0][1]),
			fmt.Sprintf("%s: %s", namesSlice[1][0], namesSlice[1][1]),
		)
		assert.Equal(t, expectedOut, buf.String())
	})

	t.Run("EmptyList", func(t *testing.T) {
		if os.Getenv(subProcEnv) == "1" {
			st := newListSuite(t, os.Stdout)
			namesSlice := [][2]string{}

			st.service.On(serviceMethod, st.ctx).Return(namesSlice, nil)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newListSuite(t, buf)
		defer st.PrettyPanic()
		expectedExitCode := 0
		expectedOut := "there are no saved passwords\nPASS\n"

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)
		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)
		require.Equal(t, expectedExitCode, cmd.ProcessState.ExitCode())

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)

	})

	t.Run("FailedListInternalErr", func(t *testing.T) {
		expectedErr := errors.New("something happened in service")

		if os.Getenv(subProcEnv) == "1" {
			st := newListSuite(t, os.Stdout)
			namesSlice := [][2]string{}

			st.service.On(
				serviceMethod, st.ctx,
			).Return(namesSlice, expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newListSuite(t, buf)
		defer st.PrettyPanic()
		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"the application completed with an error: %s\n",
			expectedErr.Error(),
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)
		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})
}

// Test *EditPassword* handler

type editSuite struct {
	ctx         context.Context
	log         logger.Logger
	service     *pwdEditService
	valueGetter *valueGetter
	handler     *pwdhandler.PwdEditHandler
	w           io.Writer
	t           *testing.T
}

func newEditSuite(t *testing.T, w io.Writer) *editSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	service := &pwdEditService{}
	valueGetter := &valueGetter{}
	handler := pwdhandler.NewEditHandler(log, service, w)
	return &editSuite{ctx, log, service, valueGetter, handler, w, t}
}

func (st *editSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestEdit(t *testing.T) {
	const subProcEnv = "GOPHKEEPER_TEST_PWDHDR_EDIT"
	const serviceMethod = "Edit"
	t.Run("Ordinary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		st := newEditSuite(t, buf)
		defer st.PrettyPanic()

		key := "testKey"
		id := 1
		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}
		expectedOut := fmt.Sprintf(
			"the password under the record number %d was edited\n",
			id,
		)

		st.valueGetter.On(
			"GetString", pwdcommand.MasterKeyFlag,
		).Return(key, nil)

		st.valueGetter.On(
			"GetInt", pwdcommand.EntryNumFlag,
		).Return(id, nil)

		st.valueGetter.On(
			"GetString", pwdcommand.NameFlag,
		).Return(obj.Name, nil)

		st.valueGetter.On(
			"GetString", pwdcommand.PasswordFlag,
		).Return(obj.Password, nil)

		st.valueGetter.On(
			"GetString", pwdcommand.LoginFlag,
		).Return(obj.Login, nil)

		st.service.On(
			serviceMethod, st.ctx, key, id, obj.Name, obj,
		).Return(nil)

		st.handler.Handle(st.ctx, st.valueGetter)
		assert.Equal(t, expectedOut, buf.String())
	})

	t.Run("RequiredFlagsNotSpecified", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newEditSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, errors.New(""))

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, errors.New(""))

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, errors.New(""))

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, errors.New(""))

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, nil)

			st.service.On(
				serviceMethod, st.ctx, key, id, obj.Name, obj,
			).Return(nil)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newEditSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n--%s\n--%s\n--%s\n",
			pwdcommand.MasterKeyFlag,
			pwdcommand.NameFlag,
			pwdcommand.EntryNumFlag,
			pwdcommand.PasswordFlag,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("RequiredFlagsSpaceStr", func(t *testing.T) {
		key := "    "
		id := 1

		obj := dto.PWD{
			Name:     "  ",
			Login:    "testLogin",
			Password: "   ",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newEditSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, nil)

			st.service.On(serviceMethod, st.ctx, key, id, obj.Name, obj).Return(nil)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newReadSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n--%s\n--%s\n",
			pwdcommand.MasterKeyFlag,
			pwdcommand.NameFlag,
			pwdcommand.PasswordFlag,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("InvalidKey", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newEditSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, nil)

			st.service.On(
				serviceMethod, st.ctx, key, id, obj.Name, obj,
			).Return(service.ErrInvalidKey)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newEditSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintln(service.ErrInvalidKey)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("NotExistsPassword", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv(subProcEnv) == "1" {
			st := newEditSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, nil)

			st.service.On(
				serviceMethod, st.ctx, key, id, obj.Name, obj,
			).Return(service.ErrNotExists)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newEditSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 0
		expectedOut := fmt.Sprintf(
			"the password with entry number %d is not exists\nPASS\n",
			id,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)
		require.Equal(t, expectedExitCode, cmd.ProcessState.ExitCode())

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("EditFailedInternalErr", func(t *testing.T) {
		key := "testKey"
		id := 1

		obj := dto.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		expectedErr := errors.New("something happened with service")

		if os.Getenv(subProcEnv) == "1" {
			st := newEditSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetString", pwdcommand.MasterKeyFlag,
			).Return(key, nil)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.NameFlag,
			).Return(obj.Name, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.PasswordFlag,
			).Return(obj.Password, nil)

			st.valueGetter.On(
				"GetString", pwdcommand.LoginFlag,
			).Return(obj.Login, nil)

			st.service.On(
				serviceMethod, st.ctx, key, id, obj.Name, obj,
			).Return(expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newEditSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"the application completed with an error: %s\n",
			expectedErr.Error(),
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})
}

// Test *RemovePassword* handler

type removeSuite struct {
	ctx         context.Context
	log         logger.Logger
	service     *pwdRemoveService
	valueGetter *valueGetter
	handler     *pwdhandler.PwdRemoveHandler
	w           io.Writer
	t           *testing.T
}

func newRemoveSuite(t *testing.T, w io.Writer) *removeSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	service := &pwdRemoveService{}
	valueGetter := &valueGetter{}
	handler := pwdhandler.NewRemoveHandler(log, service, w)
	return &removeSuite{ctx, log, service, valueGetter, handler, w, t}
}

func (st *removeSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestRemove(t *testing.T) {
	const subProcEnv = "GOPHKEEPER_TEST_PWDHDR_REMOVE"
	const serviceMethod = "Remove"
	t.Run("Ordinary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		st := newRemoveSuite(t, buf)
		defer st.PrettyPanic()

		id := 1
		expectedOut := fmt.Sprintf(
			"the password under the record number %d was removed\n",
			id,
		)

		st.valueGetter.On(
			"GetInt", pwdcommand.EntryNumFlag,
		).Return(id, nil)

		st.service.On(serviceMethod, st.ctx, id).Return(nil)

		st.handler.Handle(st.ctx, st.valueGetter)
		assert.Equal(t, expectedOut, buf.String())
	})

	t.Run("RequiredFlagsNotSpecified", func(t *testing.T) {
		if os.Getenv(subProcEnv) == "1" {
			st := newRemoveSuite(t, os.Stdout)
			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(0, errors.New(""))

			st.service.On(serviceMethod, st.ctx, 0).Return(nil)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newRemoveSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n",
			pwdcommand.EntryNumFlag,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("NotExistsPassword", func(t *testing.T) {
		id := 1
		if os.Getenv(subProcEnv) == "1" {
			st := newRemoveSuite(t, os.Stdout)

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.service.On(
				serviceMethod, st.ctx, id,
			).Return(service.ErrNotExists)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newRemoveSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 0
		expectedOut := fmt.Sprintf(
			"the password with entry number %d is not exists\nPASS\n",
			id,
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)
		require.Equal(t, expectedExitCode, cmd.ProcessState.ExitCode())

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("ReadFailedInternalErr", func(t *testing.T) {
		expectedErr := errors.New("something happened with service")
		if os.Getenv(subProcEnv) == "1" {
			st := newRemoveSuite(t, os.Stdout)
			id := 1

			st.valueGetter.On(
				"GetInt", pwdcommand.EntryNumFlag,
			).Return(id, nil)

			st.service.On(serviceMethod, st.ctx, id).Return(expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
			return
		}

		buf := new(bytes.Buffer)
		st := newRemoveSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"the application completed with an error: %s\n",
			expectedErr.Error(),
		)

		cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
		cmd.Env = append(
			os.Environ(), subProcEnv+"=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})
}

func prettyPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Log(r)
		t.FailNow()
	}
}
