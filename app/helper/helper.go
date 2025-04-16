package helper

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"time"
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

func TimestampedLog(s string) string {
	// 🧠 In Go, time.Format uses a specific reference time (Mon Jan 2 15:04:05 MST 2006) to define the layout — we need to pass an example time with the exact formatting we want.
	return fmt.Sprintf("(%s) %s", time.Now().Format("15:04:05"), s)
}
