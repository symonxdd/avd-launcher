package app

import (
	"avd-launcher/app/helper"
	"avd-launcher/app/models"
	"context"
	"fmt"

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

// Opens the folder containing the app's config file and selects the file
func (a *App) OpenConfigFolder() {
	configPath := helper.GetConfigPath()
	// /select follows the pattern: explorer.exe /select,C:\path\to\file
	cmd := helper.NewCommand("explorer", "/select,"+configPath)

	// Suppress error because explorer.exe /select often returns exit status 1
	// even when it successfully opens the window.
	_ = cmd.Run()
}

// Opens a directory selection dialog to choose the Android SDK path and saves it
func (a *App) SelectAndSaveSdkPath() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Android SDK Location",
	})
	if err != nil {
		if err.Error() == "shellItem is nil" {
			return "", nil // User cancelled (Windows)
		}
		return "", err
	}

	if path == "" {
		return "", nil // User cancelled
	}

	if !helper.IsValidSdkPath(path) {
		return "", fmt.Errorf("invalid Android SDK path: missing platform-tools or emulator")
	}

	err = helper.SaveSdkPath(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
