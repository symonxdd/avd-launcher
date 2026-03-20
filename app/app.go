package app

import (
	"context"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) setContext(ctx context.Context) {
	a.ctx = ctx
}

// SetContext provides a way for the main package to inject the context without exposing it to Wails bindings.
func SetContext(a *App, ctx context.Context) {
	a.setContext(ctx)
}
