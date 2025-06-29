package info

import (
	"fmt"
	"os/exec"
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
	exec.Command("cmd", "/c", "start", Repository).Start()
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
