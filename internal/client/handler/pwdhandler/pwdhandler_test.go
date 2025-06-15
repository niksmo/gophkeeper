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
	"github.com/niksmo/gophkeeper/internal/client/handler/pwdhandler"
	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type valueGetter struct {
	mock.Mock
}

func (v *valueGetter) GetString(name string) (string, error) {
	args := v.Called(name)
	return args.String(0), args.Error(1)
}

type pwdAddService struct {
	mock.Mock
}

func (s *pwdAddService) Add(
	ctx context.Context, key string, obj objects.PWD,
) (int, error) {
	args := s.Called(ctx, key, obj)
	return args.Int(0), args.Error(1)
}

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
	prettyPanic(st.t)
}

func TestAdd(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		st := newAddSuite(t, buf)
		defer st.PrettyPanic()

		expectedEntryNo := 1
		var expectedErr error
		expectedOut := fmt.Sprintf(
			"the password is saved under the record number: %d\n",
			expectedEntryNo,
		)

		masterKey := "testKey"

		obj := objects.PWD{
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

		st.service.On("Add", st.ctx, masterKey, obj).Return(expectedEntryNo, expectedErr)

		st.handler.Handle(st.ctx, st.valueGetter)
		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("RequiredFlagsNotSpecified", func(t *testing.T) {
		var expectedErr error
		expectedEntryNo := 0
		masterKey := "testKey"
		obj := objects.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv("GOPHKEEPER_TEST_PWDADDHDR_REQUIREDFLAGS") == "1" {
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
				"Add", st.ctx, masterKey, obj,
			).Return(expectedEntryNo, expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
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

		cmd := exec.Command(
			os.Args[0], "-test.run=TestAdd/RequiredFlagsNotSpecified",
		)
		cmd.Env = append(
			os.Environ(), "GOPHKEEPER_TEST_PWDADDHDR_REQUIREDFLAGS=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)
		require.Equal(t, expectedExitCode, cmd.ProcessState.ExitCode())

		actualOut := buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("FailedSave", func(t *testing.T) {
		expectedEntryNo := 0
		expectedErr := errors.New("something happened with database")
		masterKey := "testKey"
		obj := objects.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

		if os.Getenv("GOPHKEEPER_TEST_PWDADDHDR_FAILEDSAVE") == "1" {
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
				"Add", st.ctx, masterKey, obj,
			).Return(expectedEntryNo, expectedErr)

			st.handler.Handle(st.ctx, st.valueGetter)
		}

		buf := new(bytes.Buffer)
		st := newAddSuite(t, buf)
		defer st.PrettyPanic()

		expectedExitCode := 1
		expectedOut := fmt.Sprintf(
			"the application completed with an error: %s\n",
			expectedErr.Error(),
		)

		cmd := exec.Command(
			os.Args[0], "-test.run=TestAdd/FailedSave",
		)
		cmd.Env = append(
			os.Environ(), "GOPHKEEPER_TEST_PWDADDHDR_FAILEDSAVE=1",
		)

		cmd.Stderr = os.Stderr
		cmd.Stdout = buf
		err := cmd.Run()
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr))
		require.Equal(t, exitErr.ExitCode(), expectedExitCode)
		require.Equal(t, expectedExitCode, cmd.ProcessState.ExitCode())

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
