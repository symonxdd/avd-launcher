package models

import "os/exec"

// AVD represents an Android Virtual Device (AVD) with its process and other data
type AVD struct {
	Name    string
	Process *exec.Cmd
}
