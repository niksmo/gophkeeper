package pwdhandler_test

import (
	"context"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/objects"
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

func TestAdd(t *testing.T) {

}
