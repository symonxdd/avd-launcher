package helper

import (
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

// Resolves ANDROID_HOME or returns a default Windows path
func GetAndroidSdkPath() string {
	sdkPath := os.Getenv("ANDROID_HOME")
	if sdkPath != "" {
		return sdkPath
	}
	return ""
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
