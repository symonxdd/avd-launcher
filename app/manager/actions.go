package manager

import (
	"avd-launcher/app/helper"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (m *AvdManager) RenameAVD(oldName, newName string) error {
	newID := sanitizeAvdID(newName)
	if newID == "" {
		return fmt.Errorf("invalid AVD name")
	}

	avdDir, _ := helper.GetAvdDirectory()
	oldIniPath := filepath.Join(avdDir, oldName+".ini")
	newIniPath := filepath.Join(avdDir, newID+".ini")

	if oldName != newID {
		if _, err := os.Stat(newIniPath); err == nil {
			return fmt.Errorf("an AVD with ID '%s' already exists", newID)
		}
	}

	info, err := m.GetAvdInfo(oldName)
	if err != nil {
		return err
	}
	newAvdPath := filepath.Join(avdDir, newID+".avd")

	if info.Path != newAvdPath {
		if _, err := os.Stat(newAvdPath); err == nil {
			return fmt.Errorf("destination folder already exists")
		}
		if err := os.Rename(info.Path, newAvdPath); err != nil {
			return err
		}
	}

	if oldName != newID {
		if err := m.updateIniFile(oldIniPath, newIniPath, newAvdPath, newID); err != nil {
			return err
		}
		_ = os.Remove(oldIniPath)
	}

	return m.updateConfigFile(filepath.Join(newAvdPath, "config.ini"), newName, newID)
}

func (m *AvdManager) DeleteAVD(name string) error {
	path, _ := helper.GetAvdManagerPath()
	output, err := helper.NewCommand(path, "delete", "avd", "-n", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete: %w, %s", err, string(output))
	}
	return nil
}

func (m *AvdManager) OpenAvdFolder(path string) {
	_ = helper.NewCommand("explorer", path).Run()
}

// --- Internal Helpers ---

func sanitizeAvdID(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	return strings.Trim(regexp.MustCompile(`\s+`).ReplaceAllString(reg.ReplaceAllString(strings.TrimSpace(name), " "), "_"), "_")
}

func stripEmojis(s string) string {
	return strings.TrimSpace(regexp.MustCompile(`[^\x00-\x7F]+`).ReplaceAllString(s, ""))
}

func (m *AvdManager) updateIniFile(oldPath, newPath, avdPath, newID string) error {
	data, err := os.ReadFile(oldPath)
	if err != nil { return err }
	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "path=") { newLines = append(newLines, "path="+avdPath)
		} else if strings.HasPrefix(trimmed, "path.rel=") { newLines = append(newLines, "path.rel=avd\\"+newID+".avd")
		} else { newLines = append(newLines, trimmed) }
	}
	return os.WriteFile(newPath, []byte(strings.Join(newLines, "\r\n")+"\r\n"), 0644)
}

func (m *AvdManager) updateConfigFile(path, newName, newID string) error {
	data, err := os.ReadFile(path)
	if err != nil { return err }
	lines := strings.Split(string(data), "\n")
	var newConfigLines []string
	tags := map[string]string{"avd.ini.displayname": stripEmojis(newName), "avd.launcher.displayname": newName, "AvdId": newID}
	found := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if val, ok := tags[key]; ok {
				newConfigLines = append(newConfigLines, key+"="+val)
				found[key] = true
				continue
			}
		}
		newConfigLines = append(newConfigLines, trimmed)
	}
	for key, val := range tags {
		if !found[key] { newConfigLines = append(newConfigLines, key+"="+val) }
	}
	return os.WriteFile(path, []byte(strings.Join(newConfigLines, "\r\n")+"\r\n"), 0644)
}
