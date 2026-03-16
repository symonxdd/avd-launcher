package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Alias for exec.Cmd
type Command = exec.Cmd

func TimestampedLog(s string) string {
	// 🧠 In Go, time.Format uses a specific reference time (Mon Jan 2 15:04:05 MST 2006) to define the layout — we need to pass an example time with the exact formatting we want.
	return fmt.Sprintf("(%s) %s", time.Now().Format("15:04:05"), s)
}

type AppConfig struct {
	CustomSdkPath string `json:"custom_sdk_path"`
}

type SdkInfo struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}

func GetConfigPath() string {
	// 🧠 On Windows, os.UserConfigDir() returns AppData\Roaming, but we want AppData\Local for this app's settings.
	dir := os.Getenv("LOCALAPPDATA")
	if dir == "" {
		// Fallback for non-Windows or if LOCALAPPDATA is missing
		var err error
		dir, err = os.UserConfigDir()
		if err != nil {
			dir = os.TempDir()
		}
	}

	appDir := filepath.Join(dir, "avd-launcher")
	os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, "config.json")
}

func SaveSdkPath(path string) error {
	cfgPath := GetConfigPath()
	cfg := AppConfig{CustomSdkPath: path}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0644)
}

// Checks if a path looks like a valid Android SDK (contains platform-tools or emulator)
func IsValidSdkPath(path string) bool {
	if path == "" {
		return false
	}
	// Check for platform-tools or emulator directory
	if _, err := os.Stat(filepath.Join(path, "platform-tools")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(path, "emulator")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(path, "platforms")); err == nil {
		return true
	}
	return false
}

// Resolves the Android SDK location and returns the path along with how it was resolved
func GetAndroidSdkPath() SdkInfo {
	// 1. Check custom config first
	cfgPath := GetConfigPath()
	if data, err := os.ReadFile(cfgPath); err == nil {
		var cfg AppConfig
		if err := json.Unmarshal(data, &cfg); err == nil && cfg.CustomSdkPath != "" {
			if IsValidSdkPath(cfg.CustomSdkPath) {
				return SdkInfo{Path: cfg.CustomSdkPath, Source: "user selected path"}
			}
		}
	}

	// 2. Check ANDROID_HOME env var
	sdkPath := os.Getenv("ANDROID_HOME")
	if IsValidSdkPath(sdkPath) {
		return SdkInfo{Path: sdkPath, Source: "ANDROID_HOME environment variable"}
	}

	// 3. Check common Windows default locations
	localAppData := os.Getenv("LOCALAPPDATA")
	commonPaths := []struct {
		path   string
		source string
	}{
		{filepath.Join(localAppData, "Android", "Sdk"), "default location (local appdata)"},
		{`C:\Android\Sdk`, "default location (C:\\Android\\Sdk)"},
		{`C:\Program Files (x86)\Android\android-sdk`, "default location (program files)"},
	}

	var info SdkInfo
	for _, cp := range commonPaths {
		if IsValidSdkPath(cp.path) {
			info = SdkInfo{Path: cp.path, Source: cp.source}
			break
		}
	}

	if info.Path == "" {
		return SdkInfo{Path: "", Source: "Not found"}
	}

	return info
}

// Returns the avdmanager executable path
func GetAvdManagerPath() (string, error) {
	sdkPath := GetAndroidSdkPath().Path
	cmdlineToolsPath := filepath.Join(sdkPath, "cmdline-tools", "latest", "bin", "avdmanager.bat")
	if _, err := os.Stat(cmdlineToolsPath); err == nil {
		return cmdlineToolsPath, nil
	}

	// Fallback to older tools/bin
	toolsPath := filepath.Join(sdkPath, "tools", "bin", "avdmanager.bat")
	if _, err := os.Stat(toolsPath); err == nil {
		return toolsPath, nil
	}

	return "", fmt.Errorf("avdmanager not found. Searched in cmdline-tools/latest/bin and tools/bin")
}

// Returns the adb executable path
func GetAdbPath() (string, error) {
	sdkPath := GetAndroidSdkPath().Path
	adbPath := filepath.Join(sdkPath, "platform-tools", "adb.exe")

	if _, err := os.Stat(adbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("adb not found at: %s", adbPath)
	}
	return adbPath, nil
}

// Returns the emulator executable path
func GetEmulatorPath() (string, error) {
	sdkPath := GetAndroidSdkPath().Path
	emulatorPath := filepath.Join(sdkPath, "emulator", "emulator.exe")
	if _, err := os.Stat(emulatorPath); os.IsNotExist(err) {
		return "", fmt.Errorf("emulator not found at: %s", emulatorPath)
	}
	return emulatorPath, nil
}

// Returns the platform-specific AVD storage directory
func GetAvdDirectory() (string, error) {
	// Usually ~/.android/avd or %USERPROFILE%\.android\avd
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}

	avdPath := filepath.Join(home, ".android", "avd")
	if _, err := os.Stat(avdPath); os.IsNotExist(err) {
		return "", fmt.Errorf("AVD directory not found at: %s", avdPath)
	}

	return avdPath, nil
}

func ResolvePortForAVD(avdName string) (int, error) {
	fmt.Println("Resolving port for AVD:", avdName)

	adbPath, err := GetAdbPath()
	if err != nil {
		return 0, err
	}

	output, err := NewCommand(adbPath, "devices").Output()

	if err != nil {
		return 0, fmt.Errorf("failed to list adb devices: %w", err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "emulator-") && strings.Contains(line, "device") {
			deviceID := strings.Fields(line)[0]
			nameOut, err := NewCommand(adbPath, "-s", deviceID, "emu", "avd", "name").Output()
			if err != nil {
				continue
			}

			// Just grab the first line before "OK"
			actualName := strings.TrimSpace(strings.SplitN(string(nameOut), "\n", 2)[0])
			if actualName == avdName {
				portStr := strings.TrimPrefix(deviceID, "emulator-")
				port, err := strconv.Atoi(portStr)
				if err != nil {
					return 0, fmt.Errorf("invalid port in %s", deviceID)
				}
				fmt.Printf("Resolved %s to port %d\n", avdName, port)
				return port, nil
			}
		}
	}

	return 0, fmt.Errorf("AVD %s not found among running devices", avdName)
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func FormatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
