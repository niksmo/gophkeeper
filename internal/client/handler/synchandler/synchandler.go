package synchandler

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/synccommand"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type SyncRunner interface {
	Run(ctx context.Context, token string)
}

type StartHandler struct {
	logger  logger.Logger
	service SyncRunner
	writer  io.Writer
}

func NewStart(l logger.Logger, s SyncRunner, w io.Writer) *StartHandler {
	return &StartHandler{l, s, w}
}

func (h *StartHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "SignupHandler.Handle"

	h.logger.Debug().Str("op", op).Msg("handle start")

	token := h.getToken(v)

	ctx, stop := signal.NotifyContext(
		ctx, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
	)
	defer stop()

	h.service.Run(ctx, token)
}

func (h *StartHandler) getToken(v command.ValueGetter) string {
	const op = "StartHandler.getToken"
	log := h.logger.WithOp(op)

	token, err := v.GetString(synccommand.TokenFlag)
	if err != nil {
		log.Debug().Msg("token not specified")
		h.printOutput("--%s flag not specified", synccommand.TokenFlag)
		os.Exit(1)
	}
	return token
}

func (h *StartHandler) printOutput(formated string, args ...any) {
	fmt.Fprintf(h.writer, formated, args...)
}
