package genservice_test

import "github.com/stretchr/testify/mock"

type encoder struct {
	mock.Mock
}

func (e *encoder) Encode(src any) ([]byte, error) {
	args := e.Called(src)
	return args.Get(0).([]byte), args.Error(1)
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
