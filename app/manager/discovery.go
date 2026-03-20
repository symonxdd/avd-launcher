package manager

import (
	"avd-launcher/app/helper"
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func (m *AvdManager) ListAVDs() ([]string, error) {
	avdDir, err := helper.GetAvdDirectory()
	if err == nil {
		if files, err := os.ReadDir(avdDir); err == nil {
			var avds []string
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".ini") {
					avds = append(avds, strings.TrimSuffix(file.Name(), ".ini"))
				}
			}
			if len(avds) > 0 {
				return avds, nil
			}
		}
	}

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
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			avds = append(avds, line)
		}
	}
	return avds, scanner.Err()
}

func (m *AvdManager) ListRunningAVDs() ([]string, error) {
	adbPath, err := helper.GetAdbPath()
	if err != nil {
		return nil, err
	}

	output, err := helper.NewCommand(adbPath, "devices").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run adb devices: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var runningAVDs []string

	for _, line := range lines {
		if strings.HasPrefix(line, "emulator-") && strings.Contains(line, "device") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				if name, err := m.getAvdNameForSerial(adbPath, parts[0]); err == nil && name != "" {
					runningAVDs = append(runningAVDs, name)
				}
			}
		}
	}
	return runningAVDs, nil
}

func (m *AvdManager) getAvdNameForSerial(adbPath, serial string) (string, error) {
	out, err := helper.NewCommand(adbPath, "-s", serial, "emu", "avd", "name").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(strings.TrimSpace(string(out)), "\n")[0]), nil
}
