package app

import (
	"avd-launcher/app/helper"
	"avd-launcher/app/models"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// The App struct (think of it like an object / class in other languages)
type App struct {
	ctx context.Context

	// Change to map of AVDs
	runningAVDs map[string]*models.AVD
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		runningAVDs: make(map[string]*models.AVD),
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
// (a *App) means the function belongs to that struct.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	fmt.Println("Shutting down app...")

	// Loop through any running AVDs and kill them
	for name, avd := range a.runningAVDs {
		fmt.Printf("Killing emulator for AVD: %s\n", name)
		_ = avd.Process.Process.Kill() // Or call adb kill, whatever you prefer
	}
}

// StartAVD starts the emulator for a given AVD
func (a *App) StartAVD(avdName string, coldBoot bool) string {
	emulatorPath, err := getEmulatorPath()
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
				runtime.EventsEmit(a.ctx, "avd-booted", helper.TimestampedLog(avdName))
			}
			// Detect shutdown or exit
			if strings.Contains(line, "Saving with gfxstream=1") {
				runtime.EventsEmit(a.ctx, "avd-shutdown", helper.TimestampedLog(avdName))
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

// StopAVD attempts to gracefully shut down a running AVD by name
func (a *App) StopAVD(name string) error {
	_, exists := a.runningAVDs[name] // ✅ Only use what you need
	if !exists {
		return fmt.Errorf("AVD not running: %s", name)
	}

	port, err := resolvePortForAVD(name)
	if err != nil {
		return fmt.Errorf("could not find emulator port: %w", err)
	}

	adbPath, err := getAdbPath()
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

// Retrieves list installed AVDs
func (a *App) ListAVDs() ([]string, error) {
	// Step 1: Get the path to the Android emulator executable.
	emulatorPath, err := getEmulatorPath()
	if err != nil {
		// If we fail to get the emulator path, return the error immediately.
		return nil, err
	}

	// Step 2: Prepare the command to list available AVDs.
	cmd := helper.NewCommand(emulatorPath, "-list-avds")

	// Ensure the command inherits the environment variables of the current process.
	cmd.Env = os.Environ()

	// Step 3: Execute the command and capture its output.
	out, err := cmd.Output()
	if err != nil {
		// If there's an error running the command, return a descriptive error.
		return nil, fmt.Errorf("error running emulator command: %s", err.Error())
	}

	// Step 4: Initialize a slice to store the list of AVD names.
	var avds []string

	// Step 5: Use a scanner to process the output line by line.
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		// Retrieve the current line and trim any leading/trailing whitespace (including \r, \n, or spaces).
		line := strings.TrimSpace(scanner.Text())

		// If the line is not empty, add it to the list of AVDs.
		if line != "" {
			avds = append(avds, line)
		}
	}

	// Step 6: Check if any error occurred while scanning the output.
	if err := scanner.Err(); err != nil {
		// If an error occurred during scanning, return it.
		return nil, fmt.Errorf("error reading command output: %s", err.Error())
	}

	// Step 7: If no AVDs were found, return a helpful error message.
	if len(avds) == 0 {
		return nil, fmt.Errorf("no AVDs found. Please create an AVD first")
	}

	// Step 8: Return the list of AVD names.
	return avds, nil
}

func getEmulatorPath() (string, error) {
	sdkPath := os.Getenv("ANDROID_SDK_ROOT")
	if sdkPath == "" {
		sdkPath = os.Getenv("ANDROID_HOME")
	}
	if sdkPath == "" {
		// Default fallback
		sdkPath = fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\Android\\Sdk", os.Getenv("USERNAME"))
	}

	emulatorPath := filepath.Join(sdkPath, "emulator", "emulator.exe")

	// Check if emulator exists
	if _, err := os.Stat(emulatorPath); os.IsNotExist(err) {
		return "", fmt.Errorf("emulator not found at: %s", emulatorPath)
	}

	return emulatorPath, nil
}

// GetEnvVariables returns important environment variables needed by the tool
func (a *App) GetEnvVariables() map[string]string {
	result := make(map[string]string)

	// Check Android SDK path
	androidHome := os.Getenv("ANDROID_HOME")
	fmt.Println("ANDROID_HOME:", androidHome)
	if androidHome == "" {
		androidHome = os.Getenv("ANDROID_SDK_ROOT")
		fmt.Println("ANDROID_SDK_ROOT:", androidHome)
	}

	// Fallback for common location
	if androidHome == "" {
		androidHome = fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\Android\\Sdk", os.Getenv("USERNAME"))
	}

	result["ANDROID_HOME"] = androidHome

	return result
}

// Retrieves the list of currently running Android Virtual Devices (AVDs)
func (a *App) ListRunningAVDs() ([]string, error) {
	// Get the path to the adb executable
	adbPath, err := getAdbPath()
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

// Helper to find adb
func getAdbPath() (string, error) {
	sdkPath := os.Getenv("ANDROID_SDK_ROOT")
	if sdkPath == "" {
		sdkPath = os.Getenv("ANDROID_HOME")
	}
	if sdkPath == "" {
		sdkPath = fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\Android\\Sdk", os.Getenv("USERNAME"))
	}

	adbPath := filepath.Join(sdkPath, "platform-tools", "adb.exe")
	if _, err := os.Stat(adbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("adb not found at: %s", adbPath)
	}
	return adbPath, nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func resolvePortForAVD(avdName string) (int, error) {
	fmt.Println("Resolving port for AVD:", avdName)

	adbPath, err := getAdbPath()
	if err != nil {
		return 0, err
	}

	output, err := helper.NewCommand(adbPath, "devices").Output()

	if err != nil {
		return 0, fmt.Errorf("failed to list adb devices: %w", err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "emulator-") && strings.Contains(line, "device") {
			deviceID := strings.Fields(line)[0]
			nameOut, err := helper.NewCommand(adbPath, "-s", deviceID, "emu", "avd", "name").Output()
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
