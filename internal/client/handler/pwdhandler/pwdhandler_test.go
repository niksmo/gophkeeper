package pwdhandler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/handler/pwdhandler"
	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	buf         *bytes.Buffer
	t           *testing.T
}

func NewAddSuite(t *testing.T) *addSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	service := &pwdAddService{}
	valueGetter := &valueGetter{}
	buf := new(bytes.Buffer)
	handler := pwdhandler.NewAddHandler(log, service, buf)
	return &addSuite{ctx, log, service, valueGetter, handler, buf, t}
}

func (st *addSuite) SetT(t *testing.T) {
	st.t = t
	st.t.Cleanup(func() {
		st.service.ExpectedCalls = nil
		st.valueGetter.ExpectedCalls = nil
		st.buf.Reset()
	})
}

func (st *addSuite) PrettyPanic() {
	prettyPanic(st.t)
}

func TestAdd(t *testing.T) {
	st := NewAddSuite(t)
	t.Run("Ordinary", func(t *testing.T) {
		st.SetT(t)
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
		actualOut := st.buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})

	t.Run("RequiredFlagsNotSpecified", func(t *testing.T) {
		st.SetT(t)
		defer st.PrettyPanic()

		expectedEntryNo := 1
		var expectedErr error
		expectedOut := fmt.Sprintf(
			"required flags are not specified:\n--%s\n--%s\n--%s\n",
			pwdcommand.MasterKeyFlag,
			pwdcommand.NameFlag,
			pwdcommand.PasswordFlag,
		)

		masterKey := "testKey"

		obj := objects.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}

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

		st.service.On("Add", st.ctx, masterKey, obj).Return(expectedEntryNo, expectedErr)

		st.handler.Handle(st.ctx, st.valueGetter)
		actualOut := st.buf.String()
		assert.Equal(t, expectedOut, actualOut)
	})
}

func prettyPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Log(r)
		t.FailNow()
	}
}
