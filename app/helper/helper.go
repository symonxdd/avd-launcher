package helper

import (
	"os/exec"
	"runtime"
	"syscall"
)

// Command is just an alias for exec.Cmd
type Command = exec.Cmd

// NewCommand creates an exec.Cmd and hides the terminal window on Windows
func NewCommand(name string, arg ...string) *Command {
	cmd := exec.Command(name, arg...)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	return cmd
}
