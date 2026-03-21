package services

import (
	"avd-launcher/app/helper"
	"avd-launcher/app/models"
	"context"
	"fmt"
	"strings"

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

// CheckAcceleration runs `emulator -accel-check` and returns the hypervisor status.
func (s *SystemService) CheckAcceleration() models.AccelInfo {
	emulatorPath, err := helper.GetEmulatorPath()
	if err != nil {
		return models.AccelInfo{
			Status:     "unavailable",
			Hypervisor: "Unknown",
			Details:    "Emulator binary not found. Install the Android Emulator via SDK Manager.",
		}
	}

	// emulator -accel-check outputs to both stdout and stderr; CombinedOutput captures both
	output, _ := helper.NewCommand(emulatorPath, "-accel-check").CombinedOutput()
	raw := string(output)

	return parseAccelOutput(raw)
}

// parseAccelOutput interprets the output of `emulator -accel-check`.
// Success example (Windows): "accel:\n0\nWHPX(10.0.22631) is installed and usable.\naccel"
// Success example (Linux):   "accel:\n0\nKVM (version 12) is installed and usable.\naccel"
// Success example (macOS):   "accel:\n0\nHypervisor.Framework OS X Version 13.2\naccel"
func parseAccelOutput(raw string) models.AccelInfo {
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")

	// Find the detail line (the one that mentions the hypervisor)
	var detailLine string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "accel" || trimmed == "accel:" {
			continue
		}
		// Skip the status code line (just a number)
		if len(trimmed) <= 3 {
			if _, isNum := isNumericString(trimmed); isNum {
				continue
			}
		}
		detailLine = trimmed
		break
	}

	if detailLine == "" {
		return models.AccelInfo{
			Status:     "unavailable",
			Hypervisor: "Unknown",
			Details:    "Could not determine acceleration status.",
		}
	}

	// Determine status
	status := "unavailable"
	if strings.Contains(strings.ToLower(detailLine), "is installed and usable") ||
		strings.Contains(strings.ToLower(detailLine), "acceleration can be used") {
		status = "available"
	}

	// Extract the hypervisor name from the detail line
	hypervisor := extractHypervisorName(detailLine)

	return models.AccelInfo{
		Status:     status,
		Hypervisor: hypervisor,
		Details:    detailLine,
	}
}

func extractHypervisorName(detail string) string {
	upper := strings.ToUpper(detail)

	switch {
	case strings.Contains(upper, "WHPX"):
		return "WHPX"
	case strings.Contains(upper, "HAXM"):
		return "HAXM"
	case strings.Contains(upper, "AEHD"):
		return "AEHD"
	case strings.Contains(upper, "KVM"):
		return "KVM"
	case strings.Contains(upper, "HYPERVISOR.FRAMEWORK"):
		return "Hypervisor.Framework"
	case strings.Contains(upper, "HVF"):
		return "HVF"
	default:
		return "Unknown"
	}
}

func isNumericString(s string) (int, bool) {
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
	}
	return 0, true
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
