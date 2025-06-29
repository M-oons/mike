package info

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

const (
	AppName    = "mike"
	Author     = "Moons"
	Repository = "https://github.com/m-oons/mike"
)

var (
	Version = "dev"
	Commit  = "-"
	Date    = "-"
)

func OpenRepository() {
	cmd := exec.Command("cmd", "/c", "start", "", Repository)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	cmd.Start()
}

func VersionString() string {
	if Version == "dev" {
		return "dev"
	}

	return fmt.Sprintf("%s (%s)", Version, Commit)
}

func DateString() string {
	if Version == "dev" {
		return time.Now().Format("2006-01-02T15:04:05Z")
	}

	return Date
}
