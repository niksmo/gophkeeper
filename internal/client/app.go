package client

import (
	"context"

	"github.com/niksmo/gophkeeper/pkg/logger"
)

type CMD interface {
	Execute()
}

type App struct {
	l logger.Logger
	c CMD
}

func New(l logger.Logger, c CMD) *App {
	return &App{l, c}
}

func (a *App) Run(ctx context.Context) {
	a.c.Execute()
}
