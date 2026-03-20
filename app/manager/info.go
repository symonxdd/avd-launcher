package manager

import (
	"avd-launcher/app/models"
	"avd-launcher/app/helper"
	"fmt"
	"path/filepath"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func (m *AvdManager) GetAvdInfo(name string) (models.AvdInfo, error) {
	path, err := m.readAvdPathFromIni(name)
	if err != nil {
		return models.AvdInfo{}, err
	}

	avd := models.AvdInfo{Name: name, Path: path}
	m.parseConfigIni(path, &avd)

	if avd.ApiLevel != "" {
		avd.AndroidVersion, avd.AndroidCodename = getAndroidVersionInfo(avd.ApiLevel)
	}
	if avd.DisplayName == "" {
		avd.DisplayName = name
	}

	avd.Running = m.isAvdRunningOnDisk(path)
	return avd, nil
}

func (m *AvdManager) GetAvdDiskUsage(name string) (string, error) {
	path, err := m.readAvdPathFromIni(name)
	if err != nil {
		return "", err
	}
	size, err := helper.DirSize(path)
	if err != nil {
		return "", err
	}
	return helper.FormatSize(size), nil
}

// --- Internal Helpers ---

func (m *AvdManager) parseConfigIni(path string, avd *models.AvdInfo) {
	file, err := os.Open(filepath.Join(path, "config.ini"))
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		m.applyConfigParam(key, val, avd)
	}
}

func (m *AvdManager) applyConfigParam(key, val string, avd *models.AvdInfo) {
	switch key {
	case "avd.launcher.displayname", "avd.ini.displayname":
		if avd.DisplayName == "" || key == "avd.launcher.displayname" {
			avd.DisplayName = val
		}
	case "abi.type": avd.Abi = val
	case "hw.ramSize": avd.RamSize = val
	case "PlayStore.enabled": avd.HasGooglePlay = (val == "true" || val == "yes")
	case "target":
		if strings.HasPrefix(val, "android-") { avd.ApiLevel = strings.TrimPrefix(val, "android-") }
	case "image.sysdir.1":
		for _, part := range strings.FieldsFunc(val, func(r rune) bool { return r == '/' || r == '\\' }) {
			if strings.HasPrefix(part, "android-") { avd.ApiLevel = strings.TrimPrefix(part, "android-") }
		}
	case "hw.lcd.width": 
		if avd.Resolution == "" { avd.Resolution = val } else { avd.Resolution = val + "x" + avd.Resolution }
	case "hw.lcd.height":
		if avd.Resolution == "" { avd.Resolution = val } else { avd.Resolution = avd.Resolution + "x" + val }
	}
}

func (m *AvdManager) isAvdRunningOnDisk(path string) bool {
	lockPath := filepath.Join(path, "hardware-qemu.ini.lock")
	if info, err := os.Stat(lockPath); err == nil && info.IsDir() {
		if pidData, err := os.ReadFile(filepath.Join(lockPath, "pid")); err == nil {
			if pid, err := strconv.Atoi(strings.TrimSpace(string(pidData))); err == nil {
				return helper.IsProcessAlive(pid)
			}
		}
	}
	return false
}

func (m *AvdManager) readAvdPathFromIni(name string) (string, error) {
	avdDir, err := helper.GetAvdDirectory()
	if err != nil {
		return "", err
	}
	file, err := os.Open(filepath.Join(avdDir, name+".ini"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); strings.HasPrefix(line, "path=") {
			return strings.TrimPrefix(line, "path="), nil
		}
	}
	return "", fmt.Errorf("path not found in ini")
}
