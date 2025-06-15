package pwdcommand_test

import (
	"context"
	"testing"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/stretchr/testify/mock"
)

type addHandler struct {
	mock.Mock
}

func (h *addHandler) Handle(ctx context.Context, v command.ValueGetter) {
	h.Called(ctx, v)
}

func TestAdd(t *testing.T) {

}
