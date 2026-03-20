package services

import (
	"avd-launcher/app/helper"
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SystemService struct {
	ctx context.Context
}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (s *SystemService) setContext(ctx context.Context) {
	s.ctx = ctx
}

// SetContext provides a way for other packages to inject the context without exposing it to Wails bindings.
func SetContext(s *SystemService, ctx context.Context) {
	s.setContext(ctx)
}

func (s *SystemService) OpenEnvironmentVariables() error {
	return helper.NewCommand("rundll32", "sysdm.cpl,EditEnvironmentVariables").Run()
}

func (s *SystemService) GetAndroidSdkEnv() helper.SdkInfo {
	return helper.GetAndroidSdkPath()
}

func (s *SystemService) OpenConfigFolder() {
	configPath := helper.GetConfigPath()
	// /select follows the pattern: explorer.exe /select,C:\path\to\file
	_ = helper.NewCommand("explorer", "/select,"+configPath).Run()
}

func (s *SystemService) SelectAndSaveSdkPath() (string, error) {
	path, err := runtime.OpenDirectoryDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Select Android SDK Location",
	})
	if err != nil {
		if err.Error() == "shellItem is nil" {
			return "", nil
		}
		return "", err
	}

	if path == "" {
		return "", nil
	}

	if !helper.IsValidSdkPath(path) {
		return "", fmt.Errorf("invalid Android SDK path: missing platform-tools or emulator")
	}

	if err := helper.SaveSdkPath(path); err != nil {
		return "", err
	}

	return path, nil
}
