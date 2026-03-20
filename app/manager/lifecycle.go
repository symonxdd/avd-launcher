package manager

import (
	"avd-launcher/app/helper"
	"avd-launcher/app/models"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (m *AvdManager) StartAVD(avdName string, coldBoot bool) string {
	emulatorPath, err := helper.GetEmulatorPath()
	if err != nil {
		return "Failed to find emulator: " + err.Error()
	}

	avdName = strings.TrimSpace(avdName)
	cmd, stdout, stderr, err := m.launchEmulator(emulatorPath, avdName, coldBoot)
	
	if err != nil && strings.Contains(err.Error(), "already running") {
		m.cleanupLockFiles(avdName)
		cmd, stdout, stderr, err = m.launchEmulator(emulatorPath, avdName, coldBoot)
	}

	if err != nil {
		return "Failed to start emulator: " + err.Error()
	}

	m.setRunningAvd(avdName, &models.AVD{Name: avdName, Process: cmd})

	go m.streamLogs(avdName, stdout, stderr)

	return "Emulator started"
}

func (m *AvdManager) StopAVD(name string) error {
	if !m.isAvdRunningInProcess(name) {
		// Even if not in our map, let's try to resolve the port to be sure
	}

	port, err := helper.ResolvePortForAVD(name)
	if err != nil {
		return err
	}

	adbPath, err := helper.GetAdbPath()
	if err != nil {
		return err
	}

	if output, err := helper.NewCommand(adbPath, "-s", fmt.Sprintf("emulator-%d", port), "emu", "kill").CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop AVD '%s': %w, output: %s", name, err, string(output))
	}

	m.deleteRunningAvd(name)
	return nil
}

func (m *AvdManager) launchEmulator(path, name string, cold bool) (*helper.Command, io.ReadCloser, io.ReadCloser, error) {
	args := []string{"-avd", name}
	if cold {
		args = append(args, "-no-snapshot-load")
	}

	cmd := helper.NewCommand(path, args...)
	cmd.Env = os.Environ()
	cmd.Dir = filepath.Dir(path)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	return cmd, stdout, stderr, cmd.Start()
}

func (m *AvdManager) cleanupLockFiles(name string) {
	lockPath := filepath.Join(os.Getenv("USERPROFILE"), ".android", "avd", name+".avd")
	if files, err := filepath.Glob(filepath.Join(lockPath, "*.lock")); err == nil {
		for _, file := range files {
			_ = os.Remove(file)
		}
	}
}

func (m *AvdManager) streamLogs(name string, stdout, stderr io.ReadCloser) {
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		line := scanner.Text()
		runtime.EventsEmit(m.ctx, "avd-log", helper.TimestampedLog(line))

		if strings.Contains(line, "Successfully loaded snapshot") || strings.Contains(line, "Boot completed") {
			runtime.EventsEmit(m.ctx, "avd-booted", name)
		}
		if strings.Contains(line, "killing emulator, bye bye") {
			runtime.EventsEmit(m.ctx, "avd-shutdown", name)
			m.deleteRunningAvd(name)
			break
		}
	}
}
