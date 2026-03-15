package app

import (
	"avd-launcher/app/helper"
	"avd-launcher/app/models"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Retrieves list of installed AVDs
func (a *App) ListAVDs() ([]string, error) {
	// Step 1: Attempt fast discovery by reading .ini files directly from the AVD directory.
	// This avoids the overhead of spawning a Java CLI tool.
	avdDir, err := helper.GetAvdDirectory()
	if err == nil {
		files, err := os.ReadDir(avdDir)
		if err == nil {
			var avds []string
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".ini") {
					name := strings.TrimSuffix(file.Name(), ".ini")
					avds = append(avds, name)
				}
			}
			if len(avds) > 0 {
				return avds, nil
			}
		}
	}

	// Step 2: Fallback to the original emulator command if filesystem discovery fails or finds nothing.
	emulatorPath, err := helper.GetEmulatorPath()
	if err != nil {
		return nil, err
	}

	cmd := helper.NewCommand(emulatorPath, "-list-avds")
	cmd.Env = os.Environ()

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running emulator command: %s", err.Error())
	}

	var avds []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			avds = append(avds, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading command output: %s", err.Error())
	}

	if len(avds) == 0 {
		return nil, fmt.Errorf("no AVDs found. Please create an AVD first")
	}

	return avds, nil
}

// Retrieves the list of currently running Android Virtual Devices (AVDs)
func (a *App) ListRunningAVDs() ([]string, error) {
	// Get the path to the adb executable
	adbPath, err := helper.GetAdbPath()
	if err != nil {
		// If adb path can't be found, return an error
		return nil, err
	}

	// Prepare the command to list connected adb devices
	cmd := helper.NewCommand(adbPath, "devices")
	// Inherit environment variables from the current process
	cmd.Env = os.Environ()

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		// If the adb devices command fails, return an error
		return nil, fmt.Errorf("failed to run adb devices: %v", err)
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")

	// Prepare a slice to store the names of running AVDs
	var runningAVDs []string

	// Iterate over each line of adb devices output
	for _, line := range lines {
		// Look for lines that represent emulator instances
		// A running emulator will have a line starting with "emulator-" and containing "device"
		if strings.HasPrefix(line, "emulator-") && strings.Contains(line, "device") {
			// Split the line into fields (usually: device serial, state, etc.)
			parts := strings.Fields(line)
			if len(parts) > 0 {
				// Extract the device identifier (e.g., emulator-5554)
				device := parts[0]

				// Prepare a command to get the AVD name for this emulator device.
				nameCmd := helper.NewCommand(adbPath, "-s", device, "emu", "avd", "name")

				// Run the command and capture its output
				nameOutput, err := nameCmd.Output()
				if err != nil {
					// If fetching the AVD name fails, skip this emulator
					continue
				}

				// Clean up the command output by trimming spaces and splitting by lines.
				// Typically, the first line is the name, and subsequent lines might contain "OK" or other info
				nameLines := strings.Split(strings.TrimSpace(string(nameOutput)), "\n")

				// Extract and trim the actual AVD name from the first line
				name := strings.TrimSpace(nameLines[0])

				// If the name is not empty, add it to the list of running AVDs
				if name != "" {
					runningAVDs = append(runningAVDs, name)
				}
			}
		}
	}

	// Debug print: list all running AVD names found.
	fmt.Println("Running AVDs:", runningAVDs)

	// Return the list of running AVD names.
	return runningAVDs, nil
}

