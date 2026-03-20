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
	"regexp"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func sanitizeAvdID(name string) string {
	// 🧠 Trim and replace all non-allowed chars with spaces first, 
	// then collapse those spaces into single underscores.
	name = strings.TrimSpace(name)
	regInvalid := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	name = regInvalid.ReplaceAllString(name, " ")
	
	regSpace := regexp.MustCompile(`\s+`)
	name = regSpace.ReplaceAllString(name, "_")
	
	// Final trim of underscores from the ends
	return strings.Trim(name, "_")
}

func stripEmojis(s string) string {
	// 🧠 Keep only alphanumeric, spaces, and common safe punctuation for the standard field.
	// This ensures Android Studio/SDK tools don't crash or error out.
	reg := regexp.MustCompile(`[^\x00-\x7F]+`)
	return strings.TrimSpace(reg.ReplaceAllString(s, ""))
}

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

// Renames an existing AVD by manually updating its file references and metadata
func (a *App) RenameAVD(oldName, newName string) error {
	newID := sanitizeAvdID(newName)
	if newID == "" {
		return fmt.Errorf("invalid AVD name")
	}

	avdDir, err := helper.GetAvdDirectory()
	if err != nil {
		return err
	}

	oldIniPath := filepath.Join(avdDir, oldName+".ini")
	newIniPath := filepath.Join(avdDir, newID+".ini")

	// 1. Check if new AVD ID already exists (and it's not the same as old)
	if oldName != newID {
		if _, err := os.Stat(newIniPath); err == nil {
			return fmt.Errorf("an AVD with ID '%s' already exists", newID)
		}
	}

	// 2. Identify the current AVD folder path from the .ini file
	oldAvdInfo, err := a.GetAvdInfo(oldName)
	if err != nil {
		return fmt.Errorf("could not find current AVD info: %w", err)
	}
	oldAvdPath := oldAvdInfo.Path
	newAvdPath := filepath.Join(avdDir, newID+".avd")

	// 3. Rename the .avd folder
	if oldAvdPath != newAvdPath {
		if _, err := os.Stat(newAvdPath); err == nil {
			return fmt.Errorf("destination folder '%s' already exists", newID+".avd")
		}
		if err := os.Rename(oldAvdPath, newAvdPath); err != nil {
			return fmt.Errorf("failed to rename AVD folder: %w", err)
		}
	}

	// 4. Rename and Update the root .ini file
	if oldName != newID {
		iniData, err := os.ReadFile(oldIniPath)
		if err != nil {
			return fmt.Errorf("failed to read .ini file: %w", err)
		}

		lines := strings.Split(string(iniData), "\n")
		var newLines []string
		foundPath := false
		foundPathRel := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				newLines = append(newLines, trimmed)
				continue
			}

			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) < 2 {
				newLines = append(newLines, trimmed)
				continue
			}

			key := strings.TrimSpace(parts[0])
			if key == "path" {
				if !foundPath {
					newLines = append(newLines, "path="+newAvdPath)
					foundPath = true
				}
			} else if key == "path.rel" {
				if !foundPathRel {
					newLines = append(newLines, "path.rel=avd\\"+newID+".avd")
					foundPathRel = true
				}
			} else {
				newLines = append(newLines, trimmed)
			}
		}

		if err := os.WriteFile(newIniPath, []byte(strings.Join(newLines, "\r\n")+"\r\n"), 0644); err != nil {
			return fmt.Errorf("failed to create new .ini file: %w", err)
		}

		// Delete old .ini file
		if err := os.Remove(oldIniPath); err != nil {
			fmt.Printf("Warning: failed to remove old .ini file '%s': %v\n", oldIniPath, err)
		}
	}

	// 5. Update config.ini
	configPath := filepath.Join(newAvdPath, "config.ini")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config.ini: %w", err)
	}

	lines := strings.Split(string(configData), "\n")
	var newConfigLines []string
	foundDisplayName := false
	foundLauncherName := false
	foundAvdId := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			newConfigLines = append(newConfigLines, trimmed)
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) < 2 {
			newConfigLines = append(newConfigLines, trimmed)
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "avd.ini.displayname" {
			if !foundDisplayName {
				newConfigLines = append(newConfigLines, "avd.ini.displayname="+stripEmojis(newName))
				foundDisplayName = true
			}
		} else if key == "avd.launcher.displayname" {
			if !foundLauncherName {
				newConfigLines = append(newConfigLines, "avd.launcher.displayname="+newName)
				foundLauncherName = true
			}
		} else if key == "AvdId" {
			if !foundAvdId {
				newConfigLines = append(newConfigLines, "AvdId="+newID)
				foundAvdId = true
			}
		} else {
			newConfigLines = append(newConfigLines, trimmed)
		}
	}

	if !foundDisplayName {
		newConfigLines = append(newConfigLines, "avd.ini.displayname="+stripEmojis(newName))
	}
	if !foundLauncherName {
		newConfigLines = append(newConfigLines, "avd.launcher.displayname="+newName)
	}
	if !foundAvdId {
		newConfigLines = append(newConfigLines, "AvdId="+newID)
	}

	if err := os.WriteFile(configPath, []byte(strings.Join(newConfigLines, "\r\n")+"\r\n"), 0644); err != nil {
		return fmt.Errorf("failed to update config.ini: %w", err)
	}

	fmt.Printf("Renamed AVD '%s' to '%s' (ID: %s)\n", oldName, newName, newID)
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

	var avd models.AvdInfo
	avd.Name = avdName
	avd.Path = avdPath

	// Read config.ini for more metadata
	configPath := filepath.Join(avdPath, "config.ini")
	if configFile, err := os.Open(configPath); err == nil {
		defer configFile.Close()
		configScanner := bufio.NewScanner(configFile)
		for configScanner.Scan() {
			line := strings.TrimSpace(configScanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) < 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			switch key {
			case "avd.launcher.displayname":
				avd.DisplayName = val
			case "avd.ini.displayname":
				if avd.DisplayName == "" {
					avd.DisplayName = val
				}
			case "abi.type":
				avd.Abi = val
			case "hw.ramSize":
				avd.RamSize = val
			case "PlayStore.enabled":
				if val == "true" || val == "yes" {
					avd.HasGooglePlay = true
				}
			case "target":
				if strings.HasPrefix(val, "android-") {
					avd.ApiLevel = strings.TrimPrefix(val, "android-")
				}
			case "image.sysdir.1":
				sysDir := val
				p := strings.FieldsFunc(sysDir, func(r rune) bool {
					return r == '/' || r == '\\'
				})
				for _, part := range p {
					if strings.HasPrefix(part, "android-") {
						avd.ApiLevel = strings.TrimPrefix(part, "android-")
						break
					}
				}
			case "hw.lcd.width":
				if avd.Resolution == "" {
					avd.Resolution = val
				} else {
					avd.Resolution = val + "x" + avd.Resolution
				}
			case "hw.lcd.height":
				if avd.Resolution == "" {
					avd.Resolution = val
				} else {
					avd.Resolution = avd.Resolution + "x" + val
				}
			}
		}
	}

	// Post-process metadata
	if avd.ApiLevel != "" {
		avd.AndroidVersion, avd.AndroidCodename = getAndroidVariantInfo(avd.ApiLevel)
	}

	if avd.DisplayName == "" {
		avd.DisplayName = avd.Name
	}

	// Check for lock folder to determine if AVD is currently running
	// ... (rest of the logic)
	lockPath := filepath.Join(avdPath, "hardware-qemu.ini.lock")
	isRunning := false

	if info, statErr := os.Stat(lockPath); statErr == nil && info.IsDir() {
		// ...
		pidFile := filepath.Join(lockPath, "pid")
		pidData, readErr := os.ReadFile(pidFile)
		if readErr == nil {
			pidStr := strings.TrimSpace(string(pidData))
			pid, parseErr := strconv.Atoi(pidStr)
			if parseErr == nil {
				if helper.IsProcessAlive(pid) {
					isRunning = true
				} else {
					_ = os.RemoveAll(lockPath)
				}
			} else {
				_ = os.RemoveAll(lockPath)
			}
		} else {
			_ = os.RemoveAll(lockPath)
		}
	}

	avd.Running = isRunning
	return avd, nil
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

func getAndroidVariantInfo(apiLevel string) (string, string) {
	// versionMap: API level -> [Commercial Version, Codename]
	versionMap := map[string]string{
		"35": "15|Vanilla Ice Cream",
		"34": "14|Upside Down Cake",
		"33": "13|Tiramisu",
		"32": "12L|S-V2",
		"31": "12|Snow Cone",
		"30": "11|Red Velvet Cake",
		"29": "10|Quince Tart",
		"28": "9.0|Pie",
		"27": "8.1|Oreo",
		"26": "8.0|Oreo",
		"25": "7.1|Nougat",
		"24": "7.0|Nougat",
		"23": "6.0|Marshmallow",
		"22": "5.1|Lollipop",
		"21": "5.0|Lollipop",
		"20": "4.4W|KitKat",
		"19": "4.4|KitKat",
		"18": "4.3|Jelly Bean",
		"17": "4.2|Jelly Bean",
		"16": "4.1|Jelly Bean",
	}

	if info, ok := versionMap[apiLevel]; ok {
		parts := strings.Split(info, "|")
		return "Android " + parts[0], parts[1]
	}
	return "Android " + apiLevel, ""
}
