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

func GetConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
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

// Resolves ANDROID_HOME or returns a default Windows path
func GetAndroidSdkPath() string {
	// 1. Check custom config first
	cfgPath := GetConfigPath()
	if data, err := os.ReadFile(cfgPath); err == nil {
		var cfg AppConfig
		if err := json.Unmarshal(data, &cfg); err == nil && cfg.CustomSdkPath != "" {
			return cfg.CustomSdkPath
		}
	}

	sdkPath := os.Getenv("ANDROID_HOME")
	if sdkPath != "" {
		return sdkPath
	}
	return ""
}

// Returns the avdmanager executable path
func GetAvdManagerPath() (string, error) {
	sdkPath := GetAndroidSdkPath()
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
	sdkPath := GetAndroidSdkPath()
	adbPath := filepath.Join(sdkPath, "platform-tools", "adb.exe")

	if _, err := os.Stat(adbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("adb not found at: %s", adbPath)
	}
	return adbPath, nil
}

// Returns the emulator executable path
func GetEmulatorPath() (string, error) {
	sdkPath := GetAndroidSdkPath()
	emulatorPath := filepath.Join(sdkPath, "emulator", "emulator.exe")
	if _, err := os.Stat(emulatorPath); os.IsNotExist(err) {
		return "", fmt.Errorf("emulator not found at: %s", emulatorPath)
	}
	return emulatorPath, nil
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