// Starts the emulator for a given AVD
func (a *App) StartAVD(avdName string, coldBoot bool) string {
	emulatorPath, err := helper.GetEmulatorPath()
	if err != nil {
		return "Failed to find emulator: " + err.Error()
	}

	avdName = strings.TrimSpace(avdName)

	start := func() (*helper.Command, io.ReadCloser, io.ReadCloser, error) {
		args := []string{"-avd", avdName}
		if coldBoot {
			args = append(args, "-no-snapshot-load")
		}

		cmd := helper.NewCommand(emulatorPath, args...)
		cmd.Env = os.Environ()
		cmd.Dir = filepath.Dir(emulatorPath)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, nil, nil, err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, nil, nil, err
		}

		if err := cmd.Start(); err != nil {
			return nil, nil, nil, err
		}
		return cmd, stdout, stderr, nil
	}

	cmd, stdout, stderr, err := start()
	if err != nil && strings.Contains(err.Error(), "already running") {
		// Lock issue detected, attempt to delete lock files and retry
		lockPath := filepath.Join(os.Getenv("USERPROFILE"), ".android", "avd", avdName+".avd")
		files, _ := filepath.Glob(filepath.Join(lockPath, "*.lock"))
		for _, file := range files {
			_ = os.Remove(file)
		}
		// Retry
		cmd, stdout, stderr, err = start()
		if err != nil {
			return "Failed to start emulator even after deleting lock files: " + err.Error()
		}
	} else if err != nil {
		return "Failed to start emulator: " + err.Error()
	}

	avd := &models.AVD{
		Name:    avdName,
		Process: cmd,
	}
	a.runningAVDs[avdName] = avd

	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			runtime.EventsEmit(a.ctx, "avd-log", helper.TimestampedLog(line))

			// Detect successful boot
			if strings.Contains(line, "Successfully loaded snapshot") || strings.Contains(line, "Boot completed") {
				fmt.Println("AVD booted successfully")
				runtime.EventsEmit(a.ctx, "avd-booted", avdName)
			}
			// Detect shutdown or exit
			if strings.Contains(line, "Saving with gfxstream=1") || strings.Contains(line, "killing emulator, bye bye") {
				runtime.EventsEmit(a.ctx, "avd-shutdown", avdName)
				delete(a.runningAVDs, avdName)
				break
			}

			// Detect multiple AVDs warning and handle here (optional logging)
			if strings.Contains(line, "Running multiple emulators with the same AVD") {
				runtime.EventsEmit(a.ctx, "avd-log", helper.TimestampedLog("Detected multiple emulator conflict. Lock file issue suspected."))
			}
		}
	}()

	return "Emulator started"
}

// Attempts to gracefully shut down a running AVD by name
func (a *App) StopAVD(name string) error {
	fmt.Println("[Go backend] Trying to stop:", name)

	gg, exists := a.runningAVDs[name] // ✅ Only use what you need

	fmt.Println("AVD exists:", exists)
	fmt.Println("gg:", gg)

	if !exists {
		return fmt.Errorf("AVD not running: %s", name)
	}

	port, err := helper.ResolvePortForAVD(name)
	if err != nil {
		return fmt.Errorf("could not find emulator port: %w", err)
	}

	adbPath, err := helper.GetAdbPath()
	if err != nil {
		return fmt.Errorf("failed to find adb: %w", err)
	}

	emulatorID := fmt.Sprintf("emulator-%d", port)
	fmt.Printf("Stopping AVD '%s' using %s\n", name, emulatorID)

	cmd := helper.NewCommand(adbPath, "-s", emulatorID, "emu", "kill")
	fmt.Printf("Executing command: %v\n", cmd.Args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop AVD '%s': %w, output: %s", name, err, string(output))
	}

	fmt.Printf("Stopped AVD '%s': %s\n", name, string(output))
	delete(a.runningAVDs, name)
	return nil
}

// Retrieves the Android SDK environment variable and its resolution source
func (a *App) GetAndroidSdkEnv() helper.SdkInfo {
	sdkInfo := helper.GetAndroidSdkPath()
	fmt.Printf("Resolved Android SDK: %s (via %s)\n", sdkInfo.Path, sdkInfo.Source)

	return sdkInfo
}

