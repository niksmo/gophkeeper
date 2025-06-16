package pwdservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/service/pwdservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type encoder struct {
	mock.Mock
}

func (e *encoder) Encode(src any) ([]byte, error) {
	args := e.Called(src)
	return args.Get(0).([]byte), args.Error(1)
}

type decoder struct {
	mock.Mock
}

func (d *decoder) Decode(dst any, src []byte) error {
	args := d.Called(dst, src)
	if obj, ok := dst.(*objects.PWD); ok {
		obj.Name = "testName"
		obj.Login = "testLogin"
		obj.Password = "testPassword"
	}
	return args.Error(0)
}

type encrypter struct {
	mock.Mock
}

func (e *encrypter) SetKey(k string) {
	e.Called(k)
}

func (e *encrypter) Encrypt(data []byte) []byte {
	args := e.Called(data)
	return args.Get(0).([]byte)
}

type decrypter struct {
	mock.Mock
}

func (d *decrypter) SetKey(k string) {
	d.Called(k)
}

func (d *decrypter) Decrypt(data []byte) ([]byte, error) {
	args := d.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

type pwdAddRepository struct {
	mock.Mock
}

func (r *pwdAddRepository) Add(
	ctx context.Context, name string, data []byte,
) (int, error) {
	args := r.Called(ctx, name, data)
	return args.Int(0), args.Error(1)
}

type pwdReadRepository struct {
	mock.Mock
}

func (r *pwdReadRepository) ReadByID(ctx context.Context, id int) ([]byte, error) {
	args := r.Called(ctx, id)
	return args.Get(0).([]byte), args.Error(1)
}

type pwdListRepository struct {
	mock.Mock
}

func (r *pwdListRepository) ListNames(ctx context.Context) ([][2]string, error) {
	args := r.Called(ctx)
	return args.Get(0).([][2]string), args.Error(1)
}

type addSuite struct {
	t         *testing.T
	ctx       context.Context
	log       logger.Logger
	addRepo   *pwdAddRepository
	encoder   *encoder
	encrypter *encrypter
	service   *pwdservice.PwdService
}

func newAddSuite(t *testing.T) *addSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	encoder := &encoder{}
	encrypter := &encrypter{}
	addRepo := &pwdAddRepository{}
	service := pwdservice.New(
		log,
		addRepo, &pwdReadRepository{}, &pwdListRepository{},
		encoder, &decoder{},
		encrypter, &decrypter{},
	)
	st := &addSuite{
		t, ctx, log,
		addRepo,
		encoder,
		encrypter,
		service,
	}
	return st
}

