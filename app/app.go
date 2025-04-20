package app

import (
	"avd-launcher/app/helper"
	"avd-launcher/app/models"
	"context"
)

// The App struct (think of it like an object/class in other languages)
type App struct {
	ctx         context.Context
	runningAVDs map[string]*models.AVD
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		runningAVDs: make(map[string]*models.AVD),
	}
}

// Called when the app starts.
// `(a *App)` means the function belongs to that (`App`) struct
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenEnvironmentVariables() error {
	cmd := helper.NewCommand("rundll32", "sysdm.cpl,EditEnvironmentVariables")
	return cmd.Run()
}