// Renames an existing AVD using avdmanager
func (a *App) RenameAVD(oldName, newName string) error {
	avdManagerStr, err := helper.GetAvdManagerPath()
	if err != nil {
		return err
	}

	cmd := helper.NewCommand(avdManagerStr, "move", "avd", "-n", oldName, "-r", newName)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rename AVD '%s' to '%s': %w, output: %s", oldName, newName, err, string(output))
	}

	fmt.Printf("Renamed AVD '%s' to '%s': %s\n", oldName, newName, string(output))
	return nil
}

// Deletes an existing AVD using avdmanager
func (a *App) DeleteAVD(avdName string) error {
	avdManagerStr, err := helper.GetAvdManagerPath()
	if err != nil {
		return err
	}

	cmd := helper.NewCommand(avdManagerStr, "delete", "avd", "-n", avdName)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete AVD '%s': %w, output: %s", avdName, err, string(output))
	}

	fmt.Printf("Deleted AVD '%s': %s\n", avdName, string(output))
	return nil
}

// Retrieves info about an AVD: path and running state (via lock-folder check).
// If a stale lock folder is detected (emulator process no longer alive), it is cleaned up.
// Does NOT calculate disk usage — use GetAvdDiskUsage for that.
func (a *App) GetAvdInfo(avdName string) (models.AvdInfo, error) {
	avdDir, err := helper.GetAvdDirectory()
	if err != nil {
		return models.AvdInfo{}, err
	}

	iniPath := filepath.Join(avdDir, avdName+".ini")
	file, err := os.Open(iniPath)
	if err != nil {
		return models.AvdInfo{}, err
	}
	defer file.Close()

	var avdPath string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "path=") {
			avdPath = strings.TrimPrefix(line, "path=")
			break
		}
	}

	if avdPath == "" {
		return models.AvdInfo{}, fmt.Errorf("could not find path in AVD ini file")
	}

	// Check for lock folder to determine if AVD is currently running
	lockPath := filepath.Join(avdPath, "hardware-qemu.ini.lock")
	isRunning := false

	if info, statErr := os.Stat(lockPath); statErr == nil && info.IsDir() {
		// Lock folder exists — check if the emulator process is actually alive
		pidFile := filepath.Join(lockPath, "pid")
		pidData, readErr := os.ReadFile(pidFile)
		if readErr == nil {
			pidStr := strings.TrimSpace(string(pidData))
			pid, parseErr := strconv.Atoi(pidStr)
			if parseErr == nil {
				if helper.IsProcessAlive(pid) {
					isRunning = true
				} else {
					// Stale lock: process is dead, clean up
					_ = os.RemoveAll(lockPath)
				}
			} else {
				// Can't parse PID, treat as stale
				_ = os.RemoveAll(lockPath)
			}
		} else {
			// No pid file in lock folder, treat as stale
			_ = os.RemoveAll(lockPath)
		}
	}

	return models.AvdInfo{
		Name:    avdName,
		Path:    avdPath,
		Running: isRunning,
	}, nil
}

// GetAvdDiskUsage calculates disk usage for an AVD (can be slow for large AVDs).
// Intended to be called in the background after cards are already displayed.
func (a *App) GetAvdDiskUsage(avdName string) (string, error) {
	avdDir, err := helper.GetAvdDirectory()
	if err != nil {
		return "", err
	}

	iniPath := filepath.Join(avdDir, avdName+".ini")
	file, err := os.Open(iniPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var avdPath string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "path=") {
			avdPath = strings.TrimPrefix(line, "path=")
			break
		}
	}

	if avdPath == "" {
		return "", fmt.Errorf("could not find path in AVD ini file")
	}

	size, err := helper.DirSize(avdPath)
	if err != nil {
		return "", err
	}

	return helper.FormatSize(size), nil
}

// Opens the AVD's directory in Windows Explorer
func (a *App) OpenAvdFolder(path string) {
	// 🧠 /select,path would select the folder in its parent. 
	// But usually users want to OPEN the folder content if it's a directory.
	// Or select it. Let's just open it.
	cmd := helper.NewCommand("explorer", path)
	_ = cmd.Run()
}
