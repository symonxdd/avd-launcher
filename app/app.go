package app

import (
	"avd-launcher/app/helper"
	"avd-launcher/app/models"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (a *App) SelectSdkPath() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Android SDK Path",
	})
	if err != nil || path == "" {
		return "", err
	}

	// Save to config
	// cfgPath := helper.GetConfigPath() // Removed: unused variable
	// cfg := helper.AppConfig{CustomSdkPath: path} // Removed: unused variable

	// Since json.Marshal is in encoding/json, and os is already imported, we need to make sure we deal with imports.
	// But it's easier to just do it in helper where it's already there
	if err := helper.SaveSdkPath(path); err != nil {
		return "", err
	}

	return path, nil
}
