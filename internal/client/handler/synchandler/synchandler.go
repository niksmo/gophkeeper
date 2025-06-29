package synchandler

import (
	"context"
	"io"
	"os/signal"
	"syscall"

	"github.com/niksmo/gophkeeper/internal/client/command/synccommand"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type SyncRunner interface {
	Run(ctx context.Context, token string)
}

type StartHandler struct {
	l logger.Logger
	s SyncRunner
	w io.Writer
}

func NewStart(l logger.Logger, s SyncRunner, w io.Writer) *StartHandler {
	return &StartHandler{l, s, w}
}

func (h *StartHandler) Handle(ctx context.Context, fv synccommand.StartCmdFlags) {
	const op = "SignupHandler.Handle"

	h.l.Debug().Str("op", op).Msg("handle start")

	ctx, stop := signal.NotifyContext(
		ctx, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
	)
	defer stop()

	h.s.Run(ctx, fv.Token)
}
