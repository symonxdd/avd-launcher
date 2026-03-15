//go:build !windows

package helper

import (
	"os"
	"syscall"
)

// IsProcessAlive checks whether a process with the given PID is still running.
// On Unix-like systems, sending signal 0 checks for process existence without side effects.
func IsProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}
