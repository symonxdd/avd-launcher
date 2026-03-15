//go:build windows

package helper

import (
	"golang.org/x/sys/windows"
)

// IsProcessAlive checks whether a process with the given PID is still running.
// On Windows, os.FindProcess always returns a non-nil process, so we use
// OpenProcess with PROCESS_QUERY_LIMITED_INFORMATION to verify it truly exists.
func IsProcessAlive(pid int) bool {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	windows.CloseHandle(handle)
	return true
}