func (st *addSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestAdd(t *testing.T) {

	key := "testMasterKey"

	obj := objects.PWD{
		Name:     "testName",
		Login:    "testLogin",
		Password: "TestPassword",
	}

	encodedData := []byte("encodedData")
	encryptedData := []byte("encryptedData")

	t.Run("Ordinary", func(t *testing.T) {
		st := newAddSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr error
		expected := 1

		st.encoder.On("Encode", obj).Return(encodedData, encodeErr)
		st.encrypter.On("SetKey", key)
		st.encrypter.On("Encrypt", encodedData).Return(encryptedData)
		st.addRepo.On(
			"Add", st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("EncodeFailed", func(t *testing.T) {
		st := newAddSuite(t)
		defer st.PrettyPanic()

		var repoAddErr error
		var encodeErr = errors.New("encode failed")
		expected := 0

		st.encoder.On("Encode", obj).Return(encodedData, encodeErr)
		st.encrypter.On("SetKey", key)
		st.encrypter.On("Encrypt", encodedData).Return(encryptedData)
		st.addRepo.On(
			"Add", st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj)
		require.ErrorIs(t, err, encodeErr)
		assert.Equal(t, expected, actual)
	})

	t.Run("RepoAddFailed", func(t *testing.T) {
		st := newAddSuite(t)
		defer st.PrettyPanic()

		var encodeErr error
		var repoAddErr = errors.New("repo add failed")
		expected := 0

		st.encoder.On("Encode", obj).Return(encodedData, encodeErr)
		st.encrypter.On("SetKey", key)
		st.encrypter.On("Encrypt", encodedData).Return(encryptedData)
		st.addRepo.On(
			"Add", st.ctx, obj.Name, encryptedData,
		).Return(expected, repoAddErr)

		actual, err := st.service.Add(st.ctx, key, obj)
		require.ErrorIs(t, err, repoAddErr)
		assert.Equal(t, expected, actual)
	})
}

type readSuite struct {
	t         *testing.T
	ctx       context.Context
	log       logger.Logger
	readRepo  *pwdReadRepository
	decoder   *decoder
	decrypter *decrypter
	service   *pwdservice.PwdService
}

func newReadSuite(t *testing.T) *readSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	decoder := &decoder{}
	decrypter := &decrypter{}
	readRepo := &pwdReadRepository{}
	service := pwdservice.New(
		log,
		&pwdAddRepository{}, readRepo, &pwdListRepository{},
		&encoder{}, decoder,
		&encrypter{}, decrypter,
	)
	st := &readSuite{
		t, ctx, log,
		readRepo,
		decoder,
		decrypter,
		service,
	}
	return st
}

func (st *readSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestRead(t *testing.T) {
	key := "testMasterKey"
	id := 1
	encryptedData := []byte("encryptedData")
	encodedData := []byte("encodedData")
	var obj objects.PWD
	t.Run("Ordinary", func(t *testing.T) {
		st := newReadSuite(t)
		defer st.PrettyPanic()

		st.readRepo.On("ReadByID", st.ctx, id).Return(encryptedData, nil)
		st.decrypter.On("SetKey", key)
		st.decrypter.On("Decrypt", encryptedData).Return(encodedData, nil)
		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.NoError(t, err)
		expectedObj := objects.PWD{
			Name:     "testName",
			Login:    "testLogin",
			Password: "testPassword",
		}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("NotExists", func(t *testing.T) {
		st := newReadSuite(t)
		defer st.PrettyPanic()

		st.readRepo.On(
			"ReadByID", st.ctx, id,
		).Return(encryptedData, repository.ErrNotExists)

		st.decrypter.On("SetKey", key)
		st.decrypter.On("Decrypt", encryptedData).Return(encodedData, nil)
		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, pwdservice.ErrPwdNotExists)
		expectedObj := objects.PWD{}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("RepoReadFailed", func(t *testing.T) {
		st := newReadSuite(t)
		defer st.PrettyPanic()

		repoErr := errors.New("something happened with repo")

		st.readRepo.On(
			"ReadByID", st.ctx, id,
		).Return(encryptedData, repoErr)

		st.decrypter.On("SetKey", key)
		st.decrypter.On("Decrypt", encryptedData).Return(encodedData, nil)
		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, repoErr)
		expectedObj := objects.PWD{}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("DecryptFailedInvalidKey", func(t *testing.T) {
		st := newReadSuite(t)
		defer st.PrettyPanic()

		st.readRepo.On(
			"ReadByID", st.ctx, id,
		).Return(encryptedData, nil)

		st.decrypter.On("SetKey", key)

		st.decrypter.On(
			"Decrypt", encryptedData,
		).Return(encodedData, errors.New("invalid key"))

		st.decoder.On("Decode", &obj, encodedData).Return(nil)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, pwdservice.ErrInvalidKey)
		expectedObj := objects.PWD{}
		assert.Equal(t, expectedObj, obj)
	})

	t.Run("DecodeFailed", func(t *testing.T) {
		st := newReadSuite(t)
		defer st.PrettyPanic()

		decodeErr := errors.New("failed to decode")

		st.readRepo.On(
			"ReadByID", st.ctx, id,
		).Return(encryptedData, nil)

		st.decrypter.On("SetKey", key)

		st.decrypter.On(
			"Decrypt", encryptedData,
		).Return(encodedData, nil)

		st.decoder.On("Decode", &obj, encodedData).Return(decodeErr)

		obj, err := st.service.Read(st.ctx, key, id)
		require.ErrorIs(t, err, decodeErr)
		expectedObj := objects.PWD{}
		assert.Equal(t, expectedObj, obj)
	})
}

type listSuite struct {
	t        *testing.T
	ctx      context.Context
	log      logger.Logger
	listRepo *pwdListRepository
	service  *pwdservice.PwdService
}

func newListSuite(t *testing.T) *listSuite {
	ctx := context.Background()
	log := logger.NewPretty("debug")
	listRepo := &pwdListRepository{}
	service := pwdservice.New(
		log,
		&pwdAddRepository{}, &pwdReadRepository{}, listRepo,
		&encoder{}, &decoder{},
		&encrypter{}, &decrypter{},
	)
	st := &listSuite{
		t, ctx, log,
		listRepo,
		service,
	}
	return st
}

func (st *listSuite) PrettyPanic() {
	st.t.Helper()
	prettyPanic(st.t)
}

func TestList(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newListSuite(t)
		defer st.PrettyPanic()
		data := [][2]string{
			{"1", "testName1"},
			{"2", "testName2"},
		}

		st.listRepo.On("ListNames", st.ctx).Return(data, nil)
		actual, err := st.service.List(st.ctx)
		require.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	t.Run("EmptyList", func(t *testing.T) {
		st := newListSuite(t)
		defer st.PrettyPanic()
		data := [][2]string{}

		st.listRepo.On("ListNames", st.ctx).Return(data, nil)
		actual, err := st.service.List(st.ctx)
		require.ErrorIs(t, err, pwdservice.ErrEmptyList)
		assert.Nil(t, actual)
	})

	t.Run("RepoFailed", func(t *testing.T) {
		st := newListSuite(t)
		defer st.PrettyPanic()
		data := [][2]string{
			{"1", "testName1"},
			{"2", "testName2"},
		}
		listNamesErr := errors.New("something happened with repo")

		st.listRepo.On("ListNames", st.ctx).Return(data, listNamesErr)
		actual, err := st.service.List(st.ctx)
		require.Error(t, err)
		assert.Nil(t, actual)
	})
}

func prettyPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Log(r)
		t.FailNow()
	}
}
